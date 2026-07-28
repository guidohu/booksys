package handlers

import (
	"net/http/httptest"
	"testing"
	"time"

	"server/database"
	"server/database/dbtest"
)

func TestGetSession(t *testing.T) {
	start := time.Date(2024, time.May, 4, 10, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)
	session := database.Session{
		ID:         3,
		StartTime:  start,
		EndTime:    end,
		Title:      "Morning session",
		Comment:    "flat water",
		FreeSpaces: 4,
		CreatorID:  1,
		Creator:    database.User{ID: 1, FirstName: "Ann", LastName: "Rider"},
	}

	t.Run("returns the session with its riders", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetSessionFn: func(uint) (database.Session, error) { return session, nil },
			GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) {
				return []database.UserToSession{
					{UserID: 5, User: database.User{ID: 5, FirstName: "Bob", LastName: "Sailor"}},
				}, nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetSession(rec, r, GetSessionRequest{SessionID: 3}, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data SessionResponse
		decodeData(t, resp, &data)
		if data.ID != 3 || data.Title != "Morning session" {
			t.Errorf("session = %+v, does not match the database entry", data)
		}
		if data.Duration != 7200 {
			t.Errorf("duration = %d, want 7200", data.Duration)
		}
		if data.CreatorFirstName != "Ann" {
			t.Errorf("creator = %q, want Ann", data.CreatorFirstName)
		}
		if len(data.Riders) != 1 || data.Riders[0].Name != "Bob Sailor" {
			t.Errorf("riders = %+v, want Bob Sailor", data.Riders)
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{
				name: "session does not exist",
				db: &dbtest.FakeDB{
					GetSessionFn: func(uint) (database.Session, error) { return database.Session{}, errNotFound },
				},
			},
			{
				name: "riders cannot be read",
				db: &dbtest.FakeDB{
					GetSessionFn: func(uint) (database.Session, error) { return session, nil },
					GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) {
						return nil, errNotFound
					},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.GetSession(rec, r, GetSessionRequest{SessionID: 3}, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestGetSessionMetadata(t *testing.T) {
	session := database.Session{
		ID:        3,
		StartTime: time.Date(2024, time.June, 21, 10, 0, 0, 0, time.UTC),
	}

	t.Run("returns sunrise and sunset for the session day", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetSessionFn: func(uint) (database.Session, error) { return session, nil },
		}
		h := newTestHandler(t, db, map[string]string{
			"location.timezone":  "Europe/Zurich",
			"location.latitude":  "47.3769",
			"location.longitude": "8.5417",
		})

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetSessionMetadata(rec, r, GetSessionRequest{SessionID: 3}, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data GetSessionMetadataResponse
		decodeData(t, resp, &data)
		if data.SunriseTime >= data.SunsetTime {
			t.Errorf("sunrise (%d) has to be before sunset (%d)", data.SunriseTime, data.SunsetTime)
		}
		// Around midsummer in Zurich the day is roughly 15h45m long.
		dayLength := time.Duration(data.SunsetTime-data.SunriseTime) * time.Second
		if dayLength < 15*time.Hour || dayLength > 16*time.Hour {
			t.Errorf("day length = %v, want roughly 15-16h for the summer solstice in Zurich", dayLength)
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name   string
			db     *dbtest.FakeDB
			values map[string]string
		}{
			{
				name: "session does not exist",
				db: &dbtest.FakeDB{
					GetSessionFn: func(uint) (database.Session, error) { return database.Session{}, errNotFound },
				},
				values: map[string]string{"location.timezone": "Europe/Zurich"},
			},
			{
				name: "invalid timezone",
				db: &dbtest.FakeDB{
					GetSessionFn: func(uint) (database.Session, error) { return session, nil },
				},
				values: map[string]string{"location.timezone": "Not/AZone"},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, tt.values)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.GetSessionMetadata(rec, r, GetSessionRequest{SessionID: 3}, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestCreateSession(t *testing.T) {
	req := CreateSessionRequest{
		Title:     "Evening session",
		Comment:   "bring a towel",
		Start:     time.Date(2024, time.May, 4, 18, 0, 0, 0, time.UTC).Unix(),
		End:       time.Date(2024, time.May, 4, 20, 0, 0, 0, time.UTC).Unix(),
		MaxRiders: 4,
		Type:      database.DefaultSessionType,
	}

	t.Run("creates a session", func(t *testing.T) {
		var created database.Session
		db := &dbtest.FakeDB{
			CreateSessionFn: func(s database.Session) (uint, error) {
				created = s
				return 42, nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(7, database.UserRoleAdmin)})
		h.CreateSession(rec, r, req, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data CreateSessionResponse
		decodeData(t, resp, &data)
		if data.SessionID != 42 {
			t.Errorf("session_id = %d, want 42", data.SessionID)
		}
		if created.CreatorID != 7 {
			t.Errorf("creator = %d, want the user of the session (7)", created.CreatorID)
		}
		if created.FreeSpaces != 4 || created.SessionTypeID != database.DefaultSessionType {
			t.Errorf("session = %+v, does not match the request", created)
		}
		if !created.StartTime.Equal(time.Unix(req.Start, 0)) {
			t.Errorf("start = %v, want %v", created.StartTime, time.Unix(req.Start, 0))
		}
	})

	t.Run("refuses colliding sessions", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetSessionsBetweenFn: func(time.Time, time.Time) ([]database.Session, error) {
				return []database.Session{{ID: 1}}, nil
			},
			CreateSessionFn: func(database.Session) (uint, error) {
				t.Error("a colliding session must not be created")
				return 0, nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(7, database.UserRoleAdmin)})
		h.CreateSession(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{
				name: "collision check fails",
				db: &dbtest.FakeDB{
					GetSessionsBetweenFn: func(time.Time, time.Time) ([]database.Session, error) {
						return nil, errNotFound
					},
				},
			},
			{
				name: "session cannot be created",
				db: &dbtest.FakeDB{
					CreateSessionFn: func(database.Session) (uint, error) { return 0, errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db, ValidSession: validSession(7, database.UserRoleAdmin)})
				h.CreateSession(rec, r, req, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestEditSession(t *testing.T) {
	existing := database.Session{ID: 3, CreatorID: 9}
	req := EditSessionRequest{
		SessionID: 3,
		Title:     "Changed",
		Start:     time.Date(2024, time.May, 4, 18, 0, 0, 0, time.UTC).Unix(),
		End:       time.Date(2024, time.May, 4, 20, 0, 0, 0, time.UTC).Unix(),
		MaxRiders: 5,
		Type:      database.DefaultSessionType,
	}

	t.Run("updates the session and keeps the original creator", func(t *testing.T) {
		var updated database.Session
		db := &dbtest.FakeDB{
			GetSessionFn: func(uint) (database.Session, error) { return existing, nil },
			GetSessionsBetweenFn: func(time.Time, time.Time) ([]database.Session, error) {
				// Only the session itself is in the new time window.
				return []database.Session{existing}, nil
			},
			UpdateSessionFn: func(s database.Session) error {
				updated = s
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(7, database.UserRoleAdmin)})
		h.EditSession(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if updated.ID != 3 || updated.Title != "Changed" || updated.FreeSpaces != 5 {
			t.Errorf("session = %+v, does not match the request", updated)
		}
		if updated.CreatorID != 9 {
			t.Errorf("creator = %d, want the original creator (9)", updated.CreatorID)
		}
	})

	t.Run("refuses a collision with a different session", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetSessionFn: func(uint) (database.Session, error) { return existing, nil },
			GetSessionsBetweenFn: func(time.Time, time.Time) ([]database.Session, error) {
				return []database.Session{{ID: 4}}, nil
			},
			UpdateSessionFn: func(database.Session) error {
				t.Error("a colliding session must not be saved")
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(7, database.UserRoleAdmin)})
		h.EditSession(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{
				name: "session does not exist",
				db: &dbtest.FakeDB{
					GetSessionFn: func(uint) (database.Session, error) { return database.Session{}, errNotFound },
				},
			},
			{
				name: "collision check fails",
				db: &dbtest.FakeDB{
					GetSessionFn: func(uint) (database.Session, error) { return existing, nil },
					GetSessionsBetweenFn: func(time.Time, time.Time) ([]database.Session, error) {
						return nil, errNotFound
					},
				},
			},
			{
				name: "session cannot be saved",
				db: &dbtest.FakeDB{
					GetSessionFn:    func(uint) (database.Session, error) { return existing, nil },
					UpdateSessionFn: func(database.Session) error { return errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db, ValidSession: validSession(7, database.UserRoleAdmin)})
				h.EditSession(rec, r, req, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestDeleteSession(t *testing.T) {
	t.Run("removes riders and the session", func(t *testing.T) {
		ridersRemoved, sessionRemoved := false, false
		db := &dbtest.FakeDB{
			GetHeatInSessionCountFn: func(uint) (int, error) { return 0, nil },
			DeleteUsersFromSessionFn: func(uint) error {
				ridersRemoved = true
				return nil
			},
			DeleteSessionFn: func(uint) error {
				if !ridersRemoved {
					t.Error("the riders have to be removed before the session")
				}
				sessionRemoved = true
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.DeleteSession(rec, r, DeleteSessionRequest{SessionID: 3}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if !sessionRemoved {
			t.Error("the session was not removed")
		}
	})

	t.Run("refuses a session that already has heats", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetHeatInSessionCountFn: func(uint) (int, error) { return 2, nil },
			DeleteSessionFn: func(uint) error {
				t.Error("a session with heats must not be removed")
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.DeleteSession(rec, r, DeleteSessionRequest{SessionID: 3}, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if resp.OK || resp.Msg != "This session already has heats. It cannot be removed." {
			t.Errorf("response = %+v, want the heat protection failure", resp)
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{
				name: "heat count fails",
				db: &dbtest.FakeDB{
					GetHeatInSessionCountFn: func(uint) (int, error) { return 0, errNotFound },
				},
			},
			{
				name: "riders cannot be removed",
				db: &dbtest.FakeDB{
					DeleteUsersFromSessionFn: func(uint) error { return errNotFound },
				},
			},
			{
				name: "session cannot be removed",
				db: &dbtest.FakeDB{
					DeleteSessionFn: func(uint) error { return errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.DeleteSession(rec, r, DeleteSessionRequest{SessionID: 3}, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestAddUserToSession(t *testing.T) {
	session := database.Session{ID: 3, FreeSpaces: 2}

	t.Run("adds the requested users", func(t *testing.T) {
		var added []database.UserToSession
		db := &dbtest.FakeDB{
			GetSessionFn:         func(uint) (database.Session, error) { return session, nil },
			GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) { return nil, nil },
			GetUserByIdFn: func(id uint) (database.User, error) {
				return database.User{ID: id, Email: "rider@example.com"}, nil
			},
			AddSessionToUserEntryFn: func(u database.UserToSession) error {
				added = append(added, u)
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.AddUserToSession(rec, r, AddSessionUserRequest{SessionID: 3, UserIDs: []uint{5, 6}}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if len(added) != 2 {
			t.Fatalf("added users = %+v, want 2", added)
		}
		if added[0].SessionID != 3 || added[0].UserID != 5 || added[1].UserID != 6 {
			t.Errorf("added users = %+v, do not match the request", added)
		}
	})

	t.Run("skips unknown and already listed users", func(t *testing.T) {
		var added []database.UserToSession
		db := &dbtest.FakeDB{
			GetSessionFn: func(uint) (database.Session, error) { return session, nil },
			GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) {
				return []database.UserToSession{{UserID: 5, User: database.User{ID: 5}}}, nil
			},
			GetUserByIdFn: func(id uint) (database.User, error) {
				if id == 99 {
					return database.User{}, errNotFound
				}
				return database.User{ID: id}, nil
			},
			AddSessionToUserEntryFn: func(u database.UserToSession) error {
				added = append(added, u)
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.AddUserToSession(rec, r, AddSessionUserRequest{SessionID: 3, UserIDs: []uint{5, 6, 99}}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if len(added) != 1 || added[0].UserID != 6 {
			t.Errorf("added users = %+v, want only user 6", added)
		}
	})

	t.Run("refuses to overbook a session", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetSessionFn:         func(uint) (database.Session, error) { return database.Session{ID: 3, FreeSpaces: 1}, nil },
			GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) { return nil, nil },
			GetUserByIdFn:        func(id uint) (database.User, error) { return database.User{ID: id}, nil },
			AddSessionToUserEntryFn: func(database.UserToSession) error {
				t.Error("no user may be added when the session is full")
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.AddUserToSession(rec, r, AddSessionUserRequest{SessionID: 3, UserIDs: []uint{5, 6}}, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if resp.OK || resp.Msg != "There is not enough space to add all the users." {
			t.Errorf("response = %+v, want the capacity failure", resp)
		}
	})

	t.Run("notifies the user that was added", func(t *testing.T) {
		// The invitation email has to be addressed to the user that is being
		// added, not to an unrelated (or empty) user record.
		emailConfig, mails := startFakeSMTP(t)
		db := &dbtest.FakeDB{
			GetSessionFn:         func(uint) (database.Session, error) { return session, nil },
			GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) { return nil, nil },
			GetUserByIdFn: func(id uint) (database.User, error) {
				return database.User{ID: id, FirstName: "Bob", LastName: "Sailor", Email: "bob@example.com"}, nil
			},
			GetEmailConfigurationFn: func() (database.EmailConfiguration, error) { return emailConfig, nil },
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.AddUserToSession(rec, r, AddSessionUserRequest{SessionID: 3, UserIDs: []uint{6}}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		sent := mails()
		if len(sent) != 1 {
			t.Fatalf("delivered mails = %d, want exactly one invitation", len(sent))
		}
		if len(sent[0].To) != 1 || sent[0].To[0] != "bob@example.com" {
			t.Errorf("recipients = %v, want [bob@example.com]", sent[0].To)
		}
		if sent[0].From != emailConfig.Sender {
			t.Errorf("sender = %q, want %q", sent[0].From, emailConfig.Sender)
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{
				name: "session does not exist",
				db: &dbtest.FakeDB{
					GetSessionFn: func(uint) (database.Session, error) { return database.Session{}, errNotFound },
				},
			},
			{
				name: "existing riders cannot be read",
				db: &dbtest.FakeDB{
					GetSessionFn: func(uint) (database.Session, error) { return session, nil },
					GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) {
						return nil, errNotFound
					},
				},
			},
			{
				name: "email configuration cannot be read",
				db: &dbtest.FakeDB{
					GetSessionFn:         func(uint) (database.Session, error) { return session, nil },
					GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) { return nil, nil },
					GetUserByIdFn:        func(id uint) (database.User, error) { return database.User{ID: id}, nil },
					GetEmailConfigurationFn: func() (database.EmailConfiguration, error) {
						return database.EmailConfiguration{}, errNotFound
					},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.AddUserToSession(rec, r, AddSessionUserRequest{SessionID: 3, UserIDs: []uint{5}}, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestRemoveUserFromSession(t *testing.T) {
	session := database.Session{ID: 3}

	t.Run("removes the user", func(t *testing.T) {
		var removedUser, removedSession uint
		db := &dbtest.FakeDB{
			GetUserByIdFn: func(id uint) (database.User, error) { return database.User{ID: id}, nil },
			GetSessionFn:  func(uint) (database.Session, error) { return session, nil },
			DeleteSessionToUserEntryFn: func(userID, sessionID uint) error {
				removedUser, removedSession = userID, sessionID
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(1, database.UserRoleAdmin)})
		h.RemoveUserFromSession(rec, r, RemoveSessionUserRequest{SessionID: 3, UserID: 5}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if removedUser != 5 || removedSession != 3 {
			t.Errorf("DeleteSessionToUserEntry(%d, %d), want (5, 3)", removedUser, removedSession)
		}
	})

	t.Run("refuses to remove a user that has heats in the session", func(t *testing.T) {
		// An admin removes user 5 while being logged in as user 1. The heats of
		// user 5 are what matters here, not the heats of the acting admin.
		db := &dbtest.FakeDB{
			GetUserByIdFn: func(id uint) (database.User, error) { return database.User{ID: id}, nil },
			GetSessionFn:  func(uint) (database.Session, error) { return session, nil },
			GetUserHeatsBySessionFn: func(_, userID uint, _ int) ([]database.Heat, error) {
				return []database.Heat{{ID: 1, UserID: userID, SessionID: 3}}, nil
			},
			DeleteSessionToUserEntryFn: func(uint, uint) error {
				t.Error("a user with heats in the session must not be removed")
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(1, database.UserRoleAdmin)})
		h.RemoveUserFromSession(rec, r, RemoveSessionUserRequest{SessionID: 3, UserID: 5}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{
				name: "user does not exist",
				db:   &dbtest.FakeDB{},
			},
			{
				name: "session does not exist",
				db: &dbtest.FakeDB{
					GetUserByIdFn: func(id uint) (database.User, error) { return database.User{ID: id}, nil },
					GetSessionFn:  func(uint) (database.Session, error) { return database.Session{}, errNotFound },
				},
			},
			{
				name: "heats cannot be read",
				db: &dbtest.FakeDB{
					GetUserByIdFn: func(id uint) (database.User, error) { return database.User{ID: id}, nil },
					GetSessionFn:  func(uint) (database.Session, error) { return session, nil },
					GetUserHeatsBySessionFn: func(uint, uint, int) ([]database.Heat, error) {
						return nil, errNotFound
					},
				},
			},
			{
				name: "entry cannot be removed",
				db: &dbtest.FakeDB{
					GetUserByIdFn:              func(id uint) (database.User, error) { return database.User{ID: id}, nil },
					GetSessionFn:               func(uint) (database.Session, error) { return session, nil },
					DeleteSessionToUserEntryFn: func(uint, uint) error { return errNotFound },
				},
			},
			{
				name: "email configuration cannot be read",
				db: &dbtest.FakeDB{
					GetUserByIdFn: func(id uint) (database.User, error) { return database.User{ID: id}, nil },
					GetSessionFn:  func(uint) (database.Session, error) { return session, nil },
					GetEmailConfigurationFn: func() (database.EmailConfiguration, error) {
						return database.EmailConfiguration{}, errNotFound
					},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db, ValidSession: validSession(1, database.UserRoleAdmin)})
				h.RemoveUserFromSession(rec, r, RemoveSessionUserRequest{SessionID: 3, UserID: 5}, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestRemoveMyUserFromSession(t *testing.T) {
	t.Run("removes the calling user", func(t *testing.T) {
		var removedUser, removedSession uint
		db := &dbtest.FakeDB{
			GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) {
				return []database.UserToSession{{UserID: 5}, {UserID: 6}}, nil
			},
			DeleteSessionToUserEntryFn: func(userID, sessionID uint) error {
				removedUser, removedSession = userID, sessionID
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(5, database.UserRoleMember)})
		h.RemoveMyUserFromSession(rec, r, RemoveSessionMyUserRequest{SessionID: 3}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if removedUser != 5 || removedSession != 3 {
			t.Errorf("DeleteSessionToUserEntry(%d, %d), want (5, 3)", removedUser, removedSession)
		}
	})

	t.Run("does not touch other riders", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) {
				return []database.UserToSession{{UserID: 6}}, nil
			},
			DeleteSessionToUserEntryFn: func(uint, uint) error {
				t.Error("a user that is not part of the session must not remove anybody")
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(5, database.UserRoleMember)})
		h.RemoveMyUserFromSession(rec, r, RemoveSessionMyUserRequest{SessionID: 3}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Errorf("unexpected failure: %s", resp.Msg)
		}
	})

	t.Run("refuses when the user already has heats", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) {
				return []database.UserToSession{{UserID: 5}}, nil
			},
			GetUserHeatsBySessionFn: func(uint, uint, int) ([]database.Heat, error) {
				return []database.Heat{{ID: 1, UserID: 5}}, nil
			},
			DeleteSessionToUserEntryFn: func(uint, uint) error {
				t.Error("a user with heats must not leave the session")
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db, ValidSession: validSession(5, database.UserRoleMember)})
		h.RemoveMyUserFromSession(rec, r, RemoveSessionMyUserRequest{SessionID: 3}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{
				name: "riders cannot be read",
				db: &dbtest.FakeDB{
					GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) { return nil, errNotFound },
				},
			},
			{
				name: "heats cannot be read",
				db: &dbtest.FakeDB{
					GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) {
						return []database.UserToSession{{UserID: 5}}, nil
					},
					GetUserHeatsBySessionFn: func(uint, uint, int) ([]database.Heat, error) {
						return nil, errNotFound
					},
				},
			},
			{
				name: "entry cannot be removed",
				db: &dbtest.FakeDB{
					GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) {
						return []database.UserToSession{{UserID: 5}}, nil
					},
					DeleteSessionToUserEntryFn: func(uint, uint) error { return errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db, ValidSession: validSession(5, database.UserRoleMember)})
				h.RemoveMyUserFromSession(rec, r, RemoveSessionMyUserRequest{SessionID: 3}, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestGetSessionHeats(t *testing.T) {
	heats := []database.Heat{
		{
			ID:              1,
			UserID:          5,
			SessionID:       3,
			DurationSeconds: 600,
			Cost:            mustDecimal(t, "28"),
			Comment:         "good run",
			User:            database.User{ID: 5, FirstName: "Bob", LastName: "Sailor", UserStatusID: database.UserStatusMember},
		},
	}

	t.Run("returns the heats with their pricing", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetHeatsInSessionFn: func(uint) ([]database.Heat, error) { return heats, nil },
			GetUserStatusToPricingsMapFn: func() (map[uint]database.Pricing, error) {
				return pricingMap(t, "2.80"), nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetSessionHeats(rec, r, GetSessionHeatsRequest{SessionID: 3}, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data []GetSessionHeatsResponse
		decodeData(t, resp, &data)
		if len(data) != 1 {
			t.Fatalf("heats = %+v, want 1", data)
		}
		if data[0].HeatID != 1 || data[0].FirstName != "Bob" || data[0].Duration != 600 {
			t.Errorf("heat = %+v, does not match the database entry", data[0])
		}
		if !data[0].Pricing.Equal(mustDecimal(t, "2.80")) {
			t.Errorf("price_per_min = %s, want 2.80", data[0].Pricing)
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{
				name: "heats cannot be read",
				db: &dbtest.FakeDB{
					GetHeatsInSessionFn: func(uint) ([]database.Heat, error) { return nil, errNotFound },
				},
			},
			{
				name: "pricing cannot be read",
				db: &dbtest.FakeDB{
					GetHeatsInSessionFn: func(uint) ([]database.Heat, error) { return heats, nil },
					GetUserStatusToPricingsMapFn: func() (map[uint]database.Pricing, error) {
						return nil, errNotFound
					},
				},
			},
			{
				name: "no pricing for the user status of a rider",
				db: &dbtest.FakeDB{
					GetHeatsInSessionFn: func(uint) ([]database.Heat, error) { return heats, nil },
					GetUserStatusToPricingsMapFn: func() (map[uint]database.Pricing, error) {
						return map[uint]database.Pricing{}, nil
					},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.GetSessionHeats(rec, r, GetSessionHeatsRequest{SessionID: 3}, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestSessionRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		target  any
		errors  map[string]string
		wantErr bool
	}{
		{
			name:    "create session",
			payload: `{"title":"t","start":1714824000,"end":1714831200,"max_riders":4,"type":1}`,
			target:  &CreateSessionRequest{},
			errors:  SessionValidationErrors,
		},
		{
			name:    "create session with an unknown type",
			payload: `{"start":1714824000,"end":1714831200,"max_riders":4,"type":99}`,
			target:  &CreateSessionRequest{},
			errors:  SessionValidationErrors,
			wantErr: true,
		},
		{
			name:    "create session without an end time",
			payload: `{"start":1714824000,"max_riders":4,"type":1}`,
			target:  &CreateSessionRequest{},
			errors:  SessionValidationErrors,
			wantErr: true,
		},
		{
			name:    "delete session requires an id",
			payload: `{}`,
			target:  &DeleteSessionRequest{},
			errors:  DeleteSessionValidationErrors,
			wantErr: true,
		},
		{
			name:    "add users requires a user list",
			payload: `{"session_id":3}`,
			target:  &AddSessionUserRequest{},
			wantErr: true,
		},
		{
			name:    "remove user requires both ids",
			payload: `{"session_id":3}`,
			target:  &RemoveSessionUserRequest{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.errors == nil {
				err = ReadBodyAndValidate(newRequest(tt.payload, nil), tt.target)
			} else {
				err = ReadBodyAndValidate(newRequest(tt.payload, nil), tt.target, tt.errors)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
