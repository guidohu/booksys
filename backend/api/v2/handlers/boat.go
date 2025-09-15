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
	AfterHours  decimal.Decimal `json:"engine_hours_after,omitempty" validate:"omitempty,numeric"`
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

type GetFuelEntriesResponse struct {
	ID                   uint             `json:"id"`
	AverageLitersPerHour decimal.Decimal  `json:"avg_liters_per_hour"`
	CostNet              *decimal.Decimal `json:"cost"`
	CostGros             *decimal.Decimal `json:"cost_brutto"`
	IsDiscounted         bool             `json:"is_discounted"`
	DiffHours            decimal.Decimal  `json:"diff_hours"`
	EngineHours          *decimal.Decimal `json:"engine_hours"`
	Liters               decimal.Decimal  `json:"liters"`
	Timestamp            int64            `json:"timestamp"`
	UserFirstName        string           `json:"user_first_name"`
	UserLastName         string           `json:"user_last_name"`
	UserID               uint             `json:"user_id"`
}

type AddFuelEntryRequest struct {
	Cost        *decimal.Decimal `json:"cost" validate:"required,numeric"`
	EngineHours *decimal.Decimal `json:"engine_hours" validate:"required,numeric"`
	Liters      *decimal.Decimal `json:"liters" validate:"required,numeric"`
	UserID      uint             `json:"user_id" validate:"required"`
}

var FuelEntryValidationErrors = map[string]string{
	"UserID":       "Please provide an user ID.",
	"Cost":         "Please provide the cost.",
	"EngineHours":  "Please provide the engine hours at time of refueling.",
	"Liters":       "Please provide the amount of fuel.",
	"ID":           "Please provide a fuel entry ID.",
	"CostNet":      "Please provide the net cost.",
	"CostGros":     "Please provide the gros cost.",
	"IsDiscounted": "Please provide whether this entry is discounted or not.",
}

type ChangeFuelEntryRequest struct {
	ID           uint             `json:"id"`
	CostNet      *decimal.Decimal `json:"cost" validate:"required,numeric"`
	CostGros     *decimal.Decimal `json:"cost_brutto,omitempty" validate:"omitempty,numeric"`
	EngineHours  *decimal.Decimal `json:"engine_hours" validate:"required,numeric"`
	Liters       *decimal.Decimal `json:"liters" validate:"required,numeric"`
	IsDiscounted bool             `json:"is_discounted"`
}

type RemoveFuelEntryRequest struct {
	ID uint `json:"id" validate:"required"`
}

type GetMaintenanceEntriesResponse struct {
	ID            uint            `json:"id"`
	Description   string          `json:"description"`
	EngineHours   decimal.Decimal `json:"engine_hours"`
	Timestamp     int64           `json:"timestamp"`
	UserFirstName string          `json:"user_first_name"`
	UserLastName  string          `json:"user_last_name"`
	UserID        uint            `json:"user_id"`
}

type AddMaintenanceEntryRequest struct {
	Description string          `json:"description" validate:"required"`
	EngineHours decimal.Decimal `json:"engine_hours" validate:"required,numeric"`
	UserID      uint            `json:"user_id" validate:"required"`
}

var AddMaintenanceEntryValidationErrors = map[string]string{
	"Description": "Please provide a description of the maintenance.",
	"EngineHours": "Please provide the engine hours at which you performed maintenance.",
	"UserID":      "Please provide a valid user ID.",
}

func (h *Handler) GetEngineHourLatest(w http.ResponseWriter, r *http.Request) {
	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database
	b, err := dbh.GetEngineHourLatest()
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
		UsageType:     b.Type.ID,
		UsageTypeName: database.DefaultSessionTypesMap[int(b.TypeID)].Name,
		UserFirstName: b.User.FirstName,
		UserLastName:  b.User.LastName,
		UserID:        b.UserID,
	}
	WriteSuccessResponse("boat engine hour entry", resp, w)
}

func (h *Handler) GetEngineHoursList(w http.ResponseWriter, r *http.Request) {
	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database
	engineHours, err := dbh.GetEngineHours()
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
			UsageType:     b.Type.ID,
			UsageTypeName: database.DefaultSessionTypesMap[int(b.TypeID)].Name,
			UserFirstName: b.User.FirstName,
			UserLastName:  b.User.LastName,
			UserID:        b.UserID,
		}
		resp = append(resp, entry)
	}
	WriteSuccessResponse("boat engine hours", resp, w)
}

