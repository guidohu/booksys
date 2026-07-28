package recaptcha

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// withVerifyServer points the package at a stub verification endpoint for the
// duration of a test.
func withVerifyServer(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	previous := recaptchaSiteVerifyUrl
	t.Cleanup(func() { recaptchaSiteVerifyUrl = previous })
	recaptchaSiteVerifyUrl = server.URL
}

func TestValid(t *testing.T) {
	t.Run("accepts a successful verification", func(t *testing.T) {
		var got url.Values
		withVerifyServer(t, func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			got = r.PostForm
			w.Write([]byte(`{"success":true,"hostname":"example.com"}`))
		})

		valid, err := Valid("the-token", "the-secret")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !valid {
			t.Error("a successful verification has to be reported as valid")
		}
		if got.Get("secret") != "the-secret" {
			t.Errorf("secret = %q, want the-secret", got.Get("secret"))
		}
		if got.Get("response") != "the-token" {
			t.Errorf("response = %q, want the-token", got.Get("response"))
		}
	})

	t.Run("rejects an unsuccessful verification", func(t *testing.T) {
		withVerifyServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"success":false,"error-codes":["invalid-input-response"]}`))
		})

		valid, err := Valid("bad-token", "the-secret")
		if err != nil {
			t.Fatalf("a rejected token is not an error: %v", err)
		}
		if valid {
			t.Error("an unsuccessful verification must not be reported as valid")
		}
	})

	t.Run("reports an unparsable response", func(t *testing.T) {
		withVerifyServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`not json`))
		})

		valid, err := Valid("the-token", "the-secret")
		if err == nil {
			t.Error("expected an error for an unparsable response")
		}
		if valid {
			t.Error("an unparsable response must not be reported as valid")
		}
	})

	t.Run("reports an unreachable endpoint", func(t *testing.T) {
		previous := recaptchaSiteVerifyUrl
		t.Cleanup(func() { recaptchaSiteVerifyUrl = previous })
		// A closed server yields a connection error.
		server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
		recaptchaSiteVerifyUrl = server.URL
		server.Close()

		valid, err := Valid("the-token", "the-secret")
		if err == nil {
			t.Error("expected an error for an unreachable endpoint")
		}
		if valid {
			t.Error("a failed request must not be reported as valid")
		}
	})
}
