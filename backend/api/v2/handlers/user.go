package handlers

import (
	"errors"
	"math"
	"math/rand"
	"net/http"
	"server/database"
	"server/notifications/email"
	"server/recaptcha"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
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
	RecaptchaToken string `json:"recaptcha_token" validate:"omitempty"`
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
	Date            string          `json:"date"` // dd.mm.YYYY representation
	DateUnixMillis  int64           `json:"date_unix_millis"`
	Cost            decimal.Decimal `json:"cost"`     // #.## representation
	DurationText    string          `json:"duration"` // ##:## representation
	DurationSeconds int64           `json:"duration_seconds"`
}

type GetMyHeatStatsResponse struct {
	HeatTimeMinutesTotal int64           `json:"heat_time_min"`
	HeatCostTotal        decimal.Decimal `json:"heat_cost"`
	HeatTimeMinutesYTD   int64           `json:"heat_time_min_ytd"`
	HeatCostYTD          decimal.Decimal `json:"heat_cost_ytd"`
}

type GetMyBalanceResponse struct {
	PaymentTotal decimal.Decimal `json:"payment_total"`
	PaybackTotal decimal.Decimal `json:"payback_total"`
	Balance      decimal.Decimal `json:"balance_current"`
}

type GetAllUsersShortResponse []UserShort

type UserShort struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserDetailed struct {
	ID               uint            `json:"id"`
	Username         string          `json:"username"`
	FirstName        string          `json:"first_name"`
	LastName         string          `json:"last_name"`
	Address          string          `json:"address"`
	City             string          `json:"city"`
	ZipCode          int             `json:"plz"`
	MobilePhoneNr    string          `json:"mobile"`
	Email            string          `json:"email"`
	License          bool            `json:"license"`
	Locked           bool            `json:"locked"`
	UserGroupID      uint            `json:"status"`
	Comment          string          `json:"comment"`
	TotalHeatCost    decimal.Decimal `json:"total_heat_cost"`
	TotalHeatSeconds int64           `json:"total_heat_seconds"`
	TotalPayment     decimal.Decimal `json:"total_payment"`
}

type GetUserGroupsResponse struct {
	PriceDescription     string          `json:"price_description"`
	PriceID              uint            `json:"price_id"`
	PricePerMinute       decimal.Decimal `json:"price_min"`
	UserGroupDescription string          `json:"user_group_description"`
	UserGroupID          uint            `json:"user_group_id"`
	UserGroupName        string          `json:"user_group_name"`
	UserRoleDescription  string          `json:"user_role_description"`
	UserRoleID           uint            `json:"user_role_id"`
	UserRoleName         string          `json:"user_role_name"`
}

type CreateUserGroupsRequest struct {
	PriceDescription     string          `json:"price_description"`
	PricePerMinute       decimal.Decimal `json:"price_min" validate:"required"`
	UserGroupDescription string          `json:"user_group_description"`
	UserGroupName        string          `json:"user_group_name" validate:"required"`
	UserRoleID           uint            `json:"user_role_id"`
}

type ChangeUserGroupRequest struct {
	PriceDescription     string          `json:"price_description" validate:"required"`
	PriceID              uint            `json:"price_id" validate:"required"`
	PricePerMinute       decimal.Decimal `json:"price_min" validate:"required"`
	UserGroupDescription string          `json:"user_group_description" validate:"required"`
	UserGroupID          uint            `json:"user_group_id" validate:"required"`
	UserGroupName        string          `json:"user_group_name" validate:"required"`
	UserRoleID           uint            `json:"user_role_id" validate:"required"`
}

type DeleteUserGroupRequest struct {
	UserGroupID uint `json:"user_group_id" validate:"required"`
}

