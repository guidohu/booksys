package handlers

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"server/database"
	"server/database/dbtest"
	"server/version"
)

func TestPing(t *testing.T) {
	db := &dbtest.FakeDB{}
	h := newTestHandler(t, db, nil)

	rec := httptest.NewRecorder()
	h.Ping(rec, newRequest("", nil))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data PingData
	decodeData(t, resp, &data)
	if data.Message == "" {
		t.Error("the ping response has to carry a message")
	}
}

func TestHealthStatus(t *testing.T) {
	dbSettings := map[string]string{
		"database.protocol": "tcp",
		"database.user":     "booksys",
		"database.host":     "127.0.0.1",
		"database.port":     "3306",
		"database.dbname":   "booksys",
		"environment":       "test",
	}

	t.Run("reports a healthy system", func(t *testing.T) {
		configFile := filepath.Join(t.TempDir(), "config.yaml")
		if err := os.WriteFile(configFile, []byte("http:\n  port: 8080\n"), 0o600); err != nil {
			t.Fatalf("cannot write the config file: %v", err)
		}
		db := &dbtest.FakeDB{
			UsersExistFn: func() (bool, error) { return true, nil },
		}
		h := NewHandler(HandlerParams{
			Database:      dbtest.NewManager(t, db),
			Configuration: newTestConfigFromFile(t, configFile, dbSettings),
		})

		rec := httptest.NewRecorder()
		h.HealthStatus(rec, newRequest("", nil))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data HealthStatusResponse
		decodeData(t, resp, &data)
		if !data.ConfigFile || !data.ConfigDB || !data.DBReachable || !data.UsersExist {
			t.Errorf("health = %+v, want everything to be healthy", data)
		}
		if data.Environment != "test" {
			t.Errorf("environment = %q, want test", data.Environment)
		}
		if data.Version != version.Release || data.Commit != version.Commit {
			t.Errorf("version info = %s/%s, want %s/%s", data.Version, data.Commit, version.Release, version.Commit)
		}
	})

	t.Run("reports a missing config file", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "does-not-exist.yaml")
		db := &dbtest.FakeDB{}
		h := NewHandler(HandlerParams{
			Database:      dbtest.NewManager(t, db),
			Configuration: newTestConfigFromFile(t, missing, dbSettings),
		})

		rec := httptest.NewRecorder()
		h.HealthStatus(rec, newRequest("", nil))

		var data HealthStatusResponse
		decodeData(t, decodeResponse(t, rec), &data)
		if data.ConfigFile {
			t.Error("configFile has to be false when the configured file does not exist")
		}
	})

	t.Run("reports a system without database configuration", func(t *testing.T) {
		db := &dbtest.FakeDB{}
		h := NewHandler(HandlerParams{
			Database:      dbtest.NewManager(t, db),
			Configuration: newTestConfigWithoutFile(t, nil),
		})

		rec := httptest.NewRecorder()
		h.HealthStatus(rec, newRequest("", nil))

		var data HealthStatusResponse
		decodeData(t, decodeResponse(t, rec), &data)
		if !data.ConfigFile {
			t.Error("configFile has to be true when no config file is expected")
		}
		if data.ConfigDB {
			t.Error("configDb has to be false without database settings")
		}
		if data.DBReachable {
			t.Error("dbReachable has to be false without database settings")
		}
	})

	t.Run("reports an unreachable database", func(t *testing.T) {
		db := &dbtest.FakeDB{
			PingFn: func() error { return errNotFound },
			UsersExistFn: func() (bool, error) {
				t.Error("users must not be queried when the database is unreachable")
				return false, nil
			},
		}
		h := NewHandler(HandlerParams{
			Database:      dbtest.NewManager(t, db),
			Configuration: newTestConfig(t, dbSettings),
		})

		rec := httptest.NewRecorder()
		h.HealthStatus(rec, newRequest("", nil))

		var data HealthStatusResponse
		decodeData(t, decodeResponse(t, rec), &data)
		if !data.ConfigDB {
			t.Error("configDb has to be true when database settings are present")
		}
		if data.DBReachable {
			t.Error("dbReachable has to be false for a failing ping")
		}
	})

	t.Run("reports an unconfigured database handle", func(t *testing.T) {
		db := &dbtest.FakeDB{
			IsConfiguredFn: func() bool { return false },
			PingFn: func() error {
				t.Error("an unconfigured database handle must not be pinged")
				return nil
			},
		}
		h := NewHandler(HandlerParams{
			Database:      dbtest.NewManager(t, db),
			Configuration: newTestConfig(t, dbSettings),
		})

		rec := httptest.NewRecorder()
		h.HealthStatus(rec, newRequest("", nil))

		var data HealthStatusResponse
		decodeData(t, decodeResponse(t, rec), &data)
		if data.DBReachable {
			t.Error("dbReachable has to be false for an unconfigured database handle")
		}
	})

	t.Run("fails without a database connection", func(t *testing.T) {
		h := NewHandler(HandlerParams{
			Database:      database.NewManager(database.Settings{}),
			Configuration: newTestConfig(t, dbSettings),
		})

		rec := httptest.NewRecorder()
		h.HealthStatus(rec, newRequest("", nil))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response without a database connection")
		}
	})

	t.Run("fails when the user check errors", func(t *testing.T) {
		db := &dbtest.FakeDB{
			UsersExistFn: func() (bool, error) { return false, errNotFound },
		}
		h := NewHandler(HandlerParams{
			Database:      dbtest.NewManager(t, db),
			Configuration: newTestConfig(t, dbSettings),
		})

		rec := httptest.NewRecorder()
		h.HealthStatus(rec, newRequest("", nil))

		if rec.Code != 500 {
			t.Errorf("status = %d, want 500", rec.Code)
		}
	})
}

