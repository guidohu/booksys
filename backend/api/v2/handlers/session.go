package handlers

import (
	"fmt"
	"net/http"
	"server/database"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
)

type GetSessionRequest struct {
	SessionID uint64 `json:"id" validate:"required"`
}

type CreateSessionRequest struct {
	Title     string `json:"title" validate:"omitempty"`
	Comment   string `json:"comment" validate:"omitempty"`
	Start     int64  `json:"start" validate:"required,numeric"`
	End       int64  `json:"end" validate:"required"`
	MaxRiders int64  `json:"max_riders" validate:"required"`
	Type      int64  `json:"type" validate:"required,sessiontype"`
}

var SessionValidationErrors = map[string]string{
	"Start":     "Start time needs to be provided and needs to be a valid time value.",
	"End":       "End time needs to be provided and needs to be a valid time value.",
	"MaxRiders": "Please provide the maximum number of allowed riders in this session.",
	"Type":      "Please provide a valid type for this session.",
}

type CreateSessionResponse struct {
	SessionID uint `json:"session_id"`
}

type EditSessionRequest struct {
	SessionID uint   `json:"id"`
	Title     string `json:"title" validate:"omitempty"`
	Comment   string `json:"comment" validate:"omitempty"`
	Start     int64  `json:"start" validate:"required,numeric"`
	End       int64  `json:"end" validate:"required"`
	MaxRiders int64  `json:"max_riders" validate:"required"`
	Type      int64  `json:"type" validate:"required,sessiontype"`
}

type DeleteSessionRequest struct {
	SessionID uint `json:"session_id" validate:"required,numeric"`
}

var DeleteSessionValidationErrors = map[string]string{
	"SessionID": "Please provide a valid session id.",
}

type AddSessionUserRequest struct {
	SessionID uint   `json:"session_id" validate:"required"`
	UserIDs   []uint `json:"user_ids" validate:"required"`
}

type RemoveSessionUserRequest struct {
	SessionID uint `json:"session_id" validate:"required,numeric"`
	UserID    uint `json:"user_id" validate:"required,numeric"`
}

type GetSessionMetadataResponse struct {
	SunriseTime int64 `json:"sunrise"`
	SunsetTime  int64 `json:"sunset"`
}

type GetSessionHeatsRequest struct {
	SessionID uint `json:"session_id" validate:"required,numeric"`
}

