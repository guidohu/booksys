package handlers

import (
	"net/http"
	"os"

	"golang.org/x/exp/slog"
)

type HealthStatusResponse struct {
	// True if no config file provided or if config file exists.
	ConfigFile bool `json:"configFile"`
	// True if database configuration exists.
	ConfigDB bool `json:"configDb"`
	// True if we can connect to the database.
	DBReachable bool `json:"dbReachable"`
	// True if users are present in the users table.
	UsersExist bool `json:"usersExist"`
	// True if myNautique is configured.
	MyNautique bool `json:"myNautiqueConfigured"`
}

func (h *Handler) HealthStatus(w http.ResponseWriter, r *http.Request) {
	healthStatus := HealthStatusResponse{}

	// Check if configuration file exists or if one
	// should exist.
	configFile, _ := h.config.GetString("config")
	if configFile == "" {
		// Config file status is 'true' in case
		// the config file is not set. We do not wait for or
		// miss a config file.
		healthStatus.ConfigFile = true
	} else {
		_, err := os.Stat(configFile)
		if err == nil {
			healthStatus.ConfigFile = true
		}
	}

	healthStatus.ConfigDB = h.config.IsDBConfigured()
	healthStatus.DBReachable = h.GetDB() != nil && h.config.IsDBConfigured() && h.GetDB().Ping() == nil
	if healthStatus.DBReachable {
		usersExist, err := h.GetDB().UsersExist()
		if err != nil {
			slog.Error("Cannot check for users:", err)
			http.Error(w, "health status response", http.StatusInternalServerError)
			return
		}
		healthStatus.UsersExist = usersExist
	}
	healthStatus.MyNautique = h.config.IsSet("mynautique.enabled")

	WriteSuccessResponse("health report", healthStatus, w)
}
