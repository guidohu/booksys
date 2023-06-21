package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	sunrise "github.com/nathan-osman/go-sunrise"
	"golang.org/x/exp/slog"
)

type GetBookingDayRequest struct {
	Start int64 `json:"start" validate:"required,number"`
	End   int64 `json:"end" validate:"required,number"`
}

type GetBookingResponse struct {
	Start            int64             `json:"window_start"`
	End              int64             `json:"window_end"`
	Timezone         string            `json:"timezone"`
	Sunrise          int64             `json:"sunrise"`
	Sunset           int64             `json:"sunset"`
	OpeningHourStart string            `json:"business_day_start"`
	OpeningHourEnd   string            `json:"business_day_end"`
	Sessions         []SessionResponse `json:"sessions"`
}

type SessionResponse struct {
	ID               uint            `json:"id"`
	Start            int64           `json:"start"`
	End              int64           `json:"end"`
	Title            string          `json:"title"`
	Comment          string          `json:"comment"`
	FreeSpaces       uint            `json:"free"`
	CreatorID        uint            `json:"creator_id"`
	CreatorFirstName string          `json:"creator_first_name"`
	CreatorLastName  string          `json:"creator_last_name"`
	Duration         int64           `json:"duration"`
	Riders           []RiderResponse `json:"riders"`
}

type RiderResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetBookingDay(w http.ResponseWriter, r *http.Request) {
	var req GetBookingDayRequest
	err := ReadBodyAndValidate(r, &req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Request payload not valid", w)
		return
	}
	b, err := h.getBooking(time.Unix(req.Start, 0), time.Unix(req.End, 0))
	if err != nil {
		slog.Warn("Cannot retrieve bookings", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get bookings", w)
		return
	}
	WriteSuccessResponse("bookings retrieved", &b, w)
}

func (h *Handler) getBooking(start time.Time, end time.Time) (GetBookingResponse, error) {
	sunrise, sunset := h.getSunriseSunset(start)
	timezone, err := h.GetDB().GetPropertyValue("location.timezone")
	if err != nil {
		slog.Warn("Cannot retrieve location.timezone from the database", slog.String("error", err.Error()))
		return GetBookingResponse{}, err
	}
	businessDayStart, err := h.GetDB().GetPropertyValue("business.day.start")
	if err != nil {
		slog.Warn("Cannot retrieve business.day.start from the database", slog.String("error", err.Error()))
		return GetBookingResponse{}, err
	}
	businessDayEnd, err := h.GetDB().GetPropertyValue("business.day.end")
	if err != nil {
		slog.Warn("Cannot retrieve business.day.end from the database", slog.String("error", err.Error()))
		return GetBookingResponse{}, err
	}
	b := &GetBookingResponse{
		Start:            start.Unix(),
		End:              end.Unix(),
		Timezone:         timezone.Value,
		Sunrise:          sunrise.Unix(),
		Sunset:           sunset.Unix(),
		OpeningHourStart: businessDayStart.Value,
		OpeningHourEnd:   businessDayEnd.Value,
		Sessions:         []SessionResponse{},
	}

	// get sessions for that timeframe
	s, err := h.GetDB().GetSessionsBetween(start, end)
	if err != nil {
		slog.Warn("Cannot retrieve the sessions from the database", slog.String("error", err.Error()))
		return GetBookingResponse{}, err
	}

	for _, session := range s {
		sr := SessionResponse{
			ID:               session.ID,
			Start:            session.Start.Unix(),
			End:              session.End.Unix(),
			Title:            session.Title,
			Comment:          session.Comment,
			FreeSpaces:       session.FreeSpaces,
			CreatorID:        session.CreatorID,
			CreatorFirstName: session.Creator.FirstName,
			CreatorLastName:  session.Creator.LastName,
			Duration:         int64(session.End.Sub(session.Start).Seconds()),
			Riders:           []RiderResponse{},
		}
		users, err := h.GetDB().GetUsersForSession(session.ID)
		if err != nil {
			slog.Warn("Cannot retrieve users for", slog.Uint64("session", uint64(session.ID)), slog.String("error", err.Error()))
			return GetBookingResponse{}, err
		}
		for _, user := range users {
			sr.Riders = append(sr.Riders, RiderResponse{
				ID:   user.ID,
				Name: fmt.Sprintf("%s %s", user.User.FirstName, user.User.LastName),
			})
		}
		b.Sessions = append(b.Sessions, sr)
	}
	return *b, nil
}

func (h *Handler) getSunriseSunset(date time.Time) (time.Time, time.Time) {
	latProperty, err := h.GetDB().GetPropertyValue("location.latitude")
	if err != nil {
		return time.Time{}, time.Time{}
	}
	lonProperty, err := h.GetDB().GetPropertyValue("location.longitude")
	if err != nil {
		return time.Time{}, time.Time{}
	}
	lat, err := strconv.ParseFloat(latProperty.Value, 64)
	if err != nil {
		return time.Time{}, time.Time{}
	}
	lon, err := strconv.ParseFloat(lonProperty.Value, 64)
	if err != nil {
		return time.Time{}, time.Time{}
	}

	return sunrise.SunriseSunset(
		lat,
		lon,
		date.Year(), date.Month(), date.Day(),
	)
}
