package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"server/config"
	"server/database"
	"server/database/dbtest"
)

// errNotFound is the error the fake database reports for records that do not
// exist. It doubles as the generic "this database call failed" error.
var errNotFound = dbtest.ErrNotFound

// newTestConfig builds a Config that is backed by an isolated (temporary)
// configuration file. All values passed in are provided through the flag layer
// which has the highest priority, so any key can be set.
func newTestConfig(t *testing.T, values map[string]string) *config.Config {
	t.Helper()
	configFile := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(configFile, []byte("http:\n  port: 8080\n"), 0o600); err != nil {
		t.Fatalf("cannot write test config file: %v", err)
	}
	return newTestConfigFromFile(t, configFile, values)
}

// newTestConfigFromFile builds a Config backed by the given configuration file.
func newTestConfigFromFile(t *testing.T, configFile string, values map[string]string) *config.Config {
	t.Helper()
	v := viper.New()
	v.Set("config", configFile)
	for key, value := range values {
		v.Set(key, value)
	}
	c, err := config.NewConfig(v)
	if err != nil {
		t.Fatalf("cannot create test config: %v", err)
	}
	return c
}

// newTestConfigWithoutFile builds a Config that has no configuration file at
// all, which is how the application runs when it is configured purely through
// flags and environment variables.
func newTestConfigWithoutFile(t *testing.T, values map[string]string) *config.Config {
	t.Helper()
	// An empty `config` value makes the config layer look for a file in the
	// default locations; point it at an empty directory so nothing is found.
	t.Chdir(t.TempDir())
	t.Setenv("HOME", t.TempDir())
	v := viper.New()
	v.Set("config", "")
	for key, value := range values {
		v.Set(key, value)
	}
	c, err := config.NewConfig(v)
	if err != nil {
		t.Fatalf("cannot create test config: %v", err)
	}
	return c
}

// failAfterFirstCall returns a function that succeeds once and fails on every
// following call. It is used to let the initial connection of a test succeed
// while a later reconnect fails.
func failAfterFirstCall() func() error {
	var calls int
	return func() error {
		calls++
		if calls == 1 {
			return nil
		}
		return errNotFound
	}
}

// newTestHandler wires a Handler with the given fake database and config
// values, the same way main.go does it.
func newTestHandler(t *testing.T, db *dbtest.FakeDB, values map[string]string) *Handler {
	t.Helper()
	manager := dbtest.NewManager(t, db)
	conf := newTestConfig(t, values)
	conf.SetDB(manager)
	return NewHandler(HandlerParams{
		Database:      manager,
		Configuration: conf,
	})
}

// cookieCleared reports whether the given cookies expire the named cookie.
func cookieCleared(cookies []*http.Cookie, name string) bool {
	for _, c := range cookies {
		if c.Name == name && c.Value == "" && c.Expires.Before(time.Now()) {
			return true
		}
	}
	return false
}

// newRequest builds a POST request carrying the given JSON body and handler
// context (both the HandlerCtx and, if a session is set, the session value).
func newRequest(body string, hCtx *HandlerCtx) *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/api/v2/test", strings.NewReader(body))
	if hCtx == nil {
		return r
	}
	ctx := context.WithValue(r.Context(), HandlerContextKey, hCtx)
	if hCtx.ValidSession != nil {
		ctx = context.WithValue(ctx, SessionContextKey, *hCtx.ValidSession)
	}
	return r.WithContext(ctx)
}

