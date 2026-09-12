package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"server/database"
	"server/database/dbtest"
	"server/util/hash"
)

// testUser builds a user whose password is `password`, hashed the same way
// the login handler does it.
func testUser(t *testing.T, id uint, password string) database.User {
	t.Helper()
	salt := 4711
	h, err := hash.CryptSha512(hash.Sha256(password), strconv.Itoa(salt))
	if err != nil {
		t.Fatalf("cannot hash test password: %v", err)
	}
	return database.User{
		ID:           id,
		Username:     "tester",
		PasswordSalt: salt,
		PasswordHash: h,
		FirstName:    "Test",
		LastName:     "User",
		Email:        "test@example.com",
		UserStatusID: database.UserStatusMember,
		UserStatus: database.UserStatus{
			ID:         database.UserStatusMember,
			UserRoleID: database.UserRoleMember,
			UserRole: database.UserRole{
				ID:   database.UserRoleMember,
				Name: "member",
			},
		},
	}
}

func TestLogin(t *testing.T) {
	const password = "Sup3rSecret!Password"

	t.Run("successful login stores a session and sets the cookie", func(t *testing.T) {
		user := testUser(t, 42, password)
		var stored database.BrowserSession
		db := &dbtest.FakeDB{
			GetUserByNameFn: func(string) (database.User, error) { return user, nil },
			GetUserByIdFn:   func(uint) (database.User, error) { return user, nil },
			AddBrowserSessionFn: func(b database.BrowserSession) (string, error) {
				stored = b
				return b.SessionSecret, nil
			},
		}
		h := newTestHandler(t, db, map[string]string{
			"http.sessioninactivitytimeout": "3600",
			"http.sessiontimeout":           "86400",
		})

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.Login(rec, r, LoginRequest{Username: "tester", Password: password}, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("login failed: %s", resp.Msg)
		}
		if stored.SessionSecret == "" {
			t.Error("no session secret was stored")
		}
		if stored.UserID != user.ID {
			t.Errorf("stored session user = %d, want %d", stored.UserID, user.ID)
		}
		if stored.UserRoleID != database.UserRoleMember {
			t.Errorf("stored session role = %d, want %d", stored.UserRoleID, database.UserRoleMember)
		}
		if stored.ValidUntil.After(stored.MaxValidUntil) {
			t.Error("the inactivity timeout must not exceed the absolute session timeout")
		}
		cookies := (&http.Response{Header: rec.Header()}).Cookies()
		if len(cookies) != 1 || cookies[0].Value != stored.SessionSecret {
			t.Errorf("session cookie does not match the stored session secret: %+v", cookies)
		}
	})

	t.Run("failures", func(t *testing.T) {
		user := testUser(t, 42, password)
		lockedUser := testUser(t, 42, password)
		lockedUser.Locked = true

		tests := []struct {
			name    string
			db      *dbtest.FakeDB
			req     LoginRequest
			wantMsg string
		}{
			{
				name:    "unknown user",
				db:      &dbtest.FakeDB{},
				req:     LoginRequest{Username: "nobody", Password: password},
				wantMsg: "invalid username/password",
			},
			{
				name: "locked user",
				db: &dbtest.FakeDB{
					GetUserByNameFn: func(string) (database.User, error) { return lockedUser, nil },
				},
				req:     LoginRequest{Username: "tester", Password: password},
				wantMsg: "user account not activated",
			},
			{
				name: "wrong password",
				db: &dbtest.FakeDB{
					GetUserByNameFn: func(string) (database.User, error) { return user, nil },
				},
				req:     LoginRequest{Username: "tester", Password: "wrong-password"},
				wantMsg: "invalid username/password",
			},
			{
				name: "user details cannot be loaded",
				db: &dbtest.FakeDB{
					GetUserByNameFn: func(string) (database.User, error) { return user, nil },
					GetUserByIdFn:   func(uint) (database.User, error) { return database.User{}, errNotFound },
				},
				req:     LoginRequest{Username: "tester", Password: password},
				wantMsg: "invalid username/password",
			},
			{
				name: "session cannot be stored",
				db: &dbtest.FakeDB{
					GetUserByNameFn: func(string) (database.User, error) { return user, nil },
					GetUserByIdFn:   func(uint) (database.User, error) { return user, nil },
					AddBrowserSessionFn: func(database.BrowserSession) (string, error) {
						return "", errNotFound
					},
				},
				req:     LoginRequest{Username: "tester", Password: password},
				wantMsg: "invalid username/password",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.Login(rec, r, tt.req, GetHandlerContext(r))

				resp := decodeResponse(t, rec)
				if resp.OK {
					t.Fatal("login must not succeed")
				}
				if resp.Msg != tt.wantMsg {
					t.Errorf("msg = %q, want %q", resp.Msg, tt.wantMsg)
				}
				if cookies := (&http.Response{Header: rec.Header()}).Cookies(); len(cookies) != 0 {
					t.Errorf("a failed login must not set a cookie, got %+v", cookies)
				}
			})
		}
	})
}

