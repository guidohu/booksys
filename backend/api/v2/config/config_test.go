package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"server/database"

	"github.com/spf13/viper"
)

// newTestConfig returns a Config that is backed by a config file in a
// temporary directory. The `config` flag is set, so that findConfigFile does
// not reach for a config file that happens to exist on the machine running the
// tests.
func newTestConfig(t *testing.T, configFile string) *Config {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(configFile), 0600); err != nil {
		t.Fatalf("cannot write config file: %v", err)
	}
	flags := viper.New()
	flags.Set("config", path)
	c, err := NewConfig(flags)
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}
	return c
}

func TestIsSecret(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		{"database.password", true},
		{"mynautique.api.key", true},
		{"mynautique.password", true},
		{"recaptcha.privatekey", true},
		{"smtp.password", true},
		// Keys are compared case insensitively, like viper compares them.
		{"Database.Password", true},
		{"SMTP.PASSWORD", true},
		// Everything else is not a credential.
		{"database.user", false},
		{"http.port", false},
		{"recaptcha.publickey", false},
		{"password", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsSecret(tt.key); got != tt.want {
			t.Errorf("IsSecret(%q) = %v, want %v", tt.key, got, tt.want)
		}
	}
}

func TestToStringFullRedactsSecrets(t *testing.T) {
	c := newTestConfig(t, `database:
  user: user
  password: filesecret
mynautique:
  api.key: apisecret
`)
	// Secrets that live in the database rather than in the file.
	c.properties = map[string]database.Configuration{
		"smtp.password":        {Property: "smtp.password", Value: "smtpsecret"},
		"recaptcha.privatekey": {Property: "recaptcha.privatekey", Value: "captchasecret"},
		"mynautique.password":  {Property: "mynautique.password", Value: "mnsecret"},
		"smtp.username":        {Property: "smtp.username", Value: "smtpuser"},
	}
	plaintext := []string{"filesecret", "apisecret", "smtpsecret", "captchasecret", "mnsecret"}

	got := c.ToStringFull()
	for _, secret := range plaintext {
		if strings.Contains(got, secret) {
			t.Errorf("ToStringFull() leaks the secret %q:\n%s", secret, got)
		}
	}
	// A redacted key is still listed, only its value is hidden. The file
	// section is the one that holds database.password here.
	if !strings.Contains(got, Redacted) {
		t.Errorf("ToStringFull() contains no %q marker:\n%s", Redacted, got)
	}
	// Non secrets stay readable, otherwise the dump is of no use.
	for _, want := range []string{"database.user", "user", "smtp.username", "smtpuser"} {
		if !strings.Contains(got, want) {
			t.Errorf("ToStringFull() does not contain %q:\n%s", want, got)
		}
	}

	unsafe := c.ToStringFullUnsafe()
	for _, secret := range plaintext {
		if !strings.Contains(unsafe, secret) {
			t.Errorf("ToStringFullUnsafe() does not contain the secret %q:\n%s", secret, unsafe)
		}
	}
}

// An unset secret is reported as unset instead of redacted, as a dump is
// consulted to find out whether a credential is configured at all.
func TestToStringFullKeepsUnsetSecretsVisible(t *testing.T) {
	c := newTestConfig(t, "database:\n  user: user\n")
	got := c.ToStringFull()

	for _, line := range strings.Split(got, "\n") {
		if !strings.HasPrefix(line, "database.password") && !strings.HasPrefix(line, "smtp.password") {
			continue
		}
		if strings.Contains(line, Redacted) {
			t.Errorf("unset secret is redacted rather than reported as unset: %q", line)
		}
	}
}
