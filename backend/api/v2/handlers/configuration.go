package handlers

import (
	"fmt"
	"net"
	"net/http"
	"server/config"
	"server/database"

	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

type GetDBConfigResponse struct {
	IsConfigured bool   `json:"is_configured"`
	DBServer     string `json:"db_server"`
	DBName       string `json:"db_name"`
	DBUser       string `json:"db_user"`
	DBPassword   string `json:"db_password,omitempty"`
}

type SetupDBConfigRequest struct {
	DBServer   string `json:"db_server" validate:"required,hostname_port"`
	DBName     string `json:"db_name" validate:"required,alphanum"`
	DBUser     string `json:"db_user" validate:"required,alphanum"`
	DBPassword string `json:"db_password" validate:"required"`
}

type SetupMyNautiqueCredentialsRequest struct {
	Enabled  bool   `json:"mynautique_enabled" validate:"boolean"`
	User     string `json:"mynautique_user" validate:"omitempty,email"`
	Password string `json:"mynautique_password"`
}

func (h *Handler) SetupDBConfig(w http.ResponseWriter, r *http.Request) {
	if config.IsDBConfigured() {
		slog.Warn("SetupDB called for already setup DB")
		WriteFailureResponse("Invalid request", w)
		return
	}

	var req SetupDBConfigRequest
	err := ReadBodyAndValidate(r, &req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid request payload", w)
		return
	}

	host, port, err := net.SplitHostPort(req.DBServer)
	if err != nil {
		slog.Warn("Request payload is not valid, invalid DBServer", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid hostname (should be host:port)", w)
		return
	}

	// test db access
	db := &database.DBMysql{
		User:     req.DBUser,
		Password: req.DBPassword,
		Protocol: "tcp",
		Host:     host,
		Port:     port,
		DBName:   req.DBName,
	}
	err = db.Connect()
	defer db.Disconnect()
	if err != nil {
		slog.Warn(fmt.Sprintf("New database parameters are not valid. Error returned from Connet(): %s", err))
		WriteFailureResponse("Cannot connect to database.", w)
		return
	}

	// if database access was successful, we store the configurationn
	viper.Set("database.user", req.DBUser)
	viper.Set("database.password", req.DBPassword)
	viper.Set("database.protocol", "tcp")
	viper.Set("database.host", host)
	viper.Set("database.port", port)
	viper.Set("database.dbname", req.DBName)
	err = viper.WriteConfig()
	if err != nil {
		slog.Warn("Cannot store new configuration", slog.String("error", err.Error()))
		WriteFailureResponse("cannot write config", w)
		return
	}
	slog.Info("New database configuration has been written to", slog.String("configfile", viper.GetString("configfile")))

	// Postprocess config change
	h.ReconnectDB()

	WriteSuccessResponse("config written", nil, w)
}

func (h *Handler) GetDBConfig(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	resp := GetDBConfigResponse{
		IsConfigured: true,
		DBServer:     fmt.Sprintf("%s:%d", viper.GetString("database.host"), viper.GetUint("database.port")),
		DBName:       viper.GetString("database.dbname"),
		DBUser:       viper.GetString("database.user"),
		DBPassword:   "",
	}
	WriteSuccessResponse("success", resp, w)
}

func (h *Handler) SetupMyNautiqueCredentials(w http.ResponseWriter, r *http.Request) {
	// in case the configuration is already present, we do not
	// allow to edit it
	session := GetSessionFromContext(r)
	if viper.IsSet("mynautique.enabled") && AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	var req SetupMyNautiqueCredentialsRequest
	err := ReadBodyAndValidate(r, &req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid request payload", w)
		return
	}

	viper.Set("mynautique.enabled", req.Enabled)
	viper.Set("mynautique.user", req.User)
	viper.Set("mynautique.password", req.Password)

	err = viper.WriteConfig()
	if err != nil {
		slog.Warn("Cannot store new configuration", slog.String("error", err.Error()))
		WriteFailureResponse("cannot write config", w)
		return
	}
	slog.Info("New mynautique configuration has been written to", slog.String("configfile", viper.GetString("configfile")))

	// Postprocess config change
	h.ReconnectDB()

	WriteSuccessResponse("config written", nil, w)
}
