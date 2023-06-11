package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type pingData struct {
	Message string `json:"message"`
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Status: Status{
			OK:  true,
			Msg: "ping received",
		},
		Data: pingData{
			Message: "I'm a ping response from the server",
		},
	}

	j, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "create pong response", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, bytes.NewReader(j))
}
