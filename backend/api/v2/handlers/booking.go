package handlers

import (
	"fmt"
	"net/http"
	"server/database"
	"strconv"
	"time"
	_ "time/tzdata"

	sunrise "github.com/nathan-osman/go-sunrise"
	"golang.org/x/exp/slog"
)

type GetBookingDayRequest struct {
	Start int64 `json:"start" validate:"required,number"`
	End   int64 `json:"end" validate:"required,number"`
}

type GetBookingSeriesRequest struct {
	TimeWindows []GetBookingDayRequest `json:"time_windows" validate:"required"`
}

type GetBookingResponse struct {
	Start            int64             `json:"window_start"`
	StartText        string            `json:"window_start_text"`
	End              int64             `json:"window_end"`
	EndText          string            `json:"window_end_text"`
	Timezone         string            `json:"timezone"`
	Sunrise          int64             `json:"sunrise"`
	SunriseText      string            `json:"sunrise_text"`
	Sunset           int64             `json:"sunset"`
	SunsetText       string            `json:"sunset_text"`
	OpeningHourStart string            `json:"business_day_start"`
	OpeningHourEnd   string            `json:"business_day_end"`
	Sessions         []SessionResponse `json:"sessions"`
}

type GetBookingSeriesResponse []GetBookingResponse

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
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *Handler) GetBookingDay(w http.ResponseWriter, r *http.Request, req GetBookingDayRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	b, err := h.getBooking(dbh, time.Unix(req.Start, 0), time.Unix(req.End, 0))
	if err != nil {
		slog.Warn("Cannot retrieve bookings", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get bookings", w)
		return
	}
	WriteSuccessResponse("bookings retrieved", &b, w)
}

func (h *Handler) GetBookingSeries(w http.ResponseWriter, r *http.Request, req GetBookingSeriesRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	resp := []GetBookingResponse{}
	for _, window := range req.TimeWindows {
		b, err := h.getBooking(dbh, time.Unix(window.Start, 0), time.Unix(window.End, 0))
		if err != nil {
			slog.Warn("Cannot retrieve bookings for window", slog.Int64("start", window.Start), slog.Int64("end", window.End), slog.String("error", err.Error()))
			b = GetBookingResponse{}
		}
		resp = append(resp, b)
	}
	var response GetBookingSeriesResponse
	response = GetBookingSeriesResponse(resp)
	WriteSuccessResponse("bookings", &response, w)
}

func (h *Handler) getBooking(db database.Database, start time.Time, end time.Time) (GetBookingResponse, error) {
	location, err := h.getLocation()
	if err != nil {
		slog.Warn("Cannot get timezone for sunrise/sunset calculations", slog.String("error", err.Error()))
		return GetBookingResponse{}, err
	}
	sunrise, sunset := h.getSunriseSunset(start, location)
	businessDayStart, err := db.GetPropertyValue("business.day.start")
	if err != nil {
		slog.Warn("Cannot retrieve business.day.start from the database", slog.String("error", err.Error()))
		return GetBookingResponse{}, err
	}
	businessDayEnd, err := db.GetPropertyValue("business.day.end")
	if err != nil {
		slog.Warn("Cannot retrieve business.day.end from the database", slog.String("error", err.Error()))
		return GetBookingResponse{}, err
	}
	b := &GetBookingResponse{
		Start:            start.Unix(),
		StartText:        start.Format(time.RFC3339),
		End:              end.Unix(),
		EndText:          end.Format(time.RFC3339),
		Timezone:         location.String(),
		Sunrise:          sunrise.Unix(),
		SunriseText:      sunrise.Format(time.RFC3339),
		Sunset:           sunset.Unix(),
		SunsetText:       sunset.Format(time.RFC3339),
		OpeningHourStart: businessDayStart.Value,
		OpeningHourEnd:   businessDayEnd.Value,
		Sessions:         []SessionResponse{},
	}

	// get sessions for that timeframe
	s, err := db.GetSessionsBetween(start, end)
	if err != nil {
		slog.Warn("Cannot retrieve the sessions from the database", slog.String("error", err.Error()))
		return GetBookingResponse{}, err
	}

	for _, session := range s {
		sr := SessionResponse{
			ID:               session.ID,
			Start:            session.StartTime.Unix(),
			End:              session.EndTime.Unix(),
			Title:            session.Title,
			Comment:          session.Comment,
			FreeSpaces:       session.FreeSpaces,
			CreatorID:        session.CreatorID,
			CreatorFirstName: session.Creator.FirstName,
			CreatorLastName:  session.Creator.LastName,
			Duration:         int64(session.EndTime.Sub(session.StartTime).Seconds()),
			Riders:           []RiderResponse{},
		}
		riders, err := h.getRiders(session.ID)
		if err != nil {
			slog.Warn("Cannot retrieve users for", slog.Uint64("session", uint64(session.ID)), slog.String("error", err.Error()))
			return GetBookingResponse{}, err
		}
		for _, rider := range riders {
			sr.Riders = append(sr.Riders, RiderResponse{
				ID:        rider.ID,
				Name:      fmt.Sprintf("%s %s", rider.FirstName, rider.LastName),
				FirstName: rider.FirstName,
				LastName:  rider.LastName,
			})
		}
		b.Sessions = append(b.Sessions, sr)
	}
	return *b, nil
}

func (h *Handler) getRiders(sessionID uint) ([]database.User, error) {
	users := []database.User{}
	dbh, done := h.Database.GetHandler()
	defer done()
	if dbh == nil {
		slog.Warn("No database connection is available.")
		return users, fmt.Errorf("no database connection available")
	}
	usersToSession, err := dbh.GetUsersForSession(sessionID)
	if err != nil {
		slog.Warn("Cannot retrieve users for", slog.Uint64("session", uint64(sessionID)), slog.String("error", err.Error()))
		return users, err
	}
	for _, entry := range usersToSession {
		users = append(users, entry.User)
	}
	return users, nil
}

func (h *Handler) getSunriseSunset(date time.Time, timezone *time.Location) (time.Time, time.Time) {
	dbh, done := h.Database.GetHandler()
	defer done()
	if dbh == nil {
		slog.Warn("No database connection is available.")
		return time.Time{}, time.Time{}
	}
	latProperty, err := dbh.GetPropertyValue("location.latitude")
	if err != nil {
		return time.Time{}, time.Time{}
	}
	lonProperty, err := dbh.GetPropertyValue("location.longitude")
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
		date.In(timezone).Year(),
		date.In(timezone).Month(),
		date.In(timezone).Day(),
	)
}

func (h *Handler) getLocation() (*time.Location, error) {
	dbh, done := h.Database.GetHandler()
	defer done()
	if dbh == nil {
		slog.Warn("No database connection is available.")
		return nil, fmt.Errorf("no database connection available")
	}
	timezone, err := dbh.GetPropertyValue("location.timezone")
	if err != nil {
		slog.Warn("Cannot retrieve location.timezone from the database", slog.String("error", err.Error()))
		return nil, err
	}
	location, err := time.LoadLocation(timezone.Value)
	if err != nil {
		slog.Warn("Invalid location.timezone retrieved from the database", slog.String("error", err.Error()))
		return nil, err
	}
	return location, nil
}
