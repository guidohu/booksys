package handlers

import (
	"net/http/httptest"
	"testing"
	"time"

	"server/database"
	"server/database/dbtest"
)

// bookingConfig is a complete configuration for the booking endpoints.
func bookingConfig() map[string]string {
	return map[string]string{
		"location.timezone":  "Europe/Zurich",
		"location.latitude":  "47.3769",
		"location.longitude": "8.5417",
		"business.day.start": "08:00",
		"business.day.end":   "20:00",
	}
}

func TestGetBookingDay(t *testing.T) {
	start := time.Date(2024, time.June, 21, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	t.Run("returns the sessions of the day with their riders", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetSessionsBetweenFn: func(gotStart, gotEnd time.Time) ([]database.Session, error) {
				if !gotStart.Equal(start) || !gotEnd.Equal(end) {
					t.Errorf("GetSessionsBetween(%v, %v), want (%v, %v)", gotStart, gotEnd, start, end)
				}
				return []database.Session{
					{
						ID:         3,
						StartTime:  start.Add(10 * time.Hour),
						EndTime:    start.Add(12 * time.Hour),
						Title:      "Morning",
						FreeSpaces: 4,
						Creator:    database.User{FirstName: "Ann", LastName: "Rider"},
					},
				}, nil
			},
			GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) {
				return []database.UserToSession{
					{UserID: 5, User: database.User{ID: 5, FirstName: "Bob", LastName: "Sailor"}},
				}, nil
			},
		}
		h := newTestHandler(t, db, bookingConfig())

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetBookingDay(rec, r, GetBookingDayRequest{Start: start.Unix(), End: end.Unix()}, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data GetBookingResponse
		decodeData(t, resp, &data)
		if data.Timezone != "Europe/Zurich" {
			t.Errorf("timezone = %q, want Europe/Zurich", data.Timezone)
		}
		if data.OpeningHourStart != "08:00" || data.OpeningHourEnd != "20:00" {
			t.Errorf("business hours = %q-%q, want 08:00-20:00", data.OpeningHourStart, data.OpeningHourEnd)
		}
		if data.Sunrise >= data.Sunset {
			t.Errorf("sunrise (%d) has to be before sunset (%d)", data.Sunrise, data.Sunset)
		}
		if len(data.Sessions) != 1 {
			t.Fatalf("sessions = %+v, want 1", data.Sessions)
		}
		if data.Sessions[0].Duration != 7200 {
			t.Errorf("duration = %d, want 7200", data.Sessions[0].Duration)
		}
		if len(data.Sessions[0].Riders) != 1 || data.Sessions[0].Riders[0].Name != "Bob Sailor" {
			t.Errorf("riders = %+v, want Bob Sailor", data.Sessions[0].Riders)
		}
	})

	t.Run("returns an empty session list for a day without sessions", func(t *testing.T) {
		db := &dbtest.FakeDB{}
		h := newTestHandler(t, db, bookingConfig())

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetBookingDay(rec, r, GetBookingDayRequest{Start: start.Unix(), End: end.Unix()}, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data GetBookingResponse
		decodeData(t, resp, &data)
		if data.Sessions == nil {
			t.Error("the session list has to be an empty array, not null")
		}
		if len(data.Sessions) != 0 {
			t.Errorf("sessions = %+v, want none", data.Sessions)
		}
	})

	t.Run("failures", func(t *testing.T) {
		noStart := bookingConfig()
		noStart["business.day.start"] = ""
		noEnd := bookingConfig()
		noEnd["business.day.end"] = ""
		badZone := bookingConfig()
		badZone["location.timezone"] = "Not/AZone"

		tests := []struct {
			name   string
			values map[string]string
			db     *dbtest.FakeDB
		}{
			{name: "business day start not configured", values: noStart, db: &dbtest.FakeDB{}},
			{name: "business day end not configured", values: noEnd, db: &dbtest.FakeDB{}},
			{name: "invalid timezone", values: badZone, db: &dbtest.FakeDB{}},
			{
				name:   "sessions cannot be read",
				values: bookingConfig(),
				db: &dbtest.FakeDB{
					GetSessionsBetweenFn: func(time.Time, time.Time) ([]database.Session, error) {
						return nil, errNotFound
					},
				},
			},
			{
				name:   "riders cannot be read",
				values: bookingConfig(),
				db: &dbtest.FakeDB{
					GetSessionsBetweenFn: func(time.Time, time.Time) ([]database.Session, error) {
						return []database.Session{{ID: 3}}, nil
					},
					GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) {
						return nil, errNotFound
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, tt.values)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.GetBookingDay(rec, r, GetBookingDayRequest{Start: start.Unix(), End: end.Unix()}, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestGetBookingSeries(t *testing.T) {
	day := time.Date(2024, time.June, 21, 0, 0, 0, 0, time.UTC)

	t.Run("returns one entry per time window", func(t *testing.T) {
		calls := 0
		db := &dbtest.FakeDB{
			GetSessionsBetweenFn: func(time.Time, time.Time) ([]database.Session, error) {
				calls++
				return nil, nil
			},
		}
		h := newTestHandler(t, db, bookingConfig())

		req := GetBookingSeriesRequest{
			TimeWindows: []GetBookingDayRequest{
				{Start: day.Unix(), End: day.Add(24 * time.Hour).Unix()},
				{Start: day.Add(24 * time.Hour).Unix(), End: day.Add(48 * time.Hour).Unix()},
			},
		}
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetBookingSeries(rec, r, req, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data GetBookingSeriesResponse
		decodeData(t, resp, &data)
		if len(data) != 2 {
			t.Fatalf("windows = %d, want 2", len(data))
		}
		if calls != 2 {
			t.Errorf("database queries = %d, want 2", calls)
		}
	})

	t.Run("a failing window yields an empty entry instead of an error", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetSessionsBetweenFn: func(time.Time, time.Time) ([]database.Session, error) {
				return nil, errNotFound
			},
		}
		h := newTestHandler(t, db, bookingConfig())

		req := GetBookingSeriesRequest{
			TimeWindows: []GetBookingDayRequest{{Start: day.Unix(), End: day.Add(24 * time.Hour).Unix()}},
		}
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetBookingSeries(rec, r, req, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("the series endpoint has to stay successful, got %+v", resp)
		}
		var data GetBookingSeriesResponse
		decodeData(t, resp, &data)
		if len(data) != 1 {
			t.Fatalf("windows = %d, want 1", len(data))
		}
		if data[0].Start != 0 || len(data[0].Sessions) != 0 {
			t.Errorf("failed window = %+v, want an empty entry", data[0])
		}
	})

	t.Run("an empty series is accepted", func(t *testing.T) {
		db := &dbtest.FakeDB{}
		h := newTestHandler(t, db, bookingConfig())
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetBookingSeries(rec, r, GetBookingSeriesRequest{}, GetHandlerContext(r))
		if resp := decodeResponse(t, rec); !resp.OK {
			t.Errorf("unexpected failure: %s", resp.Msg)
		}
	})
}

func TestGetSunriseSunset(t *testing.T) {
	h := newTestHandler(t, &dbtest.FakeDB{}, bookingConfig())
	zurich, err := time.LoadLocation("Europe/Zurich")
	if err != nil {
		t.Fatalf("cannot load the test timezone: %v", err)
	}

	sunrise, sunset := h.getSunriseSunset(time.Date(2024, time.June, 21, 12, 0, 0, 0, time.UTC), zurich)
	if sunrise.IsZero() || sunset.IsZero() {
		t.Fatal("sunrise/sunset were not calculated")
	}
	if !sunrise.Before(sunset) {
		t.Errorf("sunrise (%v) has to be before sunset (%v)", sunrise, sunset)
	}
	// Around the summer solstice the sun rises before 06:00 local time.
	if got := sunrise.In(zurich).Hour(); got > 6 {
		t.Errorf("sunrise hour = %d, want an early summer sunrise", got)
	}
}

func TestGetSunriseSunsetWithoutCoordinates(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]string
	}{
		{name: "no latitude", values: map[string]string{"location.longitude": "8.5417"}},
		{name: "no longitude", values: map[string]string{"location.latitude": "47.3769"}},
		{name: "unparsable latitude", values: map[string]string{"location.latitude": "north", "location.longitude": "8.5417"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHandler(t, &dbtest.FakeDB{}, tt.values)
			sunrise, sunset := h.getSunriseSunset(time.Now(), time.UTC)
			if !sunrise.IsZero() || !sunset.IsZero() {
				t.Errorf("sunrise/sunset = %v/%v, want zero times for an incomplete location", sunrise, sunset)
			}
		})
	}
}

func TestGetLocation(t *testing.T) {
	t.Run("loads the configured timezone", func(t *testing.T) {
		h := newTestHandler(t, &dbtest.FakeDB{}, map[string]string{"location.timezone": "Europe/Zurich"})
		loc, err := h.getLocation()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if loc.String() != "Europe/Zurich" {
			t.Errorf("location = %q, want Europe/Zurich", loc.String())
		}
	})

	t.Run("rejects an unknown timezone", func(t *testing.T) {
		h := newTestHandler(t, &dbtest.FakeDB{}, map[string]string{"location.timezone": "Not/AZone"})
		if _, err := h.getLocation(); err == nil {
			t.Error("expected an error for an unknown timezone")
		}
	})
}

func TestGetRiders(t *testing.T) {
	t.Run("returns the users of a session", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetUsersForSessionFn: func(uint) ([]database.UserToSession, error) {
				return []database.UserToSession{
					{UserID: 5, User: database.User{ID: 5, FirstName: "Bob"}},
				}, nil
			},
		}
		h := newTestHandler(t, db, nil)

		riders, err := h.getRiders(3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(riders) != 1 || riders[0].FirstName != "Bob" {
			t.Errorf("riders = %+v, want Bob", riders)
		}
	})

	t.Run("fails without a database connection", func(t *testing.T) {
		h := NewHandler(HandlerParams{
			Database:      database.NewManager(database.Settings{}),
			Configuration: newTestConfig(t, nil),
		})
		if _, err := h.getRiders(3); err == nil {
			t.Error("expected an error without a database connection")
		}
	})
}

func TestBookingRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		target  any
		wantErr bool
	}{
		{name: "booking day", payload: `{"start":1718920800,"end":1719007200}`, target: &GetBookingDayRequest{}},
		{name: "booking day without end", payload: `{"start":1718920800}`, target: &GetBookingDayRequest{}, wantErr: true},
		{
			name:    "booking series",
			payload: `{"time_windows":[{"start":1718920800,"end":1719007200}]}`,
			target:  &GetBookingSeriesRequest{},
		},
		{name: "booking series without windows", payload: `{}`, target: &GetBookingSeriesRequest{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReadBodyAndValidate(newRequest(tt.payload, nil), tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
