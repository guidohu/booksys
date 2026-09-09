package handlers

import (
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"server/database"
	"server/database/dbtest"
	"server/util/hash"
)

func TestSignUp(t *testing.T) {
	req := SignUpRequest{
		Username:  "newuser",
		Password:  "Sup3rSecret!Password",
		FirstName: "New",
		LastName:  "User",
		Address:   "Main Street 1",
		MobileNr:  "0791234567",
		ZipCode:   8000,
		City:      "Zurich",
		Email:     "new@example.com",
		License:   true,
		AcceptGTC: true,
	}

	t.Run("creates a locked guest user with a hashed password", func(t *testing.T) {
		var added database.User
		db := &dbtest.FakeDB{
			AddUserFn: func(u database.User) (uint, error) {
				added = u
				return 17, nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.SignUp(rec, r, req, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("signup failed: %s", resp.Msg)
		}
		var data SignUpResponse
		decodeData(t, resp, &data)
		if data.UserID != 17 {
			t.Errorf("user_id = %d, want 17", data.UserID)
		}
		if !added.Locked {
			t.Error("a new user has to be locked until an admin unlocks it")
		}
		if added.UserStatusID != database.UserStatusGuest {
			t.Errorf("status = %d, want guest (%d)", added.UserStatusID, database.UserStatusGuest)
		}
		if added.PasswordHash == "" || added.PasswordHash == req.Password {
			t.Error("the password has to be stored as a hash")
		}
		expected, err := hash.CryptSha512(hash.Sha256(req.Password), strconv.Itoa(added.PasswordSalt))
		if err != nil {
			t.Fatalf("cannot recompute the hash: %v", err)
		}
		if added.PasswordHash != expected {
			t.Error("the stored password hash cannot be reproduced from the submitted password")
		}
	})

	t.Run("rejects an existing user", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetUserByNameFn: func(string) (database.User, error) { return database.User{ID: 1}, nil },
			AddUserFn: func(database.User) (uint, error) {
				t.Error("an existing user must not be overwritten")
				return 0, nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.SignUp(rec, r, req, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if resp.OK || resp.Msg != "user already exists" {
			t.Errorf("response = %+v, want the 'user already exists' failure", resp)
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{
			AddUserFn: func(database.User) (uint, error) { return 0, errNotFound },
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.SignUp(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response when the user cannot be stored")
		}
	})
}

func TestSignUpRequestValidation(t *testing.T) {
	valid := `{"username":"newuser","password":"Sup3rSecret!Pwd","first_name":"New","last_name":"User",` +
		`"address":"Main Street 1","mobile":"0791234567","plz":8000,"city":"Zurich","email":"new@example.com",` +
		`"license":true,"ownRisk":true}`

	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{name: "valid payload", payload: valid},
		{name: "weak password", payload: `{"username":"u","password":"short","first_name":"N","last_name":"U","address":"A","mobile":"0791234567","plz":8000,"city":"Z","email":"a@b.co","ownRisk":true}`, wantErr: true},
		{name: "invalid email", payload: `{"username":"u","password":"Sup3rSecret!Pwd","first_name":"N","last_name":"U","address":"A","mobile":"0791234567","plz":8000,"city":"Z","email":"not-an-email","ownRisk":true}`, wantErr: true},
		{name: "gtc not accepted", payload: `{"username":"u","password":"Sup3rSecret!Pwd","first_name":"N","last_name":"U","address":"A","mobile":"0791234567","plz":8000,"city":"Z","email":"a@b.co","ownRisk":false}`, wantErr: true},
		{name: "injection characters in the name", payload: `{"username":"u","password":"Sup3rSecret!Pwd","first_name":"<script>","last_name":"U","address":"A","mobile":"0791234567","plz":8000,"city":"Z","email":"a@b.co","ownRisk":true}`, wantErr: true},
		{name: "missing zip code", payload: `{"username":"u","password":"Sup3rSecret!Pwd","first_name":"N","last_name":"U","address":"A","mobile":"0791234567","city":"Z","email":"a@b.co","ownRisk":true}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReadBodyAndValidate(newRequest(tt.payload, nil), &SignUpRequest{}, SignUpRequestValidationErrors)
			if (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMakeAdmin(t *testing.T) {
	t.Run("promotes and unlocks the first user", func(t *testing.T) {
		var statusID, statusUserID uint
		var lockUserID uint
		locked := true
		db := &dbtest.FakeDB{
			CountAdminUsersFn: func() int64 { return 0 },
			ChangeUserStatusFn: func(id, s uint) error {
				statusUserID, statusID = id, s
				return nil
			},
			ChangeLockFn: func(id uint, l bool) error {
				lockUserID, locked = id, l
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.MakeAdmin(rec, r, MakeAdminRequest{UserID: 5}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("MakeAdmin failed: %s", resp.Msg)
		}
		if statusUserID != 5 || statusID != database.UserStatusAdmin {
			t.Errorf("ChangeUserStatus(%d, %d), want (5, %d)", statusUserID, statusID, database.UserStatusAdmin)
		}
		if lockUserID != 5 || locked {
			t.Errorf("ChangeLock(%d, %v), want (5, false)", lockUserID, locked)
		}
	})

	t.Run("refuses when an admin already exists", func(t *testing.T) {
		db := &dbtest.FakeDB{
			CountAdminUsersFn: func() int64 { return 1 },
			ChangeUserStatusFn: func(uint, uint) error {
				t.Error("no further admin may be created through this endpoint")
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.MakeAdmin(rec, r, MakeAdminRequest{UserID: 5}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response when an admin already exists")
		}
	})

	t.Run("reports database failures", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{
				name: "status update fails",
				db: &dbtest.FakeDB{
					ChangeUserStatusFn: func(uint, uint) error { return errNotFound },
				},
			},
			{
				name: "unlock fails",
				db: &dbtest.FakeDB{
					ChangeLockFn: func(uint, bool) error { return errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.MakeAdmin(rec, r, MakeAdminRequest{UserID: 5}, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestDeleteUser(t *testing.T) {
	member := database.User{
		ID: 3,
		UserStatus: database.UserStatus{
			UserRoleID: database.UserRoleMember,
		},
	}
	admin := database.User{
		ID: 4,
		UserStatus: database.UserStatus{
			UserRoleID: database.UserRoleAdmin,
		},
	}

	tests := []struct {
		name       string
		req        DeleteUserRequest
		session    *database.BrowserSession
		db         *dbtest.FakeDB
		wantDelete bool
		wantMsg    string
	}{
		{
			name:    "deletes a user with a zero balance",
			req:     DeleteUserRequest{UserID: 3},
			session: validSession(1, database.UserRoleAdmin),
			db: &dbtest.FakeDB{
				GetUserByIdFn: func(uint) (database.User, error) { return member, nil },
			},
			wantDelete: true,
		},
		{
			name:    "refuses self deletion",
			req:     DeleteUserRequest{UserID: 1},
			session: validSession(1, database.UserRoleAdmin),
			db:      &dbtest.FakeDB{},
			wantMsg: "You cannot delete yourself.",
		},
		{
			name:    "refuses unknown users",
			req:     DeleteUserRequest{UserID: 3},
			session: validSession(1, database.UserRoleAdmin),
			db:      &dbtest.FakeDB{},
			wantMsg: "User not found.",
		},
		{
			name:    "refuses admins",
			req:     DeleteUserRequest{UserID: 4},
			session: validSession(1, database.UserRoleAdmin),
			db: &dbtest.FakeDB{
				GetUserByIdFn: func(uint) (database.User, error) { return admin, nil },
			},
			wantMsg: "Users with admin roles cannot be removed.",
		},
		{
			name:    "refuses users with a non-zero balance",
			req:     DeleteUserRequest{UserID: 3},
			session: validSession(1, database.UserRoleAdmin),
			db: &dbtest.FakeDB{
				GetUserByIdFn: func(uint) (database.User, error) { return member, nil },
				GetUserSessionPaymentsFn: func(uint) (decimal.Decimal, error) {
					return decimal.NewFromInt(20), nil
				},
			},
			wantMsg: "Cannot delete user with non-zero balance.",
		},
		{
			name:    "reports a failing balance lookup",
			req:     DeleteUserRequest{UserID: 3},
			session: validSession(1, database.UserRoleAdmin),
			db: &dbtest.FakeDB{
				GetUserByIdFn: func(uint) (database.User, error) { return member, nil },
				GetUserHeatStatsFn: func(uint, time.Time, time.Time) (int64, decimal.Decimal, error) {
					return 0, decimal.Zero, errNotFound
				},
			},
			wantMsg: "Cannot check for balance to be 0 for this user.",
		},
		{
			name:    "reports a failing deletion",
			req:     DeleteUserRequest{UserID: 3},
			session: validSession(1, database.UserRoleAdmin),
			db: &dbtest.FakeDB{
				GetUserByIdFn:    func(uint) (database.User, error) { return member, nil },
				DeleteUserByIdFn: func(uint) error { return errNotFound },
			},
			wantDelete: true,
			wantMsg:    "Cannot delete user.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleted := false
			inner := tt.db.DeleteUserByIdFn
			tt.db.DeleteUserByIdFn = func(id uint) error {
				deleted = true
				if inner != nil {
					return inner(id)
				}
				return nil
			}
			h := newTestHandler(t, tt.db, nil)

			rec := httptest.NewRecorder()
			r := newRequest("", &HandlerCtx{Database: tt.db, ValidSession: tt.session})
			h.DeleteUser(rec, r, tt.req, GetHandlerContext(r))

			resp := decodeResponse(t, rec)
			if tt.wantMsg == "" {
				if !resp.OK {
					t.Fatalf("unexpected failure: %s", resp.Msg)
				}
			} else if resp.OK || resp.Msg != tt.wantMsg {
				t.Errorf("response = %+v, want failure %q", resp, tt.wantMsg)
			}
			if deleted != tt.wantDelete {
				t.Errorf("user deleted = %v, want %v", deleted, tt.wantDelete)
			}
		})
	}
}

func TestUpdateMyUser(t *testing.T) {
	payload := `{"first_name":"New","last_name":"Name","address":"Street 1","mobile":"0791234567",` +
		`"plz":8000,"city":"Zurich","email":"new@example.com","license":true}`

	t.Run("updates the user of the session", func(t *testing.T) {
		var updatedID uint
		var updated database.User
		db := &dbtest.FakeDB{
			UpdateUserFn: func(id uint, u database.User) error {
				updatedID, updated = id, u
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		session := validSession(11, database.UserRoleMember)
		session.Username = "old@example.com"
		session.User = database.User{Email: "old@example.com"}

		rec := httptest.NewRecorder()
		h.UpdateMyUser(rec, newRequest(payload, &HandlerCtx{Database: db, ValidSession: session}))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("update failed: %s", resp.Msg)
		}
		if updatedID != 11 || updated.ID != 11 {
			t.Errorf("updated user id = %d/%d, want 11", updatedID, updated.ID)
		}
		if updated.Email != "new@example.com" || updated.City != "Zurich" {
			t.Errorf("updated user = %+v, does not match the request", updated)
		}
		if updated.Username != "new@example.com" {
			t.Errorf("username = %q, want the new email because the old username was the email", updated.Username)
		}
		if updated.PasswordHash != "" || updated.PasswordSalt != 0 || updated.UserStatusID != 0 {
			t.Error("privileged fields must not be set through a self service update")
		}
	})

	t.Run("keeps a username that differs from the email", func(t *testing.T) {
		var updated database.User
		db := &dbtest.FakeDB{
			UpdateUserFn: func(_ uint, u database.User) error {
				updated = u
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		session := validSession(11, database.UserRoleMember)
		session.Username = "custom-name"
		session.User = database.User{Email: "old@example.com"}

		rec := httptest.NewRecorder()
		h.UpdateMyUser(rec, newRequest(payload, &HandlerCtx{Database: db, ValidSession: session}))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("update failed: %s", resp.Msg)
		}
		if updated.Username != "custom-name" {
			t.Errorf("username = %q, want the unchanged username 'custom-name'", updated.Username)
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name    string
			payload string
			session *database.BrowserSession
			db      *dbtest.FakeDB
		}{
			{
				name:    "no user in the session",
				payload: payload,
				session: validSession(0, database.UserRoleMember),
				db:      &dbtest.FakeDB{},
			},
			{
				name:    "invalid payload",
				payload: `{"first_name":"New"}`,
				session: validSession(11, database.UserRoleMember),
				db:      &dbtest.FakeDB{},
			},
			{
				name:    "database failure",
				payload: payload,
				session: validSession(11, database.UserRoleMember),
				db: &dbtest.FakeDB{
					UpdateUserFn: func(uint, database.User) error { return errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				h.UpdateMyUser(rec, newRequest(tt.payload, &HandlerCtx{Database: tt.db, ValidSession: tt.session}))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestUpdateMyPassword(t *testing.T) {
	const oldPassword = "0ldSecret!Password"
	const newPassword = "N3wSecret!Password"
	user := testUser(t, 11, oldPassword)

	t.Run("changes the password", func(t *testing.T) {
		var stored database.User
		db := &dbtest.FakeDB{
			GetUserByIdFn: func(uint) (database.User, error) { return user, nil },
			UpdatePasswordFn: func(_ uint, u database.User) error {
				stored = u
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		payload := `{"password_old":"` + oldPassword + `","password_new":"` + newPassword + `"}`
		rec := httptest.NewRecorder()
		h.UpdateMyPassword(rec, newRequest(payload, &HandlerCtx{Database: db, ValidSession: validSession(11, database.UserRoleMember)}))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("password change failed: %s", resp.Msg)
		}
		expected, err := hash.CryptSha512(hash.Sha256(newPassword), strconv.Itoa(stored.PasswordSalt))
		if err != nil {
			t.Fatalf("cannot recompute the hash: %v", err)
		}
		if stored.PasswordHash != expected {
			t.Error("the stored hash does not match the new password")
		}
		if stored.PasswordSalt == user.PasswordSalt {
			t.Error("a new salt should be generated for a new password")
		}
	})

	t.Run("rejects a wrong old password", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetUserByIdFn: func(uint) (database.User, error) { return user, nil },
			UpdatePasswordFn: func(uint, database.User) error {
				t.Error("the password must not be changed with a wrong old password")
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		payload := `{"password_old":"Wr0ngSecret!Pwd","password_new":"` + newPassword + `"}`
		rec := httptest.NewRecorder()
		h.UpdateMyPassword(rec, newRequest(payload, &HandlerCtx{Database: db, ValidSession: validSession(11, database.UserRoleMember)}))

		resp := decodeResponse(t, rec)
		if resp.OK {
			t.Error("expected a failure response")
		}
	})

	t.Run("failures", func(t *testing.T) {
		payload := `{"password_old":"` + oldPassword + `","password_new":"` + newPassword + `"}`
		tests := []struct {
			name    string
			payload string
			session *database.BrowserSession
			db      *dbtest.FakeDB
		}{
			{
				name:    "no user in the session",
				payload: payload,
				session: validSession(0, database.UserRoleMember),
				db:      &dbtest.FakeDB{},
			},
			{
				name:    "weak new password",
				payload: `{"password_old":"` + oldPassword + `","password_new":"weak"}`,
				session: validSession(11, database.UserRoleMember),
				db:      &dbtest.FakeDB{},
			},
			{
				name:    "user cannot be found",
				payload: payload,
				session: validSession(11, database.UserRoleMember),
				db:      &dbtest.FakeDB{},
			},
			{
				name:    "password cannot be stored",
				payload: payload,
				session: validSession(11, database.UserRoleMember),
				db: &dbtest.FakeDB{
					GetUserByIdFn:    func(uint) (database.User, error) { return user, nil },
					UpdatePasswordFn: func(uint, database.User) error { return errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				h.UpdateMyPassword(rec, newRequest(tt.payload, &HandlerCtx{Database: tt.db, ValidSession: tt.session}))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestGetMySessions(t *testing.T) {
	past := database.Session{
		ID:            1,
		Title:         "past",
		StartTime:     time.Now().Add(-48 * time.Hour),
		EndTime:       time.Now().Add(-47 * time.Hour),
		SessionTypeID: database.DefaultSessionType,
		SessionType:   database.SessionType{Name: "default"},
	}
	upcoming := database.Session{
		ID:            2,
		Title:         "upcoming",
		StartTime:     time.Now().Add(24 * time.Hour),
		EndTime:       time.Now().Add(25 * time.Hour),
		SessionTypeID: database.DefaultSessionType,
		SessionType:   database.SessionType{Name: "default"},
	}

	db := &dbtest.FakeDB{
		GetSessionsByUserFn: func(uint) ([]database.Session, error) {
			return []database.Session{past, upcoming}, nil
		},
		GetUsersForSessionFn: func(id uint) ([]database.UserToSession, error) {
			return []database.UserToSession{
				{ID: 10, UserID: 5, User: database.User{ID: 5, FirstName: "Ann", LastName: "Rider"}},
			}, nil
		},
	}
	h := newTestHandler(t, db, nil)

	rec := httptest.NewRecorder()
	h.GetMySessions(rec, newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(5, database.UserRoleMember)}))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data GetMySessionsResponse
	decodeData(t, resp, &data)
	if len(data.PastSessions) != 1 || data.PastSessions[0].Title != "past" {
		t.Errorf("past sessions = %+v, want exactly the past session", data.PastSessions)
	}
	if len(data.UpcomingSessions) != 1 || data.UpcomingSessions[0].Title != "upcoming" {
		t.Errorf("upcoming sessions = %+v, want exactly the upcoming session", data.UpcomingSessions)
	}
	if len(data.UpcomingSessions[0].Riders) != 1 || data.UpcomingSessions[0].Riders[0].FirstName != "Ann" {
		t.Errorf("riders = %+v, want the rider Ann", data.UpcomingSessions[0].Riders)
	}
}

func TestGetMySessionsFailures(t *testing.T) {
	t.Run("no user in the session", func(t *testing.T) {
		db := &dbtest.FakeDB{}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		h.GetMySessions(rec, newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(0, database.UserRoleMember)}))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})

	t.Run("database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetSessionsByUserFn: func(uint) ([]database.Session, error) { return nil, errNotFound },
		}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		h.GetMySessions(rec, newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(5, database.UserRoleMember)}))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestGetMyHeats(t *testing.T) {
	timestamp := time.Date(2024, time.May, 4, 15, 30, 0, 0, time.UTC)
	db := &dbtest.FakeDB{
		GetUserHeatsFn: func(userID uint, size int) ([]database.Heat, error) {
			if userID != 5 {
				t.Errorf("GetUserHeats(%d), want user 5", userID)
			}
			return []database.Heat{
				{ID: 1, Timestamp: timestamp, DurationSeconds: 754, Cost: mustDecimal(t, "35.20")},
			}, nil
		},
	}
	h := newTestHandler(t, db, nil)

	rec := httptest.NewRecorder()
	h.GetMyHeats(rec, newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(5, database.UserRoleMember)}))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data GetMyHeatsResponse
	decodeData(t, resp, &data)
	if len(data.Heats) != 1 {
		t.Fatalf("heats = %+v, want exactly one entry", data.Heats)
	}
	got := data.Heats[0]
	if got.Date != "04.05.2024" {
		t.Errorf("date = %q, want %q", got.Date, "04.05.2024")
	}
	if got.DurationSeconds != 754 {
		t.Errorf("duration_seconds = %d, want 754", got.DurationSeconds)
	}
	if got.DurationText != "00:12:34" {
		t.Errorf("duration = %q, want %q", got.DurationText, "00:12:34")
	}
	if got.DateUnixMillis != timestamp.UnixMilli() {
		t.Errorf("date_unix_millis = %d, want %d", got.DateUnixMillis, timestamp.UnixMilli())
	}
}

func TestGetMyHeatsFailures(t *testing.T) {
	t.Run("no user in the session", func(t *testing.T) {
		db := &dbtest.FakeDB{}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		h.GetMyHeats(rec, newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(0, database.UserRoleMember)}))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})

	t.Run("database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetUserHeatsFn: func(uint, int) ([]database.Heat, error) { return nil, errNotFound },
		}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		h.GetMyHeats(rec, newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(5, database.UserRoleMember)}))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestGetMyHeatStats(t *testing.T) {
	db := &dbtest.FakeDB{
		GetUserHeatStatsFn: func(_ uint, start, _ time.Time) (int64, decimal.Decimal, error) {
			// The total query starts at the zero time, the YTD query does not.
			if start.IsZero() {
				return 3661, mustDecimal(t, "120.50"), nil
			}
			return 61, mustDecimal(t, "10.00"), nil
		},
	}
	h := newTestHandler(t, db, nil)

	rec := httptest.NewRecorder()
	h.GetMyHeatStats(rec, newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(5, database.UserRoleMember)}))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data GetMyHeatStatsResponse
	decodeData(t, resp, &data)
	// 3661s -> 61.02min -> rounded up to 62min.
	if data.HeatTimeMinutesTotal != 62 {
		t.Errorf("heat_time_min = %d, want 62", data.HeatTimeMinutesTotal)
	}
	if !data.HeatCostTotal.Equal(mustDecimal(t, "120.50")) {
		t.Errorf("heat_cost = %s, want 120.50", data.HeatCostTotal)
	}
	// 61s -> 1.02min -> rounded up to 2min.
	if data.HeatTimeMinutesYTD != 2 {
		t.Errorf("heat_time_min_ytd = %d, want 2", data.HeatTimeMinutesYTD)
	}
	if !data.HeatCostYTD.Equal(mustDecimal(t, "10.00")) {
		t.Errorf("heat_cost_ytd = %s, want 10.00", data.HeatCostYTD)
	}
}

func TestGetMyHeatStatsFailures(t *testing.T) {
	tests := []struct {
		name    string
		session *database.BrowserSession
		db      *dbtest.FakeDB
	}{
		{
			name:    "no user in the session",
			session: validSession(0, database.UserRoleMember),
			db:      &dbtest.FakeDB{},
		},
		{
			name:    "stats cannot be read",
			session: validSession(5, database.UserRoleMember),
			db: &dbtest.FakeDB{
				GetUserHeatStatsFn: func(uint, time.Time, time.Time) (int64, decimal.Decimal, error) {
					return 0, decimal.Zero, errNotFound
				},
			},
		},
		{
			name:    "timezone cannot be read",
			session: validSession(5, database.UserRoleMember),
			db: &dbtest.FakeDB{
				GetTimezoneLocationFn: func() (*time.Location, error) { return nil, errNotFound },
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHandler(t, tt.db, nil)
			rec := httptest.NewRecorder()
			h.GetMyHeatStats(rec, newRequest("", &HandlerCtx{Database: tt.db, ValidSession: tt.session}))
			if resp := decodeResponse(t, rec); resp.OK {
				t.Error("expected a failure response")
			}
		})
	}
}

func TestGetMyBalance(t *testing.T) {
	db := &dbtest.FakeDB{
		GetUserHeatStatsFn: func(uint, time.Time, time.Time) (int64, decimal.Decimal, error) {
			return 0, mustDecimal(t, "80.00"), nil
		},
		GetUserSessionPaymentsFn: func(uint) (decimal.Decimal, error) { return mustDecimal(t, "100.00"), nil },
		GetUserSessionPaybacksFn: func(uint) (decimal.Decimal, error) { return mustDecimal(t, "5.00"), nil },
	}
	h := newTestHandler(t, db, nil)

	rec := httptest.NewRecorder()
	h.GetMyBalance(rec, newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(5, database.UserRoleMember)}))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data GetMyBalanceResponse
	decodeData(t, resp, &data)
	if !data.PaymentTotal.Equal(mustDecimal(t, "100.00")) {
		t.Errorf("payment_total = %s, want 100.00", data.PaymentTotal)
	}
	if !data.PaybackTotal.Equal(mustDecimal(t, "5.00")) {
		t.Errorf("payback_total = %s, want 5.00", data.PaybackTotal)
	}
	// 100 paid - 5 paid back - 80 used = 15.
	if !data.Balance.Equal(mustDecimal(t, "15.00")) {
		t.Errorf("balance_current = %s, want 15.00", data.Balance)
	}
}

func TestGetMyBalanceFailures(t *testing.T) {
	tests := []struct {
		name    string
		session *database.BrowserSession
		db      *dbtest.FakeDB
	}{
		{
			name:    "no user in the session",
			session: validSession(0, database.UserRoleMember),
			db:      &dbtest.FakeDB{},
		},
		{
			name:    "heat stats fail",
			session: validSession(5, database.UserRoleMember),
			db: &dbtest.FakeDB{
				GetUserHeatStatsFn: func(uint, time.Time, time.Time) (int64, decimal.Decimal, error) {
					return 0, decimal.Zero, errNotFound
				},
			},
		},
		{
			name:    "payments fail",
			session: validSession(5, database.UserRoleMember),
			db: &dbtest.FakeDB{
				GetUserSessionPaymentsFn: func(uint) (decimal.Decimal, error) { return decimal.Zero, errNotFound },
			},
		},
		{
			name:    "paybacks fail",
			session: validSession(5, database.UserRoleMember),
			db: &dbtest.FakeDB{
				GetUserSessionPaybacksFn: func(uint) (decimal.Decimal, error) { return decimal.Zero, errNotFound },
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHandler(t, tt.db, nil)
			rec := httptest.NewRecorder()
			h.GetMyBalance(rec, newRequest("", &HandlerCtx{Database: tt.db, ValidSession: tt.session}))
			if resp := decodeResponse(t, rec); resp.OK {
				t.Error("expected a failure response")
			}
		})
	}
}

func TestGetAllUsersShort(t *testing.T) {
	db := &dbtest.FakeDB{
		GetUsersFn: func(includeDeleted bool) ([]database.User, error) {
			if includeDeleted {
				t.Error("deleted users must not be returned")
			}
			return []database.User{
				{ID: 1, FirstName: "Ann", LastName: "Rider", Email: "ann@example.com"},
				{ID: 2, FirstName: "Bob", LastName: "Sailor"},
			}, nil
		},
	}
	h := newTestHandler(t, db, nil)

	rec := httptest.NewRecorder()
	h.GetAllUsersShort(rec, newRequest("", &HandlerCtx{Database: db}))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data []UserShort
	decodeData(t, resp, &data)
	if len(data) != 2 || data[0].FirstName != "Ann" || data[1].ID != 2 {
		t.Errorf("users = %+v, want the two configured users", data)
	}
}

func TestGetAllUsersShortFailure(t *testing.T) {
	db := &dbtest.FakeDB{
		GetUsersFn: func(bool) ([]database.User, error) { return nil, errNotFound },
	}
	h := newTestHandler(t, db, nil)
	rec := httptest.NewRecorder()
	h.GetAllUsersShort(rec, newRequest("", &HandlerCtx{Database: db}))
	if resp := decodeResponse(t, rec); resp.OK {
		t.Error("expected a failure response")
	}
}

func TestGetAllUsersDetailed(t *testing.T) {
	db := &dbtest.FakeDB{
		GetUsersFn: func(bool) ([]database.User, error) {
			return []database.User{{ID: 1, FirstName: "Ann", LastName: "Rider", UserStatusID: database.UserStatusMember}}, nil
		},
		GetUserHeatStatsFn: func(uint, time.Time, time.Time) (int64, decimal.Decimal, error) {
			return 3600, mustDecimal(t, "60.00"), nil
		},
		GetUserSessionPaymentsFn: func(uint) (decimal.Decimal, error) { return mustDecimal(t, "100.00"), nil },
		GetUserSessionPaybacksFn: func(uint) (decimal.Decimal, error) { return mustDecimal(t, "25.00"), nil },
	}
	h := newTestHandler(t, db, nil)

	rec := httptest.NewRecorder()
	h.GetAllUsersDetailed(rec, newRequest("", &HandlerCtx{Database: db}))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data []UserDetailed
	decodeData(t, resp, &data)
	if len(data) != 1 {
		t.Fatalf("users = %+v, want exactly one entry", data)
	}
	if data[0].TotalHeatSeconds != 3600 {
		t.Errorf("total_heat_seconds = %d, want 3600", data[0].TotalHeatSeconds)
	}
	if !data[0].TotalHeatCost.Equal(mustDecimal(t, "60.00")) {
		t.Errorf("total_heat_cost = %s, want 60.00", data[0].TotalHeatCost)
	}
	// Payments minus paybacks.
	if !data[0].TotalPayment.Equal(mustDecimal(t, "75.00")) {
		t.Errorf("total_payment = %s, want 75.00", data[0].TotalPayment)
	}
}

func TestGetAllUsersDetailedFailures(t *testing.T) {
	users := []database.User{{ID: 1}}
	tests := []struct {
		name string
		db   *dbtest.FakeDB
	}{
		{
			name: "users cannot be read",
			db:   &dbtest.FakeDB{GetUsersFn: func(bool) ([]database.User, error) { return nil, errNotFound }},
		},
		{
			name: "heat stats fail",
			db: &dbtest.FakeDB{
				GetUsersFn: func(bool) ([]database.User, error) { return users, nil },
				GetUserHeatStatsFn: func(uint, time.Time, time.Time) (int64, decimal.Decimal, error) {
					return 0, decimal.Zero, errNotFound
				},
			},
		},
		{
			name: "paybacks fail",
			db: &dbtest.FakeDB{
				GetUsersFn:               func(bool) ([]database.User, error) { return users, nil },
				GetUserSessionPaybacksFn: func(uint) (decimal.Decimal, error) { return decimal.Zero, errNotFound },
			},
		},
		{
			name: "payments fail",
			db: &dbtest.FakeDB{
				GetUsersFn:               func(bool) ([]database.User, error) { return users, nil },
				GetUserSessionPaymentsFn: func(uint) (decimal.Decimal, error) { return decimal.Zero, errNotFound },
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHandler(t, tt.db, nil)
			rec := httptest.NewRecorder()
			h.GetAllUsersDetailed(rec, newRequest("", &HandlerCtx{Database: tt.db}))
			if resp := decodeResponse(t, rec); resp.OK {
				t.Error("expected a failure response")
			}
		})
	}
}

func TestGetUserGroups(t *testing.T) {
	db := &dbtest.FakeDB{
		GetPricingsFn: func() ([]database.Pricing, error) {
			return []database.Pricing{
				{
					ID:             3,
					UserStatusID:   database.UserStatusMember,
					PricePerMinute: mustDecimal(t, "2.80"),
					Comment:        "Member Price",
					UserStatus: database.UserStatus{
						ID:          database.UserStatusMember,
						Name:        "Member",
						Description: "Member Access",
						UserRoleID:  database.UserRoleMember,
						UserRole:    database.UserRole{ID: database.UserRoleMember, Name: "member", Description: "Member User Permissions"},
					},
				},
			}, nil
		},
	}
	h := newTestHandler(t, db, nil)

	rec := httptest.NewRecorder()
	h.GetUserGroups(rec, newRequest("", &HandlerCtx{Database: db}))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data []GetUserGroupsResponse
	decodeData(t, resp, &data)
	if len(data) != 1 {
		t.Fatalf("user groups = %+v, want exactly one entry", data)
	}
	got := data[0]
	if got.PriceID != 3 || got.UserGroupID != database.UserStatusMember || got.UserGroupName != "Member" {
		t.Errorf("user group = %+v, does not match the pricing", got)
	}
	if !got.PricePerMinute.Equal(mustDecimal(t, "2.80")) {
		t.Errorf("price_min = %s, want 2.80", got.PricePerMinute)
	}
	if got.UserRoleID != database.UserRoleMember || got.UserRoleName != "member" {
		t.Errorf("role = %d/%q, want member", got.UserRoleID, got.UserRoleName)
	}
}

func TestGetUserGroupsFailure(t *testing.T) {
	db := &dbtest.FakeDB{GetPricingsFn: func() ([]database.Pricing, error) { return nil, errNotFound }}
	h := newTestHandler(t, db, nil)
	rec := httptest.NewRecorder()
	h.GetUserGroups(rec, newRequest("", &HandlerCtx{Database: db}))
	if resp := decodeResponse(t, rec); resp.OK {
		t.Error("expected a failure response")
	}
}

func TestCreateUserGroup(t *testing.T) {
	t.Run("creates the group and its pricing", func(t *testing.T) {
		var gotStatus database.UserStatus
		var gotPricing database.Pricing
		db := &dbtest.FakeDB{
			CreateUserGroupFn: func(us database.UserStatus, p database.Pricing) error {
				gotStatus, gotPricing = us, p
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		req := CreateUserGroupsRequest{
			PriceDescription:     "Crew Price",
			PricePerMinute:       mustDecimal(t, "1.50"),
			UserGroupDescription: "Crew members",
			UserGroupName:        "Crew",
			UserRoleID:           database.UserRoleMember,
		}
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.CreateUserGroup(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if gotStatus.Name != "Crew" || gotStatus.UserRoleID != database.UserRoleMember {
			t.Errorf("user status = %+v, does not match the request", gotStatus)
		}
		if !gotPricing.PricePerMinute.Equal(mustDecimal(t, "1.50")) {
			t.Errorf("price = %s, want 1.50", gotPricing.PricePerMinute)
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{
			CreateUserGroupFn: func(database.UserStatus, database.Pricing) error { return errNotFound },
		}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.CreateUserGroup(rec, r, CreateUserGroupsRequest{}, GetHandlerContext(r))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestChangeUserGroup(t *testing.T) {
	t.Run("updates the group and its pricing", func(t *testing.T) {
		var gotStatus database.UserStatus
		var gotPricing database.Pricing
		db := &dbtest.FakeDB{
			ChangeUserGroupFn: func(us database.UserStatus, p database.Pricing) error {
				gotStatus, gotPricing = us, p
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		req := ChangeUserGroupRequest{
			PriceDescription:     "Crew Price",
			PriceID:              7,
			PricePerMinute:       mustDecimal(t, "1.50"),
			UserGroupDescription: "Crew members",
			UserGroupID:          4,
			UserGroupName:        "Crew",
			UserRoleID:           database.UserRoleMember,
		}
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.ChangeUserGroup(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if gotStatus.ID != 4 || gotPricing.ID != 7 || gotPricing.UserStatusID != 4 {
			t.Errorf("status/pricing = %+v/%+v, does not link the group to its pricing", gotStatus, gotPricing)
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{
			ChangeUserGroupFn: func(database.UserStatus, database.Pricing) error { return errNotFound },
		}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.ChangeUserGroup(rec, r, ChangeUserGroupRequest{}, GetHandlerContext(r))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestDeleteUserGroup(t *testing.T) {
	t.Run("deletes the group", func(t *testing.T) {
		var deletedID uint
		db := &dbtest.FakeDB{
			DeleteUserGroupFn: func(id uint) error {
				deletedID = id
				return nil
			},
		}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.DeleteUserGroup(rec, r, DeleteUserGroupRequest{UserGroupID: 6}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if deletedID != 6 {
			t.Errorf("deleted group = %d, want 6", deletedID)
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{DeleteUserGroupFn: func(uint) error { return errNotFound }}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.DeleteUserGroup(rec, r, DeleteUserGroupRequest{UserGroupID: 6}, GetHandlerContext(r))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestSetUserGroup(t *testing.T) {
	pricings := []database.Pricing{
		{ID: 1, UserStatusID: database.UserStatusGuest},
		{ID: 2, UserStatusID: database.UserStatusMember},
	}
	user := database.User{ID: 3}

	t.Run("assigns a valid group", func(t *testing.T) {
		var gotUser, gotGroup uint
		db := &dbtest.FakeDB{
			GetUserByIdFn: func(uint) (database.User, error) { return user, nil },
			GetPricingsFn: func() ([]database.Pricing, error) { return pricings, nil },
			GetAdminUsersFn: func() ([]database.User, error) {
				return []database.User{{ID: 1}, {ID: 2}}, nil
			},
			SetUserGroupFn: func(userID, groupID uint) error {
				gotUser, gotGroup = userID, groupID
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.SetUserGroup(rec, r, SetUserGroupRequest{UserID: 3, UserGroupID: database.UserStatusMember}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if gotUser != 3 || gotGroup != database.UserStatusMember {
			t.Errorf("SetUserGroup(%d, %d), want (3, %d)", gotUser, gotGroup, database.UserStatusMember)
		}
	})

	t.Run("protects the last admin", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetUserByIdFn: func(uint) (database.User, error) { return user, nil },
			GetPricingsFn: func() ([]database.Pricing, error) { return pricings, nil },
			GetAdminUsersFn: func() ([]database.User, error) {
				return []database.User{{ID: 3}}, nil
			},
			SetUserGroupFn: func(uint, uint) error {
				t.Error("the group of the last admin must not be changed")
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.SetUserGroup(rec, r, SetUserGroupRequest{UserID: 3, UserGroupID: database.UserStatusMember}, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if resp.OK || resp.Msg != "Cannot change group of the only administrator." {
			t.Errorf("response = %+v, want the last-admin protection failure", resp)
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name string
			req  SetUserGroupRequest
			db   *dbtest.FakeDB
		}{
			{
				name: "unknown user",
				req:  SetUserGroupRequest{UserID: 3, UserGroupID: database.UserStatusMember},
				db:   &dbtest.FakeDB{},
			},
			{
				name: "pricings cannot be read",
				req:  SetUserGroupRequest{UserID: 3, UserGroupID: database.UserStatusMember},
				db: &dbtest.FakeDB{
					GetUserByIdFn: func(uint) (database.User, error) { return user, nil },
					GetPricingsFn: func() ([]database.Pricing, error) { return nil, errNotFound },
				},
			},
			{
				name: "unknown group",
				req:  SetUserGroupRequest{UserID: 3, UserGroupID: 99},
				db: &dbtest.FakeDB{
					GetUserByIdFn: func(uint) (database.User, error) { return user, nil },
					GetPricingsFn: func() ([]database.Pricing, error) { return pricings, nil },
				},
			},
			{
				name: "admin lookup fails",
				req:  SetUserGroupRequest{UserID: 3, UserGroupID: database.UserStatusMember},
				db: &dbtest.FakeDB{
					GetUserByIdFn:   func(uint) (database.User, error) { return user, nil },
					GetPricingsFn:   func() ([]database.Pricing, error) { return pricings, nil },
					GetAdminUsersFn: func() ([]database.User, error) { return nil, errNotFound },
				},
			},
			{
				name: "update fails",
				req:  SetUserGroupRequest{UserID: 3, UserGroupID: database.UserStatusMember},
				db: &dbtest.FakeDB{
					GetUserByIdFn:  func(uint) (database.User, error) { return user, nil },
					GetPricingsFn:  func() ([]database.Pricing, error) { return pricings, nil },
					SetUserGroupFn: func(uint, uint) error { return errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.SetUserGroup(rec, r, tt.req, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestGetUserRoles(t *testing.T) {
	db := &dbtest.FakeDB{
		GetUserRolesFn: func() ([]database.UserRole, error) { return database.DefaultUserRoles, nil },
	}
	h := newTestHandler(t, db, nil)

	rec := httptest.NewRecorder()
	h.GetUserRoles(rec, newRequest("", &HandlerCtx{Database: db}))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data []GetUserRolesResponse
	decodeData(t, resp, &data)
	if len(data) != len(database.DefaultUserRoles) {
		t.Fatalf("roles = %+v, want %d entries", data, len(database.DefaultUserRoles))
	}
	if data[0].UserRoleID != database.UserRoleGuest || data[0].UserRoleName != "guest" {
		t.Errorf("first role = %+v, want the guest role", data[0])
	}
}

func TestGetUserRolesFailure(t *testing.T) {
	db := &dbtest.FakeDB{GetUserRolesFn: func() ([]database.UserRole, error) { return nil, errNotFound }}
	h := newTestHandler(t, db, nil)
	rec := httptest.NewRecorder()
	h.GetUserRoles(rec, newRequest("", &HandlerCtx{Database: db}))
	if resp := decodeResponse(t, rec); resp.OK {
		t.Error("expected a failure response")
	}
}

func TestSetUserLock(t *testing.T) {
	user := database.User{ID: 3}

	t.Run("locks a user", func(t *testing.T) {
		var gotID uint
		var gotLocked bool
		db := &dbtest.FakeDB{
			GetUserByIdFn:   func(uint) (database.User, error) { return user, nil },
			GetAdminUsersFn: func() ([]database.User, error) { return []database.User{{ID: 1}, {ID: 2}}, nil },
			ChangeLockFn: func(id uint, locked bool) error {
				gotID, gotLocked = id, locked
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.SetUserLock(rec, r, SetUserLockRequest{UserID: 3, Locked: true}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if gotID != 3 || !gotLocked {
			t.Errorf("ChangeLock(%d, %v), want (3, true)", gotID, gotLocked)
		}
	})

	t.Run("protects the last admin", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetUserByIdFn:   func(uint) (database.User, error) { return user, nil },
			GetAdminUsersFn: func() ([]database.User, error) { return []database.User{{ID: 3}}, nil },
			ChangeLockFn: func(uint, bool) error {
				t.Error("the last admin must not be locked out")
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.SetUserLock(rec, r, SetUserLockRequest{UserID: 3, Locked: true}, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if resp.OK || resp.Msg != "Cannot lock/unlock the only administrator." {
			t.Errorf("response = %+v, want the last-admin protection failure", resp)
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{name: "unknown user", db: &dbtest.FakeDB{}},
			{
				name: "admin lookup fails",
				db: &dbtest.FakeDB{
					GetUserByIdFn:   func(uint) (database.User, error) { return user, nil },
					GetAdminUsersFn: func() ([]database.User, error) { return nil, errNotFound },
				},
			},
			{
				name: "lock update fails",
				db: &dbtest.FakeDB{
					GetUserByIdFn: func(uint) (database.User, error) { return user, nil },
					ChangeLockFn:  func(uint, bool) error { return errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.SetUserLock(rec, r, SetUserLockRequest{UserID: 3}, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestGetPasswordResetToken(t *testing.T) {
	user := database.User{ID: 8, Email: "user@example.com", FirstName: "Ann"}

	t.Run("stores a token for a known user", func(t *testing.T) {
		var stored database.PasswordReset
		db := &dbtest.FakeDB{
			GetUserByNameFn: func(string) (database.User, error) { return user, nil },
			AddPasswordResetTokenFn: func(p database.PasswordReset) error {
				stored = p
				return nil
			},
			GetEmailConfigurationFn: func() (database.EmailConfiguration, error) {
				return database.EmailConfiguration{}, errNotFound
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetPasswordResetToken(rec, r, GetPasswordResetTokenRequest{UserEmail: user.Email}, GetHandlerContext(r))

		if stored.UserID != user.ID {
			t.Errorf("token user = %d, want %d", stored.UserID, user.ID)
		}
		if stored.Token == "" {
			t.Error("no token was generated")
		}
		if !stored.Valid {
			t.Error("a fresh token has to be valid")
		}
		if !stored.ValidUntil.After(time.Now()) {
			t.Error("the token expiry has to be in the future")
		}
		if got := stored.ValidUntil.Sub(stored.IssuedAt); got != passwordResetValidity {
			t.Errorf("token validity = %s, want %s", got, passwordResetValidity)
		}
		if len(stored.Token) != 6 {
			t.Errorf("token = %q, want six digits so the whole keyspace is used", stored.Token)
		}
		// The email configuration is broken in this case, so the request fails.
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure when the email cannot be sent")
		}
	})

	t.Run("mails the token to the user", func(t *testing.T) {
		emailConfig, mails := startFakeSMTP(t)
		var stored database.PasswordReset
		db := &dbtest.FakeDB{
			GetUserByNameFn: func(string) (database.User, error) { return user, nil },
			AddPasswordResetTokenFn: func(p database.PasswordReset) error {
				stored = p
				return nil
			},
			GetEmailConfigurationFn: func() (database.EmailConfiguration, error) { return emailConfig, nil },
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetPasswordResetToken(rec, r, GetPasswordResetTokenRequest{UserEmail: user.Email}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		sent := mails()
		if len(sent) != 1 {
			t.Fatalf("delivered mails = %d, want exactly one token mail", len(sent))
		}
		if len(sent[0].To) != 1 || sent[0].To[0] != user.Email {
			t.Errorf("recipients = %v, want [%s]", sent[0].To, user.Email)
		}
		if !strings.Contains(sent[0].Body, stored.Token) {
			t.Errorf("the mail body does not contain the generated token %q", stored.Token)
		}
	})

	t.Run("does not leak whether a user exists", func(t *testing.T) {
		db := &dbtest.FakeDB{
			AddPasswordResetTokenFn: func(database.PasswordReset) error {
				t.Error("no token may be created for an unknown user")
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetPasswordResetToken(rec, r, GetPasswordResetTokenRequest{UserEmail: "nobody@example.com"}, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Error("an unknown email address has to yield the same success response as a known one")
		}
	})

	t.Run("reports a failing token store", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetUserByNameFn:         func(string) (database.User, error) { return user, nil },
			AddPasswordResetTokenFn: func(database.PasswordReset) error { return errNotFound },
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetPasswordResetToken(rec, r, GetPasswordResetTokenRequest{UserEmail: user.Email}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestSetPasswordWithToken(t *testing.T) {
	const newPassword = "N3wSecret!Password"
	user := database.User{ID: 8, Email: "user@example.com"}
	validToken := database.PasswordReset{
		ID:         3,
		UserID:     8,
		Token:      "123456",
		IssuedAt:   time.Now(),
		ValidUntil: time.Now().Add(passwordResetValidity),
		Valid:      true,
	}

	t.Run("sets the password and invalidates the tokens", func(t *testing.T) {
		var stored database.User
		invalidated := false
		db := &dbtest.FakeDB{
			GetUserByNameFn:               func(string) (database.User, error) { return user, nil },
			GetActivePasswordResetEntryFn: func(uint) (database.PasswordReset, error) { return validToken, nil },
			UpdatePasswordFn: func(_ uint, u database.User) error {
				stored = u
				return nil
			},
			InvalidatePasswordResetEntriesFn: func(uint) error {
				invalidated = true
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		req := SetPasswordWithTokenRequest{UserEmail: user.Email, Password: newPassword, Token: validToken.Token}
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.SetPasswordWithToken(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("password reset failed: %s", resp.Msg)
		}
		expected, err := hash.CryptSha512(hash.Sha256(newPassword), strconv.Itoa(stored.PasswordSalt))
		if err != nil {
			t.Fatalf("cannot recompute the hash: %v", err)
		}
		if stored.PasswordHash != expected {
			t.Error("the stored hash does not match the new password")
		}
		if !invalidated {
			t.Error("the used reset tokens have to be invalidated")
		}
	})

	t.Run("rejects invalid tokens", func(t *testing.T) {
		expired := validToken
		expired.ValidUntil = time.Now().Add(-time.Minute)
		invalid := validToken
		invalid.Valid = false
		mismatched := validToken
		mismatched.Token = "999999"

		tests := []struct {
			name  string
			db    *dbtest.FakeDB
			token string
		}{
			{
				name:  "unknown user",
				db:    &dbtest.FakeDB{},
				token: validToken.Token,
			},
			{
				name: "token not found",
				db: &dbtest.FakeDB{
					GetUserByNameFn:               func(string) (database.User, error) { return user, nil },
					GetActivePasswordResetEntryFn: func(uint) (database.PasswordReset, error) { return database.PasswordReset{}, errNotFound },
				},
				token: validToken.Token,
			},
			{
				name: "token mismatch",
				db: &dbtest.FakeDB{
					GetUserByNameFn:               func(string) (database.User, error) { return user, nil },
					GetActivePasswordResetEntryFn: func(uint) (database.PasswordReset, error) { return mismatched, nil },
				},
				token: validToken.Token,
			},
			{
				name: "token invalidated",
				db: &dbtest.FakeDB{
					GetUserByNameFn:               func(string) (database.User, error) { return user, nil },
					GetActivePasswordResetEntryFn: func(uint) (database.PasswordReset, error) { return invalid, nil },
				},
				token: validToken.Token,
			},
			{
				name: "token expired",
				db: &dbtest.FakeDB{
					GetUserByNameFn:               func(string) (database.User, error) { return user, nil },
					GetActivePasswordResetEntryFn: func(uint) (database.PasswordReset, error) { return expired, nil },
				},
				token: validToken.Token,
			},
			{
				name: "password cannot be stored",
				db: &dbtest.FakeDB{
					GetUserByNameFn:               func(string) (database.User, error) { return user, nil },
					GetActivePasswordResetEntryFn: func(uint) (database.PasswordReset, error) { return validToken, nil },
					UpdatePasswordFn:              func(uint, database.User) error { return errNotFound },
				},
				token: validToken.Token,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				changed := false
				inner := tt.db.UpdatePasswordFn
				tt.db.UpdatePasswordFn = func(id uint, u database.User) error {
					changed = true
					if inner != nil {
						return inner(id, u)
					}
					return nil
				}
				h := newTestHandler(t, tt.db, nil)

				req := SetPasswordWithTokenRequest{UserEmail: user.Email, Password: newPassword, Token: tt.token}
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.SetPasswordWithToken(rec, r, req, GetHandlerContext(r))

				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
				if tt.name != "password cannot be stored" && changed {
					t.Error("the password must not be changed for an invalid token")
				}
			})
		}
	})
}

func TestSetPasswordWithTokenValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{name: "valid payload", payload: `{"email":"a@b.co","password":"N3wSecret!Pwd","token":"123456"}`},
		{name: "weak password", payload: `{"email":"a@b.co","password":"weak","token":"123456"}`, wantErr: true},
		{name: "missing token", payload: `{"email":"a@b.co","password":"N3wSecret!Pwd"}`, wantErr: true},
		{name: "invalid email", payload: `{"email":"nope","password":"N3wSecret!Pwd","token":"123456"}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReadBodyAndValidate(newRequest(tt.payload, nil), &SetPasswordWithTokenRequest{}, SetPasswordWithTokenValidationErrors)
			if (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPasswordResetTokenLifecycle(t *testing.T) {
	user := database.User{ID: 8, Email: "user@example.com", FirstName: "Ann"}

	// newDB returns a fake wired for a successful token request, plus the call
	// log the assertions below inspect.
	newDB := func(active database.PasswordReset, activeErr error) (*dbtest.FakeDB, *[]string, *database.PasswordReset) {
		var calls []string
		var stored database.PasswordReset
		emailConfig, _ := startFakeSMTP(t)
		db := &dbtest.FakeDB{
			GetUserByNameFn: func(string) (database.User, error) { return user, nil },
			GetActivePasswordResetEntryFn: func(uint) (database.PasswordReset, error) {
				return active, activeErr
			},
			InvalidatePasswordResetEntriesFn: func(uint) error {
				calls = append(calls, "invalidate")
				return nil
			},
			AddPasswordResetTokenFn: func(p database.PasswordReset) error {
				calls = append(calls, "add")
				stored = p
				return nil
			},
			GetEmailConfigurationFn: func() (database.EmailConfiguration, error) { return emailConfig, nil },
		}
		return db, &calls, &stored
	}

	request := func(t *testing.T, db *dbtest.FakeDB) apiResponse {
		t.Helper()
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetPasswordResetToken(rec, r, GetPasswordResetTokenRequest{UserEmail: user.Email}, GetHandlerContext(r))
		return decodeResponse(t, rec)
	}

	t.Run("retires the previous tokens before issuing a new one", func(t *testing.T) {
		stale := database.PasswordReset{ID: 1, UserID: user.ID, IssuedAt: time.Now().Add(-time.Hour), Valid: true}
		db, calls, _ := newDB(stale, nil)

		if resp := request(t, db); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		want := []string{"invalidate", "add"}
		if !reflect.DeepEqual(*calls, want) {
			t.Errorf("calls = %v, want %v (a stacked token shrinks the space an attacker has to search)", *calls, want)
		}
	})

	t.Run("keeps the live token inside the cooldown window", func(t *testing.T) {
		fresh := database.PasswordReset{ID: 1, UserID: user.ID, IssuedAt: time.Now().Add(-30 * time.Second), Valid: true}
		db, calls, _ := newDB(fresh, nil)

		resp := request(t, db)
		if !resp.OK {
			t.Fatalf("a request inside the cooldown still has to look successful, got: %s", resp.Msg)
		}
		if len(*calls) != 0 {
			t.Errorf("calls = %v, want none: invalidating here would destroy a token that was just mailed", *calls)
		}
	})

	t.Run("issues again once the cooldown has passed", func(t *testing.T) {
		old := database.PasswordReset{ID: 1, UserID: user.ID, IssuedAt: time.Now().Add(-passwordResetCooldown - time.Second), Valid: true}
		db, calls, stored := newDB(old, nil)

		if resp := request(t, db); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if len(*calls) != 2 {
			t.Fatalf("calls = %v, want invalidate then add", *calls)
		}
		if stored.Token == "" {
			t.Error("a replacement token has to be issued after the cooldown")
		}
	})

	t.Run("answers identically for unknown, cooled down and issued", func(t *testing.T) {
		fresh := database.PasswordReset{ID: 1, UserID: user.ID, IssuedAt: time.Now(), Valid: true}
		cooling, _, _ := newDB(fresh, nil)
		issuing, _, _ := newDB(database.PasswordReset{}, errNotFound)
		unknown := &dbtest.FakeDB{}

		messages := map[string]string{
			"cooldown": request(t, cooling).Msg,
			"issued":   request(t, issuing).Msg,
			"unknown":  request(t, unknown).Msg,
		}
		for name, msg := range messages {
			if msg != messages["issued"] {
				t.Errorf("%s answered %q, want %q: a differing answer reveals whether the account exists", name, msg, messages["issued"])
			}
		}
	})
}

func TestSetPasswordWithTokenAttemptBudget(t *testing.T) {
	user := database.User{ID: 8, Email: "user@example.com"}
	live := database.PasswordReset{
		ID:         3,
		UserID:     user.ID,
		Token:      "123456",
		IssuedAt:   time.Now(),
		ValidUntil: time.Now().Add(passwordResetValidity),
		Valid:      true,
	}

	redeem := func(t *testing.T, db *dbtest.FakeDB, token string) apiResponse {
		t.Helper()
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		req := SetPasswordWithTokenRequest{UserEmail: user.Email, Password: "N3wSecret!Password", Token: token}
		h.SetPasswordWithToken(rec, r, req, GetHandlerContext(r))
		return decodeResponse(t, rec)
	}

	t.Run("spends an attempt on a wrong token", func(t *testing.T) {
		var gotID, gotMax uint
		registered := false
		db := &dbtest.FakeDB{
			GetUserByNameFn:               func(string) (database.User, error) { return user, nil },
			GetActivePasswordResetEntryFn: func(uint) (database.PasswordReset, error) { return live, nil },
			RegisterFailedPasswordResetAttemptFn: func(id uint, maxAttempts uint) error {
				registered, gotID, gotMax = true, id, maxAttempts
				return nil
			},
			UpdatePasswordFn: func(uint, database.User) error {
				t.Error("the password must not change on a wrong token")
				return nil
			},
		}

		if resp := redeem(t, db, "999999"); resp.OK {
			t.Fatal("a wrong token must not reset the password")
		}
		if !registered {
			t.Fatal("a wrong token has to be counted, otherwise the keyspace can be searched")
		}
		if gotID != live.ID || gotMax != passwordResetMaxAttempts {
			t.Errorf("registered attempt against (id=%d, max=%d), want (id=%d, max=%d)", gotID, gotMax, live.ID, passwordResetMaxAttempts)
		}
	})

	t.Run("refuses a token whose attempts are used up", func(t *testing.T) {
		burned := live
		burned.Attempts = passwordResetMaxAttempts
		db := &dbtest.FakeDB{
			GetUserByNameFn:               func(string) (database.User, error) { return user, nil },
			GetActivePasswordResetEntryFn: func(uint) (database.PasswordReset, error) { return burned, nil },
			UpdatePasswordFn: func(uint, database.User) error {
				t.Error("a burned token must not reset the password, even with the right value")
				return nil
			},
		}

		if resp := redeem(t, db, burned.Token); resp.OK {
			t.Error("a token that used up its attempts has to be refused")
		}
	})

	t.Run("does not count an attempt when there is no live token", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetUserByNameFn:               func(string) (database.User, error) { return user, nil },
			GetActivePasswordResetEntryFn: func(uint) (database.PasswordReset, error) { return database.PasswordReset{}, errNotFound },
			RegisterFailedPasswordResetAttemptFn: func(uint, uint) error {
				t.Error("there is no token to count an attempt against")
				return nil
			},
		}

		if resp := redeem(t, db, "123456"); resp.OK {
			t.Error("expected a failure when no token is live")
		}
	})
}
