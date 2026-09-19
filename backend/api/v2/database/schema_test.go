package database

import (
	"testing"
	"time"
)

// TestBrowserSessionValid pins down that both timeouts of a session are
// enforced here, in the check every caller makes. The absolute lifetime used
// to be compared only in the handler that clients call to find out whether
// they are still logged in, which left it up to the client whether it applied
// at all.
func TestBrowserSessionValid(t *testing.T) {
	tests := []struct {
		name    string
		session BrowserSession
		want    bool
	}{
		{
			name: "inside both timeouts",
			session: BrowserSession{
				SessionSecret: "secret",
				ValidUntil:    time.Now().Add(time.Hour),
				MaxValidUntil: time.Now().Add(24 * time.Hour),
			},
			want: true,
		},
		{
			name: "inactivity window elapsed",
			session: BrowserSession{
				SessionSecret: "secret",
				ValidUntil:    time.Now().Add(-time.Minute),
				MaxValidUntil: time.Now().Add(24 * time.Hour),
			},
			want: false,
		},
		{
			name: "absolute lifetime reached while still active",
			session: BrowserSession{
				SessionSecret: "secret",
				ValidUntil:    time.Now().Add(time.Hour),
				MaxValidUntil: time.Now().Add(-time.Minute),
			},
			want: false,
		},
		{
			// Sessions stored before the absolute timeout existed carry no
			// lifetime, so there is nothing to bound them by.
			name: "no absolute lifetime stored",
			session: BrowserSession{
				SessionSecret: "secret",
				ValidUntil:    time.Now().Add(time.Hour),
			},
			want: false,
		},
		{
			name:    "empty session",
			session: BrowserSession{},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.session.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