func TestIsLoggedIn(t *testing.T) {
	tests := []struct {
		name         string
		cookie       bool
		session      *database.BrowserSession
		sessionErr   error
		wantLoggedIn bool
	}{
		{
			name:         "valid session",
			cookie:       true,
			session:      validSession(1, database.UserRoleMember),
			wantLoggedIn: true,
		},
		{
			name:         "no cookie",
			cookie:       false,
			wantLoggedIn: false,
		},
		{
			name:         "unknown session",
			cookie:       true,
			sessionErr:   errNotFound,
			wantLoggedIn: false,
		},
		{
			name:   "expired session",
			cookie: true,
			session: &database.BrowserSession{
				SessionSecret: "expired",
				ValidUntil:    time.Now().Add(-time.Minute),
				MaxValidUntil: time.Now().Add(time.Hour),
				UserID:        1,
			},
			wantLoggedIn: false,
		},
		{
			name:   "session beyond the absolute timeout",
			cookie: true,
			session: &database.BrowserSession{
				SessionSecret: "too old",
				ValidUntil:    time.Now().Add(time.Minute),
				MaxValidUntil: time.Now().Add(time.Minute),
				UserID:        1,
			},
			wantLoggedIn: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated := false
			db := &dbtest.FakeDB{
				GetBrowserSessionFn: func(string) (*database.BrowserSession, error) {
					return tt.session, tt.sessionErr
				},
				UpdateBrowserSessionFn: func(database.BrowserSession) error {
					updated = true
					return nil
				},
			}
			h := newTestHandler(t, db, map[string]string{"http.sessioninactivitytimeout": "3600"})

			rec := httptest.NewRecorder()
			r := newRequest("", &HandlerCtx{Database: db})
			if tt.cookie {
				r.AddCookie(&http.Cookie{Name: h.SessionCookieName(), Value: "secret"})
			}
			h.IsLoggedIn(rec, r)

			resp := decodeResponse(t, rec)
			if !resp.OK {
				t.Fatalf("IsLoggedIn always returns a success response, got %+v", resp)
			}
			var data IsLoggedInResponse
			decodeData(t, resp, &data)
			if data.LoggedIn != tt.wantLoggedIn {
				t.Errorf("loggedIn = %v, want %v", data.LoggedIn, tt.wantLoggedIn)
			}
			if tt.wantLoggedIn && !updated {
				t.Error("the session activity timestamp was not refreshed")
			}
			cookies := (&http.Response{Header: rec.Header()}).Cookies()
			if !tt.wantLoggedIn && !cookieCleared(cookies, h.SessionCookieName()) {
				t.Errorf("a request without a valid session has to clear the session cookie, got %+v", cookies)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	deleted := false
	session := validSession(3, database.UserRoleMember)
	db := &dbtest.FakeDB{
		DeleteBrowserSessionFn: func(b database.BrowserSession) error {
			deleted = true
			if b.SessionSecret != session.SessionSecret {
				t.Errorf("deleted session = %q, want %q", b.SessionSecret, session.SessionSecret)
			}
			return nil
		},
	}
	h := newTestHandler(t, db, nil)

	rec := httptest.NewRecorder()
	h.Logout(rec, newRequest("", &HandlerCtx{Database: db, ValidSession: session}))

	resp := decodeResponse(t, rec)
	if !resp.OK || resp.Msg != "logged out" {
		t.Errorf("response = %+v, want a successful logout", resp)
	}
	if !deleted {
		t.Error("the browser session was not deleted")
	}
	cookies := (&http.Response{Header: rec.Header()}).Cookies()
	if !cookieCleared(cookies, h.SessionCookieName()) {
		t.Errorf("the session cookie was not cleared: %+v", cookies)
	}
}

func TestUser(t *testing.T) {
	t.Run("returns the user of the session", func(t *testing.T) {
		user := testUser(t, 9, "irrelevant")
		user.Address = "Main Street 1"
		user.City = "Zurich"
		user.ZipCode = 8000
		user.MobilePhoneNr = "+41791234567"
		user.BoatLicense = true
		user.Comment = "a comment"

		db := &dbtest.FakeDB{
			GetUserByIdFn: func(id uint) (database.User, error) {
				if id != 9 {
					t.Errorf("GetUserByID(%d), want 9", id)
				}
				return user, nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		h.User(rec, newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(9, database.UserRoleMember)}))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data UserResponse
		decodeData(t, resp, &data)
		if data.ID != user.ID || data.Email != user.Email || data.City != user.City {
			t.Errorf("response = %+v, does not match the user %+v", data, user)
		}
		if data.UserRoleID != database.UserRoleMember || data.UserRoleName != "member" {
			t.Errorf("role = %d/%q, want %d/member", data.UserRoleID, data.UserRoleName, database.UserRoleMember)
		}
		if data.Status != database.UserStatusMember {
			t.Errorf("status = %d, want %d", data.Status, database.UserStatusMember)
		}
	})

	t.Run("fails if the user cannot be loaded", func(t *testing.T) {
		db := &dbtest.FakeDB{}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		h.User(rec, newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(9, database.UserRoleMember)}))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response for an unknown user")
		}
	})
}

func TestLoginRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{name: "username and password", payload: `{"username":"tester","password":"secret"}`},
		{name: "email as username", payload: `{"username":"test@example.com","password":"secret"}`},
		{name: "missing password", payload: `{"username":"tester"}`, wantErr: true},
		{name: "missing username", payload: `{"password":"secret"}`, wantErr: true},
		{name: "username with invalid characters", payload: `{"username":"te ster!","password":"secret"}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReadBodyAndValidate(newRequest(tt.payload, nil), &LoginRequest{})
			if (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
