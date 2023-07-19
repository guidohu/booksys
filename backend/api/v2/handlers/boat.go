package handlers

import (
	"errors"
	"net/http"
	"server/database"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

type GetEngineHourLatestResponse struct {
	ID            uint            `json:"id"`
	Timestamp     int64           `json:"timestamp"`
	BeforeHours   decimal.Decimal `json:"before_hours"`
	AfterHours    decimal.Decimal `json:"after_hours"`
	DeltaHours    decimal.Decimal `json:"delta_hours"`
	UsageType     uint            `json:"type"`
	UsageTypeName string          `json:"type_name"`
	UserFirstName string          `json:"user_first_name"`
	UserLastName  string          `json:"user_last_name"`
	UserID        uint            `json:"user_id"`
}

type GetEngineHoursResponse []GetEngineHourLatestResponse

type UpdateEngineHoursRequest struct {
	BeforeHours decimal.Decimal `json:"engine_hours_before" validate:"required,numeric"`
	AfterHours  decimal.Decimal `json:"engine_hours_after,omitempty" validate:"numeric"`
	UsageType   uint8           `json:"type" validate:"required,sessiontype"`
	UserID      uint            `json:"user_id" validate:"required,numeric"`
}

var UpdateEngineHoursValidationErrors = map[string]string{
	"BeforeHours": "Please provide engine hours before.",
	"AfterHours":  "Please provide engine hourse after.",
	"UsageType":   "Please provide a valid type.",
	"UserID":      "Please provide a valid user.",
}

type UpdateEngineHoursEntryRequest struct {
	ID        uint  `json:"id" validate:"required"`
	UsageType uint8 `json:"type" validate:"required,sessiontype"`
}

var UpdateEngineHoursEntryValidationErrors = map[string]string{
	"ID":        "Please provide an ID for the entry.",
	"UsageType": "Please provide the type of the session.",
}

func (h *Handler) GetEngineHourLatest(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	b, err := h.GetDB().GetEngineHourLatest()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Warn("Cannot get latest engine hour entry", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get latest engine hour entry.", w)
		return
	}
	resp := GetEngineHourLatestResponse{
		ID:            b.ID,
		Timestamp:     b.Timestamp.Unix(),
		BeforeHours:   b.BeforeHours,
		AfterHours:    b.AfterHours,
		DeltaHours:    b.DeltaHours,
		UsageType:     uint(b.Type),
		UsageTypeName: database.DefaultSessionTypesMap[int(b.Type)].Name,
		UserFirstName: b.User.FirstName,
		UserLastName:  b.User.LastName,
		UserID:        b.UserID,
	}
	WriteSuccessResponse("boat engine hour entry", resp, w)
}

func (h *Handler) GetEngineHours(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	engineHours, err := h.GetDB().GetEngineHours()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Warn("Cannot get engine hours", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get engine hours.", w)
		return
	}
	resp := GetEngineHoursResponse{}
	for _, b := range engineHours {
		entry := GetEngineHourLatestResponse{
			ID:            b.ID,
			Timestamp:     b.Timestamp.Unix(),
			BeforeHours:   b.BeforeHours,
			AfterHours:    b.AfterHours,
			DeltaHours:    b.DeltaHours,
			UsageType:     uint(b.Type),
			UsageTypeName: database.DefaultSessionTypesMap[int(b.Type)].Name,
			UserFirstName: b.User.FirstName,
			UserLastName:  b.User.LastName,
			UserID:        b.UserID,
		}
		resp = append(resp, entry)
	}
	WriteSuccessResponse("boat engine hours", resp, w)
}

func (h *Handler) UpdateEngineHours(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &UpdateEngineHoursRequest{}
	err := ReadBodyAndValidate(r, req, UpdateEngineHoursValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	// get latest entry
	latest, err := h.GetDB().GetEngineHourLatest()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Warn("Cannot get latest engine hour entry", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get latest engine hour entry.", w)
		return
	}

	// this is a new entry if the after hours is empty
	deltaHours := req.AfterHours.Sub(req.BeforeHours)
	if deltaHours.IsNegative() && !req.AfterHours.IsZero() {
		slog.Warn("Engine hours 'after' need to be bigger than 'before'", slog.Any("after", req.AfterHours), slog.Any("before", req.BeforeHours))
		WriteFailureResponse("Engine hours 'after' need to be bigger than 'before'", w)
		return
	}
	if deltaHours.IsNegative() {
		deltaHours = decimal.NewFromInt(0)
	}

	// new/updated entry
	entry := database.BoatEngineHour{
		ID:          0,
		Timestamp:   time.Now(),
		BeforeHours: req.BeforeHours,
		AfterHours:  req.AfterHours,
		DeltaHours:  deltaHours,
		Type:        req.UsageType,
		UserID:      req.UserID,
		Comment:     "",
	}
	// handle:
	// 1) first
	// 2) new half entry
	// 3) updated half entry
	if latest.BeforeHours.IsZero() && latest.AfterHours.IsZero() && !latest.CheckedIn {
		// this is the first entry
		entry.CheckedIn = true
		err = h.GetDB().AddEngineHours(entry)
	} else if !req.BeforeHours.IsZero() && req.AfterHours.IsZero() && !latest.CheckedIn {
		// start new half entry
		entry.CheckedIn = true
		err = h.GetDB().AddEngineHours(entry)
	} else if latest.AfterHours.IsZero() && !req.AfterHours.IsZero() && latest.CheckedIn {
		// finish half entry
		entry.ID = latest.ID
		entry.CheckedIn = false
		err = h.GetDB().UpdateEngineHours(entry)
	} else {
		slog.Warn("Cannot add new engine hours, not a valid entry.", slog.Any("latest", latest), slog.Any("req", req))
		WriteFailureResponse("Cannot add new engine hours, please check your input.", w)
		return
	}

	if err != nil {
		slog.Warn("Cannot add new engine hours", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot add new engine hours", w)
		return
	}
	WriteSuccessResponse("engine hours added", nil, w)
}

func (h *Handler) UpdateEngineHoursEntry(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &UpdateEngineHoursEntryRequest{}
	err := ReadBodyAndValidate(r, req, UpdateEngineHoursEntryValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	entry, err := h.GetDB().GetEngineHoursEntry(req.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Warn("Cannot get engine hour entry", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot find engine hour entry.", w)
		return
	}

	entry.Type = req.UsageType
	err = h.GetDB().UpdateEngineHours(entry)
	if err != nil {
		slog.Warn("Cannot update engine hour entry", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot update engine hour entry.", w)
		return
	}
	WriteSuccessResponse("entry updated", nil, w)
}
