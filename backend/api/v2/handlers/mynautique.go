package handlers

import (
	"net/http"
	"server/mynautique"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
)

type GetBoatInfoRequest struct {
	BoatID int64 `json:"boat_id" validate:"required"`
}

type GetBoatInfoResponse struct {
	Telemetry    Telemetry       `json:"telemetry"`
	FuelCapacity decimal.Decimal `json:"fuel_capacity"`
}

type Telemetry struct {
	FuelLevel   int64           `json:"fuel_level"`
	EngineHours decimal.Decimal `json:"engine_hours"`
}

func (h *Handler) GetBoatTelemetry(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	req := &GetBoatInfoRequest{}
	err := ReadBodyAndValidate(r, req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Please provide a boat ID.", w)
		return
	}

	config, err := h.GetDB().GetMyNautiqueConfiguration()
	if err != nil {
		slog.Warn("cannot lookup mynautique configuration", slog.String("error", err.Error()))
		WriteFailureResponse("myNautique is not properly configured", w)
		return
	}
	if !config.Enabled {
		slog.Warn("mynautique is not configured but GetBoatInfo was called")
		WriteFailureResponse("myNautique is not configured", w)
		return
	}

	client := h.GetMyNautiqueClient()
	// Prefer specific config over database config.
	// Note: We intend to implement an abstraction for this.
	apiKey := h.config.GetString("mynautique.api.key")
	if apiKey == "" {
		apiKey = config.APIKey
	}
	if client == nil {
		slog.Info("Creating new myNautique Client")
		client = mynautique.NewMyNautiqueClient(&mynautique.Options{
			User:       config.User,
			Password:   config.Password,
			AuthAPIKey: apiKey,
		})
		h.SetMyNautiqueClient(*client)
	}

	t, err := client.GetBoatTelemetry(config.BoatID)
	if err != nil {
		slog.Warn("cannot get boat telemetry from my nautique", slog.String("error", err.Error()))
		WriteFailureResponse("cannot get boat information from myNautique", w)
		return
	}

	resp := &GetBoatInfoResponse{
		Telemetry: Telemetry{
			FuelLevel:   t.FuelLevelLinc,
			EngineHours: t.EngineTotalHoursOfOperation,
		},
		FuelCapacity: config.FuelCapacity,
	}
	WriteSuccessResponse("boat telemetry", resp, w)

}

func (h *Handler) isMyNautiqueConfigured() bool {
	config, err := h.GetDB().GetPropertyValue("mynautique.enabled")
	if err != nil {
		slog.Error("Cannot determine if mynautique.enabled is set orr not")
		return false
	}
	return config.Value == "true"
}
