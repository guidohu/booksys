package email

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"server/database"
)

// delivery is one message as the stub server received it.
type delivery struct {
	AuthUser string
	AuthPass string
	From     string
	To       []string
	Body     string
}

// stubSMTP is an in-process SMTP server. It speaks just enough of the protocol
// for net/smtp to authenticate and deliver a message, and records what it was
// sent. It listens on 127.0.0.1 because net/smtp refuses to send PLAIN
// credentials over an unencrypted connection to anything but localhost.
type stubSMTP struct {
	Host string
	Port string

	mu         sync.Mutex
	deliveries []delivery
}

// newStubSMTP starts a stub server that is shut down when the test ends.
func newStubSMTP(t *testing.T) *stubSMTP {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("cannot listen: %v", err)
	}
	t.Cleanup(func() { listener.Close() })

	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("cannot split listener address: %v", err)
	}
	s := &stubSMTP{Host: host, Port: port}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				// The listener was closed, the test is over.
				return
			}
			go s.serve(conn)
		}
	}()
	return s
}

// Deliveries returns the messages received so far.
func (s *stubSMTP) Deliveries() []delivery {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]delivery(nil), s.deliveries...)
}

func (s *stubSMTP) serve(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	reply := func(lines ...string) error {
		for _, line := range lines {
			if _, err := fmt.Fprintf(conn, "%s\r\n", line); err != nil {
				return err
			}
		}
		return nil
	}

	if reply("220 stub.test ESMTP") != nil {
		return
	}
	var d delivery
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		verb, rest, _ := strings.Cut(strings.TrimRight(line, "\r\n"), " ")
		switch strings.ToUpper(verb) {
		case "EHLO":
			err = reply("250-stub.test greets you", "250 AUTH PLAIN")
		case "HELO":
			err = reply("250 stub.test greets you")
		case "AUTH":
			mech, credentials, _ := strings.Cut(rest, " ")
			if !strings.EqualFold(mech, "PLAIN") {
				err = reply("504 unsupported mechanism")
				break
			}
			// PLAIN credentials are "\x00username\x00password", base64 encoded.
			raw, decodeErr := base64.StdEncoding.DecodeString(credentials)
			if decodeErr != nil {
				err = reply("535 malformed credentials")
				break
			}
			if parts := strings.Split(string(raw), "\x00"); len(parts) == 3 {
				d.AuthUser, d.AuthPass = parts[1], parts[2]
			}
			err = reply("235 2.7.0 accepted")
		case "MAIL":
			d.From = addressIn(rest)
			err = reply("250 2.1.0 ok")
		case "RCPT":
			d.To = append(d.To, addressIn(rest))
			err = reply("250 2.1.5 ok")
		case "DATA":
			if err = reply("354 end data with <CR><LF>.<CR><LF>"); err != nil {
				break
			}
			body, readErr := readDATA(reader)
			if readErr != nil {
				return
			}
			d.Body = body
			s.mu.Lock()
			s.deliveries = append(s.deliveries, d)
			s.mu.Unlock()
			d = delivery{AuthUser: d.AuthUser, AuthPass: d.AuthPass}
			err = reply("250 2.0.0 ok")
		case "QUIT":
			reply("221 2.0.0 bye")
			return
		default:
			err = reply("250 2.0.0 ok")
		}
		if err != nil {
			return
		}
	}
}

// readDATA reads a message body up to the terminating dot, undoing the dot
// stuffing and the CRLF line endings the client applied.
func readDATA(reader *bufio.Reader) (string, error) {
	var body strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "." {
			return body.String(), nil
		}
		body.WriteString(strings.TrimPrefix(line, "."))
		body.WriteString("\n")
	}
}

// addressIn extracts the address from a "FROM:<addr>" or "TO:<addr>" argument.
func addressIn(s string) string {
	_, rest, found := strings.Cut(s, "<")
	if !found {
		return ""
	}
	address, _, _ := strings.Cut(rest, ">")
	return address
}

// clientFor returns a Client that talks to the stub server.
func clientFor(s *stubSMTP) *Client {
	return NewClient(database.EmailConfiguration{
		Server:   s.Host,
		Port:     s.Port,
		Username: "postmaster@example.com",
		Password: "s3cret",
	})
}

// onlyDelivery returns the single message the stub received, failing the test
// if it did not receive exactly one.
func onlyDelivery(t *testing.T, s *stubSMTP) delivery {
	t.Helper()
	deliveries := s.Deliveries()
	if len(deliveries) != 1 {
		t.Fatalf("got %d deliveries, want 1", len(deliveries))
	}
	return deliveries[0]
}

// wantContains fails the test unless the body contains every given substring.
func wantContains(t *testing.T, body string, substrings ...string) {
	t.Helper()
	for _, want := range substrings {
		if !strings.Contains(body, want) {
			t.Errorf("message body does not contain %q, body was:\n%s", want, body)
		}
	}
}

// testUser is the recipient used across the tests.
var testUser = database.User{
	FirstName: "Ada",
	LastName:  "Lovelace",
	Email:     "ada@example.com",
}

