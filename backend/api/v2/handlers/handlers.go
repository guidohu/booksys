// Package handlers implements the HTTP handlers of the v2 API together with
// the middleware that provides them with authentication, a database handle
// and a validated request body.
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/go-playground/validator/v10"
	"server/config"
	"server/database"
	"server/mynautique"
	"server/notifications/email"
	customvalidator "server/validator"
)

// HandlerContextID identifies the values this package stores in a request
// context. It is a distinct type so that the keys cannot collide with keys
// from other packages.
type HandlerContextID uint

// The keys under which the middleware stores its values in a request context.
const (
	SessionContextKey HandlerContextID = iota + 1
	HandlerContextKey
)

// HandlerCtx carries everything the middleware prepared for a handler: the
// authenticated session, a database handle and, where requested, the
// application configuration.
type HandlerCtx struct {
	ValidSession *database.BrowserSession
	Database     database.Database
	Config       *config.Config
}

// Handler holds the dependencies shared by all HTTP handlers.
type Handler struct {
	Database         *database.Manager
	config           *config.Config
	emailClient      atomic.Pointer[email.Client]
	myNautiqueClient atomic.Pointer[mynautique.Client]
}

// HandlerParams are the dependencies NewHandler needs.
type HandlerParams struct {
	Database      *database.Manager
	Configuration *config.Config
}

// NewHandler returns a Handler that uses the given dependencies.
func NewHandler(params HandlerParams) *Handler {
	h := &Handler{
		config:   params.Configuration,
		Database: params.Database,
	}
	return h
}

// SetEmailClient replaces the email client used to send notifications.
func (h *Handler) SetEmailClient(e email.Client) {
	h.emailClient.Store(&e)
}

// GetEmailClient returns the currently configured email client.
func (h *Handler) GetEmailClient() email.Client {
	return *(h.emailClient.Load())
}

// SetMyNautiqueClient replaces the client used to talk to the MyNautique API.
func (h *Handler) SetMyNautiqueClient(m mynautique.Client) {
	h.myNautiqueClient.Store(&m)
}

// GetMyNautiqueClient returns the current MyNautique client, or nil if none
// has been set.
func (h *Handler) GetMyNautiqueClient() *mynautique.Client {
	return h.myNautiqueClient.Load()
}

// MaxRequestBodyBytes is how much of a request body a handler may be made to
// read. Handlers read the body with io.ReadAll, and multipart uploads let
// net/http buffer and spill to disk, so without a cap a single unauthenticated
// caller can make the server hold arbitrarily much for as long as the read
// timeout allows.
//
// The largest legitimate request is the logo upload, which the handler itself
// caps at 512kB, so 1MiB leaves room for the multipart framing around it while
// staying three orders of magnitude below what an attacker could otherwise
// send.
const MaxRequestBodyBytes = 1 << 20

// WithMaxBodySize is middleware that caps the request body. It is applied once
// around the whole mux rather than per handler, so that routes registered
// without any other middleware, and routes added later, are covered too.
func WithMaxBodySize(next http.Handler, limit int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, limit)
		}
		next.ServeHTTP(w, r)
	})
}

// WithFlagGuarded is middleware that rejects the request unless web setup is
// enabled in the configuration.
func (h *Handler) WithFlagGuarded(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.config.GetBool("http.websetup") {
			slog.Warn("Web setup is not enabled preventing access to.", slog.String("path", r.URL.Path))
			w.WriteHeader(http.StatusForbidden)
			WriteFailureResponse("Web setup is not enabled.", w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// WithConfigContext is middleware that adds the application configuration to
// the handler context. Only handlers that validate configuration dependent
// payloads need it.
func (h *Handler) WithConfigContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hCtx := GetHandlerContext(r)
		hCtx.Config = h.config
		ctxWithConfig := context.WithValue(r.Context(), HandlerContextKey, hCtx)
		next.ServeHTTP(w, r.WithContext(ctxWithConfig))
	})
}

// WithAdminAuthentication is middleware that only lets administrators through.
func (h *Handler) WithAdminAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return h.WithAuthentication(next, database.UserRoleAdmin)
}

