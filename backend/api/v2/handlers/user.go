package handlers

import (
	"math"
	"math/rand"
	"net/http"
	"server/database"
	"strconv"

	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

type SignUpRequest struct {
	Username       string `json:"username" validate:"required,excludesall=!<>{}[]()^"`
	Password       string `json:"password" validate:"required"`
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
	"Password":       "a password has to be provided",
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
