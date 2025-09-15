package handlers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/exp/slog"

	"server/database"
	"server/util/hash"
)

type LoginRequest struct {
	Password     string `json:"password" validate:"required"`
	PasswordHash string `json:"passwordHash" validate:"omitempty,sha256"`
	Username     string `json:"username" validate:"required,email|alphanum"`
}

type IsLoggedInResponse struct {
	LoggedIn bool `json:"loggedIn"`
}

type UserResponse struct {
	ID            uint                  `json:"id"`
	Username      string                `json:"username"`
	FirstName     string                `json:"first_name"`
	LastName      string                `json:"last_name"`
	Address       string                `json:"address"`
	City          string                `json:"city"`
	ZipCode       int                   `json:"plz"`
	MobilePhoneNr string                `json:"mobile"`
	Email         string                `json:"email"`
	BoatLicense   bool                  `json:"license"`
	Status        uint                  `json:"status"`
	UserRoleId    database.UserRoleType `json:"user_role_id"`
	UserRoleName  string                `json:"user_role_name"`
	Locked        bool                  `json:"locked"`
	Comment       string                `json:"comment"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request, req LoginRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	lookupUser, err := dbh.GetUserByName(req.Username)
	if err != nil {
		slog.Warn("User not found", slog.String("username", req.Username), slog.String("error", err.Error()))
		WriteFailureResponse("invalid username/password", w)
		return
	}

	// return if user is locked
	if lookupUser.Locked {
		WriteFailureResponse("user account not activated", w)
		return
	}

	// hash the password (legacy)
	// Note: This was previously done in the UI and now and
	// the hash was used in the crypt function.
	// The UI submitted the hex representation of the Sha256 hash
	// of the password to the backend.
	passwordHash := hash.Sha256(req.Password)

	// crypt the password
	hash, err := hash.CryptSha512(passwordHash, strconv.Itoa(lookupUser.PasswordSalt))
	if err != nil {
		slog.Warn("Cannot calculate password hash for", slog.String("user", lookupUser.Username))
		WriteFailureResponse("invalid username/password", w)
		return
	}

	// compare if password is identical
	if hash != lookupUser.PasswordHash {
		slog.Warn("Password not correct for", slog.String("user", lookupUser.Username))
		WriteFailureResponse("invalid username/password", w)
		return
	}

	// generate a browser session that we store in the sessions db
	u, err := dbh.GetUserById(lookupUser.ID)
	if err != nil {
		slog.Warn("Cannot retrieve user details for", slog.String("user", lookupUser.Username))
		WriteFailureResponse("invalid username/password", w)
		return
	}

	sessionSecret := make([]byte, 256)
	_, err = rand.Read(sessionSecret)
	if err != nil {
		slog.Error("Cannot generate random number for session secret", slog.String("error", err.Error()))
		WriteFailureResponse("invalid username/password", w)
		return
	}
	sessionSecretString := fmt.Sprintf("%x", sessionSecret)
	validUntil := time.Now().Add(time.Duration(h.config.GetInt64("http.sessioninactivitytimeout")) * time.Second)
	slog.Info("Create new session for", slog.String("user", u.Username), slog.String("validUntil", validUntil.String()))
	session := database.BrowserSession{
		SessionSecret: sessionSecretString,
		ValidUntil:    validUntil,
		LastActivity:  time.Now(),
		UserID:        u.ID,
		Username:      u.Username,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		UserStatus:    u.UserStatusID,
		UserRoleID:    u.UserStatus.UserRole.ID,
	}
	_, err = dbh.AddBrowserSession(session)
	if err != nil {
		slog.Error("Cannot add browser session to database", slog.String("error", err.Error()))
		WriteFailureResponse("invalid username/password", w)
		return
	}

	// set cookie and create response
	SetSessionCookie(w, sessionSecretString, validUntil)
	WriteSuccessResponse("login successful", nil, w)
}

func (h *Handler) IsLoggedIn(w http.ResponseWriter, r *http.Request) {
	resp := IsLoggedInResponse{}
	// Get cookie SESSION
	cookie, err := r.Cookie("SESSION")
	if err != nil {
		WriteSuccessResponse("not logged in", resp, w)
		DeleteSessionCookie(w)
		return
	}

	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database

	// Check if we know of that session and whether it is not expired yet
	session, err := dbh.GetBrowserSession(cookie.Value)
	if err != nil || !session.Valid() {
		WriteSuccessResponse("not logged in", resp, w)
		DeleteSessionCookie(w)
		return
	}

	// update valid until of browser session
	session.ValidUntil = time.Now().Add(time.Duration(h.config.GetInt64("http.sessioninactivitytimeout")) * time.Second)
	// TODO only update if creation time is not older than max session time

	resp.LoggedIn = true
	err = dbh.UpdateBrowserSession(*session)
	if err != nil {
		slog.Info("Cannot update browser session", slog.String("user", session.Username), slog.String("session", session.SessionSecret), slog.String("error", err.Error()))
	}
	WriteSuccessResponse("logged in", resp, w)
}

// Logout logs out a user
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database
	err = dbh.DeleteBrowserSession(*hCtx.ValidSession)
	if err != nil {
		slog.Error("Cannot delete browser session", slog.String("error", err.Error()))
	}

	DeleteSessionCookie(w)
	WriteSuccessResponse("logged out", nil, w)
	slog.Info("User logged out", slog.String("user", hCtx.ValidSession.Username))
}

func (h *Handler) User(w http.ResponseWriter, r *http.Request) {
	hCtx, err := GetHandlerContext(w, r)
	if err != nil {
		slog.Warn("Cannot get handler context", slog.String("error", err.Error()))
		return
	}
	dbh := hCtx.Database
	user, err := dbh.GetUserById(hCtx.ValidSession.UserID)
	if err != nil {
		slog.Error("Cannot retrieve user information", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot retrieve user informaiton", w)
		return
	}
	// build response
	resp := UserResponse{
		ID:            user.ID,
		Username:      user.Username,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Address:       user.Address,
		City:          user.City,
		ZipCode:       user.ZipCode,
		MobilePhoneNr: user.MobilePhoneNr,
		Email:         user.Email,
		BoatLicense:   user.BoatLicense,
		Status:        user.UserStatusID,
		UserRoleId:    user.UserStatus.UserRoleID,
		UserRoleName:  user.UserStatus.UserRole.Name,
		Locked:        user.Locked,
		Comment:       user.Comment,
	}
	WriteSuccessResponse("success", resp, w)
}
