package recaptcha

import (
	"encoding/json"
	"net/http"
	"net/url"

	"golang.org/x/exp/slog"
)

var recaptchaSiteVerifyUrl = "https://www.google.com/recaptcha/api/siteverify"

type RecaptchaResponse struct {
	Success    bool     `json:"success"`
	Hostname   string   `json:"hostname"`
	ErrorCodes []string `json:"error-codes"`
}

func Valid(token string, secret string) (bool, error) {
	data := url.Values{
		"secret":   {secret},
		"response": {token},
	}
	resp, err := http.PostForm(recaptchaSiteVerifyUrl, data)
	if err != nil {
		slog.Warn("Cannot verify recaptcha token:", slog.String("error", err.Error()))
		return false, err
	}

	response := RecaptchaResponse{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		slog.Warn("Cannot verify recaptcha token:", slog.String("error", err.Error()))
		return false, err
	}

	if !response.Success {
		slog.Warn("Invalid recaptcha token provided.")
		return false, nil
	}

	return true, nil
}
