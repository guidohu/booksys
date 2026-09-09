package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"server/database"
	"server/database/dbtest"
)

// pricingMap returns a user status to pricing map with the given price for
// the member status.
func pricingMap(t *testing.T, price string) map[uint]database.Pricing {
	t.Helper()
	return map[uint]database.Pricing{
		database.UserStatusMember: {
			ID:             2,
			UserStatusID:   database.UserStatusMember,
			PricePerMinute: mustDecimal(t, price),
		},
	}
}

func TestCalculateHeatCost(t *testing.T) {
	h := newTestHandler(t, &dbtest.FakeDB{}, nil)

	tests := []struct {
		name            string
		durationSeconds uint64
		price           string
		want            string
	}{
		{name: "zero duration", durationSeconds: 0, price: "2.80", want: "0"},
		{name: "exactly one minute", durationSeconds: 60, price: "2.80", want: "2.80"},
		{name: "ten minutes", durationSeconds: 600, price: "2.80", want: "28"},
		// 90s -> 1.5min * 2.80 = 4.20, already a multiple of 0.05.
		{name: "one and a half minutes", durationSeconds: 90, price: "2.80", want: "4.20"},
		// 100s -> 1.6667min * 2.80 = 4.6667, rounded up to the next 0.05 -> 4.70.
		{name: "rounds up to the next five cents", durationSeconds: 100, price: "2.80", want: "4.70"},
		// 61s -> 1.0167min * 1.30 = 1.3217, rounded up -> 1.35.
		{name: "rounds up with a different price", durationSeconds: 61, price: "1.30", want: "1.35"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := h.calculateHeatCost(tt.durationSeconds, mustDecimal(t, tt.price))
			if !got.Equal(mustDecimal(t, tt.want)) {
				t.Errorf("calculateHeatCost(%d, %s) = %s, want %s", tt.durationSeconds, tt.price, got, tt.want)
			}
			// The result has to be payable in five cent steps.
			if !got.Mod(mustDecimal(t, "0.05")).IsZero() {
				t.Errorf("calculateHeatCost(%d, %s) = %s is not a multiple of 0.05", tt.durationSeconds, tt.price, got)
			}
		})
	}
}