var UserGroupsValidationErrors = map[string]string{
	"PriceDescription":     "Please add a short description for the pricing.",
	"PriceID":              "Invalid price_id provided.",
	"PricePerMinute":       "Please provide pricing information for a price per minute (e.g. 2.30).",
	"UserGroupDescription": "Please provide a description for the user group.",
	"UserGroupID":          "The user group ID is invalid.",
	"UserGroupName":        "Please provide a unique user group name.",
	"UserRoleDescription":  "Please provide a description for the user role.",
	"UserRoleID":           "The user role ID is invalid.",
	"UserRoleName":         "Please provide a user role name.",
}

type SetUserGroupRequest struct {
	UserID      uint `json:"user_id" validate:"required"`
	UserGroupID uint `json:"status_id" validate:"required"`
}

var SetUserGroupsValidationErrors = map[string]string{
	"UserID":      "The user ID is invalid.",
	"UserGroupID": "The user group ID is invalid.",
}

type GetUserRolesResponse struct {
	UserRoleDescription string `json:"user_role_description"`
	UserRoleID          uint   `json:"user_role_id"`
	UserRoleName        string `json:"user_role_name"`
}

type SetUserLockRequest struct {
	UserID uint `json:"user_id" validate:"required"`
	Locked bool `json:"locked"`
}

var SetUserLockValidationErrors = map[string]string{
	"UserID": "The user ID is invalid.",
	"Locked": "The value for 'Locked' is invalid.",
}

type DeleteUserRequest struct {
	UserID uint `json:"user_id" validate:"required"`
}

var DeleteUserValidationErrors = map[string]string{
	"UserID": "The user ID is invalid.",
}

type GetPasswordResetTokenRequest struct {
	UserEmail      string `json:"email" validate:"email,required"`
	RecaptchaToken string `json:"recaptcha_token"`
}

var GetPasswordResetTokenValidationErrors = map[string]string{
	"UserEmail":      "Please provide a valid email address.",
	"RecaptchaToken": "Please provide a recaptcha token to prove you are not a silly robot.",
}

type SetPasswordWithTokenRequest struct {
	UserEmail string `json:"email" validate:"email,required"`
	Password  string `json:"password" validate:"required,strongpassword"`
	Token     string `json:"token" validate:"required"`
}

