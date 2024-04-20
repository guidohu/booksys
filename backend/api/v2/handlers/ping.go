package handlers

import (
	"net/http"
)

type PingData struct {
	Message string `json:"message"`
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	resp := &PingData{
		Message: "I'm a ping response from the server",
	}
	WriteSuccessResponse("ping message", resp, w)
}
