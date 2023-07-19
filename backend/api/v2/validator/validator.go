package validator

import (
	"regexp"
	"server/database"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

// GoogleMapsURL is a custom validator for `googlemapsurl` tag
func GoogleMapsURL(fl validator.FieldLevel) bool {
	r, err := regexp.Compile(`^https:\/\/www\.google\.com\/maps\/embed\?pb=[^"\s]+$`)
	if err != nil {
		slog.Error("cannot compile GoogleMapsURL regex")
		return false
	}
	return r.Match([]byte(fl.Field().String()))
}

// RecaptchaKey is a custom validator for the `recaptchakey` tag
func RecaptchaKey(fl validator.FieldLevel) bool {
	r, err := regexp.Compile(`^[0-9a-zA-Z_-]{40}$`)
	if err != nil {
		slog.Error("cannot compile RecaptchaKey regex")
		return false
	}
	return r.Match([]byte(fl.Field().String()))
}

// PasswordStrength validates if the password has at least:
// - 12 characters
// - contains a capital letter A-Z
// - contains a lower case letter A-Z
// - contains a number or special character
func PasswordStrength(fl validator.FieldLevel) bool {
	pwd := fl.Field().String()
	if len(pwd) < 12 {
		return false
	}
	if !strings.ContainsAny(pwd, "abcdefghijklmnopqrstuvwxyz") {
		return false
	}
	if !strings.ContainsAny(pwd, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false
	}
	if !strings.ContainsAny(pwd, "{}[]@!#$%^&*();:,./?<>") {
		return false
	}
	return true
}

// SessionType checks if the provided session type is valid.
func SessionType(fl validator.FieldLevel) bool {
	_, found := database.DefaultSessionTypesMap[int(fl.Field().Uint())]
	return found
}
