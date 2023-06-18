package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"time"

	// "github.com/GehirnInc/crypt/sha512_crypt"

	"github.com/spf13/viper"
	crypt "github.com/tredoe/osutil/v2/userutil/crypt/sha512_crypt"
	"golang.org/x/exp/slog"

	"server/database"
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
	ID            uint   `json:"id"`
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Address       string `json:"address"`
	City          string `json:"city"`
	ZipCode       int    `json:"plz"`
	MobilePhoneNr string `json:"mobile"`
	Email         string `json:"email"`
	BoatLicense   bool   `json:"license"`
	Status        uint   `json:"status"`
	UserRoleId    uint   `json:"user_role_id"`
	UserRoleName  string `json:"user_role_name"`
	Locked        bool   `json:"locked"`
	Comment       string `json:"comment"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := ReadBodyAndValidate(r, &req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid request payload", w)
		return
	}

	// get user from database
	lookupUser, err := h.GetDB().GetUserByName(req.Username)
	if err != nil {
		slog.Warn("User not found", slog.String("username", req.Username))
		WriteFailureResponse("invalid username/password", w)
		return
	}

	// return if user is locked
	if lookupUser.Locked {
		WriteFailureResponse("user account is locked", w)
		return
	}

	// hash the password (legacy)
	// Note: This was previously done in the UI and now and
	// the hash was used in the crypt function.
	passwordHash := hashSha256(req.Password)

	// crypt the password
	hash, err := cryptSha512(passwordHash, strconv.Itoa(lookupUser.PasswordSalt))
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
	u, err := h.GetDB().GetUserById(lookupUser.ID)
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
	validUntil := time.Now().Add(time.Duration(viper.GetInt64("http.sessioninactivitytimeout")) * time.Second)
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
	_, err = h.GetDB().AddBrowserSession(session)
	if err != nil {
		slog.Error("Cannot add browser session to database", slog.String("error", err.Error()))
		WriteFailureResponse("invalid username/password", w)
		return
	}

	// set cookie and create response
	SetSessionCookie(w, sessionSecretString)
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

	// Check if we know of that session and whether it is not expired yet
	session, err := h.GetDB().GetBrowserSession(cookie.Value)
	if err != nil || !session.Valid() {
		WriteSuccessResponse("not logged in", resp, w)
		DeleteSessionCookie(w)
		return
	}

	// update valid until of browser session
	session.ValidUntil = time.Now().Add(time.Duration(viper.GetInt64("http.sessioninactivitytimeout")) * time.Second)
	// TODO only update if creation time is not older than max session time

	resp.LoggedIn = true
	err = h.GetDB().UpdateBrowserSession(*session)
	if err != nil {
		slog.Info("Cannot update browser session", slog.String("user", session.Username), slog.String("session", session.SessionSecret), slog.String("error", err.Error()))
	}
	WriteSuccessResponse("logged in", resp, w)
}

// Logout logs out a user
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedOrFailure(session, w) != nil {
		slog.Warn("Logout without authentication rejected")
		return
	}

	err := h.GetDB().DeleteBrowserSession(session)
	if err != nil {
		slog.Error("Cannot delete browser session", slog.String("error", err.Error()))
	}

	DeleteSessionCookie(w)
	WriteSuccessResponse("logged out", nil, w)
	slog.Info("User logged out", slog.String("user", session.Username))
}

func (h *Handler) User(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedOrFailure(session, w) != nil {
		return
	}

	user, err := h.GetDB().GetUserById(session.UserID)
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

// hashPassword returns a sha256 hashed version of the given string
func hashSha256(p string) string {
	hasher := sha256.New()
	hasher.Write([]byte(p))
	hash := string(hasher.Sum(nil)[:])
	return fmt.Sprintf("%x", hash)
}

// cryptSha512 generates a sha512 hashes password representation that
// uses 5000 rounds and the given salt.
func cryptSha512(p string, s string) (string, error) {
	c := crypt.New()
	cryptConfig := fmt.Sprintf("$6$rounds=5000$%s$", s)
	return c.Generate([]byte(p), []byte(cryptConfig))
}