// WithAnyAuthentication is middleware that lets any authenticated user
// through, regardless of their role.
func (h *Handler) WithAnyAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return h.WithAuthentication(next, database.UserRoleUnknown)
}

// WithNoAuthentication is middleware for endpoints that unauthenticated users
// may call. It still provides the handler with a database handle.
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

// WithAuthentication is middleware that gets the user session and calls the
// next handler. It makes sure the user is authenticated and adds this information ready for the next call.
func (h *Handler) WithAuthentication(next http.HandlerFunc, requiredRole database.UserRoleType) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get user from cookie
		cookie, err := r.Cookie(h.SessionCookieName())
		if err != nil {
			slog.Warn("Cannot read cookie", slog.Any("error", err))
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
		if err != nil {
			slog.Warn("Failed to retrieve session for given cookie", slog.Any("error", err))
			w.WriteHeader(http.StatusUnauthorized)
			WriteFailureResponse("not authenticated", w)
			return
		}
		if session == nil {
			slog.Warn("No session found for given cookie")
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
		next.ServeHTTP(w, r.WithContext(ctxWithSessionAndHandler))
	})
}

// Next is a generic handler function that gets called by the wrapper.
// It receives the parsed and validated request body of type T and the handler context.
type Next[T any] func(w http.ResponseWriter, r *http.Request, body T, hCtx *HandlerCtx)

// WithRequestBody is a generic handler wrapper that handles request body parsing, validation,
// and error handling.
func WithRequestBody[T any](next Next[T], validationErrorMessages ...map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		// Read and validate request body
		var body T
		if validationErrorMessages == nil {
			err = ReadBodyAndValidate(r, &body)
		} else {
			err = ReadBodyAndValidate(r, &body, validationErrorMessages[0])
		}
		if err != nil {
			var tooLarge *http.MaxBytesError
			if errors.As(err, &tooLarge) {
				slog.Warn("Request body exceeds the limit", slog.Int64("limit_bytes", tooLarge.Limit), slog.String("path", r.URL.Path))
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				WriteFailureResponse("Request body is too large.", w)
				return
			}
			slog.Warn("Request payload is not valid", slog.Any("error", err))
			WriteFailureResponse(err.Error(), w)
			return
		}

		// Call the actual handler
		next(w, r, body, GetHandlerContext(r))
	}
}

// WithLogging is middleware that logs the method and URL of every request
// after it has been served.
func (h *Handler) WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		slog.Info("Request served", slog.String("method", r.Method), slog.String("url", r.URL.String()))
	})
}

// GetSessionFromContext returns the browser session stored in the request
// context, or the zero session if the request is not authenticated.
func GetSessionFromContext(r *http.Request) database.BrowserSession {
	if r.Context().Value(SessionContextKey) == nil {
		return database.BrowserSession{}
	}
	return r.Context().Value(SessionContextKey).(database.BrowserSession)
}

// GetHandlerContext returns the HandlerCtx, or an empty one if none exists.
func GetHandlerContext(r *http.Request) *HandlerCtx {
	if r.Context().Value(HandlerContextKey) == nil {
		return &HandlerCtx{}
	}
	return r.Context().Value(HandlerContextKey).(*HandlerCtx)
}

