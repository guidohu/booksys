package handlers

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"server/database"
	"server/notifications/email"
	"server/recaptcha"
	"server/util/hash"

	"github.com/shopspring/decimal"
)

// SignUpRequest is the request body of SignUp, which serves
// /api/v2/user/signup.
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

// SignUpResponse is the payload returned by SignUp, which serves
// /api/v2/user/signup.
type SignUpResponse struct {
	UserID int `json:"user_id"`
}

// SignUpRequestValidationErrors maps the field names of SignUpRequest to the message the API
// returns when that field fails validation.
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

// MakeAdminRequest is the request body of MakeAdmin, which serves
// /api/v2/user/create-admin.
type MakeAdminRequest struct {
	UserID int `json:"user_id" validate:"required,number"`
}

// UpdateMyUserRequest is the request body of UpdateMyUser, which serves
// /api/v2/user/my/update.
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

// UpdateMyUserValidationErrors maps the field names of UpdateMyUserRequest to the message the API
// returns when that field fails validation.
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

// UpdateMyPasswordRequest is the request body of UpdateMyPassword, which serves
// /api/v2/user/my/password/update.
type UpdateMyPasswordRequest struct {
	PasswordOld string `json:"password_old" validate:"required"`
	PasswordNew string `json:"password_new" validate:"required,strongpassword"`
}

// UpdateMyPasswordValidationErrors maps the field names of UpdateMyPasswordRequest to the message the API
// returns when that field fails validation.
var UpdateMyPasswordValidationErrors = map[string]string{
	"PasswordOld": "an old password has to be provided",
	"PasswordNew": "A new password has to be provided. It needs to be at least 12 characters and contain a capital letter, a lower case letter and a digit or special character.",
}

// GetMySessionsResponse is the payload returned by GetMySessions, which serves
// /api/v2/user/my/sessions.
type GetMySessionsResponse struct {
	UpcomingSessions []MySessionResponse `json:"sessions"`
	PastSessions     []MySessionResponse `json:"sessions_old"`
}

// MySessionResponse is a single session in a GetMySessionsResponse.
type MySessionResponse struct {
	ID        uint        `json:"id"`
	Title     string      `json:"title"`
	Type      uint        `json:"type"`
	TypeName  string      `json:"type_name"`
	StartTime int64       `json:"start_time"`
	EndTime   int64       `json:"end_time"`
	Riders    []UserShort `json:"riders"`
}

// GetMyHeatsResponse is the payload returned by GetMyHeats, which serves
// /api/v2/user/my/heats.
type GetMyHeatsResponse struct {
	Heats []HeatResponse `json:"heats"`
}

// HeatResponse is a single ride in a GetMyHeatsResponse.
type HeatResponse struct {
	Date            string          `json:"date"` // dd.mm.YYYY representation
	DateUnixMillis  int64           `json:"date_unix_millis"`
	Cost            decimal.Decimal `json:"cost"`     // #.## representation
	DurationText    string          `json:"duration"` // ##:## representation
	DurationSeconds int64           `json:"duration_seconds"`
}

// GetMyHeatStatsResponse is the payload returned by GetMyHeatStats, which serves
// /api/v2/user/my/heats/statistics.
type GetMyHeatStatsResponse struct {
	HeatTimeMinutesTotal int64           `json:"heat_time_min"`
	HeatCostTotal        decimal.Decimal `json:"heat_cost"`
	HeatTimeMinutesYTD   int64           `json:"heat_time_min_ytd"`
	HeatCostYTD          decimal.Decimal `json:"heat_cost_ytd"`
}

// GetMyBalanceResponse is the payload returned by GetMyBalance, which serves
// /api/v2/user/my/balance.
type GetMyBalanceResponse struct {
	PaymentTotal decimal.Decimal `json:"payment_total"`
	PaybackTotal decimal.Decimal `json:"payback_total"`
	Balance      decimal.Decimal `json:"balance_current"`
}

// GetAllUsersShortResponse is the payload returned by GetAllUsersShort.
type GetAllUsersShortResponse []UserShort

// UserShort is a user reduced to ID and name.
type UserShort struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UserDetailed is a user together with their ride and payment totals.
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

// GetUserGroupsResponse is the payload returned by GetUserGroups, which serves
// /api/v2/user/groups/get.
type GetUserGroupsResponse struct {
	PriceDescription     string                `json:"price_description"`
	PriceID              uint                  `json:"price_id"`
	PricePerMinute       decimal.Decimal       `json:"price_min"`
	UserGroupDescription string                `json:"user_group_description"`
	UserGroupID          uint                  `json:"user_group_id"`
	UserGroupName        string                `json:"user_group_name"`
	UserRoleDescription  string                `json:"user_role_description"`
	UserRoleID           database.UserRoleType `json:"user_role_id"`
	UserRoleName         string                `json:"user_role_name"`
}

