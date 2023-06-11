package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"server/database"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

type HandlerContext uint

const SessionContextKey HandlerContext = 1

type Handler struct {
	db database.Database
}

type IdentificationMiddleware struct {
	db database.Database
}

func NewHandler(db database.Database) *Handler {
	return &Handler{
		db: db,
	}
}

func NewIdenticationMiddleware(db database.Database) *IdentificationMiddleware {
	return &IdentificationMiddleware{
		db: db,
	}
}

// Authentication middleware that gets the user session and calls the next Handler.
// It makes sure the user is authenticated and adds this information ready for the next call.
func (i *IdentificationMiddleware) Authenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get user from cookie
		cookie, err := r.Cookie("SESSION")
		if err != nil {
			slog.Warn("Cannot read cookie", slog.String("error", err.Error()))
			WriteFailureResponse("not authenticated", w)
			return
		}

		session, err := i.db.GetBrowserSession(cookie.Value)
		if err != nil || session == nil {
			slog.Warn("No session found for given cookie", slog.String("error", err.Error()))
			WriteFailureResponse("not authenticated", w)
			return
		}

		if session.Expired() {
			slog.Info("Session expired", slog.String("expiry", session.ValidUntil.String()), slog.String("sessionSecret", session.SessionSecret))
			WriteFailureResponse("not authenticated", w)
			return
		}

		// Add session information to context
		ctx := r.WithContext(context.WithValue(r.Context(), SessionContextKey, *session))

		next.ServeHTTP(w, ctx)
	})
}

func GetSessionFromContext(r *http.Request) database.BrowserSession {
	return r.Context().Value(SessionContextKey).(database.BrowserSession)
}

func ReadBodyAndValidate(r *http.Request, s any) error {
	// Get the content of the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		slog.Warn("Cannot read the request body", slog.String("error", err.Error()))
		return err
	}

	// Parse the content into the provided struct
	err = json.Unmarshal(body, s)
	if err != nil {
		slog.Warn("Cannot parse the request body JSON", slog.String("error", err.Error()))
		return err
	}

	// Validate the struct
	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		slog.Warn("Struct does not validate", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func WriteFailureResponse(message string, w http.ResponseWriter) {
	// Get a JSON representation of the failure response
	response := FailureResponse(message)
	j, err := json.Marshal(&response)
	if err != nil {
		slog.Error("Cannot create a FailureResponse", slog.String("error", err.Error()))
		http.Error(w, "failed to create failure response", http.StatusInternalServerError)
	}
	w.Header().Set("Conntent-Type", "application/json")
	io.Copy(w, bytes.NewReader(j))
}

// AuthenticatedOrFailure writes a failure responsne if the session is not
// authenticated. It returns an error in case the session is not authenticated.
func AuthenticatedOrFailure(b database.BrowserSession, w http.ResponseWriter) error {
	if !b.Valid() {
		slog.Warn("Unauthenticated action is not allowed")
		WriteFailureResponse("unauthenticated action not allowed", w)
		return fmt.Errorf("action is not authenticated")
	}
	return nil
}

func WriteSuccessResponse(message string, data any, w http.ResponseWriter) {
	resp := Response{
		Status: Status{
			OK:  true,
			Msg: message,
		},
		Data: data,
	}

	j, err := json.Marshal(&resp)
	if err != nil {
		slog.Error("Cannot marshal the response", slog.String("error", err.Error()))
		http.Error(w, "creating response failed", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, bytes.NewReader(j))
}

func SetSessionCookie(w http.ResponseWriter, s string) {
	cookie := http.Cookie{
		Name:     "SESSION",
		Value:    s,
		Expires:  time.Now().Add(time.Second * time.Duration(viper.GetInt64("http.sessiontimeout"))),
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
}

func DeleteSessionCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:    "SESSION",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