// ReadBodyAndValidate parses the JSON request body into s and validates it
// against its `validate` struct tags. If an errorMap is given, the message it
// holds for the first failing field is returned instead of the raw validation
// error.
func ReadBodyAndValidate(r *http.Request, s any, errorMap ...map[string]string) error {
	hCtx := GetHandlerContext(r)
	// Get the content of the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Warn("Cannot read the request body", slog.Any("error", err))
		return err
	}

	// Parse the content into the provided struct
	err = json.Unmarshal(body, s)
	if err != nil {
		slog.Warn("Cannot parse the request body JSON", slog.Any("error", err))
		return err
	}

	// Validate the struct
	validate := validator.New()
	// Add custom validators
	_ = validate.RegisterValidation("googlemapsurl", customvalidator.GoogleMapsURL)
	_ = validate.RegisterValidation("recaptchakey", customvalidator.RecaptchaKey)
	_ = validate.RegisterValidation("strongpassword", customvalidator.PasswordStrength)
	_ = validate.RegisterValidation("sessiontype", customvalidator.SessionType)
	_ = validate.RegisterValidation("expensetype", customvalidator.ExpenseType)
	_ = validate.RegisterValidation("tableid", customvalidator.TableID)
	_ = validate.RegisterValidation("currency", customvalidator.Currency)
	_ = validate.RegisterValidation("uploadfile", func(fl validator.FieldLevel) bool {
		if hCtx.Config == nil {
			slog.Error("uploadfile validator called without configuration access.")
			return false
		}
		path, _ := hCtx.Config.GetString("http.uploadpath")
		return customvalidator.UploadedFilePath(fl, path)
	})

	err = validate.Struct(s)
	if err != nil {
		slog.Warn("Struct does not validate", slog.Any("error", err))
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

// WriteFailureResponse writes a JSON failure response with the given message.
// It does not set a status code, the caller does that.
func WriteFailureResponse(message string, w http.ResponseWriter) {
	// Get a JSON representation of the failure response
	response := FailureResponse(message)
	j, err := json.Marshal(&response)
	if err != nil {
		slog.Error("Cannot create a FailureResponse", slog.Any("error", err))
		http.Error(w, "failed to create failure response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = io.Copy(w, bytes.NewReader(j))
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
		slog.Error("Cannot marshal the response", slog.Any("error", err))
		http.Error(w, "creating response failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = io.Copy(w, bytes.NewReader(j))
}

// The names of the session cookie. A cookie whose name carries the `__Host-`
// prefix is only accepted by a browser if it was set over HTTPS, without a
// `Domain` attribute and with `Path=/`, which keeps another host under the
// same registrable domain (or a network attacker on a plain HTTP sibling)
// from planting a session cookie for this application.
//
// The prefix requires the `Secure` attribute, so a deployment that serves
// plain HTTP has to fall back to the unprefixed name.
const (
	secureSessionCookieName   = "__Host-SESSION"
	insecureSessionCookieName = "SESSION"
)

// useSecureCookies reports whether the session cookie is handed out with the
// `Secure` attribute. It is only turned off to develop against a server that
// serves plain HTTP, so anything but an explicit `false` keeps the cookie
// secure.
func (h *Handler) useSecureCookies() bool {
	value, source := h.config.GetString("http.securecookie")
	if source == config.Unknown {
		return true
	}
	return value != "false"
}

// sessionCookie returns the session cookie of the configured cookie mode,
// without a value. Setting and expiring the cookie both start from here: a
// browser only replaces a cookie it already holds if the name, the domain and
// the path of the two match.
func (h *Handler) sessionCookie() http.Cookie {
	secure := h.useSecureCookies()
	name := insecureSessionCookieName
	if secure {
		name = secureSessionCookieName
	}
	return http.Cookie{
		Name:     name,
		Path:     "/",
		Domain:   "",
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// SessionCookieName returns the name the session cookie is stored under for
// the configured cookie mode.
func (h *Handler) SessionCookieName() string {
	return h.sessionCookie().Name
}

// SetSessionCookie sets the session cookie to the given secret, expiring at
// validUntil.
func (h *Handler) SetSessionCookie(w http.ResponseWriter, s string, validUntil time.Time) {
	cookie := h.sessionCookie()
	cookie.Value = s
	cookie.Expires = validUntil
	http.SetCookie(w, &cookie)
}

// DeleteSessionCookie expires the session cookie on the client.
func (h *Handler) DeleteSessionCookie(w http.ResponseWriter) {
	cookie := h.sessionCookie()
	cookie.Expires = time.Unix(0, 0)
	http.SetCookie(w, &cookie)

	if cookie.Secure {
		// A session that was handed out before the cookie was made secure
		// sits under the unprefixed name and without the `Secure` attribute,
		// so the browser would keep sending it, in the clear, until it
		// expires. Expire it here, where the client is told that it is not
		// logged in anyway.
		legacy := cookie
		legacy.Name = insecureSessionCookieName
		http.SetCookie(w, &legacy)
	}
}