// CreateUserGroupsRequest is the request body of CreateUserGroup.
type CreateUserGroupsRequest struct {
	PriceDescription     string                `json:"price_description"`
	PricePerMinute       decimal.Decimal       `json:"price_min" validate:"required"`
	UserGroupDescription string                `json:"user_group_description"`
	UserGroupName        string                `json:"user_group_name" validate:"required"`
	UserRoleID           database.UserRoleType `json:"user_role_id"`
}

// ChangeUserGroupRequest is the request body of ChangeUserGroup, which serves
// /api/v2/user/group/edit.
type ChangeUserGroupRequest struct {
	PriceDescription     string                `json:"price_description" validate:"required"`
	PriceID              uint                  `json:"price_id" validate:"required"`
	PricePerMinute       decimal.Decimal       `json:"price_min" validate:"required"`
	UserGroupDescription string                `json:"user_group_description" validate:"required"`
	UserGroupID          uint                  `json:"user_group_id" validate:"required"`
	UserGroupName        string                `json:"user_group_name" validate:"required"`
	UserRoleID           database.UserRoleType `json:"user_role_id" validate:"required"`
}

// DeleteUserGroupRequest is the request body of DeleteUserGroup, which serves
// /api/v2/user/group/delete.
type DeleteUserGroupRequest struct {
	UserGroupID uint `json:"user_group_id" validate:"required"`
}

// UserGroupsValidationErrors maps the field names of the user group requests to the message the API
// returns when that field fails validation.
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

// SetUserGroupRequest is the request body of SetUserGroup.
type SetUserGroupRequest struct {
	UserID      uint `json:"user_id" validate:"required"`
	UserGroupID uint `json:"status_id" validate:"required"`
}

// SetUserGroupsValidationErrors maps the field names of SetUserGroupRequest to the message the API
// returns when that field fails validation.
var SetUserGroupsValidationErrors = map[string]string{
	"UserID":      "The user ID is invalid.",
	"UserGroupID": "The user group ID is invalid.",
}

// GetUserRolesResponse is the payload returned by GetUserRoles, which serves
// /api/v2/user/roles/get.
type GetUserRolesResponse struct {
	UserRoleDescription string                `json:"user_role_description"`
	UserRoleID          database.UserRoleType `json:"user_role_id"`
	UserRoleName        string                `json:"user_role_name"`
}

// SetUserLockRequest is the request body of SetUserLock, which serves
// /api/v2/user/lock/set.
type SetUserLockRequest struct {
	UserID uint `json:"user_id" validate:"required"`
	Locked bool `json:"locked"`
}

// SetUserLockValidationErrors maps the field names of SetUserLockRequest to the message the API
// returns when that field fails validation.
var SetUserLockValidationErrors = map[string]string{
	"UserID": "The user ID is invalid.",
	"Locked": "The value for 'Locked' is invalid.",
}

// DeleteUserRequest is the request body of DeleteUser, which serves
// /api/v2/user/delete.
type DeleteUserRequest struct {
	UserID uint `json:"user_id" validate:"required"`
}

// DeleteUserValidationErrors maps the field names of DeleteUserRequest to the message the API
// returns when that field fails validation.
var DeleteUserValidationErrors = map[string]string{
	"UserID": "The user ID is invalid.",
}

// GetPasswordResetTokenRequest is the request body of GetPasswordResetToken, which serves
// /api/v2/user/password/token-request.
type GetPasswordResetTokenRequest struct {
	UserEmail      string `json:"email" validate:"email,required"`
	RecaptchaToken string `json:"recaptcha_token"`
}

// GetPasswordResetTokenValidationErrors maps the field names of GetPasswordResetTokenRequest to the message the API
// returns when that field fails validation.
var GetPasswordResetTokenValidationErrors = map[string]string{
	"UserEmail":      "Please provide a valid email address.",
	"RecaptchaToken": "Please provide a recaptcha token to prove you are not a silly robot.",
}

// SetPasswordWithTokenRequest is the request body of SetPasswordWithToken, which serves
// /api/v2/user/password/reset-by-token.
type SetPasswordWithTokenRequest struct {
	UserEmail string `json:"email" validate:"email,required"`
	Password  string `json:"password" validate:"required,strongpassword"`
	Token     string `json:"token" validate:"required"`
}

