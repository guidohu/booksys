package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"server/config"
	"server/database"
	"server/mynautique"
	"server/notifications/email"
	customvalidator "server/validator"
	"sync/atomic"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type HandlerContext uint

const SessionContextKey HandlerContext = 1

type Handler struct {
	db               atomic.Pointer[database.Database]
	config           *config.Config
	emailClient      atomic.Pointer[email.Client]
	myNautiqueClient atomic.Pointer[mynautique.Client]
}

type HandlerParams struct {
	Database      database.Database
	Configuration *config.Config
}

func NewHandler(params HandlerParams) *Handler {
	h := &Handler{
		config: params.Configuration,
	}
	h.db.Store(&params.Database)
	return h
}

func (h *Handler) SetDB(db database.Database) {
	h.db.Store(&db)
}

func (h *Handler) GetDB() database.Database {
	return *(h.db.Load())
}

func (h *Handler) SetEmailClient(e email.Client) {
	h.emailClient.Store(&e)
}

func (h *Handler) GetEmailClient() email.Client {
	// client := *(h.emailClient.Load())
	// if client == nil {
	// 	// Try to create a client
	// 	h.GetDB().GetMyNautiqueConfiguration()
	// }
	return *(h.emailClient.Load())
}

func (h *Handler) SetMyNautiqueClient(m mynautique.Client) {
	h.myNautiqueClient.Store(&m)
}

func (h *Handler) GetMyNautiqueClient() *mynautique.Client {
	return h.myNautiqueClient.Load()
}

// Authentication middleware that gets the user session and calls the next Handler.
// It makes sure the user is authenticated and adds this information ready for the next call.
func (h *Handler) WithAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get user from cookie
		cookie, err := r.Cookie("SESSION")
		if err != nil {
			slog.Warn("Cannot read cookie", slog.String("error", err.Error()))
			WriteFailureResponse("not authenticated", w)
			return
		}

		dbh := h.GetDB()
		session, err := dbh.GetBrowserSession(cookie.Value)
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

func (h *Handler) WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		fmt.Println(time.Now(), r.Method, r.URL)
	})
}

func GetSessionFromContext(r *http.Request) database.BrowserSession {
	if r.Context().Value(SessionContextKey) == nil {
		return database.BrowserSession{}
	}
	return r.Context().Value(SessionContextKey).(database.BrowserSession)
}

func ReadBodyAndValidate(r *http.Request, s any, errorMap ...map[string]string) error {
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
	// Add custom validators
	validate.RegisterValidation("googlemapsurl", customvalidator.GoogleMapsURL)
	validate.RegisterValidation("recaptchakey", customvalidator.RecaptchaKey)
	validate.RegisterValidation("strongpassword", customvalidator.PasswordStrength)
	validate.RegisterValidation("sessiontype", customvalidator.SessionType)
	validate.RegisterValidation("expensetype", customvalidator.ExpenseType)
	validate.RegisterValidation("tableid", customvalidator.TableID)
	err = validate.Struct(s)
	if err != nil {
		slog.Warn("Struct does not validate", slog.String("error", err.Error()))
		errs := err.(validator.ValidationErrors)

		// return plain error if no error mapping is provided
		if len(errorMap) == 0 {
			return err
		}
		// return the custom error if present
		for _, e := range errs {
			msg, found := errorMap[0][e.StructField()]
			if found {
				return fmt.Errorf("%s", msg)
			}
		}
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

func AuthenticatedAsAdminOrFailure(b database.BrowserSession, w http.ResponseWriter) error {
	if !b.Valid() || b.UserRoleID != database.UserRoleAdmin {
		slog.Warn("Non admin role not allowed to call as admin")
		WriteFailureResponse("unauthenticated action not allowed", w)
		return fmt.Errorf("admin action is not authorized")
	}
	return nil
}

// WriteSuccessResponse writes a success message that contains data.
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

func SetSessionCookie(w http.ResponseWriter, s string, validUntil time.Time) {
	cookie := http.Cookie{
		Name:     "SESSION",
		Value:    s,
		Expires:  validUntil,
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
