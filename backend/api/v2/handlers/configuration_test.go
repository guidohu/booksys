package handlers

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"server/database"
	"server/database/dbtest"
)

// fullConfigValues is a configuration where every property the endpoints read
// is set to a recognisable value.
func fullConfigValues() map[string]string {
	return map[string]string{
		"currency":                 "CHF",
		"engine.hour.format":       "hh.h",
		"fuel.payment.type":        "instant",
		"location.address":         "Lake Road 1",
		"location.latitude":        "47.376900",
		"location.longitude":       "8.541700",
		"location.map":             "https://www.google.com/maps/embed?pb=abc",
		"location.timezone":        "Europe/Zurich",
		"logo.file":                "abcdef0123.png",
		"mynautique.enabled":       "true",
		"mynautique.boat.id":       "12345",
		"mynautique.fuel.capacity": "150",
		"mynautique.user":          "boat@example.com",
		"mynautique.password":      "the-secret",
		"payment.account.bic":      "ABCDCHZZ",
		"payment.account.comment":  "pay here",
		"payment.account.iban":     "CH9300762011623852957",
		"payment.account.owner":    "The Club",
		"recaptcha.publickey":      "0123456789012345678901234567890123456789",
		"recaptcha.privatekey":     "9876543210987654321098765432109876543210",
		"smtp.sender":              "club@example.com",
		"smtp.server":              "smtp.example.com",
		"smtp.username":            "club",
		"smtp.password":            "smtp-secret",
		"url":                      "www.example.com",
	}
}

func TestGetPublicConfiguration(t *testing.T) {
	db := &dbtest.FakeDB{}
	h := newTestHandler(t, db, fullConfigValues())

	rec := httptest.NewRecorder()
	h.GetPublicConfiguration(rec, newRequest("", &HandlerCtx{Database: db}))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data PublicConfigurationMessage
	decodeData(t, resp, &data)

	if data.Currency != "CHF" || data.LocationTimeZone != "Europe/Zurich" {
		t.Errorf("configuration = %+v, does not match the settings", data)
	}
	if data.LocationLatitude != 47.3769 || data.LocationLongitude != 8.5417 {
		t.Errorf("coordinates = %v/%v, want 47.3769/8.5417", data.LocationLatitude, data.LocationLongitude)
	}
	if !data.MyNautiqueEnabled || data.MyNautiqueBoatID != 12345 || data.MyNautiqueFuelCapacity != 150 {
		t.Errorf("myNautique settings = %+v, do not match the settings", data)
	}
	if data.RecaptchaPublicKey != fullConfigValues()["recaptcha.publickey"] {
		t.Errorf("recaptcha public key = %q, does not match the setting", data.RecaptchaPublicKey)
	}

	// The public configuration must not leak any secret.
	raw := string(resp.Data)
	for _, secret := range []string{"the-secret", "smtp-secret", fullConfigValues()["recaptcha.privatekey"], "boat@example.com"} {
		if strings.Contains(raw, secret) {
			t.Errorf("the public configuration leaks %q", secret)
		}
	}
}

func TestGetConfiguration(t *testing.T) {
	db := &dbtest.FakeDB{}
	h := newTestHandler(t, db, fullConfigValues())

	rec := httptest.NewRecorder()
	h.GetConfiguration(rec, newRequest("", &HandlerCtx{Database: db}))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data ConfigurationMessage
	decodeData(t, resp, &data)

	if data.Currency != "CHF" || data.URL != "www.example.com" {
		t.Errorf("configuration = %+v, does not match the settings", data)
	}
	if data.MyNautiqueUser != "boat@example.com" {
		t.Errorf("mynautique user = %q, want boat@example.com", data.MyNautiqueUser)
	}
	// Passwords are replaced by the "hidden" marker, which is also what the
	// database layer uses to detect "do not change this value".
	if data.MyNautiquePassword != database.HiddenSecret {
		t.Errorf("mynautique password = %q, want %q", data.MyNautiquePassword, database.HiddenSecret)
	}
	if data.SMTPPassword != database.HiddenSecret {
		t.Errorf("smtp password = %q, want %q", data.SMTPPassword, database.HiddenSecret)
	}
	raw := string(resp.Data)
	for _, secret := range []string{"the-secret", "smtp-secret"} {
		if strings.Contains(raw, secret) {
			t.Errorf("the admin configuration leaks the stored password %q", secret)
		}
	}
}

