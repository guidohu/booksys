package validator

import (
	"path/filepath"
	"testing"

	"github.com/go-playground/validator/v10"
)

// validate builds a validator that has all custom validations registered and
// runs it against a single field of the given struct.
func newValidator(t *testing.T, uploadBaseDir string) *validator.Validate {
	t.Helper()
	v := validator.New()
	if err := v.RegisterValidation("googlemapsurl", GoogleMapsURL); err != nil {
		t.Fatalf("cannot register googlemapsurl: %v", err)
	}
	if err := v.RegisterValidation("recaptchakey", RecaptchaKey); err != nil {
		t.Fatalf("cannot register recaptchakey: %v", err)
	}
	if err := v.RegisterValidation("strongpassword", PasswordStrength); err != nil {
		t.Fatalf("cannot register strongpassword: %v", err)
	}
	if err := v.RegisterValidation("sessiontype", SessionType); err != nil {
		t.Fatalf("cannot register sessiontype: %v", err)
	}
	if err := v.RegisterValidation("expensetype", ExpenseType); err != nil {
		t.Fatalf("cannot register expensetype: %v", err)
	}
	if err := v.RegisterValidation("tableid", TableID); err != nil {
		t.Fatalf("cannot register tableid: %v", err)
	}
	if err := v.RegisterValidation("currency", Currency); err != nil {
		t.Fatalf("cannot register currency: %v", err)
	}
	if err := v.RegisterValidation("uploadfile", func(fl validator.FieldLevel) bool {
		return UploadedFilePath(fl, uploadBaseDir)
	}); err != nil {
		t.Fatalf("cannot register uploadfile: %v", err)
	}
	return v
}

func TestGoogleMapsURL(t *testing.T) {
	type s struct {
		URL string `validate:"googlemapsurl"`
	}
	v := newValidator(t, "")

	tests := []struct {
		name  string
		url   string
		valid bool
	}{
		{name: "embed url", url: "https://www.google.com/maps/embed?pb=!1m18!1m12", valid: true},
		{name: "plain maps url", url: "https://www.google.com/maps/place/Zurich"},
		{name: "http instead of https", url: "http://www.google.com/maps/embed?pb=abc"},
		{name: "other host", url: "https://evil.example.com/maps/embed?pb=abc"},
		{name: "host as a prefix only", url: "https://www.google.com.evil.example.com/maps/embed?pb=abc"},
		{name: "empty", url: ""},
		{name: "quote injection", url: `https://www.google.com/maps/embed?pb=abc" onload="alert(1)`},
		{name: "whitespace injection", url: "https://www.google.com/maps/embed?pb=abc def"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(s{URL: tt.url})
			if (err == nil) != tt.valid {
				t.Errorf("GoogleMapsURL(%q) valid = %v, want %v", tt.url, err == nil, tt.valid)
			}
		})
	}
}

func TestRecaptchaKey(t *testing.T) {
	type s struct {
		Key string `validate:"recaptchakey"`
	}
	v := newValidator(t, "")

	tests := []struct {
		name  string
		key   string
		valid bool
	}{
		{name: "40 characters", key: "0123456789012345678901234567890123456789", valid: true},
		{name: "with dashes and underscores", key: "6Lc_abcDEF-01234567890123456789012345678", valid: true},
		{name: "39 characters", key: "012345678901234567890123456789012345678"},
		{name: "41 characters", key: "01234567890123456789012345678901234567890"},
		{name: "invalid character", key: "0123456789012345678901234567890123456!89"},
		{name: "empty", key: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(s{Key: tt.key})
			if (err == nil) != tt.valid {
				t.Errorf("RecaptchaKey(%q) valid = %v, want %v", tt.key, err == nil, tt.valid)
			}
		})
	}
}

func TestPasswordStrength(t *testing.T) {
	type s struct {
		Password string `validate:"strongpassword"`
	}
	v := newValidator(t, "")

	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{name: "long password with all classes", password: "Sup3rSecret!Password", valid: true},
		{name: "exactly twelve characters", password: "Abcdefghijk!", valid: true},
		{name: "eleven characters", password: "Abcdefghij!"},
		{name: "no upper case", password: "abcdefghijk!"},
		{name: "no lower case", password: "ABCDEFGHIJK!"},
		{name: "no special character", password: "Abcdefghijkl"},
		{name: "empty", password: ""},
		// The documentation says "a number or special character", but only the
		// listed special characters are accepted.
		{name: "digits only as the third class", password: "Abcdefghijk1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(s{Password: tt.password})
			if (err == nil) != tt.valid {
				t.Errorf("PasswordStrength(%q) valid = %v, want %v", tt.password, err == nil, tt.valid)
			}
		})
	}
}

func TestSessionType(t *testing.T) {
	type s struct {
		Type uint8 `validate:"sessiontype"`
	}
	v := newValidator(t, "")

	tests := []struct {
		name  string
		value uint8
		valid bool
	}{
		{name: "default session", value: 1, valid: true},
		{name: "course session", value: 2, valid: true},
		{name: "zero", value: 0},
		{name: "unknown type", value: 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(s{Type: tt.value})
			if (err == nil) != tt.valid {
				t.Errorf("SessionType(%d) valid = %v, want %v", tt.value, err == nil, tt.valid)
			}
		})
	}
}

