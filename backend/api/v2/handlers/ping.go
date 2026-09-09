package handlers

import (
	"net/http"
)

// PingData is the payload returned by Ping.
type PingData struct {
	Message string `json:"message"`
}

// Ping serves a liveness check that does not touch the database.
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	resp := &PingData{
		Message: "I'm a ping response from the server",
	}
	WriteSuccessResponse("ping message", resp, w)
}