func (h *Handler) UpdateEngineHours(w http.ResponseWriter, r *http.Request) {
	req := &UpdateEngineHoursRequest{}
	err := ReadBodyAndValidate(r, req, UpdateEngineHoursValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database

	// get latest entry
	latest, err := dbh.GetEngineHourLatest()
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
		TypeID:      req.UsageType,
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
		err = dbh.AddEngineHours(entry)
	} else if !req.BeforeHours.IsZero() && req.AfterHours.IsZero() && !latest.CheckedIn {
		// start new half entry
		entry.CheckedIn = true
		err = dbh.AddEngineHours(entry)
	} else if latest.AfterHours.IsZero() && !req.AfterHours.IsZero() && latest.CheckedIn {
		// finish half entry
		entry.ID = latest.ID
		entry.CheckedIn = false
		err = dbh.UpdateEngineHours(entry)
	} else {
		slog.Warn("Cannot add new engine hours, not a valid entry.", slog.Any("req", req))
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
	req := &UpdateEngineHoursEntryRequest{}
	err := ReadBodyAndValidate(r, req, UpdateEngineHoursEntryValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database
	entry, err := dbh.GetEngineHoursEntry(req.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Warn("Cannot get engine hour entry", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot find engine hour entry.", w)
		return
	}

	entry.TypeID = req.UsageType
	err = dbh.UpdateEngineHours(entry)
	if err != nil {
		slog.Warn("Cannot update engine hour entry", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot update engine hour entry.", w)
		return
	}
	WriteSuccessResponse("entry updated", nil, w)
}

func (h *Handler) GetFuelEntries(w http.ResponseWriter, r *http.Request) {
	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database
	fuelEntries, err := dbh.GetFuelEntries()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Warn("Cannot get fuel entries", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get fuel entries.", w)
		return
	}
	resp := []GetFuelEntriesResponse{}
	var lastEngineHour decimal.Decimal
	for i := len(fuelEntries) - 1; i >= 0; i-- {
		e := fuelEntries[i]
		entry := GetFuelEntriesResponse{
			ID:            e.ID,
			CostNet:       e.Cost,
			CostGros:      e.CostBrutto,
			IsDiscounted:  e.IsDiscounted,
			EngineHours:   e.EngineHours,
			Liters:        *e.Liters,
			Timestamp:     e.Timestamp.Unix(),
			UserFirstName: e.User.FirstName,
			UserLastName:  e.User.LastName,
			UserID:        e.UserID,
		}
		if !lastEngineHour.IsZero() {
			diffHours := e.EngineHours.Sub(lastEngineHour)
			if !diffHours.IsNegative() {
				liters, _ := e.Liters.Float64()
				hours, _ := diffHours.Float64()
				var avgFuelPerHour float64 = 0
				if hours > 0.001 {
					avgFuelPerHour = liters / hours
				}
				entry.DiffHours = diffHours
				entry.AverageLitersPerHour = decimal.NewFromFloat(avgFuelPerHour)
			}
		}
		lastEngineHour = *e.EngineHours
		resp = append([]GetFuelEntriesResponse{entry}, resp...)
	}
	WriteSuccessResponse("fuel entries", resp, w)
}

func (h *Handler) AddFuelEntry(w http.ResponseWriter, r *http.Request) {
	req := &AddFuelEntryRequest{}
	err := ReadBodyAndValidate(r, req, FuelEntryValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database
	// Check that user is an admin user
	isAdmin, _ := dbh.IsAdminUser(req.UserID)
	if !isAdmin {
		slog.Warn("Non admin user tried to add fuel entry.", slog.Uint64("userID", uint64(req.UserID)))
		WriteFailureResponse("Non admin user is not allowed to change engine hours.", w)
		return
	}

	billType, err := dbh.GetPropertyValue("fuel.payment.type")
	if err != nil {
		slog.Warn("Cannot determine whether fuel is billed or paid directly", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot determine whether fuel is billed or paid directly.", w)
		return
	}
	entry := database.BoatFuel{
		ID:                  0,
		Timestamp:           time.Now(),
		UserID:              req.UserID,
		EngineHours:         req.EngineHours,
		Liters:              req.Liters,
		Cost:                req.Cost,
		CostBrutto:          nil,
		ContributeToBalance: billType.Value != "billed",
	}
	err = dbh.AddFuelEntry(entry)
	if err != nil {
		slog.Warn("Cannot add fuel entry", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot add fuel entry.", w)
		return
	}
	WriteSuccessResponse("fuel entry added", nil, w)
}

func (h *Handler) ChangeFuelEntry(w http.ResponseWriter, r *http.Request) {
	req := &ChangeFuelEntryRequest{}
	err := ReadBodyAndValidate(r, req, FuelEntryValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database
	entry, err := dbh.GetFuelEntry(req.ID)
	if err != nil {
		slog.Warn("Cannot find fuel entry", slog.Uint64("id", uint64(req.ID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot find existing fuel entry.", w)
		return
	}

	newEntry := database.BoatFuel{
		ID:                  req.ID,
		Timestamp:           entry.Timestamp,
		UserID:              entry.UserID,
		EngineHours:         req.EngineHours,
		Liters:              req.Liters,
		Cost:                req.CostNet,
		CostBrutto:          req.CostGros,
		ContributeToBalance: entry.ContributeToBalance,
		IsDiscounted:        req.IsDiscounted,
	}
	err = dbh.ChangeFuelEntry(newEntry)
	if err != nil {
		slog.Warn("Cannot change fuel entry", slog.Uint64("id", uint64(req.ID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot change existing fuel entry.", w)
		return
	}
	WriteSuccessResponse("fuel entry saved", nil, w)
}

func (h *Handler) RemoveFuelEntry(w http.ResponseWriter, r *http.Request) {
	req := &RemoveFuelEntryRequest{}
	err := ReadBodyAndValidate(r, req, FuelEntryValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database
	entry, err := dbh.GetFuelEntry(req.ID)
	if err != nil {
		slog.Warn("Cannot find fuel entry", slog.Uint64("id", uint64(req.ID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot find existing fuel entry.", w)
		return
	}

	err = dbh.RemoveFuelEntry(entry.ID)
	if err != nil {
		slog.Warn("Cannot remove fuel entry", slog.Uint64("id", uint64(req.ID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot remove existing fuel entry.", w)
		return
	}
	WriteSuccessResponse("fuel entry saved", nil, w)
}

func (h *Handler) GetMaintenanceEntries(w http.ResponseWriter, r *http.Request) {
	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database
	logs, err := dbh.GetMaintenance()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Warn("Cannot get fuel entries", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get fuel entries.", w)
		return
	}

	resp := []GetMaintenanceEntriesResponse{}
	for _, l := range logs {
		resp = append(resp, GetMaintenanceEntriesResponse{
			ID:            l.ID,
			Description:   l.Description,
			EngineHours:   l.EngineHours,
			Timestamp:     l.Timestamp.Unix(),
			UserFirstName: l.User.FirstName,
			UserLastName:  l.User.LastName,
			UserID:        l.UserID,
		})
	}
	WriteSuccessResponse("maintenance entries", &resp, w)
}

func (h *Handler) AddMaintenanceEntry(w http.ResponseWriter, r *http.Request) {
	req := &AddMaintenanceEntryRequest{}
	err := ReadBodyAndValidate(r, req, AddMaintenanceEntryValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database
	if !dbh.UserExists(req.UserID) {
		slog.Warn("Request payload is not valid", slog.String("error", "user does not exist"))
		WriteFailureResponse("Please provide a valid user ID.", w)
		return
	}

	entry := database.BoatMaintenance{
		ID:          0,
		Timestamp:   time.Now(),
		UserID:      req.UserID,
		EngineHours: req.EngineHours,
		Description: req.Description,
	}
	err = dbh.AddMaintenanceEntry(entry)
	if err != nil {
		slog.Warn("Cannot add maintenance entry", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot add maintenance entry.", w)
		return
	}
	WriteSuccessResponse("maintenance entry created", nil, w)
}