// SetPasswordWithTokenValidationErrors maps the field names of SetPasswordWithTokenRequest to the message the API
// returns when that field fails validation.
var SetPasswordWithTokenValidationErrors = map[string]string{
	"UserEmail": "Please provide a valid email address.",
	"Password":  "A new password has to be provided. It needs to be at least 12 characters and contain a capital letter, a lower case letter and a digit or special character.",
	"Token":     "Please provide the token that was sent to you.",
}

// The rules that make a six digit reset token safe. The token itself is only
// 10^6 wide, so what bounds an attacker is not its entropy but how many guesses
// they can make: one live token per user, burned after passwordResetMaxAttempts
// wrong guesses, and a new one obtainable only once per passwordResetCooldown.
// That caps a determined attacker at 60 guesses an hour, and every one of those
// tokens puts a mail in the victim's inbox.
const (
	// passwordResetValidity is how long a freshly issued token can be redeemed.
	passwordResetValidity = 15 * time.Minute
	// passwordResetCooldown is the minimum time between two token requests for
	// the same user.
	passwordResetCooldown = 2 * time.Minute
	// passwordResetMaxAttempts is how many wrong tokens may be submitted
	// against a live token before it is invalidated.
	passwordResetMaxAttempts = 2
	// passwordResetTokenSpace is the exclusive upper bound of the token value.
	// Tokens are rendered zero padded, so the whole space is six digits wide.
	passwordResetTokenSpace = 1000000
)

// passwordResetRequested is the answer GetPasswordResetToken gives on every
// path, so that an unknown address, an address in its cooldown window and a
// token that was really just sent are indistinguishable in the response body.
//
// Note that this hides the difference in the payload only. The three paths do
// different amounts of work, and the issuing one sends mail synchronously, so
// response time still tells them apart. Closing that needs the mail to move off
// the request path.
const passwordResetRequested = "Token requested, please check your email inbox."

// passwordResetRejected is the answer SetPasswordWithToken gives on every
// failure, for the same reason.
const passwordResetRejected = "Cannot reset password username or token are not valid."

// randomInt returns a uniformly distributed random integer in [0, max), drawn
// from the operating system entropy source. Password salts and reset tokens
// have to be unpredictable, so math/rand is not good enough for them.
//
// max has to be positive.
func randomInt(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

// SignUp registers a new user. The account is created locked, an administrator
// has to unlock it.
//
// It serves /api/v2/user/signup and is open to unauthenticated callers.
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request, req SignUpRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	// Get recaptcha keys
	privateKey, _ := h.config.GetString("recaptcha.privatekey")
	if privateKey != "" {
		valid, err := recaptcha.Valid(req.RecaptchaToken, privateKey)
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
	if _, err := dbh.GetUserByName(req.Username); err == nil {
		slog.Warn("Signup an already existing user", slog.String("user", req.Username))
		WriteFailureResponse("user already exists", w)
		return
	}

	// Crypt the password
	salt, err := randomInt(math.MaxUint16)
	if err != nil {
		slog.Error("Cannot generate a password salt", slog.Any("error", err))
		WriteFailureResponse("user cannot be created, please contact the administrator", w)
		return
	}
	hashedPassword, err := hash.CryptSha512(hash.Sha256(req.Password), strconv.Itoa(salt))
	if err != nil {
		slog.Warn("Password hash could no be generated", slog.Any("error", err))
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
		slog.Warn("User could not be added to database", slog.Any("error", err))
		WriteFailureResponse("user cannot be created, please contact the administrator", w)
		return
	}

	resp := &SignUpResponse{
		UserID: int(id),
	}
	WriteSuccessResponse("success", resp, w)
}

// MakeAdmin promotes a user to administrator and unlocks them. It only works while
// no administrator exists.
//
// It serves /api/v2/user/create-admin and is open to callers when web setup is enabled/unauthenticated callers.
func (h *Handler) MakeAdmin(w http.ResponseWriter, r *http.Request, req MakeAdminRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	// check if an admin user exists already
	// only the very first user can become an admin
	if dbh.CountAdminUsers() > 0 {
		slog.Warn("An admin user already exists, cannot make the user 'administrator'", slog.String("ID", strconv.Itoa(req.UserID)))
		WriteFailureResponse("An admin user already exists, cannot make the user 'administrator'", w)
		return
	}

	// change the actual user status
	err := dbh.ChangeUserStatus(uint(req.UserID), database.UserStatusAdmin)
	if err != nil {
		slog.Warn("Cannot make the user an 'administrator', database action failed", slog.String("ID", strconv.Itoa(req.UserID)), slog.Any("error", err))
		WriteFailureResponse("Cannot make user an administrator. Call to DB failed.", w)
		return
	}

	// unlock the user (users get created locked by default)
	err = dbh.ChangeLock(uint(req.UserID), false)
	if err != nil {
		slog.Warn("Cannot unlock the new user, database action failed", slog.String("ID", strconv.Itoa(req.UserID)), slog.Any("error", err))
		WriteFailureResponse("Cannot make user an administrator. Call to DB failed.", w)
		return
	}
	WriteSuccessResponse("success", nil, w)
}

// DeleteUser anonymizes a user. It refuses to delete the caller, administrators, and
// users with a non-zero balance.
//
// It serves /api/v2/user/delete and is open to administrators.
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, req DeleteUserRequest, hCtx *HandlerCtx) {
	// balance has to be zero
	// delete user (zeroing out personal info)
	if hCtx.ValidSession.UserID == req.UserID {
		slog.Info("Skip user deletion, user cannot delete itself.")
		WriteFailureResponse("You cannot delete yourself.", w)
		return
	}

	dbh := hCtx.Database
	user, err := dbh.GetUserByID(req.UserID)
	if err != nil {
		slog.Warn("Cannot find user", slog.Uint64("user_id", uint64(req.UserID)), slog.Any("error", err))
		WriteFailureResponse("User not found.", w)
		return
	}

	if user.UserStatus.UserRoleID == database.UserRoleAdmin {
		slog.Warn("Cannot delete admin ", slog.Int("user_id", int(req.UserID)))
		WriteFailureResponse("Users with admin roles cannot be removed.", w)
		return
	}

	balance, err := h.getUserBalance(dbh, req.UserID)
	if err != nil {
		slog.Warn("Cannot get balance for", slog.Int("user_id", int(req.UserID)))
		WriteFailureResponse("Cannot check for balance to be 0 for this user.", w)
		return
	}
	if !balance.Balance.IsZero() {
		slog.Warn("Cannot delete user with a non-zero balance", slog.Uint64("user_id", uint64(req.UserID)))
		WriteFailureResponse("Cannot delete user with non-zero balance.", w)
		return
	}
	err = dbh.DeleteUserByID(req.UserID)
	if err != nil {
		slog.Error("Cannot delete user", slog.Any("error", err))
		WriteFailureResponse("Cannot delete user.", w)
		return
	}
	WriteSuccessResponse("user deleted", nil, w)
}

