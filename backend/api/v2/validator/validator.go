package validator

import (
	"path/filepath"
	"regexp"
	"server/database"
	"server/validator/currency"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

var hexRegex = regexp.MustCompile("^[0-9a-fA-F]+$")

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

// SessionType checks if the provided session type is valid. It requires
// the field to be of type uint.
func SessionType(fl validator.FieldLevel) bool {
	found := false
	defer func() {
		if err := recover(); err != nil {
			slog.Error("recovered from panic, field value was not of uint family", slog.String("field", fl.Field().String()))
			found = false
		}
	}()
	_, found = database.DefaultSessionTypesMap[int(fl.Field().Uint())]
	return found
}

// ExpenseType checks if the provided expense type is valid.
func ExpenseType(fl validator.FieldLevel) bool {
	_, found := database.DefaultExpenseTypesMap[int(fl.Field().Uint())]
	return found
}

// TableID checks whether the provided table ID is valid.
func TableID(fl validator.FieldLevel) bool {
	_, found := database.TableIDMap[fl.Field().Uint()]
	return found
}

// Currency checks if the currency is an ISO4217 code.
func Currency(fl validator.FieldLevel) bool {
	return currency.IsCurrency(fl.Field().String())
}

// UploadFilePath sanitizes the file path to be safe for usage
// within the app.
func UploadedFilePath(fl validator.FieldLevel, baseDir string) bool {
	// Safe filename:
	// - filename (extension removed) only contains hexadecimal characters.
	filename := strings.TrimSuffix(fl.Field().String(), filepath.Ext(fl.Field().String()))
	if !hexRegex.MatchString(filename) {
		return false
	}
	// No directory traversal tricks.
	sanitizedPath := filepath.Join(baseDir, fl.Field().String())
	sanitizedPath = filepath.Clean(sanitizedPath)
	finalPath, err := filepath.Abs(sanitizedPath)
	if err != nil {
		return false
	}
	absBaseDir, _ := filepath.Abs(baseDir)
	return strings.HasPrefix(finalPath, absBaseDir)
}
