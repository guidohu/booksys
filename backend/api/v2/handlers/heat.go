package handlers

import (
	"fmt"
	"math"
	"net/http"
	"server/database"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
)

type AddHeatsRequest struct {
	Heats []AddHeatRequest `json:"heats" validation:"dive"`
}

type AddHeatRequest struct {
	UID             string `json:"uid" validation:"required"`
	SessionID       uint   `json:"session_id" validation:"required"`
	DurationSeconds uint64 `json:"duration_s" validation:"required"`
	UserID          uint   `json:"user_id" validation:"required"`
	Comment         string `json:"comment"`
}

type AddHeatsRequestResponse map[string]Status

type DeleteHeatRequest struct {
	HeatID uint `json:"heat_id" validation:"required"`
}

type ChangeHeatRequest struct {
	HeatID          uint   `json:"heat_id" validation:"required"`
	UserID          uint   `json:"user_id" validation:"required"`
	DurationSeconds uint64 `json:"duration_s" validation:"required"`
	Comment         string `json:"comment"`
}

func (h *Handler) AddHeats(w http.ResponseWriter, r *http.Request, req AddHeatsRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database

	resp := AddHeatsRequestResponse{}
	for _, heat := range []AddHeatRequest(req.Heats) {
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
	user, err := dbh.GetUserById(heat.UserID)
	if err != nil {
		slog.Warn("Cannot get user information for", slog.Uint64("user", uint64(heat.UserID)), slog.String("error", err.Error()))
		return fmt.Errorf("cannot get user")
	}

	// Check that session exists.
	if heat.SessionID != 0 {
		session, err := dbh.GetSession(heat.SessionID)
		if err != nil || session.ID == 0 {
			slog.Warn("Cannot find session", slog.Uint64("session", uint64(heat.SessionID)), slog.String("error", err.Error()))
			return fmt.Errorf("cannot get pricing information for user")
		}
	}

	// Get pricing
	pricing, err := dbh.GetUserStatusToPricingsMap()
	if err != nil {
		slog.Warn("Cannot get pricing information for user", slog.Uint64("user", uint64(heat.UserID)), slog.String("error", err.Error()))
		return fmt.Errorf("cannot get pricing information for user")
	}
	p, ok := pricing[user.UserStatusID]
	if !ok {
		slog.Warn("Cannot get pricing information for user, user status ID not found", slog.Uint64("user", uint64(heat.UserID)))
		return fmt.Errorf("cannot get pricing information for user")
	}
	price := p.PricePerMinute
	if price.IsZero() {
		slog.Warn("Cannot get a pricing for user, price is zero", slog.Uint64("user", uint64(heat.UserID)))
		return fmt.Errorf("cannot get a pricing for this user")
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
		slog.Warn("Cannot add heat", slog.String("error", err.Error()))
		return fmt.Errorf("cannot add the heat to the database")
	}
	return nil
}

func (h *Handler) AddHeat(w http.ResponseWriter, r *http.Request) {
	req := &AddHeatRequest{}
	err := ReadBodyAndValidate(r, req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid request.", w)
		return
	}
	hCtx := GetHandlerContext(r)
	dbh := hCtx.Database
	err = h.addHeat(dbh, *req)
	if err != nil {
		slog.Warn("Could not add heat", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}
	WriteSuccessResponse("heat added", nil, w)
}

func (h *Handler) DeleteHeat(w http.ResponseWriter, r *http.Request, req DeleteHeatRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	err := dbh.DeleteHeat(req.HeatID)
	if err != nil {
		slog.Warn("Could not delete heat", slog.String("error", err.Error()))
		WriteFailureResponse("Could not delete heat, heat not valid.", w)
		return
	}
	WriteSuccessResponse("heat deleted", nil, w)
}

func (h *Handler) ChangeHeat(w http.ResponseWriter, r *http.Request, req ChangeHeatRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	// Get user
	user, err := dbh.GetUserById(req.UserID)
	if err != nil {
		slog.Warn("Cannot get user information for", slog.Uint64("user", uint64(req.UserID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot find user.", w)
		return
	}

	// Get pricing
	pricing, err := dbh.GetUserStatusToPricingsMap()
	if err != nil {
		slog.Warn("Cannot get pricing information for user", slog.Uint64("user", uint64(user.ID)), slog.String("error", err.Error()))
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
	if err != nil || heat.ID == 0 {
		slog.Warn("Cannot find existing heat", slog.Uint64("heat", uint64(req.HeatID)), slog.String("error", err.Error()))
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
		slog.Warn("Cannot change heat", slog.String("error", err.Error()))
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
