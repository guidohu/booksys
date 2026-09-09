package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestWithMaxBodySize checks the cap that stops an unauthenticated caller from
// making the server buffer an arbitrarily large request body.
func TestWithMaxBodySize(t *testing.T) {
	const limit = 1024

	// readAll stands in for a handler that reads the whole body, which is what
	// ReadBodyAndValidate and the multipart upload both do.
	readAll := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n, err := io.Copy(io.Discard, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		fmt.Fprintf(w, "%d", n)
	})
	srv := httptest.NewServer(WithMaxBodySize(readAll, limit))
	defer srv.Close()

	post := func(t *testing.T, size int) *http.Response {
		t.Helper()
		resp, err := http.Post(srv.URL, "application/json", strings.NewReader(strings.Repeat("a", size)))
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		t.Cleanup(func() { resp.Body.Close() })
		return resp
	}

	t.Run("lets a body at the limit through", func(t *testing.T) {
		resp := post(t, limit)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want 200", resp.StatusCode)
		}
		body, _ := io.ReadAll(resp.Body)
		if got := string(body); got != fmt.Sprint(limit) {
			t.Errorf("bytes read = %s, want %d", got, limit)
		}
	})

	t.Run("cuts a body off past the limit", func(t *testing.T) {
		resp := post(t, limit*8)
		if resp.StatusCode != http.StatusRequestEntityTooLarge {
			t.Errorf("status = %d, want 413: the handler has to see a read error rather than the whole body", resp.StatusCode)
		}
	})

	t.Run("reports the limit through MaxBytesError", func(t *testing.T) {
		var seen *http.MaxBytesError
		probe := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := io.Copy(io.Discard, r.Body)
			if e, ok := err.(*http.MaxBytesError); ok {
				seen = e
			}
		})
		s := httptest.NewServer(WithMaxBodySize(probe, limit))
		defer s.Close()

		resp, err := http.Post(s.URL, "application/json", strings.NewReader(strings.Repeat("a", limit*4)))
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		resp.Body.Close()

		if seen == nil {
			t.Fatal("an oversized body has to surface as *http.MaxBytesError so callers can answer 413")
		}
		if seen.Limit != limit {
			t.Errorf("reported limit = %d, want %d", seen.Limit, limit)
		}
	})

	t.Run("the configured limit leaves room for the logo upload", func(t *testing.T) {
		// UploadLogoFile rejects files over 512kB itself; the transport cap has
		// to sit above that plus multipart framing, or valid uploads break.
		const logoLimit = 512 * 1024
		if MaxRequestBodyBytes <= logoLimit {
			t.Errorf("MaxRequestBodyBytes = %d, has to exceed the %d logo limit", MaxRequestBodyBytes, logoLimit)
		}
	})
}

// TestWithRequestBodyTooLarge checks that a body cut off by the cap is reported
// as 413 rather than surfacing the raw net/http error text.
func TestWithRequestBodyTooLarge(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}
	called := false
	handler := WithRequestBody(func(http.ResponseWriter, *http.Request, payload, *HandlerCtx) {
		called = true
	})

	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v2/test", strings.NewReader(strings.Repeat("a", 4096)))
	r.Body = http.MaxBytesReader(rec, r.Body, 128)

	handler.ServeHTTP(rec, r)

	if called {
		t.Error("the handler must not run on a body that was cut off")
	}
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
	if body := rec.Body.String(); !strings.Contains(body, "too large") {
		t.Errorf("body = %s, want a message naming the size problem", body)
	}
	if strings.Contains(rec.Body.String(), "http: request body too large") {
		t.Error("the raw net/http error text should not be handed to the client")
	}
}
