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

	// Check if database configuration is present in the config.
	healthStatus.ConfigDB = h.config.IsDBConfigured()
	if healthStatus.ConfigDB {
		// Check whether the database handler is configured and
		// able to connect.
		dbh := h.GetDB()
		if !dbh.IsConfigured() {
			slog.Warn("Database handler is not configured.")
		} else if err := dbh.Ping(); err != nil {
			slog.Warn("Database handler is configured but cannot connect", slog.String("error", err.Error()))
		} else {
			healthStatus.DBReachable = true
			usersExist, err := dbh.UsersExist()
			if err != nil {
				slog.Error("Cannot check for users:", err)
				http.Error(w, "health status response", http.StatusInternalServerError)
				return
			}
			healthStatus.UsersExist = usersExist
		}
	}

	WriteSuccessResponse("health report", healthStatus, w)
}
