package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) SetupDB(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Status: Status{
			OK:  true,
			Msg: "login received",
		},
		// Data: loginData{
		// 	Message: "I'm a login received message from the server",
		// },
	}

	j, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "create login response", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, bytes.NewReader(j))
}
