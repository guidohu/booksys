// Package recaptcha verifies Google reCAPTCHA tokens against the
// siteverify endpoint.
package recaptcha

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
)

// siteVerifyURL is the endpoint tokens are verified against. It is a
// variable so that tests can point it at a stub server.
var siteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

// Response is the reply returned by the siteverify endpoint.
type Response struct {
	Success    bool     `json:"success"`
	Hostname   string   `json:"hostname"`
	ErrorCodes []string `json:"error-codes"`
}

// Valid reports whether the given token is a valid reCAPTCHA response for
// the given secret. It returns an error if the token could not be verified.
func Valid(token string, secret string) (bool, error) {
	data := url.Values{
		"secret":   {secret},
		"response": {token},
	}
	resp, err := http.PostForm(siteVerifyURL, data)
	if err != nil {
		slog.Warn("Cannot verify recaptcha token", slog.Any("error", err))
		return false, err
	}
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.Warn("Cannot decode recaptcha verification response", slog.Any("error", err))
		return false, err
	}

	if !response.Success {
		slog.Warn("Invalid recaptcha token provided")
		return false, nil
	}

	return true, nil
}