// UpdateMyUser updates the profile of the currently logged in user.
//
// It serves /api/v2/user/my/update and is open to authenticated users.
func (h *Handler) UpdateMyUser(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	session := hCtx.ValidSession
	if session.UserID == 0 {
		slog.Warn("Unknown user accessing GetMyHeats")
		WriteFailureResponse("Not authenticated", w)
		return
	}

	req := &UpdateMyUserRequest{}
	err := ReadBodyAndValidate(r, req, UpdateMyUserValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.Any("error", err))
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

	dbh := hCtx.Database
	err = dbh.UpdateUser(session.UserID, user)
	if err != nil {
		slog.Warn("Cannot update user", slog.Uint64("userID", uint64(session.UserID)), slog.Any("error", err))
		WriteFailureResponse("Cannot update user", w)
		return
	}
	WriteSuccessResponse("user updated", nil, w)
}

// UpdateMyPassword changes the password of the currently logged in user, after checking
// the old one.
//
// It serves /api/v2/user/my/password/update and is open to authenticated users.
func (h *Handler) UpdateMyPassword(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	session := hCtx.ValidSession
	if session.UserID == 0 {
		slog.Warn("Unknown user accessing UpdateMyPassword")
		WriteFailureResponse("Not authenticated", w)
		return
	}

	req := &UpdateMyPasswordRequest{}
	err := ReadBodyAndValidate(r, req, UpdateMyPasswordValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.Any("error", err))
		WriteFailureResponse(err.Error(), w)
		return
	}

	dbh := hCtx.Database
	// verify that old password is correct
	user, err := dbh.GetUserByID(session.UserID)
	if err != nil {
		slog.Warn("Cannot find existing user by ID", slog.Any("error", err))
		WriteFailureResponse("password cannot be changed", w)
		return
	}
	hashedPassword, err := hash.CryptSha512(hash.Sha256(req.PasswordOld), strconv.Itoa(user.PasswordSalt))
	if err != nil {
		slog.Warn("Password hash could no be generated", slog.Any("error", err))
		WriteFailureResponse("password cannot be changed", w)
		return
	}
	if hashedPassword != user.PasswordHash {
		slog.Warn("Attempt to change password without correct old password")
		WriteFailureResponse("Password cannot be changed due to an incorrect old password.", w)
		return
	}

	// update new password
	salt, err := randomInt(math.MaxUint16)
	if err != nil {
		slog.Error("Cannot generate a password salt", slog.Any("error", err))
		WriteFailureResponse("password cannot be changed", w)
		return
	}
	newPasswordHash, err := hash.CryptSha512(hash.Sha256(req.PasswordNew), strconv.Itoa(salt))
	if err != nil {
		slog.Warn("Password hash could no be generated", slog.Any("error", err))
		WriteFailureResponse("password cannot be changed", w)
		return
	}
	user = database.User{
		PasswordSalt: salt,
		PasswordHash: newPasswordHash,
	}
	err = dbh.UpdatePassword(session.UserID, user)
	if err != nil {
		slog.Warn("Password could no be stored in user table", slog.Any("error", err))
		WriteFailureResponse("Password could not be changed.", w)
		return
	}
	WriteSuccessResponse("password changed", nil, w)
}

