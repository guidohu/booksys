package handlers

import (
	"net/http"
	"server/mynautique"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
)

type Config struct {
	Enabled      bool
	APIKey       string
	User         string
	Password     string
	BoatID       int64
	FuelCapacity decimal.Decimal
}

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

func (h *Handler) GetBoatTelemetry(w http.ResponseWriter, r *http.Request, req GetBoatInfoRequest, hCtx *HandlerCtx) {
	client := h.GetMyNautiqueClient()
	config := Config{}
	config.Enabled = h.config.GetBool("mynautique.enabled")
	if !config.Enabled {
		slog.Warn("mynautique is not configured but GetBoatInfo was called")
		WriteFailureResponse("myNautique is not configured", w)
		return
	}
	config.APIKey, _ = h.config.GetString("mynautique.api.key")
	if config.APIKey == "" {
		slog.Warn("mynautique API key is missing but GetBoatInfo was called")
		WriteFailureResponse("myNautique API key missing", w)
		return
	}
	config.FuelCapacity = decimal.NewFromInt(h.config.GetInt64("mynautique.fuel.capacity"))
	config.User, _ = h.config.GetString("mynautique.user")
	if config.User == "" {
		slog.Warn("mynautique user is missing but GetBoatInfo was called")
		WriteFailureResponse("myNautique user missing", w)
		return
	}
	config.Password, _ = h.config.GetString("mynautique.password")
	if config.User == "" {
		slog.Warn("mynautique password is missing but GetBoatInfo was called")
		WriteFailureResponse("myNautique password missing", w)
		return
	}

	if client == nil {
		slog.Info("Creating new myNautique Client")
		client = mynautique.NewMyNautiqueClient(&mynautique.Options{
			User:       config.User,
			Password:   config.Password,
			AuthAPIKey: config.APIKey,
		})
		h.SetMyNautiqueClient(*client)
	}

	t, err := client.GetBoatTelemetry(req.BoatID)
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