type GetSessionHeatsResponse struct {
	HeatID    uint            `json:"heat_id"`
	UserID    uint            `json:"user_id"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	SessionID uint            `json:"session_id"`
	Timestamp int64           `json:"timestamp"`
	Duration  int64           `json:"duration_s"`
	Cost      decimal.Decimal `json:"cost"`
	Pricing   decimal.Decimal `json:"price_per_min"`
	Comment   string          `json:"comment"`
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &GetSessionRequest{}
	err := ReadBodyAndValidate(r, req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Request payload is not valid.", w)
		return
	}

	s, err := h.GetDB().GetSession(uint(req.SessionID))
	if err != nil {
		slog.Warn("Cannot get session", slog.Uint64("sessionID", req.SessionID), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get session from server.", w)
		return
	}

	riders, err := h.getRiders(s.ID)
	if err != nil {
		slog.Warn("Cannot get riders for session", slog.Uint64("sessionID", req.SessionID), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get rider information for session.", w)
		return
	}
	ridersShort := []RiderResponse{}
	for _, rider := range riders {
		ridersShort = append(ridersShort, RiderResponse{
			ID:        rider.ID,
			Name:      fmt.Sprintf("%s %s", rider.FirstName, rider.LastName),
			FirstName: rider.FirstName,
			LastName:  rider.LastName,
		})
	}

	resp := SessionResponse{
		ID:               s.ID,
		Start:            s.StartTime.Unix(),
		End:              s.EndTime.Unix(),
		Title:            s.Title,
		Comment:          s.Comment,
		FreeSpaces:       s.FreeSpaces,
		CreatorID:        s.CreatorID,
		CreatorFirstName: s.Creator.FirstName,
		CreatorLastName:  s.Creator.LastName,
		Duration:         int64(s.EndTime.Sub(s.StartTime).Seconds()),
		Riders:           ridersShort,
	}
	WriteSuccessResponse("session info", &resp, w)
}

func (h *Handler) GetSessionMetadata(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &GetSessionRequest{}
	err := ReadBodyAndValidate(r, req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Request payload is not valid.", w)
		return
	}

	s, err := h.GetDB().GetSession(uint(req.SessionID))
	if err != nil {
		slog.Warn("Cannot get session", slog.Uint64("sessionID", req.SessionID), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get session from server.", w)
		return
	}

	sunrise, sunset := h.getSunriseSunset(s.StartTime)
	resp := &GetSessionMetadataResponse{
		SunriseTime: sunrise.Unix(),
		SunsetTime:  sunset.Unix(),
	}
	WriteSuccessResponse("session metadata", resp, w)
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &CreateSessionRequest{}
	err := ReadBodyAndValidate(r, req, SessionValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	start := time.Unix(req.Start, 0)
	end := time.Unix(req.End, 0)
	collidingSessions, err := h.GetDB().GetSessionsBetween(start, end)
	if err != nil {
		slog.Warn("Cannot query existing sessions and check for collisions", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	if len(collidingSessions) > 0 {
		slog.Warn("New session creation collides with existing one. Cancel request.")
		WriteFailureResponse("New session collides with an existing one. Please change the time.", w)
		return
	}

	newSession := database.Session{
		StartTime:     start,
		EndTime:       end,
		Title:         req.Title,
		Comment:       req.Comment,
		SessionTypeID: uint(req.Type),
		FreeSpaces:    uint(req.MaxRiders),
		CreatorID:     session.UserID,
	}
	id, err := h.GetDB().CreateSession(newSession)
	if err != nil {
		slog.Warn("Cannot create new session", slog.String("error", err.Error()))
		WriteFailureResponse("New session could not be created", w)
		return
	}
	resp := CreateSessionResponse{
		SessionID: id,
	}
	WriteSuccessResponse("session created", &resp, w)
}

func (h *Handler) EditSession(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &EditSessionRequest{}
	err := ReadBodyAndValidate(r, req, SessionValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	oldSession, err := h.GetDB().GetSession(req.SessionID)
	if err != nil {
		slog.Warn("Cannot check existence of session", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot find the session you like to edit.", w)
		return
	}

	start := time.Unix(req.Start, 0)
	end := time.Unix(req.End, 0)
	collidingSessions, err := h.GetDB().GetSessionsBetween(start, end)
	if err != nil {
		slog.Warn("Cannot query existing sessions and check for collisions", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	for _, s := range collidingSessions {
		if s.ID != oldSession.ID {
			slog.Warn("Updated session collides with existing one. Cancel request.")
			WriteFailureResponse("Updated session collides with an existing one. Please change the time.", w)
			return
		}
	}

	editSession := database.Session{
		ID:            req.SessionID,
		StartTime:     start,
		EndTime:       end,
		Title:         req.Title,
		Comment:       req.Comment,
		SessionTypeID: uint(req.Type),
		FreeSpaces:    uint(req.MaxRiders),
		CreatorID:     oldSession.CreatorID,
	}
	err = h.GetDB().UpdateSession(editSession)
	if err != nil {
		slog.Warn("Cannot update session", slog.String("error", err.Error()))
		WriteFailureResponse("Session could not be updated.", w)
		return
	}
	WriteSuccessResponse("session updated", nil, w)
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &DeleteSessionRequest{}
	err := ReadBodyAndValidate(r, req, DeleteSessionValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	// check if session does not have heats
	heats, err := h.GetDB().GetHeatInSessionCount(req.SessionID)
	if err != nil {
		slog.Warn("Cannot get whether heats exist for session", slog.Int("session_id", int(req.SessionID)), slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}
	if heats > 0 {
		WriteFailureResponse("This session already has heats. It cannot be removed.", w)
		return
	}
	// TODO
	// inform riders

	// remove riders from session
	err = h.GetDB().DeleteUsersFromSession(req.SessionID)
	if err != nil {
		slog.Warn("Cannot remove users for session", slog.Int("session_id", int(req.SessionID)), slog.String("error", err.Error()))
		WriteFailureResponse("Users could not be removed from the session.", w)
		return
	}

	// delete session
	err = h.GetDB().DeleteSession(req.SessionID)
	if err != nil {
		slog.Warn("Cannot delete session", slog.Int("session_id", int(req.SessionID)), slog.String("error", err.Error()))
		WriteFailureResponse("Session could not be removed.", w)
		return
	}
	WriteSuccessResponse("session removed", nil, w)
}

func (h *Handler) AddUserToSession(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &AddSessionUserRequest{}
	err := ReadBodyAndValidate(r, req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid request.", w)
		return
	}

	// check if session exists
	s, err := h.GetDB().GetSession(req.SessionID)
	if err != nil {
		slog.Warn("Session cannot be found", slog.Uint64("sessionID", uint64(req.SessionID)), slog.String("error", err.Error()))
		WriteFailureResponse("Session does not exist.", w)
		return
	}

	// check if all users exist, skip non existing
	existingUsers := []uint{}
	for _, userID := range req.UserIDs {
		user, err := h.GetDB().GetUserById(userID)
		if err != nil {
			slog.Warn("Adding user that does not exist to session is skipped", slog.Uint64("userID", uint64(userID)), slog.Uint64("sessionID", uint64(req.SessionID)), slog.String("error", err.Error()))
			continue
		}
		if user.ID == 0 {
			slog.Warn("Adding user that does not exist to session is skipped", slog.Uint64("userID", uint64(userID)), slog.Uint64("sessionID", uint64(req.SessionID)))
			continue
		}
		existingUsers = append(existingUsers, userID)
	}

	// add existing users to session
	for _, u := range existingUsers {
		entry := database.UserToSession{
			SessionID: s.ID,
			UserID:    u,
			TimeAdded: time.Now(),
		}
		err = h.GetDB().AddSessionToUserEntry(entry)
		if err != nil {
			slog.Warn("Cannot add user to session. Skipped", slog.Uint64("userID", uint64(u)), slog.Uint64("sessionID", uint64(req.SessionID)), slog.String("error", err.Error()))
		}
	}
	WriteSuccessResponse("users added", nil, w)
}

func (h *Handler) RemoveUserFromSession(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &RemoveSessionUserRequest{}
	err := ReadBodyAndValidate(r, req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid request.", w)
		return
	}

	// remove entry
	err = h.GetDB().DeleteSessionToUserEntry(req.UserID, req.SessionID)
	if err != nil {
		slog.Warn("User to Session entry not found for", slog.Uint64("sessionID", uint64(req.SessionID)), slog.Uint64("userID", uint64(req.UserID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot remove user from session.", w)
		return
	}
	WriteSuccessResponse("user removed", nil, w)
}

func (h *Handler) GetSessionHeats(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &GetSessionHeatsRequest{}
	err := ReadBodyAndValidate(r, req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid request.", w)
		return
	}

	heats, err := h.GetDB().GetHeatsInSession(req.SessionID)
	if err != nil {
		slog.Warn("Cannot get heats for session", slog.Uint64("sessionID", uint64(req.SessionID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get heats for session.", w)
		return
	}

	// get pricing for each user
	pricing, err := h.GetDB().GetPricings()
	if err != nil {
		slog.Warn("Cannot get pricing information for session", slog.Uint64("sessionID", uint64(req.SessionID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get pricing information for heats.", w)
		return
	}
	pricingMap := map[uint]decimal.Decimal{}
	for _, p := range pricing {
		pricingMap[p.UserStatusID] = p.PricePerMinute
	}

	resp := []GetSessionHeatsResponse{}
	for _, heat := range heats {
		if _, exists := pricingMap[heat.User.UserStatusID]; !exists {
			slog.Warn("Cannot get pricing information for", slog.Uint64("user", uint64(heat.User.ID)))
			WriteFailureResponse("Cannot get pricing information for users.", w)
			return
		}
		resp = append(resp, GetSessionHeatsResponse{
			HeatID:    heat.ID,
			UserID:    heat.UserID,
			FirstName: heat.User.FirstName,
			LastName:  heat.User.LastName,
			SessionID: heat.SessionID,
			Timestamp: heat.Timestamp.Unix(),
			Duration:  int64(heat.DurationSeconds),
			Cost:      heat.Cost,
			Pricing:   pricingMap[heat.User.UserStatusID],
			Comment:   heat.Comment,
		})
	}
	WriteSuccessResponse("heats for session", &resp, w)
}