// GetMySessions returns the sessions of the currently logged in user, split into
// upcoming and past ones.
//
// It serves /api/v2/user/my/sessions and is open to authenticated users.
func (h *Handler) GetMySessions(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	session := hCtx.ValidSession
	if session.UserID == 0 {
		slog.Warn("Unknown user accessing GetMyHeats")
		WriteFailureResponse("Not authenticated", w)
		return
	}

	dbh := hCtx.Database
	sessions, err := dbh.GetSessionsByUser(session.UserID)
	if err != nil {
		slog.Error("Cannot get user sessions", slog.Any("error", err))
		WriteFailureResponse("Cannot get user sessions", w)
		return
	}

	// split into past and future
	upcoming := make([]MySessionResponse, 0, len(sessions))
	past := make([]MySessionResponse, 0, len(sessions))
	now := time.Now()
	for _, s := range sessions {
		// get riders for session
		users, err := dbh.GetUsersForSession(s.ID)
		if err != nil {
			slog.Error("Cannot get users for session", slog.Any("error", err))
		}
		sessionUsers := make([]UserShort, 0, len(users))
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

// GetMyHeats returns the rides of the currently logged in user.
//
// It serves /api/v2/user/my/heats and is open to authenticated users.
func (h *Handler) GetMyHeats(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	session := hCtx.ValidSession
	if session.UserID == 0 {
		slog.Warn("Unknown user accessing GetMyHeats")
		WriteFailureResponse("Not authenticated", w)
		return
	}

	dbh := hCtx.Database
	heats, err := dbh.GetUserHeats(session.UserID, 100)
	if err != nil {
		slog.Error("Cannot get user heats", slog.Any("error", err))
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

// GetMyHeatStats returns the total and year-to-date ride time and cost of the
// currently logged in user.
//
// It serves /api/v2/user/my/heats/statistics and is open to authenticated users.
func (h *Handler) GetMyHeatStats(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	session := hCtx.ValidSession
	if session.UserID == 0 {
		slog.Warn("Unknown user accessing GetMyHeatStats")
		WriteFailureResponse("Not authenticated", w)
		return
	}

	dbh := hCtx.Database
	now := time.Now()
	duration, cost, err := dbh.GetUserHeatStats(session.UserID, time.Time{}, now)
	if err != nil {
		slog.Error("Cannot get total duration")
		WriteFailureResponse("Cannot get stats from database", w)
		return
	}
	loc, err := dbh.GetTimezoneLocation()
	if err != nil {
		slog.Error("Cannot get location", slog.Any("error", err))
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

// GetMyBalance returns what the currently logged in user paid, got back, and still
// owes.
//
// It serves /api/v2/user/my/balance and is open to authenticated users.
func (h *Handler) GetMyBalance(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	session := hCtx.ValidSession
	if session.UserID == 0 {
		slog.Warn("Unknown user accessing GetMyBalance")
		WriteFailureResponse("Not authenticated", w)
		return
	}
	dbh := hCtx.Database
	resp, err := h.getUserBalance(dbh, session.UserID)
	if err != nil {
		slog.Error("Cannot get user balance", slog.Any("error", err))
		WriteFailureResponse("Cannot get user balance.", w)
		return
	}
	WriteSuccessResponse("balance", resp, w)
}

// GetAllUsersShort returns every user reduced to ID and name.
//
// It serves /api/v2/user/list-short and is open to administrators.
func (h *Handler) GetAllUsersShort(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	dbh := hCtx.Database
	users, err := dbh.GetUsers( /*includeDeleted=*/ false)
	if err != nil {
		slog.Error("Cannot get users", slog.Any("error", err))
		WriteFailureResponse("cannot get users", w)
		return
	}

	usersShort := make([]UserShort, 0, len(users))
	for _, u := range users {
		usersShort = append(usersShort, UserShort{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		})
	}
	WriteSuccessResponse("users", usersShort, w)
}

// GetAllUsersDetailed returns every user together with their ride and payment totals.
//
// It serves /api/v2/user/list-detailed and is open to administrators.
func (h *Handler) GetAllUsersDetailed(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	dbh := hCtx.Database
	users, err := dbh.GetUsers( /*includeDeleted=*/ false)
	if err != nil {
		slog.Error("Cannot get users", slog.Any("error", err))
		WriteFailureResponse("cannot get users", w)
		return
	}

	usersDetailed := make([]UserDetailed, 0, len(users))
	for _, u := range users {
		duration, cost, err := dbh.GetUserHeatStats(u.ID, time.Time{}, time.Now())
		if err != nil {
			slog.Error("Cannot get user heat duration and cost", slog.Any("error", err))
			WriteFailureResponse("cannot get heat duration and cost", w)
			return
		}
		paybacks, err := dbh.GetUserSessionPaybacks(u.ID)
		if err != nil {
			slog.Error("Cannot get user paybacks", slog.Any("error", err))
			WriteFailureResponse("cannot get user paybacks", w)
			return
		}
		payments, err := dbh.GetUserSessionPayments(u.ID)
		if err != nil {
			slog.Error("Cannot get user payments", slog.Any("error", err))
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

// GetUserGroups returns every user group with its pricing and role.
//
// It serves /api/v2/user/groups/get and is open to administrators.
func (h *Handler) GetUserGroups(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	dbh := hCtx.Database
	pricings, err := dbh.GetPricings()
	if err != nil {
		slog.Error("Cannot get user groups", slog.Any("error", err))
		WriteFailureResponse("cannot get user groups", w)
		return
	}

	userGroups := make([]GetUserGroupsResponse, 0, len(pricings))
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

// CreateUserGroup adds a user group together with its pricing.
//
// It serves /api/v2/user/group/create and is open to administrators.
func (h *Handler) CreateUserGroup(w http.ResponseWriter, r *http.Request, req CreateUserGroupsRequest, hCtx *HandlerCtx) {
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

	dbh := hCtx.Database
	err := dbh.CreateUserGroup(u, p)
	if err != nil {
		slog.Error("Cannot create new user group", slog.Any("error", err))
		WriteFailureResponse("cannot get create new user group", w)
		return
	}
	WriteSuccessResponse("user group created", nil, w)
}

// ChangeUserGroup updates a user group and its pricing.
//
// It serves /api/v2/user/group/edit and is open to administrators.
func (h *Handler) ChangeUserGroup(w http.ResponseWriter, r *http.Request, req ChangeUserGroupRequest, hCtx *HandlerCtx) {
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
	dbh := hCtx.Database
	err := dbh.ChangeUserGroup(u, p)
	if err != nil {
		slog.Error("Cannot update user group", slog.Any("error", err))
		WriteFailureResponse("Cannot update user group.", w)
		return
	}
	WriteSuccessResponse("user group changed", nil, w)
}

// DeleteUserGroup removes a user group. It fails while users are still in it.
//
// It serves /api/v2/user/group/delete and is open to administrators.
func (h *Handler) DeleteUserGroup(w http.ResponseWriter, r *http.Request, req DeleteUserGroupRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	err := dbh.DeleteUserGroup(req.UserGroupID)
	if err != nil {
		slog.Error("Cannot delete user group", slog.Any("error", err))
		WriteFailureResponse("Cannot delete user group.", w)
		return
	}
	WriteSuccessResponse("user group deleted", nil, w)
}

// SetUserGroup moves a user into a user group.
//
// It serves /api/v2/user/group/set and is open to administrators.
func (h *Handler) SetUserGroup(w http.ResponseWriter, r *http.Request, req SetUserGroupRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	user, err := dbh.GetUserByID(req.UserID)
	if err != nil {
		slog.Error("Cannot find user", slog.Any("error", err))
		WriteFailureResponse("Cannot find user.", w)
		return
	}
	userGroups, err := dbh.GetPricings()
	if err != nil {
		slog.Error("Cannot get valid groups", slog.Any("error", err))
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
		slog.Error("Cannot get admin users", slog.Any("error", err))
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
		slog.Error("Cannot change user group for user", slog.Int("user_id", int(req.UserID)), slog.Int("user_group_id", int(req.UserGroupID)), slog.Any("error", err))
		WriteFailureResponse("Cannot change user group for user.", w)
		return
	}
	WriteSuccessResponse("user group set", nil, w)
}

// GetUserRoles returns the available access levels.
//
// It serves /api/v2/user/roles/get and is open to administrators.
func (h *Handler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	dbh := hCtx.Database
	roles, err := dbh.GetUserRoles()
	if err != nil {
		slog.Error("Cannot get user roles", slog.Any("error", err))
		WriteFailureResponse("cannot get user roles", w)
		return
	}

	userRoles := make([]GetUserRolesResponse, 0, len(roles))
	for _, u := range roles {
		userRoles = append(userRoles, GetUserRolesResponse{
			UserRoleDescription: u.Description,
			UserRoleID:          u.ID,
			UserRoleName:        u.Name,
		})
	}
	WriteSuccessResponse("user roles", userRoles, w)
}

// SetUserLock locks or unlocks a user account.
//
// It serves /api/v2/user/lock/set and is open to administrators.
func (h *Handler) SetUserLock(w http.ResponseWriter, r *http.Request, req SetUserLockRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	user, err := dbh.GetUserByID(req.UserID)
	if err != nil {
		slog.Error("Cannot find user", slog.Any("error", err))
		WriteFailureResponse("Cannot find user.", w)
		return
	}

	// check that this user is not the only remaining admin
	admins, err := dbh.GetAdminUsers()
	if err != nil {
		slog.Error("Cannot get admin users", slog.Any("error", err))
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
		slog.Error("Cannot change lock status for user", slog.Int("user_id", int(req.UserID)), slog.Bool("locked", req.Locked), slog.Any("error", err))
		WriteFailureResponse("Cannot change user group for user.", w)
		return
	}
	WriteSuccessResponse("user lock set", nil, w)
}

func (h *Handler) getUserBalance(dbh database.Database, userID uint) (*GetMyBalanceResponse, error) {
	now := time.Now()
	_, cost, err := dbh.GetUserHeatStats(userID, time.Time{}, now)
	if err != nil {
		slog.Error("Cannot get total costs", slog.Int("user", int(userID)), slog.Any("error", err))
		return nil, errors.New("cannot get balance for user")
	}

	payment, err := dbh.GetUserSessionPayments(userID)
	if err != nil {
		slog.Error("Cannot get total payments for", slog.Int("user", int(userID)), slog.Any("error", err))
		return nil, errors.New("cannot get balance for user")
	}

	payback, err := dbh.GetUserSessionPaybacks(userID)
	if err != nil {
		slog.Error("Cannot get total paybacks for", slog.Int("user", int(userID)), slog.Any("error", err))
		return nil, errors.New("cannot get balance for user")
	}

	resp := &GetMyBalanceResponse{
		PaymentTotal: payment,
		PaybackTotal: payback,
		Balance:      payment.Sub(payback).Sub(cost),
	}
	return resp, nil
}

// GetPasswordResetToken emails a password reset token to the given address. An
// unknown address is still reported as a success, so that the endpoint cannot
// be used to probe for accounts.
//
// It serves /api/v2/user/password/token-request and is open to unauthenticated callers.
func (h *Handler) GetPasswordResetToken(w http.ResponseWriter, r *http.Request, req GetPasswordResetTokenRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	// Check whether recaptcha is enabled.
	privateKey, _ := h.config.GetString("recaptcha.privatekey")
	if privateKey == "" {
		slog.Warn("Password reset token requests should be protected by recaptcha. Please setup recaptcha in the settings.")
	} else {
		valid, err := recaptcha.Valid(req.RecaptchaToken, privateKey)
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
		slog.Warn("Cannot find user for password token request", slog.Any("error", err))
		WriteSuccessResponse(passwordResetRequested, nil, w)
		return
	}

	// Enforce the cooldown before touching the existing token.
	active, err := dbh.GetActivePasswordResetEntry(user.ID)
	if err == nil && time.Since(active.IssuedAt) < passwordResetCooldown {
		slog.Info("Password reset requested inside the cooldown window, keeping the live token",
			slog.Uint64("user_id", uint64(user.ID)))
		WriteSuccessResponse(passwordResetRequested, nil, w)
		return
	}

	// A user has at most one live token, so retire the previous ones before
	// issuing a replacement. Without this an attacker could stack tokens and
	// shrink the space they have to search by the number they hold.
	if err := dbh.InvalidatePasswordResetEntries(user.ID); err != nil {
		slog.Error("Cannot invalidate the previous password reset tokens", slog.Any("error", err))
		WriteFailureResponse("Internal error, cannot send reset token.", w)
		return
	}

	// Generate token.
	token, err := randomInt(passwordResetTokenSpace)
	if err != nil {
		slog.Error("Cannot generate a password reset token", slog.Any("error", err))
		WriteFailureResponse("Internal error, cannot send reset token.", w)
		return
	}
	now := time.Now()
	tokenEntry := database.PasswordReset{
		UserID:     user.ID,
		Token:      fmt.Sprintf("%06d", token),
		IssuedAt:   now,
		ValidUntil: now.Add(passwordResetValidity),
		Valid:      true,
	}

	// Store token in database.
	err = dbh.AddPasswordResetToken(tokenEntry)
	if err != nil {
		slog.Error("Cannot store password reset token", slog.Any("error", err))
		WriteFailureResponse("Internal error, cannot send reset token.", w)
		return
	}

	// Send email with token to user.
	emailConfig, err := dbh.GetEmailConfiguration()
	if err != nil {
		slog.Error("Cannot get email configuration to reset token", slog.Any("error", err))
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
		slog.Error("Cannot send email", slog.Any("error", err))
		WriteFailureResponse("Internal error, cannot send reset token.", w)
		return
	}
	WriteSuccessResponse("Token requested, please check your email inbox.", nil, w)
}

// SetPasswordWithToken sets a new password for a user that presents a valid reset token.
//
// It serves /api/v2/user/password/reset-by-token and is open to unauthenticated callers.
func (h *Handler) SetPasswordWithToken(w http.ResponseWriter, r *http.Request, req SetPasswordWithTokenRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	// Get user by email
	user, err := dbh.GetUserByName(req.UserEmail)
	if err != nil {
		slog.Warn("User cannot be found", slog.Any("error", err))
		WriteFailureResponse(passwordResetRejected, w)
		return
	}

	// Look the token up by user rather than by the value that was submitted. A
	// wrong guess matches no row when the value is part of the query, and there
	// would be nothing left to count the failed attempt against.
	dbToken, err := dbh.GetActivePasswordResetEntry(user.ID)
	if err != nil {
		slog.Warn("No live password reset token for this user", slog.Uint64("user_id", uint64(user.ID)), slog.Any("error", err))
		WriteFailureResponse(passwordResetRejected, w)
		return
	}

	// The query already filters on all three of these. Re-checking them here is
	// cheap insurance against a later change to that query, and it enforces the
	// attempt budget even if the write that burns a spent token failed.
	if !dbToken.Valid || dbToken.ValidUntil.Before(time.Now()) || dbToken.Attempts >= passwordResetMaxAttempts {
		slog.Warn("Password reset token is no longer usable",
			slog.Uint64("user_id", uint64(user.ID)),
			slog.Bool("valid", dbToken.Valid),
			slog.Uint64("attempts", uint64(dbToken.Attempts)))
		WriteFailureResponse(passwordResetRejected, w)
		return
	}

	// Compare in constant time and spend one of the attempts on a mismatch. The
	// token is only six digits wide, so the attempt budget is what keeps it out
	// of reach rather than its entropy.
	if subtle.ConstantTimeCompare([]byte(dbToken.Token), []byte(req.Token)) != 1 {
		slog.Warn("Wrong password reset token submitted",
			slog.Uint64("user_id", uint64(user.ID)),
			slog.Uint64("attempts_used", uint64(dbToken.Attempts+1)),
			slog.Int("attempts_allowed", passwordResetMaxAttempts))
		if err := dbh.RegisterFailedPasswordResetAttempt(dbToken.ID, passwordResetMaxAttempts); err != nil {
			slog.Error("Cannot record the failed password reset attempt", slog.Any("error", err))
		}
		WriteFailureResponse(passwordResetRejected, w)
		return
	}

	// Create new password hash and
	// update new password.
	// Note: We hash the passworrd, this was previously done in the UI.
	salt, err := randomInt(math.MaxUint16)
	if err != nil {
		slog.Error("Cannot generate a password salt", slog.Any("error", err))
		WriteFailureResponse("Internal error. Password cannot be changed.", w)
		return
	}
	newPasswordHash, err := hash.CryptSha512(hash.Sha256(req.Password), strconv.Itoa(salt))
	if err != nil {
		slog.Warn("Password hash could no be generated", slog.Any("error", err))
		WriteFailureResponse("Internal error. Password cannot be changed.", w)
		return
	}
	userPassword := database.User{
		PasswordSalt: salt,
		PasswordHash: newPasswordHash,
	}
	err = dbh.UpdatePassword(user.ID, userPassword)
	if err != nil {
		slog.Warn("Password could no be stored in user table", slog.Any("error", err))
		WriteFailureResponse("Internal error. Password could not be changed.", w)
		return
	}

	// Invalidate password reset tokens for this user. A failure here leaves the
	// just used token redeemable until it expires, so it must not stay silent.
	if err := dbh.InvalidatePasswordResetEntries(user.ID); err != nil {
		slog.Error("Cannot invalidate the used password reset token", slog.Uint64("user_id", uint64(user.ID)), slog.Any("error", err))
	}

	WriteSuccessResponse("password reset", nil, w)
}
