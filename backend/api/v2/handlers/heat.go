package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
	"server/database"
)

// AddHeatsRequest is the request body of AddHeats, which serves
// /api/v2/heats/create.
type AddHeatsRequest struct {
	Heats []AddHeatRequest `json:"heats" validation:"dive"`
}

// AddHeatRequest is a single ride in an AddHeatsRequest.
type AddHeatRequest struct {
	UID             string `json:"uid" validation:"required"`
	SessionID       uint   `json:"session_id" validation:"required"`
	DurationSeconds uint64 `json:"duration_s" validation:"required"`
	UserID          uint   `json:"user_id" validation:"required"`
	Comment         string `json:"comment"`
}

// AddHeatsRequestResponse reports, per submitted ride, whether it could be
// recorded. It is keyed by the UID the client sent with that ride.
type AddHeatsRequestResponse map[string]Status

// DeleteHeatRequest is the request body of DeleteHeat, which serves
// /api/v2/heat/delete.
type DeleteHeatRequest struct {
	HeatID uint `json:"heat_id" validation:"required"`
}

// ChangeHeatRequest is the request body of ChangeHeat, which serves
// /api/v2/heat/change.
type ChangeHeatRequest struct {
	HeatID          uint   `json:"heat_id" validation:"required"`
	UserID          uint   `json:"user_id" validation:"required"`
	DurationSeconds uint64 `json:"duration_s" validation:"required"`
	Comment         string `json:"comment"`
}

// AddHeats records the rides of a session.
//
// It serves /api/v2/heats/create and is open to administrators.
func (h *Handler) AddHeats(w http.ResponseWriter, r *http.Request, req AddHeatsRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database

	resp := AddHeatsRequestResponse{}
	for _, heat := range req.Heats {
		err := h.addHeat(dbh, heat)
		if err != nil {
			resp[heat.UID] = Status{
				OK:  false,
				Msg: err.Error(),
			}
			continue
		}
		resp[heat.UID] = Status{
			OK:  true,
			Msg: "heat added",
		}
	}
	WriteSuccessResponse("added heats, check individual responses", &resp, w)
}

func (h *Handler) addHeat(dbh database.Database, heat AddHeatRequest) error {
	// Get user
	user, err := dbh.GetUserByID(heat.UserID)
	if err != nil {
		slog.Warn("Cannot get user information for", slog.Uint64("user", uint64(heat.UserID)), slog.Any("error", err))
		return errors.New("cannot get user")
	}

	// Check that session exists.
	if heat.SessionID != 0 {
		session, err := dbh.GetSession(heat.SessionID)
		if err != nil {
			slog.Warn("Failed to retrieve session", slog.Uint64("session", uint64(heat.SessionID)), slog.Any("error", err))
			return fmt.Errorf("cannot get pricing information for user: %w", err)
		}
		if session.ID == 0 {
			slog.Warn("Session not found", slog.Uint64("session", uint64(heat.SessionID)))
			return errors.New("cannot get pricing information for user: session not found")
		}
	}

	// Get pricing
	pricing, err := dbh.GetUserStatusToPricingsMap()
	if err != nil {
		slog.Warn("Cannot get pricing information for user", slog.Uint64("user", uint64(heat.UserID)), slog.Any("error", err))
		return errors.New("cannot get pricing information for user")
	}
	p, ok := pricing[user.UserStatusID]
	if !ok {
		slog.Warn("Cannot get pricing information for user, user status ID not found", slog.Uint64("user", uint64(heat.UserID)))
		return errors.New("cannot get pricing information for user")
	}
	price := p.PricePerMinute
	if price.IsZero() {
		slog.Warn("Cannot get a pricing for user, price is zero", slog.Uint64("user", uint64(heat.UserID)))
		return errors.New("cannot get a pricing for this user")
	}

	newHeat := database.Heat{
		ID:              0,
		UserID:          user.ID,
		SessionID:       heat.SessionID,
		Timestamp:       time.Now(),
		DurationSeconds: int32(heat.DurationSeconds),
		Cost:            h.calculateHeatCost(heat.DurationSeconds, price),
		Comment:         heat.Comment,
	}
	err = dbh.AddHeat(&newHeat)
	if err != nil {
		slog.Warn("Cannot add heat", slog.Any("error", err))
		return errors.New("cannot add the heat to the database")
	}
	return nil
}