// SessionType recovers from a panic when the field is not an unsigned integer.
func TestSessionTypeOnNonUintField(t *testing.T) {
	type s struct {
		Type string `validate:"sessiontype"`
	}
	v := newValidator(t, "")

	err := v.Struct(s{Type: "default"})
	if err == nil {
		t.Error("a non numeric session type has to be rejected")
	}
}

func TestExpenseType(t *testing.T) {
	type s struct {
		Type uint64 `validate:"expensetype"`
	}
	v := newValidator(t, "")

	tests := []struct {
		name  string
		value uint64
		valid bool
	}{
		{name: "fuel direct", value: 1, valid: true},
		{name: "session", value: 5, valid: true},
		{name: "fuel bill", value: 10, valid: true},
		{name: "zero", value: 0},
		{name: "unknown type", value: 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(s{Type: tt.value})
			if (err == nil) != tt.valid {
				t.Errorf("ExpenseType(%d) valid = %v, want %v", tt.value, err == nil, tt.valid)
			}
		})
	}
}

func TestTableID(t *testing.T) {
	type s struct {
		TableID uint64 `validate:"tableid"`
	}
	v := newValidator(t, "")

	tests := []struct {
		name  string
		value uint64
		valid bool
	}{
		{name: "expenditure", value: 0, valid: true},
		{name: "payment", value: 1, valid: true},
		{name: "boat fuel", value: 2, valid: true},
		{name: "unknown table", value: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(s{TableID: tt.value})
			if (err == nil) != tt.valid {
				t.Errorf("TableID(%d) valid = %v, want %v", tt.value, err == nil, tt.valid)
			}
		})
	}
}

func TestCurrency(t *testing.T) {
	type s struct {
		Currency string `validate:"currency"`
	}
	v := newValidator(t, "")

	tests := []struct {
		name     string
		currency string
		valid    bool
	}{
		{name: "swiss franc", currency: "CHF", valid: true},
		{name: "euro", currency: "EUR", valid: true},
		{name: "us dollar", currency: "USD", valid: true},
		{name: "unknown code", currency: "XYZ"},
		{name: "lower case", currency: "chf", valid: true},
		{name: "empty", currency: ""},
		{name: "sql injection attempt", currency: "CHF'; DROP TABLE heat; --"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(s{Currency: tt.currency})
			if (err == nil) != tt.valid {
				t.Errorf("Currency(%q) valid = %v, want %v", tt.currency, err == nil, tt.valid)
			}
		})
	}
}

func TestUploadedFilePath(t *testing.T) {
	base := t.TempDir()
	type s struct {
		Path string `validate:"uploadfile"`
	}
	v := newValidator(t, base)

	tests := []struct {
		name  string
		path  string
		valid bool
	}{
		{name: "hex name with extension", path: "0123456789abcdef.png", valid: true},
		{name: "hex name without extension", path: "0123456789abcdef", valid: true},
		{name: "upper case hex", path: "0123456789ABCDEF.jpg", valid: true},
		{name: "non hex name", path: "logo.png"},
		{name: "relative traversal", path: "../0123456789abcdef.png"},
		{name: "absolute path", path: "/etc/passwd"},
		{name: "subdirectory", path: "sub/0123456789abcdef.png"},
		{name: "empty", path: ""},
		{name: "null byte", path: "0123456789abcdef\x00.png"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(s{Path: tt.path})
			if (err == nil) != tt.valid {
				t.Errorf("UploadedFilePath(%q) valid = %v, want %v", tt.path, err == nil, tt.valid)
			}
		})
	}
}

// An accepted upload path must always resolve to a file inside the base
// directory.
func TestUploadedFilePathStaysInsideTheBaseDirectory(t *testing.T) {
	base := t.TempDir()
	type s struct {
		Path string `validate:"uploadfile"`
	}
	v := newValidator(t, base)

	candidates := []string{
		"0123456789abcdef.png",
		"../0123456789abcdef.png",
		"../../etc/0123456789abcdef",
		"/0123456789abcdef",
		"./0123456789abcdef",
	}
	for _, candidate := range candidates {
		if err := v.Struct(s{Path: candidate}); err != nil {
			continue
		}
		resolved, err := filepath.Abs(filepath.Clean(filepath.Join(base, candidate)))
		if err != nil {
			t.Fatalf("cannot resolve %q: %v", candidate, err)
		}
		absBase, _ := filepath.Abs(base)
		if rel, err := filepath.Rel(absBase, resolved); err != nil || rel == ".." || filepath.IsAbs(rel) ||
			(len(rel) >= 3 && rel[:3] == "../") {
			t.Errorf("accepted path %q resolves to %q which is outside of %q", candidate, resolved, absBase)
		}
	}
}