var SetPasswordWithTokenValidationErrors = map[string]string{
	"UserEmail": "Please provide a valid email address.",
	"Password":  "A new password has to be provided. It needs to be at least 12 characters and contain a capital letter, a lower case letter and a digit or special character.",
	"Token":     "Please provide the token that was sent to you.",
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	req := &SignUpRequest{}
	err := ReadBodyAndValidate(r, req, SignUpRequestValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	dbh := h.GetDB()
	// Get recaptcha keys (resp, entire configuration).
	config, err := dbh.GetAllPropertyValuesMap()
	if err != nil {
		slog.Warn("Cannot load internal configuration properties", slog.String("error", err.Error()))
		WriteFailureResponse("Internal error, cannot sign up user.", w)
		return
	}
	v, exists := config["recaptcha.privatekey"]
	if exists && v.Value != "" {
		valid, err := recaptcha.Valid(req.RecaptchaToken, v.Value)
		if err != nil {
			WriteFailureResponse("Recaptcha check failed, cannot verify token.", w)
			return
		}
		if !valid {
			WriteFailureResponse("Recaptcha check failed, provided recaptcha token is invalid.", w)
			return
		}
	}

	// Check if user already exists to not overwrite it
	if _, err = dbh.GetUserByName(req.Username); err == nil {
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
	id, err := dbh.AddUser(user)
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

	dbh := h.GetDB()
	// check if an admin user exists already
	// only the very first user can become an admin
	if dbh.CountAdminUsers() > 0 {
		slog.Warn("An admin user already exists, cannot make the user 'administrator'", slog.String("ID", strconv.Itoa(req.UserID)))
		WriteFailureResponse("An admin user already exists, cannot make the user 'administrator'", w)
		return
	}

	// change the actual user status
	err = dbh.ChangeUserStatus(uint(req.UserID), database.UserStatusAdmin)
	if err != nil {
		slog.Warn("Cannot make the user an 'administrator', database action failed", slog.String("ID", strconv.Itoa(req.UserID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot make user an administrator. Call to DB failed.", w)
		return
	}

	// unlock the user (users get created locked by default)
	err = dbh.ChangeLock(uint(req.UserID), false)
	if err != nil {
		slog.Warn("Cannot unlock the new user, database action failed", slog.String("ID", strconv.Itoa(req.UserID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot make user an administrator. Call to DB failed.", w)
		return
	}
	WriteSuccessResponse("success", nil, w)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// balance has to be zero
	// delete user (zeroing out personal info)
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	req := &DeleteUserRequest{}
	err := ReadBodyAndValidate(r, req, DeleteUserValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	if session.UserID == req.UserID {
		slog.Info("Skip user deletion, user cannot delete itself.")
		WriteFailureResponse("You cannot delete yourself.", w)
		return
	}

	dbh := h.GetDB()
	user, err := dbh.GetUserById(req.UserID)
	if err != nil {
		slog.Warn("Cannot find", slog.Int("user_id", int(req.UserID)), ":", err.Error())
		WriteFailureResponse("User not found.", w)
		return
	}

	if user.UserStatus.UserRoleID == database.UserRoleAdmin {
		slog.Warn("Cannot delete admin ", slog.Int("user_id", int(req.UserID)))
		WriteFailureResponse("Users with admin roles cannot be removed.", w)
		return
	}

	balance, err := h.getUserBalance(req.UserID)
	if err != nil {
		slog.Warn("Cannot get balance for", slog.Int("user_id", int(req.UserID)))
		WriteFailureResponse("Cannot check for balance to be 0 for this user.", w)
		return
	}
	if !balance.Balance.IsZero() {
		slog.Warn("Balance is not 0 for", slog.Int("user_id", int(req.UserID)), ": Cannot delete user.")
		WriteFailureResponse("Cannot delete user with non-zero balance.", w)
		return
	}
	err = dbh.DeleteUserById(req.UserID)
	if err != nil {
		slog.Error("Cannot delete user", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot delete user.", w)
		return
	}
	WriteSuccessResponse("user deleted", nil, w)
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

	// Handle special case where there are multiple users with
	// the same email but different usernames.
	username := req.Email
	if session.Username != session.User.Email {
		username = session.Username
	}
	user := database.User{
		ID:            session.UserID,
		Username:      username,
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

	dbh := h.GetDB()
	err = dbh.UpdateUser(session.UserID, user)
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

	dbh := h.GetDB()
	// verify that old password is correct
	user, err := dbh.GetUserById(session.UserID)
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
	err = dbh.UpdatePassword(session.UserID, user)
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

	dbh := h.GetDB()
	sessions, err := dbh.GetSessionsByUser(session.UserID)
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
		users, err := dbh.GetUsersForSession(s.ID)
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

	dbh := h.GetDB()
	heats, err := dbh.GetUserHeats(session.UserID, 100)
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
			Cost:            heat.Cost,
			DurationText:    t.Add(duration).Format("15:04:05"),
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

	dbh := h.GetDB()
	now := time.Now()
	duration, cost, err := dbh.GetUserHeatStats(session.UserID, time.Time{}, now)
	if err != nil {
		slog.Error("Cannot get total duration")
		WriteFailureResponse("Cannot get stats from database", w)
		return
	}
	loc, err := dbh.GetTimezoneLocation()
	if err != nil {
		slog.Error("Cannot get location", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get stats from database", w)
		return
	}
	beginOfYear := time.Date(now.Year(), time.January, 0, 0, 0, 0, 0, loc)
	durationYTD, costYTD, err := dbh.GetUserHeatStats(session.UserID, beginOfYear, time.Now())
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

	resp, err := h.getUserBalance(session.UserID)
	if err != nil {
		slog.Error("Cannot get user balance", err.Error())
		WriteFailureResponse("Cannot get user balance.", w)
		return
	}
	WriteSuccessResponse("balance", resp, w)
}

func (h *Handler) GetAllUsersShort(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	dbh := h.GetDB()
	users, err := dbh.GetUsers( /*includeDeleted=*/ false)
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

func (h *Handler) GetAllUsersDetailed(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	dbh := h.GetDB()
	users, err := dbh.GetUsers( /*includeDeleted=*/ false)
	if err != nil {
		slog.Error("Cannot get users", slog.String("error", err.Error()))
		WriteFailureResponse("cannot get users", w)
		return
	}

	usersDetailed := []UserDetailed{}
	for _, u := range users {
		duration, cost, err := dbh.GetUserHeatStats(u.ID, time.Time{}, time.Now())
		if err != nil {
			slog.Error("Cannot get user heat duration and cost", slog.String("error", err.Error()))
			WriteFailureResponse("cannot get heat duration and cost", w)
			return
		}
		paybacks, err := dbh.GetUserSessionPaybacks(u.ID)
		if err != nil {
			slog.Error("Cannot get user paybacks", slog.String("error", err.Error()))
			WriteFailureResponse("cannot get user paybacks", w)
			return
		}
		payments, err := dbh.GetUserSessionPayments(u.ID)
		if err != nil {
			slog.Error("Cannot get user payments", slog.String("error", err.Error()))
			WriteFailureResponse("cannot get user payments", w)
			return
		}
		usersDetailed = append(usersDetailed, UserDetailed{
			ID:               u.ID,
			Username:         u.Username,
			FirstName:        u.FirstName,
			LastName:         u.LastName,
			Address:          u.Address,
			City:             u.City,
			ZipCode:          u.ZipCode,
			MobilePhoneNr:    u.MobilePhoneNr,
			Email:            u.Email,
			License:          u.BoatLicense,
			Locked:           u.Locked,
			UserGroupID:      u.UserStatusID,
			Comment:          u.Comment,
			TotalHeatCost:    cost,
			TotalHeatSeconds: duration,
			TotalPayment:     payments.Sub(paybacks),
		})
	}
	WriteSuccessResponse("users detailed", usersDetailed, w)
}

func (h *Handler) GetUserGroups(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	dbh := h.GetDB()
	pricings, err := dbh.GetPricings()
	if err != nil {
		slog.Error("Cannot get user groups", slog.String("error", err.Error()))
		WriteFailureResponse("cannot get user groups", w)
		return
	}

	userGroups := []GetUserGroupsResponse{}
	for _, p := range pricings {
		userGroups = append(userGroups, GetUserGroupsResponse{
			PriceDescription:     p.Comment,
			PriceID:              p.ID,
			PricePerMinute:       p.PricePerMinute,
			UserGroupDescription: p.UserStatus.Description,
			UserGroupID:          p.UserStatusID,
			UserGroupName:        p.UserStatus.Name,
			UserRoleDescription:  p.UserStatus.UserRole.Description,
			UserRoleID:           p.UserStatus.UserRoleID,
			UserRoleName:         p.UserStatus.UserRole.Name,
		})
	}
	WriteSuccessResponse("user groups", userGroups, w)
}

func (h *Handler) CreateUserGroup(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	req := &CreateUserGroupsRequest{}
	err := ReadBodyAndValidate(r, req, UserGroupsValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	// this request consists of creating a user pricing
	// and a user group in a single request.
	u := database.UserStatus{
		Name:        req.UserGroupName,
		Description: req.UserGroupDescription,
		UserRoleID:  req.UserRoleID,
	}
	p := database.Pricing{
		PricePerMinute: req.PricePerMinute,
		Comment:        req.PriceDescription,
	}

	dbh := h.GetDB()
	err = dbh.CreateUserGroup(u, p)
	if err != nil {
		slog.Error("Cannot create new user group", slog.String("error", err.Error()))
		WriteFailureResponse("cannot get create new user group", w)
		return
	}
	WriteSuccessResponse("user group created", nil, w)
}

func (h *Handler) ChangeUserGroup(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	req := &ChangeUserGroupRequest{}
	err := ReadBodyAndValidate(r, req, UserGroupsValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	// this request consists of creating a user pricing
	// and a user group in one go
	u := database.UserStatus{
		ID:          req.UserGroupID,
		Name:        req.UserGroupName,
		Description: req.UserGroupDescription,
		UserRoleID:  req.UserRoleID,
	}
	p := database.Pricing{
		ID:             req.PriceID,
		UserStatusID:   req.UserGroupID,
		PricePerMinute: req.PricePerMinute,
		Comment:        req.PriceDescription,
	}
	dbh := h.GetDB()
	err = dbh.ChangeUserGroup(u, p)
	if err != nil {
		slog.Error("Cannot update user group", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot update user group.", w)
		return
	}
	WriteSuccessResponse("user group changed", nil, w)
}

func (h *Handler) DeleteUserGroup(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	req := &DeleteUserGroupRequest{}
	err := ReadBodyAndValidate(r, req, UserGroupsValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	dbh := h.GetDB()
	err = dbh.DeleteUserGroup(req.UserGroupID)
	if err != nil {
		slog.Error("Cannot delete user group", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot delete user group.", w)
		return
	}
	WriteSuccessResponse("user group deleted", nil, w)
}

func (h *Handler) SetUserGroup(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	req := &SetUserGroupRequest{}
	err := ReadBodyAndValidate(r, req, SetUserGroupsValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	dbh := h.GetDB()
	user, err := dbh.GetUserById(req.UserID)
	if err != nil {
		slog.Error("Cannot find user", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot find user.", w)
		return
	}
	userGroups, err := dbh.GetPricings()
	if err != nil {
		slog.Error("Cannot get valid groups", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get valid groups.", w)
		return
	}

	groupFound := false
	for _, group := range userGroups {
		if group.UserStatusID == req.UserGroupID {
			groupFound = true
			break
		}
	}
	if !groupFound {
		slog.Warn("Invalid group provided", slog.Int("user_group_id", int(req.UserGroupID)))
		WriteFailureResponse("Invalid group ID provided.", w)
		return
	}

	// check that his user is not the only remaining admin
	admins, err := dbh.GetAdminUsers()
	if err != nil {
		slog.Error("Cannot get admin users", slog.String("error", err.Error()))
		WriteFailureResponse("Internal error, cannot lock/unlock user.", w)
		return
	}
	if len(admins) == 1 && admins[0].ID == user.ID {
		slog.Warn("Cannot change group of the only admin ", slog.Uint64("user_id", uint64(user.ID)))
		WriteFailureResponse("Cannot change group of the only administrator.", w)
		return
	}

	err = dbh.SetUserGroup(user.ID, req.UserGroupID)
	if err != nil {
		slog.Error("Cannot change user group for user", slog.Int("user_id", int(req.UserID)), slog.Int("user_group_id", int(req.UserGroupID)), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot change user group for user.", w)
		return
	}
	WriteSuccessResponse("user group set", nil, w)
}

func (h *Handler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	dbh := h.GetDB()
	roles, err := dbh.GetUserRoles()
	if err != nil {
		slog.Error("Cannot get user roles", slog.String("error", err.Error()))
		WriteFailureResponse("cannot get user roles", w)
		return
	}

	userRoles := []GetUserRolesResponse{}
	for _, u := range roles {
		userRoles = append(userRoles, GetUserRolesResponse{
			UserRoleDescription: u.Description,
			UserRoleID:          u.ID,
			UserRoleName:        u.Name,
		})
	}
	WriteSuccessResponse("user roles", userRoles, w)
}

func (h *Handler) SetUserLock(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	req := &SetUserLockRequest{}
	err := ReadBodyAndValidate(r, req, SetUserLockValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	dbh := h.GetDB()
	user, err := dbh.GetUserById(req.UserID)
	if err != nil {
		slog.Error("Cannot find user", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot find user.", w)
		return
	}

	// check that this user is not the only remaining admin
	admins, err := dbh.GetAdminUsers()
	if err != nil {
		slog.Error("Cannot get admin users", slog.String("error", err.Error()))
		WriteFailureResponse("Internal error, cannot lock/unlock user.", w)
		return
	}
	if len(admins) == 1 && admins[0].ID == user.ID {
		slog.Warn("Cannot lock/unlock the only admin ", slog.Uint64("user_id", uint64(user.ID)))
		WriteFailureResponse("Cannot lock/unlock the only administrator.", w)
		return
	}

	err = dbh.ChangeLock(user.ID, req.Locked)
	if err != nil {
		slog.Error("Cannot change lock status for user", slog.Int("user_id", int(req.UserID)), slog.Bool("locked", req.Locked), slog.String("error", err.Error()))
		WriteFailureResponse("Cannot change user group for user.", w)
		return
	}
	WriteSuccessResponse("user lock set", nil, w)
}

func (h *Handler) getUserBalance(userID uint) (*GetMyBalanceResponse, error) {
	dbh := h.GetDB()
	now := time.Now()
	_, cost, err := dbh.GetUserHeatStats(userID, time.Time{}, now)
	if err != nil {
		slog.Error("Cannot get total costs", slog.Int("user", int(userID)), slog.String("error", err.Error()))
		return nil, errors.New("Cannot get balance for user.")
	}

	payment, err := dbh.GetUserSessionPayments(userID)
	if err != nil {
		slog.Error("Cannot get total payments for", slog.Int("user", int(userID)), slog.String("error", err.Error()))
		return nil, errors.New("Cannot get balance for user.")
	}

	payback, err := dbh.GetUserSessionPaybacks(userID)
	if err != nil {
		slog.Error("Cannot get total paybacks for", slog.Int("user", int(userID)), slog.String("error", err.Error()))
		return nil, errors.New("Cannot get balance for user.")
	}

	resp := &GetMyBalanceResponse{
		PaymentTotal: payment,
		PaybackTotal: payback,
		Balance:      payment.Sub(payback).Sub(cost),
	}
	return resp, nil
}

func (h *Handler) GetPasswordResetToken(w http.ResponseWriter, r *http.Request) {
	req := &GetPasswordResetTokenRequest{}
	err := ReadBodyAndValidate(r, req, GetPasswordResetTokenValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	dbh := h.GetDB()
	// Check whether recaptcha is enabled.
	config, err := dbh.GetAllPropertyValuesMap()
	if err != nil {
		slog.Warn("Cannot load internal configuration properties", slog.String("error", err.Error()))
		WriteFailureResponse("Internal error, cannot send reset token.", w)
		return
	}
	prop, exists := config["recaptcha.privatekey"]
	if !exists || prop.Value == "" {
		slog.Warn("Password reset token requests should be protected by recaptcha. Please setup recaptcha in the settings.")
	} else {
		valid, err := recaptcha.Valid(req.RecaptchaToken, prop.Value)
		if err != nil {
			WriteFailureResponse("Recaptcha check failed, cannot verify token.", w)
			return
		}
		if !valid {
			WriteFailureResponse("Recaptcha check failed, provided recaptcha token is invalid.", w)
			return
		}
	}

	// Check whether user exists.
	user, err := dbh.GetUserByName(req.UserEmail)
	if err != nil {
		slog.Warn("Cannot find user for password token request", slog.String("error", err.Error()))
		WriteSuccessResponse("Token requested, please check your email inbox.", nil, w)
		return
	}

	// Generate token.
	tokenEntry := database.PasswordReset{
		UserID:    user.ID,
		Token:     strconv.Itoa(rand.Intn(999999)),
		Timestamp: time.Now().Add(1 * time.Hour),
		Valid:     true,
	}

	// Store token in database.
	err = dbh.AddPasswordResetToken(tokenEntry)
	if err != nil {
		slog.Error("Cannot store password reset token", slog.String("error", err.Error()))
		WriteFailureResponse("Internal error, cannot send reset token.", w)
		return
	}

	// Send email with token to user.
	emailConfig, err := dbh.GetEmailConfiguration()
	if err != nil {
		slog.Error("Cannot get email configuration to reset token", slog.String("error", err.Error()))
		WriteFailureResponse("Internal error, cannot send reset token.", w)
		return
	}
	client := email.NewClient(emailConfig)
	err = client.SendTokenResetMessage(
		user,
		tokenEntry.Token,
		emailConfig.Sender,
		"Password Reset Token")
	if err != nil {
		slog.Error("Cannot send email", slog.String("error", err.Error()))
		WriteFailureResponse("Internal error, cannot send reset token.", w)
		return
	}
	WriteSuccessResponse("Token requested, please check your email inbox.", nil, w)
}

func (h *Handler) SetPasswordWithToken(w http.ResponseWriter, r *http.Request) {
	req := &SetPasswordWithTokenRequest{}
	err := ReadBodyAndValidate(r, req, SetPasswordWithTokenValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	dbh := h.GetDB()
	// Get user by email
	user, err := dbh.GetUserByName(req.UserEmail)
	if err != nil {
		slog.Warn("User cannot be found", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot reset password username or token are not valid.", w)
		return
	}

	// Check that we have a password reset token for this email.
	// That did not expire and is valid.
	dbToken, err := dbh.GetPasswordResetEntry(user.ID, req.Token)
	if err != nil {
		slog.Warn("Token cannot be found", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot reset password username or token are not valid.", w)
		return
	}

	// Verify that token is valid (we did this in the SQL statement already, though
	// better be sure)
	if dbToken.Token != req.Token {
		slog.Warn("Token do not match", slog.String("db", dbToken.Token), slog.String("request", req.Token))
		WriteFailureResponse("Cannot reset password username or token are not valid.", w)
		return
	}
	if !dbToken.Valid {
		slog.Warn("Token is not valid.", slog.String("user", req.UserEmail), slog.String("token", dbToken.Token))
		WriteFailureResponse("Cannot reset password username or token are not valid.", w)
		return
	}
	if dbToken.Timestamp.Before(time.Now()) {
		slog.Warn("Token is expired", slog.String("timestamp", dbToken.Timestamp.String()), slog.String("now", time.Now().String()))
		WriteFailureResponse("Cannot reset password username or token are not valid.", w)
		return
	}

	// Create new password hash and
	// update new password.
	// Note: We hash the passworrd, this was previously done in the UI.
	salt := rand.Intn(math.MaxUint16)
	newPasswordHash, err := cryptSha512(hashSha256(req.Password), strconv.Itoa(salt))
	if err != nil {
		slog.Warn("Password hash could no be generated", slog.String("error", err.Error()))
		WriteFailureResponse("Internal error. Password cannot be changed.", w)
		return
	}
	userPassword := database.User{
		PasswordSalt: salt,
		PasswordHash: newPasswordHash,
	}
	err = dbh.UpdatePassword(user.ID, userPassword)
	if err != nil {
		slog.Warn("Password could no be stored in user table", slog.String("error", err.Error()))
		WriteFailureResponse("Internal error. Password could not be changed.", w)
		return
	}

	// Invalidate password reset tokens for this user.
	dbh.InvalidatePasswordResetEntries(user.ID)

	WriteSuccessResponse("password reset", nil, w)
}
