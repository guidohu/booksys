package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"server/config"

	"golang.org/x/exp/slog"

	"github.com/spf13/viper"
)

type HealthStatusData struct {
	ConfigFile  bool `json:"configFile"`
	ConfigDB    bool `json:"configDb"`
	DBReachable bool `json:"dbReachable"`
	UsersExist  bool `json:"usersExist"`
}

func (h *Handler) HealthStatus(w http.ResponseWriter, r *http.Request) {
	healthStatus := HealthStatusData{}
	resp := Response{
		Status: Status{
			OK:  true,
			Msg: "",
		},
		Data: healthStatus,
	}

	// check if configuration file exists
	_, err := os.Stat(viper.GetString("configfile"))
	if err == nil {
		healthStatus.ConfigFile = true
		resp.Data = healthStatus
	}

	// check if db configured
	if config.IsDBConfigured() {
		healthStatus.ConfigDB = true
	}

	// check if db reachable
	if h.db.Ping() == nil {
		healthStatus.DBReachable = true
		resp.Data = healthStatus
	}

	// check if admin user exists
	if h.db.Ping() == nil {
		u, err := h.db.UsersExist()
		if err != nil {
			slog.Error("Cannot check for users:", err)
			http.Error(w, "health status response", http.StatusInternalServerError)
			return
		}
		healthStatus.UsersExist = u
		resp.Data = healthStatus
	}

	j, err := json.Marshal(&resp)
	if err != nil {
		slog.Error("Cannot build HealthStatus response:", err)
		http.Error(w, "health status response", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, bytes.NewReader(j))
}