func TestAddHeats(t *testing.T) {
	member := database.User{ID: 5, UserStatusID: database.UserStatusMember}

	t.Run("adds heats and reports the result per heat", func(t *testing.T) {
		var stored []database.Heat
		db := &dbtest.FakeDB{
			GetUserByIdFn: func(id uint) (database.User, error) {
				if id == 99 {
					return database.User{}, errNotFound
				}
				return member, nil
			},
			GetSessionFn: func(id uint) (database.Session, error) {
				return database.Session{ID: id}, nil
			},
			GetUserStatusToPricingsMapFn: func() (map[uint]database.Pricing, error) {
				return pricingMap(t, "2.80"), nil
			},
			AddHeatFn: func(h *database.Heat) error {
				stored = append(stored, *h)
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		req := AddHeatsRequest{
			Heats: []AddHeatRequest{
				{UID: "a", SessionID: 3, DurationSeconds: 600, UserID: 5, Comment: "good run"},
				{UID: "b", SessionID: 3, DurationSeconds: 300, UserID: 99},
			},
		}
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.AddHeats(rec, r, req, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data AddHeatsRequestResponse
		decodeData(t, resp, &data)
		if !data["a"].OK {
			t.Errorf("heat a = %+v, want a success", data["a"])
		}
		if data["b"].OK {
			t.Errorf("heat b = %+v, want a failure for the unknown user", data["b"])
		}
		if len(stored) != 1 {
			t.Fatalf("stored heats = %d, want 1", len(stored))
		}
		if !stored[0].Cost.Equal(mustDecimal(t, "28")) {
			t.Errorf("cost = %s, want 28 (10min * 2.80)", stored[0].Cost)
		}
		if stored[0].DurationSeconds != 600 || stored[0].SessionID != 3 || stored[0].UserID != 5 {
			t.Errorf("heat = %+v, does not match the request", stored[0])
		}
		if stored[0].Comment != "good run" {
			t.Errorf("comment = %q, want %q", stored[0].Comment, "good run")
		}
	})

	t.Run("rejects heats that cannot be priced or assigned", func(t *testing.T) {
		tests := []struct {
			name string
			db   *dbtest.FakeDB
			req  AddHeatRequest
		}{
			{
				name: "unknown user",
				db:   &dbtest.FakeDB{},
				req:  AddHeatRequest{UID: "x", DurationSeconds: 60, UserID: 5},
			},
			{
				name: "session does not exist",
				db: &dbtest.FakeDB{
					GetUserByIdFn: func(uint) (database.User, error) { return member, nil },
					GetSessionFn:  func(uint) (database.Session, error) { return database.Session{}, errNotFound },
				},
				req: AddHeatRequest{UID: "x", SessionID: 3, DurationSeconds: 60, UserID: 5},
			},
			{
				name: "session id refers to a missing row",
				db: &dbtest.FakeDB{
					GetUserByIdFn: func(uint) (database.User, error) { return member, nil },
					// No error, but also no session: the row does not exist.
					GetSessionFn: func(uint) (database.Session, error) { return database.Session{}, nil },
				},
				req: AddHeatRequest{UID: "x", SessionID: 3, DurationSeconds: 60, UserID: 5},
			},
			{
				name: "pricing cannot be read",
				db: &dbtest.FakeDB{
					GetUserByIdFn: func(uint) (database.User, error) { return member, nil },
					GetUserStatusToPricingsMapFn: func() (map[uint]database.Pricing, error) {
						return nil, errNotFound
					},
				},
				req: AddHeatRequest{UID: "x", DurationSeconds: 60, UserID: 5},
			},
			{
				name: "no pricing for the user status",
				db: &dbtest.FakeDB{
					GetUserByIdFn: func(uint) (database.User, error) {
						return database.User{ID: 5, UserStatusID: 42}, nil
					},
					GetUserStatusToPricingsMapFn: func() (map[uint]database.Pricing, error) {
						return pricingMap(t, "2.80"), nil
					},
				},
				req: AddHeatRequest{UID: "x", DurationSeconds: 60, UserID: 5},
			},
			{
				name: "price is zero",
				db: &dbtest.FakeDB{
					GetUserByIdFn: func(uint) (database.User, error) { return member, nil },
					GetUserStatusToPricingsMapFn: func() (map[uint]database.Pricing, error) {
						return pricingMap(t, "0"), nil
					},
				},
				req: AddHeatRequest{UID: "x", DurationSeconds: 60, UserID: 5},
			},
			{
				name: "heat cannot be stored",
				db: &dbtest.FakeDB{
					GetUserByIdFn: func(uint) (database.User, error) { return member, nil },
					GetUserStatusToPricingsMapFn: func() (map[uint]database.Pricing, error) {
						return pricingMap(t, "2.80"), nil
					},
					AddHeatFn: func(*database.Heat) error { return errNotFound },
				},
				req: AddHeatRequest{UID: "x", DurationSeconds: 60, UserID: 5},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.AddHeats(rec, r, AddHeatsRequest{Heats: []AddHeatRequest{tt.req}}, GetHandlerContext(r))

				resp := decodeResponse(t, rec)
				var data AddHeatsRequestResponse
				decodeData(t, resp, &data)
				if data["x"].OK {
					t.Errorf("heat x = %+v, want a failure", data["x"])
				}
			})
		}
	})

	t.Run("an empty request is accepted", func(t *testing.T) {
		db := &dbtest.FakeDB{}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.AddHeats(rec, r, AddHeatsRequest{}, GetHandlerContext(r))
		if resp := decodeResponse(t, rec); !resp.OK {
			t.Errorf("unexpected failure: %s", resp.Msg)
		}
	})
}

func TestAddHeat(t *testing.T) {
	member := database.User{ID: 5, UserStatusID: database.UserStatusMember}

	t.Run("adds a single heat", func(t *testing.T) {
		added := false
		db := &dbtest.FakeDB{
			GetUserByIdFn: func(uint) (database.User, error) { return member, nil },
			GetUserStatusToPricingsMapFn: func() (map[uint]database.Pricing, error) {
				return pricingMap(t, "2.80"), nil
			},
			AddHeatFn: func(*database.Heat) error {
				added = true
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		h.AddHeat(rec, newRequest(`{"uid":"a","duration_s":600,"user_id":5}`, &HandlerCtx{Database: db}))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if !added {
			t.Error("the heat was not added")
		}
	})

	t.Run("rejects an invalid payload", func(t *testing.T) {
		db := &dbtest.FakeDB{}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		h.AddHeat(rec, newRequest(`{`, &HandlerCtx{Database: db}))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})

	t.Run("reports a failing heat", func(t *testing.T) {
		db := &dbtest.FakeDB{}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		h.AddHeat(rec, newRequest(`{"uid":"a","duration_s":600,"user_id":5}`, &HandlerCtx{Database: db}))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response for an unknown user")
		}
	})
}

func TestDeleteHeat(t *testing.T) {
	t.Run("deletes the heat", func(t *testing.T) {
		var deletedID uint
		db := &dbtest.FakeDB{
			DeleteHeatFn: func(id uint) error {
				deletedID = id
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.DeleteHeat(rec, r, DeleteHeatRequest{HeatID: 12}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if deletedID != 12 {
			t.Errorf("deleted heat = %d, want 12", deletedID)
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{DeleteHeatFn: func(uint) error { return errNotFound }}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.DeleteHeat(rec, r, DeleteHeatRequest{HeatID: 12}, GetHandlerContext(r))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestChangeHeat(t *testing.T) {
	member := database.User{ID: 5, UserStatusID: database.UserStatusMember}
	existing := database.Heat{
		ID:              12,
		UserID:          5,
		SessionID:       3,
		DurationSeconds: 600,
		Cost:            mustDecimal(t, "28"),
	}

	t.Run("recalculates the cost and keeps the session", func(t *testing.T) {
		var stored database.Heat
		db := &dbtest.FakeDB{
			GetUserByIdFn: func(uint) (database.User, error) { return member, nil },
			GetUserStatusToPricingsMapFn: func() (map[uint]database.Pricing, error) {
				return pricingMap(t, "2.80"), nil
			},
			GetHeatFn: func(uint) (database.Heat, error) { return existing, nil },
			ChangeHeatFn: func(h *database.Heat) error {
				stored = *h
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		req := ChangeHeatRequest{HeatID: 12, UserID: 5, DurationSeconds: 300, Comment: "shorter"}
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.ChangeHeat(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if stored.ID != 12 || stored.SessionID != 3 {
			t.Errorf("heat = %+v, want the existing id and session preserved", stored)
		}
		if stored.DurationSeconds != 300 {
			t.Errorf("duration = %d, want 300", stored.DurationSeconds)
		}
		if !stored.Cost.Equal(mustDecimal(t, "14")) {
			t.Errorf("cost = %s, want 14 (5min * 2.80)", stored.Cost)
		}
		if stored.Comment != "shorter" {
			t.Errorf("comment = %q, want %q", stored.Comment, "shorter")
		}
	})

	t.Run("failures", func(t *testing.T) {
		withPricing := func(db *dbtest.FakeDB) *dbtest.FakeDB {
			db.GetUserByIdFn = func(uint) (database.User, error) { return member, nil }
			if db.GetUserStatusToPricingsMapFn == nil {
				db.GetUserStatusToPricingsMapFn = func() (map[uint]database.Pricing, error) {
					return pricingMap(t, "2.80"), nil
				}
			}
			return db
		}

		tests := []struct {
			name string
			db   *dbtest.FakeDB
		}{
			{name: "unknown user", db: &dbtest.FakeDB{}},
			{
				name: "pricing cannot be read",
				db: withPricing(&dbtest.FakeDB{
					GetUserStatusToPricingsMapFn: func() (map[uint]database.Pricing, error) { return nil, errNotFound },
				}),
			},
			{
				name: "no pricing for the user status",
				db: &dbtest.FakeDB{
					GetUserByIdFn: func(uint) (database.User, error) {
						return database.User{ID: 5, UserStatusID: 42}, nil
					},
					GetUserStatusToPricingsMapFn: func() (map[uint]database.Pricing, error) {
						return pricingMap(t, "2.80"), nil
					},
				},
			},
			{
				name: "price is zero",
				db: withPricing(&dbtest.FakeDB{
					GetUserStatusToPricingsMapFn: func() (map[uint]database.Pricing, error) {
						return pricingMap(t, "0"), nil
					},
				}),
			},
			{
				name: "heat does not exist",
				db: withPricing(&dbtest.FakeDB{
					GetHeatFn: func(uint) (database.Heat, error) { return database.Heat{}, errNotFound },
				}),
			},
			{
				name: "heat id refers to a missing row",
				db: withPricing(&dbtest.FakeDB{
					// No error, but also no heat: the row does not exist.
					GetHeatFn: func(uint) (database.Heat, error) { return database.Heat{}, nil },
				}),
			},
			{
				name: "heat cannot be stored",
				db: withPricing(&dbtest.FakeDB{
					GetHeatFn:    func(uint) (database.Heat, error) { return existing, nil },
					ChangeHeatFn: func(*database.Heat) error { return errNotFound },
				}),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.ChangeHeat(rec, r, ChangeHeatRequest{HeatID: 12, UserID: 5, DurationSeconds: 300}, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestHeatCostNeverExceedsAFullMinuteRoundingError(t *testing.T) {
	h := newTestHandler(t, &dbtest.FakeDB{}, nil)
	price := mustDecimal(t, "2.80")
	for seconds := uint64(1); seconds < 3600; seconds += 7 {
		got := h.calculateHeatCost(seconds, price)
		exact := price.Mul(decimal.NewFromInt(int64(seconds))).Div(decimal.NewFromInt(60))
		if got.LessThan(exact) {
			t.Fatalf("cost for %ds = %s is below the exact price %s", seconds, got, exact)
		}
		if got.Sub(exact).GreaterThan(mustDecimal(t, "0.05")) {
			t.Fatalf("cost for %ds = %s overcharges the exact price %s by more than five cents", seconds, got, exact)
		}
	}
}
