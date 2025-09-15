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

type HandlerContextID uint

const (
	SessionContextKey HandlerContextID = iota + 1
	HandlerContextKey
)

type HandlerCtx struct {
	ValidSession *database.BrowserSession
	Database     database.Database
}

type Handler struct {
	Database         *database.Manager
	config           *config.Config
	emailClient      atomic.Pointer[email.Client]
	myNautiqueClient atomic.Pointer[mynautique.Client]
}

type HandlerParams struct {
	Database      *database.Manager
	Configuration *config.Config
}

func NewHandler(params HandlerParams) *Handler {
	h := &Handler{
		config:   params.Configuration,
		Database: params.Database,
	}
	return h
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

func (h *Handler) WithAdminAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return h.WithAuthentication(next, database.UserRoleAdmin)
}

func (h *Handler) WithAnyAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return h.WithAuthentication(next, database.UserRoleUnknown)
}

func (h *Handler) WithNoAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbh, done := h.Database.GetHandler()
		defer done()
		if dbh == nil {
			slog.Warn("No database connection is available.")
			w.WriteHeader(http.StatusInternalServerError)
			WriteFailureResponse("Operation cannot be performed. Database connection is not established properly.", w)
			return
		}
		ctxWithDatabase := context.WithValue(r.Context(), HandlerContextKey, &HandlerCtx{
			Database: dbh,
		})
		next.ServeHTTP(w, r.WithContext(ctxWithDatabase))
	})
}

// Authentication middleware that gets the user session and calls the next Handler.
// It makes sure the user is authenticated and adds this information ready for the next call.
func (h *Handler) WithAuthentication(next http.HandlerFunc, requiredRole database.UserRoleType) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get user from cookie
		cookie, err := r.Cookie("SESSION")
		if err != nil {
			slog.Warn("Cannot read cookie", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			WriteFailureResponse("not authenticated", w)
			return
		}
		dbh, done := h.Database.GetHandler()
		defer done()
		if dbh == nil {
			slog.Warn("No database connection is available.")
			w.WriteHeader(http.StatusInternalServerError)
			WriteFailureResponse("Operation cannot be performed. Database connection is not established properly.", w)
			return
		}
		session, err := dbh.GetBrowserSession(cookie.Value)
		if err != nil || session == nil {
			slog.Warn("No session found for given cookie", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			WriteFailureResponse("not authenticated", w)
			return
		}
		if !session.Valid() {
			slog.Info("Session expired", slog.String("expiry", session.ValidUntil.String()), slog.String("sessionSecret", session.SessionSecret))
			w.WriteHeader(http.StatusUnauthorized)
			WriteFailureResponse("session expired or invalid", w)
			return
		}
		// Check for proper authentication in case we require a specific role.
		if requiredRole != database.UserRoleUnknown {
			if session.UserRoleID != requiredRole {
				slog.Warn("User does not have the required role", slog.String("user", session.Username), slog.Uint64("role", uint64(session.UserRoleID)), slog.Uint64("required_role", uint64(requiredRole)))
				w.WriteHeader(http.StatusUnauthorized)
				WriteFailureResponse("Operation cannot be performed. Database connection is not established properly.", w)
				return
			}
		}

		// Add session information to context
		ctxWithSession := context.WithValue(r.Context(), SessionContextKey, *session)
		ctxWithSessionAndHandler := context.WithValue(ctxWithSession, HandlerContextKey, &HandlerCtx{
			ValidSession: session,
			Database:     dbh,
		})

		// ctx := r.WithContext(context.WithValue(r.Context(), SessionContextKey, *session))

		next.ServeHTTP(w, r.WithContext(ctxWithSessionAndHandler))
	})
}

// Next is a generic handler function that gets called by the wrapper.
// It receives the parsed and validated request body of type T and the handler context.
type Next[T any] func(w http.ResponseWriter, r *http.Request, body T, hCtx *HandlerCtx)

// WithBody is a generic handler wrapper that handles request body parsing, validation,
// and error handling.
func WithRequestBody[T any](next Next[T], validationErrorMessages ...map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get handler context
		hCtx, err := GetHandlerContext(w, r)
		if err != nil {
			slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			WriteFailureResponse("Internal error.", w)
			return
		}

		// Read and validate request body
		var body T
		if validationErrorMessages == nil {
			err = ReadBodyAndValidate(r, &body)
		} else {
			err = ReadBodyAndValidate(r, &body, validationErrorMessages[0])
		}
		if err != nil {
			slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
			WriteFailureResponse(err.Error(), w)
			return
		}

		// Call the actual handler
		next(w, r, body, hCtx)
	}
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

// GetHandlerContext returns the HandlerCtx in case it exists or an error in case it does not.
// It directly sets a FailureResponse and sets the HTTP Status code to 500.
func GetHandlerContext(w http.ResponseWriter, r *http.Request) (*HandlerCtx, error) {
	if r.Context().Value(HandlerContextKey) == nil {
		slog.Warn("No handler context is available.")
		w.WriteHeader(http.StatusInternalServerError)
		WriteFailureResponse("Internal error.", w)
		return nil, fmt.Errorf("no handler context is available")
	}
	return r.Context().Value(HandlerContextKey).(*HandlerCtx), nil
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