// apiResponse mirrors the JSON envelope written by WriteSuccessResponse and
// WriteFailureResponse.
type apiResponse struct {
	OK   bool            `json:"ok"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// decodeResponse parses the recorded response body into the API envelope.
func decodeResponse(t *testing.T, rec *httptest.ResponseRecorder) apiResponse {
	t.Helper()
	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("cannot read response body: %v", err)
	}
	var resp apiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("cannot parse response body %q: %v", string(body), err)
	}
	return resp
}

// decodeData parses the `data` member of the response into target.
func decodeData(t *testing.T, resp apiResponse, target any) {
	t.Helper()
	if len(resp.Data) == 0 {
		t.Fatalf("response does not contain any data: %+v", resp)
	}
	if err := json.Unmarshal(resp.Data, target); err != nil {
		t.Fatalf("cannot parse response data %q: %v", string(resp.Data), err)
	}
}

// mustDecimal parses a decimal or fails the test.
func mustDecimal(t *testing.T, s string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(s)
	if err != nil {
		t.Fatalf("cannot parse decimal %q: %v", s, err)
	}
	return d
}

// decimalPtr returns a pointer to the parsed decimal.
func decimalPtr(t *testing.T, s string) *decimal.Decimal {
	t.Helper()
	d := mustDecimal(t, s)
	return &d
}

// sentMail is a message captured by the SMTP stub.
type sentMail struct {
	From string
	To   []string
	Body string
}

// startFakeSMTP starts a minimal SMTP server on localhost and returns the
// email configuration pointing at it plus an accessor for the messages that
// were delivered. The server does not advertise AUTH, so net/smtp skips
// authentication.
func startFakeSMTP(t *testing.T) (database.EmailConfiguration, func() []sentMail) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("cannot start the SMTP stub: %v", err)
	}
	t.Cleanup(func() { listener.Close() })

	var mu sync.Mutex
	var mails []sentMail

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func() {
				defer conn.Close()
				var mail sentMail
				var body strings.Builder
				inData := false
				w := bufio.NewWriter(conn)
				reader := bufio.NewReader(conn)
				fmt.Fprint(w, "220 localhost ESMTP stub\r\n")
				w.Flush()
				for {
					line, err := reader.ReadString('\n')
					if err != nil {
						return
					}
					line = strings.TrimRight(line, "\r\n")
					if inData {
						if line == "." {
							inData = false
							mail.Body = body.String()
							mu.Lock()
							mails = append(mails, mail)
							mu.Unlock()
							fmt.Fprint(w, "250 Ok\r\n")
							w.Flush()
							continue
						}
						body.WriteString(line + "\n")
						continue
					}
					switch {
					case strings.HasPrefix(line, "EHLO"), strings.HasPrefix(line, "HELO"):
						// AUTH PLAIN over a plain connection is accepted by
						// net/smtp because the host is a loopback address.
						fmt.Fprint(w, "250-localhost\r\n250 AUTH PLAIN\r\n")
					case strings.HasPrefix(line, "AUTH"):
						fmt.Fprint(w, "235 Authentication successful\r\n")
					case strings.HasPrefix(line, "MAIL FROM:"):
						mail.From = strings.Trim(strings.TrimPrefix(line, "MAIL FROM:"), "<> ")
						fmt.Fprint(w, "250 Ok\r\n")
					case strings.HasPrefix(line, "RCPT TO:"):
						mail.To = append(mail.To, strings.Trim(strings.TrimPrefix(line, "RCPT TO:"), "<> "))
						fmt.Fprint(w, "250 Ok\r\n")
					case strings.HasPrefix(line, "DATA"):
						inData = true
						fmt.Fprint(w, "354 End data with <CR><LF>.<CR><LF>\r\n")
					case strings.HasPrefix(line, "QUIT"):
						fmt.Fprint(w, "221 Bye\r\n")
						w.Flush()
						return
					default:
						fmt.Fprint(w, "250 Ok\r\n")
					}
					w.Flush()
				}
			}()
		}
	}()

	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("cannot determine the SMTP stub port: %v", err)
	}
	config := database.EmailConfiguration{
		Sender:     "club@example.com",
		Server:     "127.0.0.1",
		ServerPort: listener.Addr().String(),
		Port:       port,
		Username:   "club",
		Password:   "secret",
	}
	return config, func() []sentMail {
		mu.Lock()
		defer mu.Unlock()
		out := make([]sentMail, len(mails))
		copy(out, mails)
		return out
	}
}

// validSession returns a browser session that is valid for the next hour.
func validSession(userID uint, role database.UserRoleType) *database.BrowserSession {
	return &database.BrowserSession{
		SessionSecret: "test-session-secret",
		ValidUntil:    time.Now().Add(time.Hour),
		MaxValidUntil: time.Now().Add(24 * time.Hour),
		LastActivity:  time.Now(),
		UserID:        userID,
		Username:      "tester",
		UserRoleID:    role,
	}
}
