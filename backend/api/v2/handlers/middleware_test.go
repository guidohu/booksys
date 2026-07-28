package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"server/database"
	"server/database/dbtest"
	"server/mynautique"
	"server/notifications/email"
)

func TestWithFlagGuarded(t *testing.T) {
	tests := []struct {
		name       string
		webSetup   string
		wantCalled bool
		wantStatus int
	}{
		{name: "web setup enabled", webSetup: "true", wantCalled: true, wantStatus: http.StatusOK},
		{name: "web setup disabled", webSetup: "false", wantCalled: false, wantStatus: http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHandler(t, &dbtest.FakeDB{}, map[string]string{"http.websetup": tt.webSetup})
			called := false
			guarded := h.WithFlagGuarded(func(w http.ResponseWriter, r *http.Request) {
				called = true
			})

			rec := httptest.NewRecorder()
			guarded.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v2/database/setup", nil))

			if called != tt.wantCalled {
				t.Errorf("next handler called = %v, want %v", called, tt.wantCalled)
			}
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestWithConfigContext(t *testing.T) {
	h := newTestHandler(t, &dbtest.FakeDB{}, nil)
	var got *HandlerCtx
	next := h.WithConfigContext(func(w http.ResponseWriter, r *http.Request) {
		got = GetHandlerContext(r)
	})

	next.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", nil))

	if got == nil {
		t.Fatal("handler context is nil")
	}
	if got.Config == nil {
		t.Error("config was not added to the handler context")
	}
}

func TestWithNoAuthentication(t *testing.T) {
	t.Run("provides a database handle", func(t *testing.T) {
		db := &dbtest.FakeDB{}
		h := newTestHandler(t, db, nil)
		var got *HandlerCtx
		next := h.WithNoAuthentication(func(w http.ResponseWriter, r *http.Request) {
			got = GetHandlerContext(r)
		})

		next.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", nil))

		if got == nil || got.Database == nil {
			t.Fatal("no database handle was added to the handler context")
		}
		if got.ValidSession != nil {
			t.Error("an unauthenticated request must not carry a session")
		}
	})

	t.Run("fails without a database connection", func(t *testing.T) {
		h := NewHandler(HandlerParams{
			Database:      database.NewManager(database.Settings{}),
			Configuration: newTestConfig(t, nil),
		})
		called := false
		next := h.WithNoAuthentication(func(w http.ResponseWriter, r *http.Request) { called = true })

		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", nil))

		if called {
			t.Error("next handler must not be called without a database")
		}
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})
}

func TestWithAuthentication(t *testing.T) {
	expiredSession := &database.BrowserSession{
		SessionSecret: "expired",
		ValidUntil:    time.Now().Add(-time.Hour),
		MaxValidUntil: time.Now().Add(time.Hour),
		UserID:        7,
		UserRoleID:    database.UserRoleAdmin,
	}

	tests := []struct {
		name         string
		cookie       *http.Cookie
		session      *database.BrowserSession
		sessionErr   error
		requiredRole database.UserRoleType
		wantCalled   bool
		wantStatus   int
	}{
		{
			name:         "valid admin session",
			cookie:       &http.Cookie{Name: "SESSION", Value: "secret"},
			session:      validSession(7, database.UserRoleAdmin),
			requiredRole: database.UserRoleAdmin,
			wantCalled:   true,
			wantStatus:   http.StatusOK,
		},
		{
			name:         "any role accepted",
			cookie:       &http.Cookie{Name: "SESSION", Value: "secret"},
			session:      validSession(7, database.UserRoleGuest),
			requiredRole: database.UserRoleUnknown,
			wantCalled:   true,
			wantStatus:   http.StatusOK,
		},
		{
			name:         "insufficient role",
			cookie:       &http.Cookie{Name: "SESSION", Value: "secret"},
			session:      validSession(7, database.UserRoleMember),
			requiredRole: database.UserRoleAdmin,
			wantCalled:   false,
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:         "no cookie",
			cookie:       nil,
			session:      validSession(7, database.UserRoleAdmin),
			requiredRole: database.UserRoleAdmin,
			wantCalled:   false,
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:         "unknown session",
			cookie:       &http.Cookie{Name: "SESSION", Value: "secret"},
			sessionErr:   errNotFound,
			requiredRole: database.UserRoleAdmin,
			wantCalled:   false,
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:         "session not found without error",
			cookie:       &http.Cookie{Name: "SESSION", Value: "secret"},
			session:      nil,
			sessionErr:   nil,
			requiredRole: database.UserRoleAdmin,
			wantCalled:   false,
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:         "expired session",
			cookie:       &http.Cookie{Name: "SESSION", Value: "secret"},
			session:      expiredSession,
			requiredRole: database.UserRoleAdmin,
			wantCalled:   false,
			wantStatus:   http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &dbtest.FakeDB{
				GetBrowserSessionFn: func(string) (*database.BrowserSession, error) {
					return tt.session, tt.sessionErr
				},
			}
			h := newTestHandler(t, db, nil)
			called := false
			next := h.WithAuthentication(func(w http.ResponseWriter, r *http.Request) {
				called = true
				hCtx := GetHandlerContext(r)
				if hCtx.ValidSession == nil {
					t.Error("authenticated request is missing the session in the handler context")
				}
				if GetSessionFromContext(r).UserID != tt.session.UserID {
					t.Error("authenticated request is missing the session value in the context")
				}
			}, tt.requiredRole)

			rec := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			if tt.cookie != nil {
				r.AddCookie(tt.cookie)
			}
			next.ServeHTTP(rec, r)

			if called != tt.wantCalled {
				t.Errorf("next handler called = %v, want %v", called, tt.wantCalled)
			}
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestWithAdminAndAnyAuthenticationRoles(t *testing.T) {
	db := &dbtest.FakeDB{
		GetBrowserSessionFn: func(string) (*database.BrowserSession, error) {
			return validSession(1, database.UserRoleMember), nil
		},
	}
	h := newTestHandler(t, db, nil)

	adminCalled := false
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.AddCookie(&http.Cookie{Name: "SESSION", Value: "secret"})
	h.WithAdminAuthentication(func(http.ResponseWriter, *http.Request) { adminCalled = true }).ServeHTTP(rec, r)
	if adminCalled {
		t.Error("a member must not pass the admin authentication")
	}

	anyCalled := false
	rec = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/", nil)
	r.AddCookie(&http.Cookie{Name: "SESSION", Value: "secret"})
	h.WithAnyAuthentication(func(http.ResponseWriter, *http.Request) { anyCalled = true }).ServeHTTP(rec, r)
	if !anyCalled {
		t.Error("a member has to pass the 'any role' authentication")
	}
}

type testBody struct {
	Name string `json:"name" validate:"required"`
}

func TestWithRequestBody(t *testing.T) {
	tests := []struct {
		name       string
		payload    string
		errorMap   map[string]string
		wantCalled bool
		wantMsg    string
	}{
		{name: "valid body", payload: `{"name":"foo"}`, wantCalled: true},
		{name: "invalid json", payload: `{`, wantCalled: false},
		{
			name:       "validation error with custom message",
			payload:    `{}`,
			errorMap:   map[string]string{"Name": "a name is required"},
			wantCalled: false,
			wantMsg:    "a name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			var got testBody
			var handler http.HandlerFunc
			next := func(w http.ResponseWriter, r *http.Request, body testBody, hCtx *HandlerCtx) {
				called = true
				got = body
			}
			if tt.errorMap == nil {
				handler = WithRequestBody(next)
			} else {
				handler = WithRequestBody(next, tt.errorMap)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, newRequest(tt.payload, &HandlerCtx{Database: &dbtest.FakeDB{}}))

			if called != tt.wantCalled {
				t.Fatalf("next handler called = %v, want %v", called, tt.wantCalled)
			}
			if tt.wantCalled {
				if got.Name != "foo" {
					t.Errorf("body.Name = %q, want %q", got.Name, "foo")
				}
				return
			}
			resp := decodeResponse(t, rec)
			if resp.OK {
				t.Error("response must not be a success response")
			}
			if tt.wantMsg != "" && resp.Msg != tt.wantMsg {
				t.Errorf("msg = %q, want %q", resp.Msg, tt.wantMsg)
			}
		})
	}
}

func TestGetHandlerContextWithoutContextValue(t *testing.T) {
	hCtx := GetHandlerContext(httptest.NewRequest(http.MethodPost, "/", nil))
	if hCtx == nil {
		t.Fatal("GetHandlerContext returned nil")
	}
	if hCtx.Database != nil || hCtx.ValidSession != nil || hCtx.Config != nil {
		t.Errorf("expected an empty handler context, got %+v", hCtx)
	}
}

func TestGetSessionFromContextWithoutContextValue(t *testing.T) {
	s := GetSessionFromContext(httptest.NewRequest(http.MethodPost, "/", nil))
	if !s.IsEmpty() {
		t.Errorf("expected an empty session, got %+v", s)
	}
}

func TestReadBodyAndValidateUploadFileValidator(t *testing.T) {
	type logoBody struct {
		Logo string `json:"logo" validate:"omitempty,uploadfile"`
	}

	uploadDir := t.TempDir()
	conf := newTestConfig(t, map[string]string{"http.uploadpath": uploadDir})

	tests := []struct {
		name    string
		payload string
		hCtx    *HandlerCtx
		wantErr bool
	}{
		{name: "hex filename is accepted", payload: `{"logo":"abcdef0123.png"}`, hCtx: &HandlerCtx{Config: conf}},
		{name: "directory traversal is rejected", payload: `{"logo":"../../etc/passwd"}`, hCtx: &HandlerCtx{Config: conf}, wantErr: true},
		{name: "non hex filename is rejected", payload: `{"logo":"logo.png"}`, hCtx: &HandlerCtx{Config: conf}, wantErr: true},
		{name: "no config means no upload path", payload: `{"logo":"abcdef0123.png"}`, hCtx: &HandlerCtx{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &logoBody{}
			err := ReadBodyAndValidate(newRequest(tt.payload, tt.hCtx), body)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBodyAndValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadBodyAndValidateFallsBackToPlainError(t *testing.T) {
	body := &testBody{}
	// The error map does not contain an entry for the failing field, so the
	// original validation error has to be returned.
	err := ReadBodyAndValidate(newRequest(`{}`, nil), body, map[string]string{"Other": "not relevant"})
	if err == nil {
		t.Fatal("expected a validation error")
	}
	if err.Error() == "not relevant" {
		t.Error("an unrelated error mapping must not be used")
	}
}

func TestWriteSuccessResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteSuccessResponse("all good", PingData{Message: "hello"}, rec)

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}
	resp := decodeResponse(t, rec)
	if !resp.OK || resp.Msg != "all good" {
		t.Errorf("response = %+v, want ok with message 'all good'", resp)
	}
	var data PingData
	decodeData(t, resp, &data)
	if data.Message != "hello" {
		t.Errorf("data.Message = %q, want %q", data.Message, "hello")
	}
}

func TestWriteFailureResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteFailureResponse("it broke", rec)

	resp := decodeResponse(t, rec)
	if resp.OK {
		t.Error("failure response must not be ok")
	}
	if resp.Msg != "it broke" {
		t.Errorf("msg = %q, want %q", resp.Msg, "it broke")
	}
	if len(resp.Data) != 0 {
		t.Errorf("failure response must not contain data, got %q", string(resp.Data))
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}
}

func TestSetAndDeleteSessionCookie(t *testing.T) {
	rec := httptest.NewRecorder()
	validUntil := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	SetSessionCookie(rec, "the-secret", validUntil)

	cookies := (&http.Response{Header: rec.Header()}).Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected exactly one cookie, got %d", len(cookies))
	}
	c := cookies[0]
	if c.Name != "SESSION" || c.Value != "the-secret" {
		t.Errorf("cookie = %s=%s, want SESSION=the-secret", c.Name, c.Value)
	}
	if !c.HttpOnly {
		t.Error("session cookie has to be HttpOnly")
	}
	if c.Path != "/" {
		t.Errorf("cookie path = %q, want %q", c.Path, "/")
	}
	if !c.Expires.Equal(validUntil) {
		t.Errorf("cookie expiry = %v, want %v", c.Expires, validUntil)
	}

	rec = httptest.NewRecorder()
	DeleteSessionCookie(rec)
	cookies = (&http.Response{Header: rec.Header()}).Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected exactly one cookie, got %d", len(cookies))
	}
	if cookies[0].Value != "" {
		t.Errorf("deleted cookie value = %q, want empty", cookies[0].Value)
	}
	if !cookies[0].Expires.Before(time.Now()) {
		t.Error("deleted cookie has to be expired")
	}
}

func TestWithLogging(t *testing.T) {
	h := newTestHandler(t, &dbtest.FakeDB{}, nil)
	called := false
	wrapped := h.WithLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	wrapped.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if !called {
		t.Error("the wrapped handler was not called")
	}
}

func TestEmailAndMyNautiqueClientAccessors(t *testing.T) {
	h := newTestHandler(t, &dbtest.FakeDB{}, nil)
	if h.GetMyNautiqueClient() != nil {
		t.Error("expected no myNautique client on a fresh handler")
	}
	h.SetMyNautiqueClient(mynautique.Client{})
	if h.GetMyNautiqueClient() == nil {
		t.Error("the stored myNautique client was not returned")
	}

	h.SetEmailClient(email.Client{Host: "smtp.example.com", Port: "25"})
	if got := h.GetEmailClient(); got.Host != "smtp.example.com" {
		t.Errorf("email client host = %q, want smtp.example.com", got.Host)
	}
}

// GetEmailClient dereferences the stored pointer, so it must not be called
// before a client has been set.
func TestGetEmailClientWithoutAClientPanics(t *testing.T) {
	h := newTestHandler(t, &dbtest.FakeDB{}, nil)
	defer func() {
		if recover() == nil {
			t.Error("expected GetEmailClient to panic when no client was set")
		}
	}()
	h.GetEmailClient()
}

// errorReader is used to make io.ReadAll on the request body fail.
type errorReader struct{}

func (errorReader) Read([]byte) (int, error) { return 0, errors.New("broken body") }
func (errorReader) Close() error             { return nil }

func TestReadBodyAndValidateUnreadableBody(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Body = errorReader{}
	if err := ReadBodyAndValidate(r, &testBody{}); err == nil {
		t.Error("expected an error for an unreadable body")
	}
}