func TestSetConfiguration(t *testing.T) {
	t.Run("writes every property", func(t *testing.T) {
		var stored []database.Configuration
		db := &dbtest.FakeDB{
			UpdateOrInsertPropertyValuesFn: func(conf []database.Configuration) error {
				stored = conf
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		req := ConfigurationMessage{
			Currency:               "EUR",
			EngineHourFormat:       "hh:mm",
			FuelPaymentType:        "billed",
			LocationAddress:        "Lake Road 2",
			LocationLatitude:       47.5,
			LocationLongitude:      8.75,
			LocationTimeZone:       "Europe/Berlin",
			MyNautiqueBoatID:       999,
			MyNautiqueEnabled:      true,
			MyNautiqueFuelCapacity: 200,
			MyNautiquePassword:     "np",
			MyNautiqueUser:         "b@example.com",
			SMTPSender:             "a@example.com",
			URL:                    "www.example.org",
		}
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.SetConfiguration(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}

		got := map[string]string{}
		for _, c := range stored {
			got[c.Property] = c.Value
		}
		want := map[string]string{
			"currency":                 "EUR",
			"engine.hour.format":       "hh:mm",
			"fuel.payment.type":        "billed",
			"location.address":         "Lake Road 2",
			"location.latitude":        "47.500000",
			"location.longitude":       "8.750000",
			"location.timezone":        "Europe/Berlin",
			"mynautique.boat.id":       "999",
			"mynautique.enabled":       "true",
			"mynautique.fuel.capacity": "200",
			"url":                      "www.example.org",
		}
		for key, value := range want {
			if got[key] != value {
				t.Errorf("property %q = %q, want %q", key, got[key], value)
			}
		}

		// Every property that is written has to be accepted by the database
		// layer, otherwise it is silently dropped.
		for _, c := range stored {
			if _, allowed := database.AllowedProperties[c.Property]; !allowed {
				t.Errorf("property %q is written but not in database.AllowedProperties", c.Property)
			}
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{
			UpdateOrInsertPropertyValuesFn: func([]database.Configuration) error { return errNotFound },
		}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.SetConfiguration(rec, r, ConfigurationMessage{}, GetHandlerContext(r))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestConfigurationMessageValidation(t *testing.T) {
	uploadDir := t.TempDir()
	base := map[string]any{
		"currency":            "CHF",
		"engine_hour_format":  "hh.h",
		"fuel_payment_type":   "instant",
		"location_latitude":   47.3769,
		"location_longitude":  8.5417,
		"location_time_zone":  "Europe/Zurich",
		"mynautique_enabled":  false,
		"payment_account_bic": "ABCDCHZZ",
	}
	payload := func(overrides map[string]any) string {
		m := map[string]any{}
		for k, v := range base {
			m[k] = v
		}
		for k, v := range overrides {
			m[k] = v
		}
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatalf("cannot build the payload: %v", err)
		}
		return string(b)
	}

	tests := []struct {
		name      string
		overrides map[string]any
		wantErr   bool
	}{
		{name: "minimal valid configuration"},
		{name: "invalid currency", overrides: map[string]any{"currency": "XYZ"}, wantErr: true},
		{name: "invalid engine hour format", overrides: map[string]any{"engine_hour_format": "h"}, wantErr: true},
		{name: "invalid fuel payment type", overrides: map[string]any{"fuel_payment_type": "later"}, wantErr: true},
		{name: "latitude out of range", overrides: map[string]any{"location_latitude": 100.0}, wantErr: true},
		{name: "longitude out of range", overrides: map[string]any{"location_longitude": 200.0}, wantErr: true},
		{name: "missing timezone", overrides: map[string]any{"location_time_zone": ""}, wantErr: true},
		{name: "non google maps url", overrides: map[string]any{"location_map": "https://example.com/map"}, wantErr: true},
		{name: "google maps embed url", overrides: map[string]any{"location_map": "https://www.google.com/maps/embed?pb=abc"}},
		{name: "short recaptcha key", overrides: map[string]any{"recaptcha_publickey": "tooshort"}, wantErr: true},
		{
			name: "recaptcha public key without private key",
			overrides: map[string]any{
				"recaptcha_publickey": "0123456789012345678901234567890123456789",
			},
			wantErr: true,
		},
		{
			name: "recaptcha key pair",
			overrides: map[string]any{
				"recaptcha_publickey":  "0123456789012345678901234567890123456789",
				"recaptcha_privatekey": "9876543210987654321098765432109876543210",
			},
		},
		{
			name: "mynautique enabled without credentials",
			overrides: map[string]any{
				"mynautique_enabled": true,
			},
			wantErr: true,
		},
		{
			name: "mynautique enabled with credentials",
			overrides: map[string]any{
				"mynautique_enabled":       true,
				"mynautique_user":          "boat@example.com",
				"mynautique_password":      "secret",
				"mynautique_boat_id":       12345,
				"mynautique_fuel_capacity": 150,
			},
		},
		{name: "smtp sender without a server", overrides: map[string]any{"smtp_sender": "a@example.com"}, wantErr: true},
		{
			name: "smtp configuration",
			overrides: map[string]any{
				"smtp_sender":   "a@example.com",
				"smtp_server":   "smtp.example.com",
				"smtp_username": "club",
			},
		},
		{name: "invalid url", overrides: map[string]any{"url": "http://example.com/path"}, wantErr: true},
		{name: "logo path with traversal", overrides: map[string]any{"logo_file": "../secrets.png"}, wantErr: true},
		{name: "hex logo path", overrides: map[string]any{"logo_file": "abcdef0123.png"}},
		{name: "address with markup", overrides: map[string]any{"location_address": "<script>"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := newTestConfig(t, map[string]string{"http.uploadpath": uploadDir})
			hCtx := &HandlerCtx{Config: conf}
			err := ReadBodyAndValidate(newRequest(payload(tt.overrides), hCtx), &ConfigurationMessage{}, ConfigurationMessageValidationErrors)
			if (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetupDBConfig(t *testing.T) {
	validPayload := `{"db_server":"db.example.com:3306","db_name":"booksys","db_user":"booksys","db_password":"secret"}`

	t.Run("stores the settings in the config file", func(t *testing.T) {
		configFile := filepath.Join(t.TempDir(), "config.yaml")
		if err := os.WriteFile(configFile, []byte("http:\n  port: 8080\n"), 0o600); err != nil {
			t.Fatalf("cannot write the config file: %v", err)
		}
		db := &dbtest.FakeDB{}
		h := NewHandler(HandlerParams{
			Database:      dbtest.NewManager(t, db),
			Configuration: newTestConfigFromFile(t, configFile, nil),
		})

		rec := httptest.NewRecorder()
		h.SetupDBConfig(rec, newRequest(validPayload, &HandlerCtx{Database: db}))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		content, err := os.ReadFile(configFile)
		if err != nil {
			t.Fatalf("cannot read the config file: %v", err)
		}
		for _, want := range []string{"db.example.com", "3306", "booksys"} {
			if !strings.Contains(string(content), want) {
				t.Errorf("the config file does not contain %q:\n%s", want, content)
			}
		}
	})

	t.Run("failures", func(t *testing.T) {
		tests := []struct {
			name    string
			payload string
			// values are applied on top of an isolated config file.
			values     map[string]string
			noConfFile bool
			db         *dbtest.FakeDB
		}{
			{
				name:       "no config file configured",
				payload:    validPayload,
				noConfFile: true,
				db:         &dbtest.FakeDB{},
			},
			{
				name:    "database already configured",
				payload: validPayload,
				values: map[string]string{
					"database.protocol": "tcp",
					"database.user":     "u",
					"database.host":     "h",
					"database.port":     "3306",
					"database.dbname":   "d",
				},
				db: &dbtest.FakeDB{},
			},
			{
				name:    "invalid payload",
				payload: `{"db_server":"db.example.com:3306"}`,
				db:      &dbtest.FakeDB{},
			},
			{
				name:    "server without a port",
				payload: `{"db_server":"db.example.com","db_name":"booksys","db_user":"booksys","db_password":"secret"}`,
				db:      &dbtest.FakeDB{},
			},
			{
				name:    "database is not reachable",
				payload: validPayload,
				// The initial connection succeeds, the reconnect with the
				// submitted credentials fails.
				db: &dbtest.FakeDB{ConnectFn: failAfterFirstCall()},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var conf = newTestConfig(t, tt.values)
				if tt.noConfFile {
					conf = newTestConfigWithoutFile(t, tt.values)
				}
				h := NewHandler(HandlerParams{
					Database:      dbtest.NewManager(t, tt.db),
					Configuration: conf,
				})

				rec := httptest.NewRecorder()
				h.SetupDBConfig(rec, newRequest(tt.payload, &HandlerCtx{Database: tt.db}))

				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestSetupDBConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{name: "valid payload", payload: `{"db_server":"db:3306","db_name":"booksys","db_user":"u","db_password":"p"}`},
		{name: "server without a port", payload: `{"db_server":"db","db_name":"booksys","db_user":"u","db_password":"p"}`, wantErr: true},
		{name: "non alphanumeric database name", payload: `{"db_server":"db:3306","db_name":"book-sys","db_user":"u","db_password":"p"}`, wantErr: true},
		{name: "missing password", payload: `{"db_server":"db:3306","db_name":"booksys","db_user":"u"}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReadBodyAndValidate(newRequest(tt.payload, nil), &SetupDBConfigRequest{})
			if (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetLogoPath(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]string
		want   string
	}{
		{
			name:   "logo configured",
			values: map[string]string{"logo.file": "abcdef.png", "http.uploadpath": "/srv/uploads"},
			want:   "/srv/uploads/abcdef.png",
		},
		{
			name:   "no logo configured",
			values: map[string]string{"logo.file": "", "http.uploadpath": "/srv/uploads"},
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &dbtest.FakeDB{}
			h := newTestHandler(t, db, tt.values)

			rec := httptest.NewRecorder()
			h.GetLogoPath(rec, newRequest("", &HandlerCtx{Database: db}))

			resp := decodeResponse(t, rec)
			if !resp.OK {
				t.Fatalf("unexpected failure: %s", resp.Msg)
			}
			var data GetLogoPathResponse
			decodeData(t, resp, &data)
			if data.URI != tt.want {
				t.Errorf("uri = %q, want %q", data.URI, tt.want)
			}
		})
	}
}

func TestGetRecaptchaKey(t *testing.T) {
	key := "0123456789012345678901234567890123456789"
	db := &dbtest.FakeDB{}
	h := newTestHandler(t, db, map[string]string{"recaptcha.publickey": key})

	rec := httptest.NewRecorder()
	h.GetRecaptchaKey(rec, newRequest("", &HandlerCtx{Database: db}))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data GetRecaptchaKeyResponse
	decodeData(t, resp, &data)
	if data.Key != key {
		t.Errorf("key = %q, want %q", data.Key, key)
	}
}

// pngUpload builds a multipart request that carries a small PNG in the given
// form field.
func pngUpload(t *testing.T, field, filename string, size int) *http.Request {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	var raw bytes.Buffer
	if err := png.Encode(&raw, img); err != nil {
		t.Fatalf("cannot encode the test image: %v", err)
	}
	content := raw.Bytes()
	if size > len(content) {
		content = append(content, bytes.Repeat([]byte{0}, size-len(content))...)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile(field, filename)
	if err != nil {
		t.Fatalf("cannot create the form file: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("cannot write the form file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("cannot close the multipart writer: %v", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/api/v2/admin/upload/logo", &body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	return r
}

func TestUploadLogoFile(t *testing.T) {
	t.Run("stores an image under a content based name", func(t *testing.T) {
		uploadDir := t.TempDir()
		db := &dbtest.FakeDB{}
		h := newTestHandler(t, db, map[string]string{"http.uploadpath": uploadDir})

		rec := httptest.NewRecorder()
		h.UploadLogoFile(rec, pngUpload(t, "logo", "club-logo.png", 0))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data UploadLogoFileResponse
		decodeData(t, resp, &data)
		if filepath.Ext(data.FileName) != ".png" {
			t.Errorf("filename = %q, want a .png extension", data.FileName)
		}
		if _, err := os.Stat(filepath.Join(uploadDir, data.FileName)); err != nil {
			t.Errorf("the uploaded file was not stored: %v", err)
		}
		// The stored name has to pass the validator that guards the config.
		if err := ReadBodyAndValidate(
			newRequest(`{"logo_file":"`+data.FileName+`"}`, &HandlerCtx{Config: newTestConfig(t, map[string]string{"http.uploadpath": uploadDir})}),
			&struct {
				LogoFilePath string `json:"logo_file" validate:"omitempty,uploadfile"`
			}{},
		); err != nil {
			t.Errorf("the generated filename %q is rejected by the uploadfile validator: %v", data.FileName, err)
		}
	})

	t.Run("rejects a file that is not an image", func(t *testing.T) {
		uploadDir := t.TempDir()
		db := &dbtest.FakeDB{}
		h := newTestHandler(t, db, map[string]string{"http.uploadpath": uploadDir})

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, err := writer.CreateFormFile("logo", "evil.png")
		if err != nil {
			t.Fatalf("cannot create the form file: %v", err)
		}
		part.Write([]byte(strings.Repeat("<?php echo 'hi'; ?>", 40)))
		writer.Close()
		r := httptest.NewRequest(http.MethodPost, "/", &body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		rec := httptest.NewRecorder()
		h.UploadLogoFile(rec, r)

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("a non image upload has to be rejected")
		}
		entries, _ := os.ReadDir(uploadDir)
		if len(entries) != 0 {
			t.Errorf("a rejected upload must not be stored, found %d files", len(entries))
		}
	})

	t.Run("rejects a file larger than 512kB", func(t *testing.T) {
		uploadDir := t.TempDir()
		db := &dbtest.FakeDB{}
		h := newTestHandler(t, db, map[string]string{"http.uploadpath": uploadDir})

		rec := httptest.NewRecorder()
		h.UploadLogoFile(rec, pngUpload(t, "logo", "big.png", 513*1024))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("an oversized upload has to be rejected")
		}
		entries, _ := os.ReadDir(uploadDir)
		if len(entries) != 0 {
			t.Errorf("a rejected upload must not be stored, found %d files", len(entries))
		}
	})

	t.Run("rejects a request without the logo field", func(t *testing.T) {
		db := &dbtest.FakeDB{}
		h := newTestHandler(t, db, map[string]string{"http.uploadpath": t.TempDir()})

		rec := httptest.NewRecorder()
		h.UploadLogoFile(rec, pngUpload(t, "other", "club-logo.png", 0))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("a request without the logo field has to be rejected")
		}
	})
}

func TestSetupMyNautiqueCredentials(t *testing.T) {
	t.Run("stores the credentials when they are not configured yet", func(t *testing.T) {
		var stored []database.Configuration
		db := &dbtest.FakeDB{
			UpdateOrInsertPropertyValuesFn: func(conf []database.Configuration) error {
				stored = append(stored, conf...)
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		req := SetupMyNautiqueCredentialsRequest{Enabled: true, User: "boat@example.com", Password: "secret"}
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.SetupMyNautiqueCredentials(rec, r, req, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		got := map[string]string{}
		for _, c := range stored {
			got[c.Property] = c.Value
		}
		want := map[string]string{
			"mynautique.enabled":  "true",
			"mynautique.user":     "boat@example.com",
			"mynautique.password": "secret",
		}
		for key, value := range want {
			if got[key] != value {
				t.Errorf("property %q = %q, want %q", key, got[key], value)
			}
		}
	})

	t.Run("refuses to overwrite an existing configuration", func(t *testing.T) {
		db := &dbtest.FakeDB{
			UpdateOrInsertPropertyValuesFn: func([]database.Configuration) error {
				t.Error("an existing myNautique configuration must not be overwritten")
				return nil
			},
		}
		h := newTestHandler(t, db, map[string]string{"mynautique.enabled": "true"})

		req := SetupMyNautiqueCredentialsRequest{Enabled: false, User: "other@example.com"}
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.SetupMyNautiqueCredentials(rec, r, req, GetHandlerContext(r))

		if rec.Body.Len() != 0 {
			t.Errorf("expected no response body for a rejected setup, got %q", rec.Body.String())
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{
			UpdateOrInsertPropertyValuesFn: func([]database.Configuration) error { return errNotFound },
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.SetupMyNautiqueCredentials(rec, r, SetupMyNautiqueCredentialsRequest{Enabled: true}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestSetupMyNautiqueCredentialsValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{name: "valid payload", payload: `{"mynautique_enabled":true,"mynautique_user":"a@b.co","mynautique_password":"p"}`},
		{name: "empty user is allowed", payload: `{"mynautique_enabled":false}`},
		{name: "invalid user", payload: `{"mynautique_enabled":true,"mynautique_user":"not-an-email"}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReadBodyAndValidate(newRequest(tt.payload, nil), &SetupMyNautiqueCredentialsRequest{})
			if (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