// AddHeat records a single ride. It is not routed, AddHeats is the endpoint
// the API exposes.
func (h *Handler) AddHeat(w http.ResponseWriter, r *http.Request) {
	req := &AddHeatRequest{}
	err := ReadBodyAndValidate(r, req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.Any("error", err))
		WriteFailureResponse("Invalid request.", w)
		return
	}
	hCtx := GetHandlerContext(r)
	dbh := hCtx.Database
	err = h.addHeat(dbh, *req)
	if err != nil {
		slog.Warn("Could not add heat", slog.Any("error", err))
		WriteFailureResponse(err.Error(), w)
		return
	}
	WriteSuccessResponse("heat added", nil, w)
}

// DeleteHeat removes a ride.
//
// It serves /api/v2/heat/delete and is open to administrators.
func (h *Handler) DeleteHeat(w http.ResponseWriter, r *http.Request, req DeleteHeatRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	err := dbh.DeleteHeat(req.HeatID)
	if err != nil {
		slog.Warn("Could not delete heat", slog.Any("error", err))
		WriteFailureResponse("Could not delete heat, heat not valid.", w)
		return
	}
	WriteSuccessResponse("heat deleted", nil, w)
}

// ChangeHeat updates a ride and recalculates its cost.
//
// It serves /api/v2/heat/change and is open to administrators.
func (h *Handler) ChangeHeat(w http.ResponseWriter, r *http.Request, req ChangeHeatRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	// Get user
	user, err := dbh.GetUserByID(req.UserID)
	if err != nil {
		slog.Warn("Cannot get user information for", slog.Uint64("user", uint64(req.UserID)), slog.Any("error", err))
		WriteFailureResponse("Cannot find user.", w)
		return
	}

	// Get pricing
	pricing, err := dbh.GetUserStatusToPricingsMap()
	if err != nil {
		slog.Warn("Cannot get pricing information for user", slog.Uint64("user", uint64(user.ID)), slog.Any("error", err))
		WriteFailureResponse("Cannot get pricing information.", w)
		return
	}
	p, ok := pricing[user.UserStatusID]
	if !ok {
		slog.Warn("Cannot get pricing information for user, user status ID not found", slog.Uint64("user", uint64(user.ID)))
		WriteFailureResponse("Cannot get pricing information.", w)
		return
	}
	price := p.PricePerMinute
	if price.IsZero() {
		slog.Warn("Cannot get a pricing for user, price is zero", slog.Uint64("user", uint64(user.ID)))
		WriteFailureResponse("Cannot get pricing information.", w)
		return
	}

	// get existing heat
	heat, err := dbh.GetHeat(req.HeatID)
	if err != nil {
		slog.Warn("Failed to retrieve existing heat", slog.Uint64("heat", uint64(req.HeatID)), slog.Any("error", err))
		WriteFailureResponse("Cannot find existing heat.", w)
		return
	}
	if heat.ID == 0 {
		slog.Warn("Existing heat not found", slog.Uint64("heat", uint64(req.HeatID)))
		WriteFailureResponse("Cannot find existing heat.", w)
		return
	}

	cost := h.calculateHeatCost(req.DurationSeconds, price)
	newHeat := database.Heat{
		ID:              heat.ID,
		UserID:          user.ID,
		SessionID:       heat.SessionID,
		Timestamp:       heat.Timestamp,
		DurationSeconds: int32(req.DurationSeconds),
		Cost:            cost,
		Comment:         req.Comment,
	}
	err = dbh.ChangeHeat(&newHeat)
	if err != nil {
		slog.Warn("Cannot change heat", slog.Any("error", err))
		WriteFailureResponse("Cannot update existing heat.", w)
		return
	}
	WriteSuccessResponse("heat updated", nil, w)
}

func (h *Handler) calculateHeatCost(durationSeconds uint64, price decimal.Decimal) decimal.Decimal {
	cost := price.InexactFloat64() * (float64(durationSeconds) / 60.0)
	costRound := decimal.NewFromFloat(math.Ceil(cost/0.05) * 0.05).RoundCash(5)
	return costRound
}
