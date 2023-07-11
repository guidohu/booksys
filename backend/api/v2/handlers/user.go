package handlers

import (
	"math"
	"math/rand"
	"net/http"
	"server/database"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

type SignUpRequest struct {
	Username       string `json:"username" validate:"required,excludesall=!<>{}[]()^"`
	Password       string `json:"password" validate:"required,strongpassword"`
	FirstName      string `json:"first_name" validate:"required,excludesall=!<>{}[]()^"`
	LastName       string `json:"last_name" validate:"required,excludesall=!<>{}[]()^"`
	Address        string `json:"address" validate:"required,excludesall=!<>{}[]()^"`
	MobileNr       string `json:"mobile" validate:"required,number|containsany=()+"`
	ZipCode        int    `json:"plz" validate:"required,number"`
	City           string `json:"city" validate:"required,excludesall=!<>{}[]()^"`
	Email          string `json:"email" validate:"required,email"`
	License        bool   `json:"license" validate:"boolean"`
	AcceptGTC      bool   `json:"ownRisk" validate:"required,boolean"`
	RecaptchaToken string `json:"recaptcha_token" validate:"omitempty,alphanum"`
}

type SignUpResponse struct {
	UserID int `json:"user_id"`
}

var SignUpRequestValidationErrors = map[string]string{
	"Username":       "a username has to be provided (e.g. email)",
	"Password":       "A password has to be provided. It needs to be at least 12 characters and contain a capital letter, a lower case letter and a digit or special character.",
	"FirstName":      "first name is missing or uses invalid characters",
	"LastName":       "last name is missing or uses invalid characters",
	"Address":        "the address is missing or uses invalid characters",
	"MobileNr":       "the mobile number is missing or not a valid phone number",
	"ZipCode":        "the zip code is either missing or not a number",
	"City":           "the city is missing or contains invalid characters",
	"Email":          "the email is missing or not a valid email address",
	"License":        "license information is missing",
	"AcceptGTC":      "Please accept our General Terms and Conditions",
	"RecaptchaToken": "recaptcha token has an invalid format",
}

type MakeAdminRequest struct {
	UserID int `json:"user_id" validate:"required,number"`
}

type UpdateMyUserRequest struct {
	FirstName string `json:"first_name" validate:"required,excludesall=!<>{}[]()^"`
	LastName  string `json:"last_name" validate:"required,excludesall=!<>{}[]()^"`
	Address   string `json:"address" validate:"required,excludesall=!<>{}[]()^"`
	MobileNr  string `json:"mobile" validate:"required,number|containsany=()+"`
	ZipCode   int    `json:"plz" validate:"required,number"`
	City      string `json:"city" validate:"required,excludesall=!<>{}[]()^"`
	Email     string `json:"email" validate:"required,email"`
	License   bool   `json:"license" validate:"boolean"`
}

var UpdateMyUserValidationErrors = map[string]string{
	"FirstName": "first name is missing or uses invalid characters",
	"LastName":  "last name is missing or uses invalid characters",
	"Address":   "the address is missing or uses invalid characters",
	"MobileNr":  "the mobile number is missing or not a valid phone number",
	"ZipCode":   "the zip code is either missing or not a number",
	"City":      "the city is missing or contains invalid characters",
	"Email":     "the email is missing or not a valid email address",
	"License":   "license information is missing",
}

type UpdateMyPasswordRequest struct {
	PasswordOld string `json:"password_old" validate:"required"`
	PasswordNew string `json:"password_new" validate:"required,strongpassword"`
}

var UpdateMyPasswordValidationErrors = map[string]string{
	"PasswordOld": "an old password has to be provided",
	"PasswordNew": "A new password has to be provided. It needs to be at least 12 characters and contain a capital letter, a lower case letter and a digit or special character.",
}

type GetMySessionsResponse struct {
	UpcomingSessions []MySessionResponse `json:"sessions"`
	PastSessions     []MySessionResponse `json:"sessions_old"`
}

type MySessionResponse struct {
	ID        uint        `json:"id"`
	Title     string      `json:"title"`
	Type      uint        `json:"type"`
	TypeName  string      `json:"type_name"`
	StartTime int64       `json:"start_time"`
	EndTime   int64       `json:"end_time"`
	Riders    []UserShort `json:"riders"`
}

type GetMyHeatsResponse struct {
	Heats []HeatResponse `json:"heats"`
}

type HeatResponse struct {
	Date            string  `json:"date"` // dd.mm.YYYY representation
	DateUnixMillis  int64   `json:"date_unix_millis"`
	Cost            float64 `json:"cost"`     // #.## representation
	DurationText    string  `json:"duration"` // ##:## representation
	DurationSeconds int64   `json:"duration_seconds"`
}

type GetMyHeatStatsResponse struct {
	HeatTimeMinutesTotal int64   `json:"heat_time_min"`
	HeatCostTotal        float64 `json:"heat_cost"`
	HeatTimeMinutesYTD   int64   `json:"heat_time_min_ytd"`
	HeatCostYTD          float64 `json:"heat_cost_ytd"`
}

type GetMyBalanceResponse struct {
	PaymentTotal float64 `json:"payment_total"`
	PaybackTotal float64 `json:"payback_total"`
	Balance      float64 `json:"balance_current"`
}

type GetAllUsersShortResponse []UserShort

type UserShort struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	req := &SignUpRequest{}
	err := ReadBodyAndValidate(r, req, SignUpRequestValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	// TODO recaptcha validation
	if viper.GetString("recaptcha.privatekey") != "" {
		// TODO verify recaptcha
		panic("not implemented recaptcha validation")
	}

	// Check if user already exists to not overwrite it
	if _, err = h.GetDB().GetUserByName(req.Username); err == nil {
		slog.Warn("Signup an already existing user", slog.String("user", req.Username))
		WriteFailureResponse("user already exists", w)
		return
	}

	// Crypt the password
	salt := rand.Intn(math.MaxUint16)
	hashedPassword, err := cryptSha512(hashSha256(req.Password), strconv.Itoa(salt))
	if err != nil {
		slog.Warn("Password hash could no be generated", slog.String("error", err.Error()))
		WriteFailureResponse("user cannot be created, please contact the administrator", w)
		return
	}

	// Add user to database
	user := database.User{
		Username:      req.Username,
		PasswordSalt:  salt,
		PasswordHash:  hashedPassword,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Address:       req.Address,
		City:          req.City,
		ZipCode:       req.ZipCode,
		MobilePhoneNr: req.MobileNr,
		Email:         req.Email,
		BoatLicense:   req.License,
		UserStatusID:  database.UserStatusGuest,
		Locked:        true,
	}
	id, err := h.GetDB().AddUser(user)
	if err != nil {
		slog.Warn("User could not be added to database", slog.String("error", err.Error()))
		WriteFailureResponse("user cannot be created, please contact the administrator", w)
		return
	}

	resp := &SignUpResponse{
		UserID: int(id),
	}
	WriteSuccessResponse("success", resp, w)
}

func (h *Handler) MakeAdmin(w http.ResponseWriter, r *http.Request) {
	req := &MakeAdminRequest{}
	err := ReadBodyAndValidate(r, req, nil)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Request payload not valid", w)
		return
	}

	// check if an admin user exists already
	// only the very first user can become an admin
	if h.GetDB().CountAdminUsers() > 0 {
		slog.Warn("An admin user already exists, cannot make the user 'administrator'", slog.String("ID", strconv.Itoa(req.UserID)))
		WriteFailureResponse("An admin user already exists, cannot make the user 'administrator'", w)
		return
	}

	// change the actual user status
	err = h.GetDB().ChangeUserStatus(uint(req.UserID), database.UserStatusAdmin)
	if err != nil {
		slog.Warn("Cannot make the user an 'administrator', database action failed", slog.String("ID", strconv.Itoa(req.UserID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot make user an administrator. Call to DB failed.", w)
		return
	}

	// unlock the user (users get created locked by default)
	err = h.GetDB().ChangeLock(uint(req.UserID), false)
	if err != nil {
		slog.Warn("Cannot unlock the new user, database action failed", slog.String("ID", strconv.Itoa(req.UserID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot make user an administrator. Call to DB failed.", w)
		return
	}
	WriteSuccessResponse("success", nil, w)
}

func (h *Handler) UpdateMyUser(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if session.UserID == 0 {
		slog.Warn("Unknown user accessing GetMyHeats")
		WriteFailureResponse("Not authenticated", w)
		return
	}

	req := &UpdateMyUserRequest{}
	err := ReadBodyAndValidate(r, req, UpdateMyUserValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	user := database.User{
		ID:            session.UserID,
		Username:      req.Email,
		PasswordSalt:  0,  // will not be set
		PasswordHash:  "", // will not be set
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Address:       req.Address,
		City:          req.City,
		ZipCode:       req.ZipCode,
		MobilePhoneNr: req.MobileNr,
		Email:         req.Email,
		BoatLicense:   req.License,
		UserStatusID:  0,     // will not be set
		Locked:        false, // will not be set
		Comment:       "",    // will not be set
		IsDeleted:     false, // will not be set
	}

	err = h.GetDB().UpdateUser(session.UserID, user)
	if err != nil {
		slog.Warn("Cannot update user", slog.Uint64("userID", uint64(session.UserID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot update user", w)
		return
	}
	WriteSuccessResponse("user updated", nil, w)
}

func (h *Handler) UpdateMyPassword(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if session.UserID == 0 {
		slog.Warn("Unknown user accessing UpdateMyPassword")
		WriteFailureResponse("Not authenticated", w)
		return
	}

	req := &UpdateMyPasswordRequest{}
	err := ReadBodyAndValidate(r, req, UpdateMyPasswordValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	// verify that old password is correct
	user, err := h.GetDB().GetUserById(session.UserID)
	if err != nil {
		slog.Warn("Cannot find existing user by ID", slog.String("error", err.Error()))
		WriteFailureResponse("password cannot be changed", w)
		return
	}
	hashedPassword, err := cryptSha512(hashSha256(req.PasswordOld), strconv.Itoa(user.PasswordSalt))
	if err != nil {
		slog.Warn("Password hash could no be generated", slog.String("error", err.Error()))
		WriteFailureResponse("password cannot be changed", w)
		return
	}
	if hashedPassword != user.PasswordHash {
		slog.Warn("Attempt to change password without correct old password")
		WriteFailureResponse("Password cannot be changed due to an incorrect old password.", w)
		return
	}

	// update new password
	salt := rand.Intn(math.MaxUint16)
	newPasswordHash, err := cryptSha512(hashSha256(req.PasswordNew), strconv.Itoa(salt))
	if err != nil {
		slog.Warn("Password hash could no be generated", slog.String("error", err.Error()))
		WriteFailureResponse("password cannot be changed", w)
		return
	}
	user = database.User{
		PasswordSalt: salt,
		PasswordHash: newPasswordHash,
	}
	err = h.GetDB().UpdatePassword(session.UserID, user)
	if err != nil {
		slog.Warn("Password could no be stored in user table", slog.String("error", err.Error()))
		WriteFailureResponse("Password could not be changed.", w)
		return
	}
	WriteSuccessResponse("password changed", nil, w)
}

func (h *Handler) GetMySessions(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if session.UserID == 0 {
		slog.Warn("Unknown user accessing GetMyHeats")
		WriteFailureResponse("Not authenticated", w)
		return
	}

	sessions, err := h.GetDB().GetSessionsByUser(session.UserID)
	if err != nil {
		slog.Error("Cannot get user sessions", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get user sessions", w)
		return
	}

	// split into past and future
	upcoming := []MySessionResponse{}
	past := []MySessionResponse{}
	now := time.Now()
	for _, s := range sessions {
		// get riders for session
		users, err := h.GetDB().GetUsersForSession(s.ID)
		if err != nil {
			slog.Error("Cannot get users for session", slog.String("error", err.Error()))
		}
		sessionUsers := []UserShort{}
		for _, u := range users {
			sessionUsers = append(sessionUsers, UserShort{
				ID:        u.ID,
				FirstName: u.User.FirstName,
				LastName:  u.User.LastName,
			})
		}

		session := MySessionResponse{
			ID:        s.ID,
			Title:     s.Title,
			Type:      s.SessionTypeID,
			TypeName:  s.SessionType.Name,
			StartTime: s.StartTime.Unix(),
			EndTime:   s.EndTime.Unix(),
			Riders:    sessionUsers,
		}
		if now.After(s.EndTime) {
			past = append(past, session)
		} else {
			upcoming = append(upcoming, session)
		}
	}
	resp := &GetMySessionsResponse{
		UpcomingSessions: upcoming,
		PastSessions:     past,
	}

	WriteSuccessResponse("user sessions", resp, w)
}

func (h *Handler) GetMyHeats(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if session.UserID == 0 {
		slog.Warn("Unknown user accessing GetMyHeats")
		WriteFailureResponse("Not authenticated", w)
		return
	}

	heats, err := h.GetDB().GetUserHeats(session.UserID, 100)
	if err != nil {
		slog.Error("Cannot get user heats", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get user heats", w)
		return
	}

	resp := &GetMyHeatsResponse{}
	resp.Heats = make([]HeatResponse, 0)
	for _, heat := range heats {
		duration := time.Duration(int64(heat.DurationSeconds) * int64(time.Second))
		var t time.Time
		resp.Heats = append(resp.Heats, HeatResponse{
			Date:            heat.Timestamp.Format("02.01.2006"),
			DateUnixMillis:  heat.Timestamp.UnixMilli(),
			Cost:            math.Floor(float64(heat.Cost)*100) / 100,
			DurationText:    t.Add(duration).Format("15:04"),
			DurationSeconds: int64(heat.DurationSeconds),
		})
	}

	WriteSuccessResponse("heats", resp, w)
}

func (h *Handler) GetMyHeatStats(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if session.UserID == 0 {
		slog.Warn("Unknown user accessing GetMyHeatStats")
		WriteFailureResponse("Not authenticated", w)
		return
	}

	now := time.Now()
	duration, cost, err := h.GetDB().GetUserHeatStats(session.UserID, time.Time{}, now)
	if err != nil {
		slog.Error("Cannot get total duration")
		WriteFailureResponse("Cannot get stats from database", w)
		return
	}
	loc, err := h.GetDB().GetTimezoneLocation()
	if err != nil {
		slog.Error("Cannot get location", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get stats from database", w)
		return
	}
	beginOfYear := time.Date(now.Year(), 0, 0, 0, 0, 0, 0, loc)
	durationYTD, costYTD, err := h.GetDB().GetUserHeatStats(session.UserID, beginOfYear, time.Now())
	if err != nil {
		slog.Error("Cannot get total duration")
		WriteFailureResponse("Cannot get stats from database", w)
		return
	}

	resp := &GetMyHeatStatsResponse{
		HeatTimeMinutesTotal: int64(math.Ceil(float64(duration) / 60.0)),
		HeatCostTotal:        cost,
		HeatTimeMinutesYTD:   int64(math.Ceil(float64(durationYTD) / 60.0)),
		HeatCostYTD:          costYTD,
	}
	WriteSuccessResponse("stats", resp, w)
}

func (h *Handler) GetMyBalance(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if session.UserID == 0 {
		slog.Warn("Unknown user accessing GetMyBalance")
		WriteFailureResponse("Not authenticated", w)
		return
	}

	now := time.Now()
	_, cost, err := h.GetDB().GetUserHeatStats(session.UserID, time.Time{}, now)
	if err != nil {
		slog.Error("Cannot get total costs", slog.String("user", session.User.Username), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get balance for user.", w)
		return
	}

	payment, err := h.GetDB().GetUserSessionPayments(session.UserID)
	if err != nil {
		slog.Error("Cannot get total payments for", slog.String("user", session.User.Username), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get balance for user", w)
		return
	}

	payback, err := h.GetDB().GetUserSessionPaybacks(session.UserID)
	if err != nil {
		slog.Error("Cannot get total paybacks for", slog.String("user", session.User.Username), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get balance for user", w)
		return
	}

	resp := &GetMyBalanceResponse{
		PaymentTotal: math.Round(payment*100) / 100,
		PaybackTotal: math.Round(payback*100) / 100,
		Balance:      math.Round((payment-payback-cost)*100) / 100,
	}

	WriteSuccessResponse("balance", resp, w)
}

func (h *Handler) GetAllUsersShort(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	users, err := h.GetDB().GetUsers()
	if err != nil {
		slog.Error("Cannot get users", slog.String("error", err.Error()))
		WriteFailureResponse("cannot get users", w)
		return
	}

	usersShort := []UserShort{}
	for _, u := range users {
		usersShort = append(usersShort, UserShort{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		})
	}
	WriteSuccessResponse("users", usersShort, w)

}
