package handlers

import (
	"log/slog"
	"net/http"
	"os"

	"server/version"
)

// HealthStatusResponse reports how far the application is set up and whether
// its dependencies are reachable.
type HealthStatusResponse struct {
	// True if no config file provided or if config file exists.
	ConfigFile bool `json:"configFile"`
	// True if database configuration exists.
	ConfigDB bool `json:"configDb"`
	// True if we can connect to the database.
	DBReachable bool `json:"dbReachable"`
	// True if users are present in the users table.
	UsersExist bool `json:"usersExist"`
	// Version is the backend version.
	Version string `json:"version"`
	// Commit is the git commit.
	Commit string `json:"commit"`
	// BuildDate is the date of this build.
	Built string `json:"buildDate"`
	// Environment is the environment this application is running under (e.g. prod, dev)
	Environment string `json:"environment"`
}

// HealthStatus serves the setup and dependency health of the application.
func (h *Handler) HealthStatus(w http.ResponseWriter, r *http.Request) {
	environment, _ := h.config.GetString("environment")
	healthStatus := HealthStatusResponse{
		Version:     version.Release,
		Commit:      version.Commit,
		Built:       version.BuildDate,
		Environment: environment,
	}

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
		dbh, done := h.Database.GetHandler()
		defer done()
		if dbh == nil {
			slog.Warn("No database connection is available.")
			WriteFailureResponse("Operation cannot be performed. Database connection is not established properly.", w)
			return
		}
		if !dbh.IsConfigured() {
			slog.Warn("Database handler is not configured.")
		} else if err := dbh.Ping(); err != nil {
			slog.Warn("Database handler is configured but cannot connect", slog.Any("error", err))
		} else {
			healthStatus.DBReachable = true
			usersExist, err := dbh.UsersExist()
			if err != nil {
				slog.Error("Cannot check for users", slog.Any("error", err))
				http.Error(w, "health status response", http.StatusInternalServerError)
				return
			}
			healthStatus.UsersExist = usersExist
		}
	}

	WriteSuccessResponse("health report", healthStatus, w)
}
