package validator

import (
	"regexp"

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
		slog.Error("cannot coompile RecaptchaKey regex")
		return false
	}
	return r.Match([]byte(fl.Field().String()))
}