func TestGetLogs(t *testing.T) {
	t.Run("returns the logs for the configured currency", func(t *testing.T) {
		var gotCurrency string
		db := &dbtest.FakeDB{
			GetLogsFn: func(currency string) ([]database.Log, error) {
				gotCurrency = currency
				return []database.Log{{ID: 1, Type: "heat", LogMessage: "Ann was riding"}}, nil
			},
		}
		h := newTestHandler(t, db, map[string]string{"currency": "CHF"})

		rec := httptest.NewRecorder()
		h.GetLogs(rec, newRequest("", &HandlerCtx{Database: db}))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if gotCurrency != "CHF" {
			t.Errorf("currency = %q, want CHF", gotCurrency)
		}
		var data GetLogsResponse
		decodeData(t, resp, &data)
		if len(data) != 1 || data[0].LogMessage != "Ann was riding" {
			t.Errorf("logs = %+v, want the single log entry", data)
		}
	})

	t.Run("fails without a configured currency", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetLogsFn: func(string) ([]database.Log, error) {
				t.Error("the logs must not be queried without a currency")
				return nil, nil
			},
		}
		h := newTestHandler(t, db, map[string]string{"currency": ""})

		rec := httptest.NewRecorder()
		h.GetLogs(rec, newRequest("", &HandlerCtx{Database: db}))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetLogsFn: func(string) ([]database.Log, error) { return nil, errNotFound },
		}
		h := newTestHandler(t, db, map[string]string{"currency": "CHF"})

		rec := httptest.NewRecorder()
		h.GetLogs(rec, newRequest("", &HandlerCtx{Database: db}))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestGetBoatTelemetryConfigurationChecks(t *testing.T) {
	// The happy path talks to the myNautique cloud service, so only the
	// configuration guards in front of it are covered here.
	complete := map[string]string{
		"mynautique.enabled":       "true",
		"mynautique.api.key":       "the-api-key",
		"mynautique.user":          "boat@example.com",
		"mynautique.password":      "secret",
		"mynautique.fuel.capacity": "150",
	}
	without := func(key string) map[string]string {
		values := map[string]string{}
		for k, v := range complete {
			values[k] = v
		}
		values[key] = ""
		return values
	}

	tests := []struct {
		name    string
		values  map[string]string
		wantMsg string
	}{
		{name: "disabled", values: without("mynautique.enabled"), wantMsg: "myNautique is not configured"},
		{name: "no api key", values: without("mynautique.api.key"), wantMsg: "myNautique API key missing"},
		{name: "no user", values: without("mynautique.user"), wantMsg: "myNautique user missing"},
		{name: "no password", values: without("mynautique.password"), wantMsg: "myNautique password missing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &dbtest.FakeDB{}
			h := newTestHandler(t, db, tt.values)

			rec := httptest.NewRecorder()
			r := newRequest("", &HandlerCtx{Database: db})
			h.GetBoatTelemetry(rec, r, GetBoatInfoRequest{BoatID: 12345}, GetHandlerContext(r))

			resp := decodeResponse(t, rec)
			if resp.OK {
				t.Fatal("expected a failure response")
			}
			if resp.Msg != tt.wantMsg {
				t.Errorf("msg = %q, want %q", resp.Msg, tt.wantMsg)
			}
		})
	}
}

func TestGetBoatTelemetryValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{name: "valid payload", payload: `{"boat_id":12345}`},
		{name: "missing boat id", payload: `{}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReadBodyAndValidate(newRequest(tt.payload, nil), &GetBoatInfoRequest{})
			if (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFailureResponseShape(t *testing.T) {
	resp := FailureResponse("nope")
	if resp.OK {
		t.Error("a failure response must not be ok")
	}
	if resp.Msg != "nope" {
		t.Errorf("msg = %q, want nope", resp.Msg)
	}
	if resp.Data != nil {
		t.Errorf("data = %v, want nil", resp.Data)
	}
}
