package handlers

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"server/database"
	"server/database/dbtest"
)

func TestGetEngineHourLatest(t *testing.T) {
	t.Run("returns the latest entry", func(t *testing.T) {
		timestamp := time.Date(2024, time.May, 4, 15, 0, 0, 0, time.UTC)
		db := &dbtest.FakeDB{
			GetEngineHourLatestFn: func() (database.BoatEngineHour, error) {
				return database.BoatEngineHour{
					ID:          4,
					Timestamp:   timestamp,
					BeforeHours: mustDecimal(t, "400.0"),
					AfterHours:  mustDecimal(t, "402.5"),
					DeltaHours:  mustDecimal(t, "2.5"),
					TypeID:      database.DefaultSessionType,
					UserID:      3,
					User:        database.User{ID: 3, FirstName: "Ann", LastName: "Rider"},
				}, nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		h.GetEngineHourLatest(rec, newRequest("", &HandlerCtx{Database: db}))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data GetEngineHourLatestResponse
		decodeData(t, resp, &data)
		if data.ID != 4 || data.UserID != 3 || data.UserFirstName != "Ann" {
			t.Errorf("entry = %+v, does not match the database entry", data)
		}
		if data.Timestamp != timestamp.Unix() {
			t.Errorf("timestamp = %d, want %d", data.Timestamp, timestamp.Unix())
		}
		if data.UsageTypeName != "default" {
			t.Errorf("type_name = %q, want %q", data.UsageTypeName, "default")
		}
		if !data.DeltaHours.Equal(mustDecimal(t, "2.5")) {
			t.Errorf("delta_hours = %s, want 2.5", data.DeltaHours)
		}
	})

	t.Run("tolerates an empty table", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetEngineHourLatestFn: func() (database.BoatEngineHour, error) {
				return database.BoatEngineHour{}, gorm.ErrRecordNotFound
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		h.GetEngineHourLatest(rec, newRequest("", &HandlerCtx{Database: db}))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Errorf("an empty table has to yield an empty success response, got %+v", resp)
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetEngineHourLatestFn: func() (database.BoatEngineHour, error) {
				return database.BoatEngineHour{}, errNotFound
			},
		}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		h.GetEngineHourLatest(rec, newRequest("", &HandlerCtx{Database: db}))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestGetEngineHoursList(t *testing.T) {
	t.Run("returns all entries", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetEngineHoursFn: func() ([]database.BoatEngineHour, error) {
				return []database.BoatEngineHour{
					{ID: 1, TypeID: database.DefaultSessionType},
					{ID: 2, TypeID: database.CourseSessionType},
				}, nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		h.GetEngineHoursList(rec, newRequest("", &HandlerCtx{Database: db}))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data GetEngineHoursResponse
		decodeData(t, resp, &data)
		if len(data) != 2 {
			t.Fatalf("entries = %+v, want 2", data)
		}
		if data[0].UsageTypeName != "default" || data[1].UsageTypeName != "course" {
			t.Errorf("type names = %q/%q, want default/course", data[0].UsageTypeName, data[1].UsageTypeName)
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetEngineHoursFn: func() ([]database.BoatEngineHour, error) { return nil, errNotFound },
		}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		h.GetEngineHoursList(rec, newRequest("", &HandlerCtx{Database: db}))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestUpdateEngineHours(t *testing.T) {
	// The endpoint models a check-in/check-out flow: an entry is created with
	// the hours before the ride and completed with the hours after the ride.
	tests := []struct {
		name       string
		latest     database.BoatEngineHour
		payload    string
		wantAdd    bool
		wantUpdate bool
		wantDelta  string
		wantOK     bool
	}{
		{
			name:      "first entry ever",
			latest:    database.BoatEngineHour{},
			payload:   `{"user_id":1,"engine_hours_before":400,"type":1}`,
			wantAdd:   true,
			wantDelta: "0",
			wantOK:    true,
		},
		{
			name: "check in after a completed entry",
			latest: database.BoatEngineHour{
				ID:          1,
				BeforeHours: mustDecimal(t, "400"),
				AfterHours:  mustDecimal(t, "402"),
				CheckedIn:   false,
			},
			payload:   `{"user_id":1,"engine_hours_before":402,"type":1}`,
			wantAdd:   true,
			wantDelta: "0",
			wantOK:    true,
		},
		{
			name: "check out completes the open entry",
			latest: database.BoatEngineHour{
				ID:          9,
				BeforeHours: mustDecimal(t, "402"),
				CheckedIn:   true,
			},
			payload:    `{"user_id":1,"engine_hours_before":402,"engine_hours_after":404.5,"type":1}`,
			wantUpdate: true,
			wantDelta:  "2.5",
			wantOK:     true,
		},
		{
			name: "after hours smaller than before hours",
			latest: database.BoatEngineHour{
				ID:          9,
				BeforeHours: mustDecimal(t, "402"),
				CheckedIn:   true,
			},
			payload: `{"user_id":1,"engine_hours_before":402,"engine_hours_after":400,"type":1}`,
			wantOK:  false,
		},
		{
			name: "check out without an open entry",
			latest: database.BoatEngineHour{
				ID:          1,
				BeforeHours: mustDecimal(t, "400"),
				AfterHours:  mustDecimal(t, "402"),
				CheckedIn:   false,
			},
			payload: `{"user_id":1,"engine_hours_before":402,"engine_hours_after":404,"type":1}`,
			wantOK:  false,
		},
		{
			name:    "invalid session type",
			latest:  database.BoatEngineHour{},
			payload: `{"user_id":1,"engine_hours_before":400,"type":99}`,
			wantOK:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			added, updated := false, false
			var stored database.BoatEngineHour
			db := &dbtest.FakeDB{
				GetEngineHourLatestFn: func() (database.BoatEngineHour, error) { return tt.latest, nil },
				AddEngineHoursFn: func(b database.BoatEngineHour) error {
					added, stored = true, b
					return nil
				},
				UpdateEngineHoursFn: func(b database.BoatEngineHour) error {
					updated, stored = true, b
					return nil
				},
			}
			h := newTestHandler(t, db, nil)

			rec := httptest.NewRecorder()
			h.UpdateEngineHours(rec, newRequest(tt.payload, &HandlerCtx{Database: db}))

			resp := decodeResponse(t, rec)
			if resp.OK != tt.wantOK {
				t.Fatalf("response = %+v, want ok = %v", resp, tt.wantOK)
			}
			if added != tt.wantAdd {
				t.Errorf("entry added = %v, want %v", added, tt.wantAdd)
			}
			if updated != tt.wantUpdate {
				t.Errorf("entry updated = %v, want %v", updated, tt.wantUpdate)
			}
			if !tt.wantOK {
				return
			}
			if !stored.DeltaHours.Equal(mustDecimal(t, tt.wantDelta)) {
				t.Errorf("delta hours = %s, want %s", stored.DeltaHours, tt.wantDelta)
			}
			if tt.wantAdd && !stored.CheckedIn {
				t.Error("a new entry has to be marked as checked in")
			}
			if tt.wantUpdate {
				if stored.CheckedIn {
					t.Error("a completed entry must not stay checked in")
				}
				if stored.ID != tt.latest.ID {
					t.Errorf("updated entry id = %d, want %d", stored.ID, tt.latest.ID)
				}
			}
		})
	}
}

func TestUpdateEngineHoursFailures(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		db      *dbtest.FakeDB
	}{
		{
			name:    "invalid payload",
			payload: `{"engine_hours_before":400,"type":1}`,
			db:      &dbtest.FakeDB{},
		},
		{
			name:    "latest entry cannot be read",
			payload: `{"user_id":1,"engine_hours_before":400,"type":1}`,
			db: &dbtest.FakeDB{
				GetEngineHourLatestFn: func() (database.BoatEngineHour, error) {
					return database.BoatEngineHour{}, errNotFound
				},
			},
		},
		{
			name:    "entry cannot be stored",
			payload: `{"user_id":1,"engine_hours_before":400,"type":1}`,
			db: &dbtest.FakeDB{
				AddEngineHoursFn: func(database.BoatEngineHour) error { return errNotFound },
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHandler(t, tt.db, nil)
			rec := httptest.NewRecorder()
			h.UpdateEngineHours(rec, newRequest(tt.payload, &HandlerCtx{Database: tt.db}))
			if resp := decodeResponse(t, rec); resp.OK {
				t.Error("expected a failure response")
			}
		})
	}
}

func TestUpdateEngineHoursEntry(t *testing.T) {
	t.Run("changes the session type of an entry", func(t *testing.T) {
		var stored database.BoatEngineHour
		db := &dbtest.FakeDB{
			GetEngineHoursEntryFn: func(id uint) (database.BoatEngineHour, error) {
				return database.BoatEngineHour{ID: id, TypeID: database.DefaultSessionType, UserID: 3}, nil
			},
			UpdateEngineHoursFn: func(b database.BoatEngineHour) error {
				stored = b
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		h.UpdateEngineHoursEntry(rec, newRequest(`{"id":7,"type":2}`, &HandlerCtx{Database: db}))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if stored.ID != 7 {
			t.Errorf("updated entry = %d, want 7", stored.ID)
		}
		if stored.TypeID != database.CourseSessionType {
			t.Errorf("type = %d, want %d", stored.TypeID, database.CourseSessionType)
		}
		if stored.UserID != 3 {
			t.Error("the existing entry data has to be preserved")
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name    string
			payload string
			db      *dbtest.FakeDB
		}{
			{name: "missing id", payload: `{"type":1}`, db: &dbtest.FakeDB{}},
			{name: "invalid type", payload: `{"id":7,"type":99}`, db: &dbtest.FakeDB{}},
			{
				name:    "entry cannot be read",
				payload: `{"id":7,"type":1}`,
				db: &dbtest.FakeDB{
					GetEngineHoursEntryFn: func(uint) (database.BoatEngineHour, error) {
						return database.BoatEngineHour{}, errNotFound
					},
				},
			},
			{
				name:    "entry cannot be stored",
				payload: `{"id":7,"type":1}`,
				db: &dbtest.FakeDB{
					UpdateEngineHoursFn: func(database.BoatEngineHour) error { return errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				h.UpdateEngineHoursEntry(rec, newRequest(tt.payload, &HandlerCtx{Database: tt.db}))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestGetFuelEntries(t *testing.T) {
	fuelEntry := func(id uint, engineHours, liters string) database.BoatFuel {
		eh := mustDecimal(t, engineHours)
		l := mustDecimal(t, liters)
		cost := mustDecimal(t, "100.00")
		return database.BoatFuel{
			ID:          id,
			Timestamp:   time.Date(2024, time.May, int(id), 12, 0, 0, 0, time.UTC),
			UserID:      1,
			EngineHours: &eh,
			Liters:      &l,
			Cost:        &cost,
		}
	}

	t.Run("computes the average consumption between refuellings", func(t *testing.T) {
		// The handler walks the list backwards, so it is given newest first.
		db := &dbtest.FakeDB{
			GetFuelEntriesFn: func() ([]database.BoatFuel, error) {
				return []database.BoatFuel{
					fuelEntry(3, "420", "60"),
					fuelEntry(2, "410", "50"),
					fuelEntry(1, "400", "40"),
				}, nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		h.GetFuelEntries(rec, newRequest("", &HandlerCtx{Database: db}))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data []GetFuelEntriesResponse
		decodeData(t, resp, &data)
		if len(data) != 3 {
			t.Fatalf("entries = %+v, want 3", data)
		}
		// The response keeps the input order (newest first).
		if data[0].ID != 3 || data[2].ID != 1 {
			t.Errorf("entry order = %d,%d,%d, want 3,2,1", data[0].ID, data[1].ID, data[2].ID)
		}
		// The oldest entry has no predecessor, so there is nothing to compare to.
		if !data[2].DiffHours.IsZero() || !data[2].AverageLitersPerHour.IsZero() {
			t.Errorf("oldest entry = %+v, want no consumption data", data[2])
		}
		// Entry 2: 10 engine hours since entry 1, 50 liters -> 5 l/h.
		if !data[1].DiffHours.Equal(mustDecimal(t, "10")) {
			t.Errorf("diff_hours = %s, want 10", data[1].DiffHours)
		}
		if !data[1].AverageLitersPerHour.Equal(mustDecimal(t, "5")) {
			t.Errorf("avg_liters_per_hour = %s, want 5", data[1].AverageLitersPerHour)
		}
		// Entry 3: 10 engine hours since entry 2, 60 liters -> 6 l/h.
		if !data[0].AverageLitersPerHour.Equal(mustDecimal(t, "6")) {
			t.Errorf("avg_liters_per_hour = %s, want 6", data[0].AverageLitersPerHour)
		}
	})

	t.Run("tolerates entries without liters or engine hours", func(t *testing.T) {
		// These columns are nullable in the schema, so rows without values
		// must not crash the endpoint.
		db := &dbtest.FakeDB{
			GetFuelEntriesFn: func() ([]database.BoatFuel, error) {
				return []database.BoatFuel{{ID: 1, UserID: 1}}, nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		h.GetFuelEntries(rec, newRequest("", &HandlerCtx{Database: db}))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Errorf("unexpected failure: %s", resp.Msg)
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetFuelEntriesFn: func() ([]database.BoatFuel, error) { return nil, errNotFound },
		}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		h.GetFuelEntries(rec, newRequest("", &HandlerCtx{Database: db}))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestAddFuelEntry(t *testing.T) {
	req := AddFuelEntryRequest{
		Cost:        decimalPtr(t, "120.50"),
		EngineHours: decimalPtr(t, "412.5"),
		Liters:      decimalPtr(t, "55.3"),
		UserID:      3,
	}

	t.Run("instantly paid fuel contributes to the balance", func(t *testing.T) {
		var stored database.BoatFuel
		db := &dbtest.FakeDB{
			AddFuelEntryFn: func(e database.BoatFuel) error {
				stored = e
				return nil
			},
		}
		h := newTestHandler(t, db, map[string]string{"fuel.payment.type": "instant"})

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.AddFuelEntry(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if !stored.ContributeToBalance {
			t.Error("instantly paid fuel has to contribute to the balance")
		}
		if stored.UserID != 3 || !stored.Cost.Equal(mustDecimal(t, "120.50")) {
			t.Errorf("entry = %+v, does not match the request", stored)
		}
	})

	t.Run("billed fuel does not contribute to the balance", func(t *testing.T) {
		var stored database.BoatFuel
		db := &dbtest.FakeDB{
			AddFuelEntryFn: func(e database.BoatFuel) error {
				stored = e
				return nil
			},
		}
		h := newTestHandler(t, db, map[string]string{"fuel.payment.type": "billed"})

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.AddFuelEntry(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if stored.ContributeToBalance {
			t.Error("billed fuel must not contribute to the balance")
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name   string
			values map[string]string
			db     *dbtest.FakeDB
		}{
			{
				name:   "payment type not configured",
				values: map[string]string{"fuel.payment.type": ""},
				db:     &dbtest.FakeDB{},
			},
			{
				name:   "entry cannot be stored",
				values: map[string]string{"fuel.payment.type": "instant"},
				db: &dbtest.FakeDB{
					AddFuelEntryFn: func(database.BoatFuel) error { return errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, tt.values)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.AddFuelEntry(rec, r, req, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestChangeFuelEntry(t *testing.T) {
	existing := func() database.BoatFuel {
		eh := mustDecimal(t, "400")
		l := mustDecimal(t, "50")
		c := mustDecimal(t, "100")
		return database.BoatFuel{
			ID:                  4,
			Timestamp:           time.Date(2024, time.May, 4, 12, 0, 0, 0, time.UTC),
			UserID:              9,
			EngineHours:         &eh,
			Liters:              &l,
			Cost:                &c,
			ContributeToBalance: true,
		}
	}

	t.Run("keeps timestamp, user and balance flag of the existing entry", func(t *testing.T) {
		var stored database.BoatFuel
		db := &dbtest.FakeDB{
			GetFuelEntryFn: func(uint) (database.BoatFuel, error) { return existing(), nil },
			ChangeFuelEntryFn: func(e database.BoatFuel) error {
				stored = e
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		req := ChangeFuelEntryRequest{
			ID:           4,
			CostNet:      decimalPtr(t, "110.00"),
			CostGros:     decimalPtr(t, "118.00"),
			EngineHours:  decimalPtr(t, "405"),
			Liters:       decimalPtr(t, "55"),
			IsDiscounted: true,
		}
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.ChangeFuelEntry(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if stored.ID != 4 || stored.UserID != 9 {
			t.Errorf("entry = %+v, want the existing id and user preserved", stored)
		}
		if !stored.Timestamp.Equal(existing().Timestamp) {
			t.Errorf("timestamp = %v, want the existing timestamp", stored.Timestamp)
		}
		if !stored.ContributeToBalance {
			t.Error("the balance flag of the existing entry has to be preserved")
		}
		if !stored.Cost.Equal(mustDecimal(t, "110.00")) || !stored.IsDiscounted {
			t.Errorf("entry = %+v, does not match the request", stored)
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{
				name: "entry does not exist",
				db: &dbtest.FakeDB{
					GetFuelEntryFn: func(uint) (database.BoatFuel, error) { return database.BoatFuel{}, errNotFound },
				},
			},
			{
				name: "entry cannot be stored",
				db: &dbtest.FakeDB{
					GetFuelEntryFn:    func(uint) (database.BoatFuel, error) { return existing(), nil },
					ChangeFuelEntryFn: func(database.BoatFuel) error { return errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.ChangeFuelEntry(rec, r, ChangeFuelEntryRequest{
					ID:          4,
					CostNet:     decimalPtr(t, "1"),
					EngineHours: decimalPtr(t, "1"),
					Liters:      decimalPtr(t, "1"),
				}, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestRemoveFuelEntry(t *testing.T) {
	t.Run("removes an existing entry", func(t *testing.T) {
		var removedID uint
		db := &dbtest.FakeDB{
			GetFuelEntryFn: func(id uint) (database.BoatFuel, error) { return database.BoatFuel{ID: id}, nil },
			RemoveFuelEntryFn: func(id uint) error {
				removedID = id
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.RemoveFuelEntry(rec, r, RemoveFuelEntryRequest{ID: 8}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if removedID != 8 {
			t.Errorf("removed entry = %d, want 8", removedID)
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{
				name: "entry does not exist",
				db: &dbtest.FakeDB{
					GetFuelEntryFn: func(uint) (database.BoatFuel, error) { return database.BoatFuel{}, errNotFound },
				},
			},
			{
				name: "entry cannot be removed",
				db: &dbtest.FakeDB{
					GetFuelEntryFn:    func(id uint) (database.BoatFuel, error) { return database.BoatFuel{ID: id}, nil },
					RemoveFuelEntryFn: func(uint) error { return errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.RemoveFuelEntry(rec, r, RemoveFuelEntryRequest{ID: 8}, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestGetMaintenanceEntries(t *testing.T) {
	t.Run("returns the maintenance log", func(t *testing.T) {
		timestamp := time.Date(2024, time.May, 4, 12, 0, 0, 0, time.UTC)
		db := &dbtest.FakeDB{
			GetMaintenanceFn: func() ([]database.BoatMaintenance, error) {
				return []database.BoatMaintenance{
					{
						ID:          2,
						Timestamp:   timestamp,
						UserID:      3,
						User:        database.User{FirstName: "Ann", LastName: "Rider"},
						EngineHours: mustDecimal(t, "410.5"),
						Description: "oil change",
					},
				}, nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		h.GetMaintenanceEntries(rec, newRequest("", &HandlerCtx{Database: db}))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data []GetMaintenanceEntriesResponse
		decodeData(t, resp, &data)
		if len(data) != 1 {
			t.Fatalf("entries = %+v, want 1", data)
		}
		if data[0].Description != "oil change" || data[0].UserFirstName != "Ann" {
			t.Errorf("entry = %+v, does not match the database entry", data[0])
		}
		if data[0].Timestamp != timestamp.Unix() {
			t.Errorf("timestamp = %d, want %d", data[0].Timestamp, timestamp.Unix())
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetMaintenanceFn: func() ([]database.BoatMaintenance, error) { return nil, errNotFound },
		}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		h.GetMaintenanceEntries(rec, newRequest("", &HandlerCtx{Database: db}))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestAddMaintenanceEntry(t *testing.T) {
	req := AddMaintenanceEntryRequest{
		Description: "oil change",
		EngineHours: mustDecimal(t, "410.5"),
		UserID:      3,
	}

	t.Run("stores the entry", func(t *testing.T) {
		var stored database.BoatMaintenance
		db := &dbtest.FakeDB{
			UserExistsFn: func(uint) bool { return true },
			AddMaintenanceEntryFn: func(m database.BoatMaintenance) error {
				stored = m
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.AddMaintenanceEntry(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if stored.UserID != 3 || stored.Description != "oil change" {
			t.Errorf("entry = %+v, does not match the request", stored)
		}
		if !stored.EngineHours.Equal(mustDecimal(t, "410.5")) {
			t.Errorf("engine hours = %s, want 410.5", stored.EngineHours)
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{
				name: "unknown user",
				db: &dbtest.FakeDB{
					UserExistsFn: func(uint) bool { return false },
					AddMaintenanceEntryFn: func(database.BoatMaintenance) error {
						t.Error("no entry may be stored for an unknown user")
						return nil
					},
				},
			},
			{
				name: "entry cannot be stored",
				db: &dbtest.FakeDB{
					UserExistsFn:          func(uint) bool { return true },
					AddMaintenanceEntryFn: func(database.BoatMaintenance) error { return errNotFound },
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.AddMaintenanceEntry(rec, r, req, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestBoatRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		target  any
		errors  map[string]string
		wantErr bool
	}{
		{
			name:    "fuel entry requires cost, hours, liters and user",
			payload: `{"cost":"12.5","engine_hours":"400","liters":"50","user_id":1}`,
			target:  &AddFuelEntryRequest{},
			errors:  FuelEntryValidationErrors,
		},
		{
			name:    "fuel entry without liters",
			payload: `{"cost":"12.5","engine_hours":"400","user_id":1}`,
			target:  &AddFuelEntryRequest{},
			errors:  FuelEntryValidationErrors,
			wantErr: true,
		},
		{
			name:    "maintenance entry requires a description",
			payload: `{"engine_hours":"400","user_id":1}`,
			target:  &AddMaintenanceEntryRequest{},
			errors:  AddMaintenanceEntryValidationErrors,
			wantErr: true,
		},
		{
			name:    "remove fuel entry requires an id",
			payload: `{}`,
			target:  &RemoveFuelEntryRequest{},
			errors:  FuelEntryValidationErrors,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReadBodyAndValidate(newRequest(tt.payload, nil), tt.target, tt.errors)
			if (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateEngineHoursRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{
			name:    "valid payload with null engine_hours_after",
			payload: `{"user_id":1,"engine_hours_before":400,"engine_hours_after":null,"type":1}`,
			wantErr: false,
		},
		{
			name:    "valid payload with null engine_hours_after as string",
			payload: `{"user_id":1,"engine_hours_before":"400","engine_hours_after":null,"type":1}`,
			wantErr: false,
		},
		{
			name:    "valid payload with non-null engine_hours_after",
			payload: `{"user_id":1,"engine_hours_before":400,"engine_hours_after":405.5,"type":1}`,
			wantErr: false,
		},
		{
			name:    "missing required user_id",
			payload: `{"engine_hours_before":400,"type":1}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &UpdateEngineHoursRequest{}
			err := ReadBodyAndValidate(newRequest(tt.payload, nil), req, UpdateEngineHoursValidationErrors)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBodyAndValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// decimal values that arrive as JSON numbers must keep their precision.
func TestEngineHourDecimalPrecision(t *testing.T) {
	req := &UpdateEngineHoursRequest{}
	if err := ReadBodyAndValidate(newRequest(`{"user_id":1,"engine_hours_before":400.125,"type":1}`, nil), req, UpdateEngineHoursValidationErrors); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if !req.BeforeHours.Equal(decimal.RequireFromString("400.125")) {
		t.Errorf("engine_hours_before = %s, want 400.125", req.BeforeHours)
	}
}
