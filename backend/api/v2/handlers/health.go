package handlers

import (
	"net/http"
	"os"
	"server/config"

	"golang.org/x/exp/slog"

	"github.com/spf13/viper"
)

type HealthStatusResponse struct {
	ConfigFile  bool `json:"configFile"`
	ConfigDB    bool `json:"configDb"`
	DBReachable bool `json:"dbReachable"`
	UsersExist  bool `json:"usersExist"`
	MyNautique  bool `json:"myNautiqueConfigured"`
}

func (h *Handler) HealthStatus(w http.ResponseWriter, r *http.Request) {
	healthStatus := HealthStatusResponse{}

	// check if configuration file exists
	_, err := os.Stat(viper.GetString("configfile"))
	if err == nil {
		healthStatus.ConfigFile = true
	}

	healthStatus.ConfigDB = config.IsDBConfigured()
	healthStatus.DBReachable = config.IsDBConfigured() && h.GetDB().Ping() == nil
	if healthStatus.DBReachable {
		healthStatus.UsersExist, err = h.GetDB().UsersExist()
		if err != nil {
			slog.Error("Cannot check for users:", err)
			http.Error(w, "health status response", http.StatusInternalServerError)
			return
		}
	}
	healthStatus.MyNautique = viper.IsSet("mynautique.enabled")

	WriteSuccessResponse("health report", healthStatus, w)
}