// testSession starts and ends at fixed times so that the rendered dates are
// stable.
var testSession = database.Session{
	StartTime: time.Date(2026, time.July, 14, 9, 30, 0, 0, time.UTC),
	EndTime:   time.Date(2026, time.July, 14, 11, 0, 0, 0, time.UTC),
}

func TestNewClient(t *testing.T) {
	client := NewClient(database.EmailConfiguration{
		Sender:     "noreply@example.com",
		Server:     "smtp.example.com",
		Port:       "587",
		ServerPort: "smtp.example.com:587",
		Username:   "postmaster@example.com",
		Password:   "s3cret",
	})

	// Host comes from Server, not from the combined ServerPort field.
	if client.Host != "smtp.example.com" {
		t.Errorf("Host = %q, want smtp.example.com", client.Host)
	}
	if client.Port != "587" {
		t.Errorf("Port = %q, want 587", client.Port)
	}
	if client.Username != "postmaster@example.com" {
		t.Errorf("Username = %q, want postmaster@example.com", client.Username)
	}
	if client.Password != "s3cret" {
		t.Errorf("Password = %q, want s3cret", client.Password)
	}
}

func TestSend(t *testing.T) {
	t.Run("delivers the message to the configured server", func(t *testing.T) {
		server := newStubSMTP(t)
		client := clientFor(server)

		err := client.Send(&Email{
			Recipient: "ada@example.com",
			Sender:    "noreply@example.com",
			Subject:   "Hello",
			Message:   "To: ada@example.com\nSubject: Hello\n\nBody text.",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got := onlyDelivery(t, server)
		if got.From != "noreply@example.com" {
			t.Errorf("MAIL FROM = %q, want noreply@example.com", got.From)
		}
		if len(got.To) != 1 || got.To[0] != "ada@example.com" {
			t.Errorf("RCPT TO = %v, want [ada@example.com]", got.To)
		}
		wantContains(t, got.Body, "Subject: Hello", "Body text.")
	})

	t.Run("authenticates with the configured credentials", func(t *testing.T) {
		server := newStubSMTP(t)
		client := clientFor(server)

		if err := client.Send(&Email{
			Recipient: "ada@example.com",
			Sender:    "noreply@example.com",
			Message:   "hi",
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got := onlyDelivery(t, server)
		if got.AuthUser != "postmaster@example.com" {
			t.Errorf("auth user = %q, want postmaster@example.com", got.AuthUser)
		}
		if got.AuthPass != "s3cret" {
			t.Errorf("auth password = %q, want s3cret", got.AuthPass)
		}
	})

	t.Run("reports an unreachable server", func(t *testing.T) {
		// Port 0 cannot be dialed, so this fails without touching the network.
		client := NewClient(database.EmailConfiguration{
			Server: "127.0.0.1",
			Port:   "0",
		})

		err := client.Send(&Email{
			Recipient: "ada@example.com",
			Sender:    "noreply@example.com",
			Message:   "hi",
		})
		if err == nil {
			t.Error("expected an error for an unreachable server")
		}
	})
}

func TestSendTokenResetMessage(t *testing.T) {
	server := newStubSMTP(t)
	client := clientFor(server)

	err := client.SendTokenResetMessage(testUser, "424242", "noreply@example.com", "Password Reset Token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := onlyDelivery(t, server)
	if got.From != "noreply@example.com" {
		t.Errorf("MAIL FROM = %q, want noreply@example.com", got.From)
	}
	if len(got.To) != 1 || got.To[0] != "ada@example.com" {
		t.Errorf("RCPT TO = %v, want [ada@example.com]", got.To)
	}
	wantContains(t, got.Body,
		"To: ada@example.com",
		"From: noreply@example.com",
		"Subject: Password Reset Token",
		"Dear Ada Lovelace",
		"Token: 424242",
	)
}

func TestSendUserAddedToSessionMessage(t *testing.T) {
	server := newStubSMTP(t)
	client := clientFor(server)

	err := client.SendUserAddedToSessionMessage(
		testUser, testSession, "noreply@example.com", "Added to session", "https://booksys.example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := onlyDelivery(t, server)
	wantContains(t, got.Body,
		"To: ada@example.com",
		"Subject: Added to session",
		"Dear Ada Lovelace",
		"You have been added to the following session",
		"Date: Tuesday 14.07.2026",
		"09:30 to 11:00",
		"https://booksys.example.com",
	)
}

func TestSendUserRemovedFromSessionMessage(t *testing.T) {
	server := newStubSMTP(t)
	client := clientFor(server)

	err := client.SendUserRemovedFromSessionMessage(
		testUser, testSession, "noreply@example.com", "Session cancelled", "https://booksys.example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := onlyDelivery(t, server)
	wantContains(t, got.Body,
		"To: ada@example.com",
		"Subject: Session cancelled",
		"Dear Ada Lovelace",
		"The following session has been cancelled",
		"Date: Tuesday 14.07.2026",
		"09:30 to 11:00",
		"https://booksys.example.com",
	)
}
