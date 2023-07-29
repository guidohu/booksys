package handlers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"server/config"
	"server/database"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

type GetDBConfigResponse struct {
	IsConfigured bool   `json:"is_configured"`
	DBServer     string `json:"db_server"`
	DBName       string `json:"db_name"`
	DBUser       string `json:"db_user"`
	DBPassword   string `json:"db_password,omitempty"`
}

type ConfigurationMessage struct {
	Currency               string  `json:"currency" validate:"required,excludesall={} []!()<>"`
	EngineHourFormat       string  `json:"engine_hour_format" validate:"required,oneof=hh.h hh:mm"`
	FuelPaymentType        string  `json:"fuel_payment_type" validate:"required,oneof=billed instant"`
	LocationAddress        string  `json:"location_address" validate:"excludesall={}[]!><"`
	LocationLatitude       float32 `json:"location_latitude" validate:"required,latitude"`
	LocationLongitude      float32 `json:"location_longitude" validate:"required,longitude"`
	LocationMap            string  `json:"location_map" validate:"omitempty,googlemapsurl"`
	LocationTimeZone       string  `json:"location_time_zone" validate:"required"`
	LogoFilePath           string  `json:"logo_file" validate:"omitempty,file"`
	MyNautiqueBoatID       int     `json:"mynautique_boat_id" validate:"required_if=MyNautiqueEnabled true,omitempty,number,gt=10"`
	MyNautiqueEnabled      bool    `json:"mynautique_enabled" validate:"omitempty,boolean"`
	MyNautiqueFuelCapacity int     `json:"mynautique_fuel_capacity" validate:"required_if=MyNautiqueEnabled true,omitempty,number,gt=10"`
	MyNautiquePassword     string  `json:"mynautique_password" validate:"required_if=MyNautiqueEnabled true,omitempty,gt=1"`
	MyNautiqueUser         string  `json:"mynautique_user" validate:"required_if=MyNautiqueEnabled true,omitempty,email"`
	PaymentAccountBIC      string  `json:"payment_account_bic" validate:"omitempty,printascii"`
	PaymentAccountComment  string  `json:"payment_account_comment"`
	PaymentAccountIBAN     string  `json:"payment_account_iban" validate:"omitempty,printascii"`
	PaymentAccountOwner    string  `json:"payment_account_owner"`
	RecaptchaPrivateKey    string  `json:"recaptcha_privatekey" validate:"omitempty,recaptchakey,required_with=RecaptchaPublicKey"`
	RecaptchaPublicKey     string  `json:"recaptcha_publickey" validate:"omitempty,recaptchakey,required_with=RecaptchaPrivateKey"`
	SMTPPassword           string  `json:"smtp_password" validate:"required_with=SMTPSender"`
	SMTPSender             string  `json:"smtp_sender" validate:"omitempty,required_with=SMTPSender,email"`
	SMTPServer             string  `json:"smtp_server" validate:"required_with=SMTPSender"`
	SMTPUsername           string  `json:"smtp_username" validate:"required_with=SMTPSender"`
}

var ConfigurationMessageValidationErrors = map[string]string{
	"Currency":               "Use the 3 letter currency representation. E.g., USD, EUR, CHF.",
	"EngineHourFormat":       "Engine hour format can only be hh.m or hh:mm.",
	"FuelPaymentType":        "Fuel payment type is either 'instant' or 'billed'",
	"LocationAddress":        "Location address cannot contain invalid characters such as []{}<> or similar.",
	"LocationLatitude":       "Latitude is not a valid value.",
	"LocationLongitude":      "Longitude is not a valid value.",
	"LocationMap":            "Map URL needs to be a google embeded maps URL of the form https://www.google.com/maps/embeded?pb=...",
	"LocationTimeZone":       "The timezone needs to be a valid representation such as Europe/Zurich.",
	"LogoFilePath":           "Needs to be the path to the logo file on the server.",
	"MyNautiqueBoatID":       "The boat ID from the my Nautique App should be a number.",
	"MyNautiqueEnabled":      "MyNautique enabled needs to be true or false",
	"MyNautiqueFuelCapacity": "The MyNautique fuel capacity needs to be a number.",
	"MyNautiqueUser":         "The user for MyNautique needs to be an Email address.",
	"MyNautiquePassword":     "A password for the MyNautique user needs to be set.",
	"PaymentAccountBIC":      "BIC needs to be set to a valid value.",
	"PaymentAccountComment":  "Comment cannot be parsed.",
	"PaymentAccountIBAN":     "The IBAN number needs to be valid.",
	"PaymentAccountOwner":    "The account owner cannot be parsed.",
	"RecaptchaPrivateKey":    "Recaptcha private key is not valid.",
	"RecaptchaPublicKey":     "Recaptcha public key is not valid.",
	"SMTPPassword":           "SMTP Password needs to be set.",
	"SMTPSender":             "SMTP sender is not a valid email address.",
	"SMTPServer":             "SMTP server is not a valid  address.",
	"SMTPUsername":           "No or invalid SMTP username provided",
}

type SetupDBConfigRequest struct {
	DBServer   string `json:"db_server" validate:"required,hostname_port"`
	DBName     string `json:"db_name" validate:"required,alphanum"`
	DBUser     string `json:"db_user" validate:"required,alphanum"`
	DBPassword string `json:"db_password" validate:"required"`
}

type SetupMyNautiqueCredentialsRequest struct {
	Enabled  bool   `json:"mynautique_enabled" validate:"boolean"`
	User     string `json:"mynautique_user" validate:"omitempty,email"`
	Password string `json:"mynautique_password"`
}

type UploadLogoFileResponse struct {
	URI string `json:"uri"`
}

type GetLogoPathResponse struct {
	URI string `json:"uri"`
}

func (h *Handler) SetupDBConfig(w http.ResponseWriter, r *http.Request) {
	if config.IsDBConfigured() {
		slog.Warn("SetupDB called for already setup DB")
		WriteFailureResponse("Invalid request", w)
		return
	}

	var req SetupDBConfigRequest
	err := ReadBodyAndValidate(r, &req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid request payload", w)
		return
	}

	host, port, err := net.SplitHostPort(req.DBServer)
	if err != nil {
		slog.Warn("Request payload is not valid, invalid DBServer", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid hostname (should be host:port)", w)
		return
	}

	// test db access
	db := &database.DBMysql{
		User:     req.DBUser,
		Password: req.DBPassword,
		Protocol: "tcp",
		Host:     host,
		Port:     port,
		DBName:   req.DBName,
	}
	err = db.Connect()
	defer db.Disconnect()
	if err != nil {
		slog.Warn(fmt.Sprintf("New database parameters are not valid. Error returned from Connet(): %s", err))
		WriteFailureResponse("Cannot connect to database.", w)
		return
	}

	// if database access was successful, we store the configurationn
	viper.Set("database.user", req.DBUser)
	viper.Set("database.password", req.DBPassword)
	viper.Set("database.protocol", "tcp")
	viper.Set("database.host", host)
	viper.Set("database.port", port)
	viper.Set("database.dbname", req.DBName)
	err = viper.WriteConfig()
	if err != nil {
		slog.Warn("Cannot store new configuration", slog.String("error", err.Error()))
		WriteFailureResponse("cannot write config", w)
		return
	}
	slog.Info("New database configuration has been written to", slog.String("configfile", viper.GetString("configfile")))

	// Postprocess config change
	h.ReconnectDB()

	WriteSuccessResponse("config written", nil, w)
}

func (h *Handler) GetDBConfig(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	resp := GetDBConfigResponse{
		IsConfigured: true,
		DBServer:     fmt.Sprintf("%s:%d", viper.GetString("database.host"), viper.GetUint("database.port")),
		DBName:       viper.GetString("database.dbname"),
		DBUser:       viper.GetString("database.user"),
		DBPassword:   "",
	}
	WriteSuccessResponse("success", resp, w)
}

func (h *Handler) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if !session.Valid() {
		slog.Warn("Call to GetConfiguration without authentication")
		WriteFailureResponse("Not authenticated", w)
		return
	}

	properties, err := h.GetDB().GetAllPropertyValues()
	if err != nil {
		slog.Error("Cannot get configuration properties from database", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get configuration", w)
		return
	}

	pMap := make(map[string]string)
	for _, p := range properties {
		pMap[p.Property] = p.Value
	}

	lat, err := strconv.ParseFloat(pMap["location.latitude"], 32)
	if err != nil {
		slog.Error("Cannot convert location.latitude to float", slog.String("error", err.Error()))
		lat = 0
	}
	lon, err := strconv.ParseFloat(pMap["location.longitude"], 32)
	if err != nil {
		slog.Error("Cannot convert location.longitude to float", slog.String("error", err.Error()))
		lon = 0
	}
	boatid, err := strconv.Atoi(pMap["mynautique.boat.id"])
	if pMap["mynautique.boat.id"] == "" {
		boatid = 0
	} else if err != nil {
		slog.Error("Cannot convert mynautique.boat.id to int", slog.String("error", err.Error()))
		boatid = 0
	}
	mynautiqueEnabled, err := strconv.ParseBool(pMap["mynautique.enabled"])
	if pMap["mynautique.enabled"] == "" {
		mynautiqueEnabled = false
	} else if err != nil {
		slog.Error("Cannot convert mynautique.enabled to bool", slog.String("error", err.Error()))
		mynautiqueEnabled = false
	}
	mynautiqueFuelCapacity, err := strconv.Atoi(pMap["mynautique.fuel.capacity"])
	if pMap["mynautique.fuel.capacity"] == "" {
		mynautiqueFuelCapacity = 0
	} else if err != nil {
		slog.Error("Cannot convert boat.fuel.capacity to int", slog.String("error", err.Error()))
		mynautiqueFuelCapacity = 0
	}
	resp := &ConfigurationMessage{
		Currency:               pMap["currency"],
		EngineHourFormat:       pMap["engine.hour.format"],
		FuelPaymentType:        pMap["fuel.payment.type"],
		LocationAddress:        pMap["location.address"],
		LocationLatitude:       float32(lat),
		LocationLongitude:      float32(lon),
		LocationMap:            pMap["location.map"],
		LocationTimeZone:       pMap["location.timezone"],
		LogoFilePath:           pMap["logo.file"],
		MyNautiqueBoatID:       boatid,
		MyNautiqueEnabled:      mynautiqueEnabled,
		MyNautiqueFuelCapacity: mynautiqueFuelCapacity,
		MyNautiquePassword:     "hidden",
		MyNautiqueUser:         pMap["mynautique.user"],
		PaymentAccountBIC:      pMap["payment.account.bic"],
		PaymentAccountComment:  pMap["payment.account.comment"],
		PaymentAccountIBAN:     pMap["payment.account.iban"],
		PaymentAccountOwner:    pMap["payment.account.owner"],
		RecaptchaPrivateKey:    pMap["recaptcha.privatekey"],
		RecaptchaPublicKey:     pMap["recaptcha.publickey"],
		SMTPPassword:           "hidden",
		SMTPSender:             pMap["smtp.sender"],
		SMTPServer:             pMap["smtp.server"],
		SMTPUsername:           pMap["smtp.username"],
	}
	WriteSuccessResponse("configuration", resp, w)
}

func (h *Handler) SetConfiguration(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	var req ConfigurationMessage
	err := ReadBodyAndValidate(r, &req, ConfigurationMessageValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	props := []database.Configuration{
		{
			Property: "currency",
			Value:    req.Currency,
		},
		{
			Property: "engine.hour.format",
			Value:    req.EngineHourFormat,
		},
		{
			Property: "fuel.payment.type",
			Value:    req.FuelPaymentType,
		},
		{
			Property: "location.address",
			Value:    req.LocationAddress,
		},
		{
			Property: "location.latitude",
			Value:    strconv.FormatFloat(float64(req.LocationLatitude), 'f', 6, 32),
		},
		{
			Property: "location.longitude",
			Value:    strconv.FormatFloat(float64(req.LocationLongitude), 'f', 6, 32),
		},
		{
			Property: "location.map",
			Value:    req.LocationMap,
		},
		{
			Property: "location.time.zone",
			Value:    req.LocationTimeZone,
		},
		{
			Property: "logo.file",
			Value:    req.LogoFilePath,
		},
		{
			Property: "mynautique.boat.id",
			Value:    strconv.FormatInt(int64(req.MyNautiqueBoatID), 10),
		},
		{
			Property: "mynautique.enabled",
			Value:    strconv.FormatBool(req.MyNautiqueEnabled),
		},
		{
			Property: "mynautique.fuel.capacity",
			Value:    strconv.FormatInt(int64(req.MyNautiqueFuelCapacity), 10),
		},
		{
			Property: "mynautique.password",
			Value:    req.MyNautiquePassword,
		},
		{
			Property: "mynautique.user",
			Value:    req.MyNautiqueUser,
		},
		{
			Property: "payment.account.bic",
			Value:    req.PaymentAccountBIC,
		},
		{
			Property: "payment.account.comment",
			Value:    req.PaymentAccountComment,
		},
		{
			Property: "payment.account.iban",
			Value:    req.PaymentAccountIBAN,
		},
		{
			Property: "payment.account.owner",
			Value:    req.PaymentAccountOwner,
		},
		{
			Property: "recaptcha.privatekey",
			Value:    req.RecaptchaPrivateKey,
		},
		{
			Property: "recaptcha.publickey",
			Value:    req.RecaptchaPublicKey,
		},
		{
			Property: "smtp.password",
			Value:    req.SMTPPassword,
		},
		{
			Property: "smtp.sender",
			Value:    req.SMTPSender,
		},
		{
			Property: "smtp.server",
			Value:    req.SMTPServer,
		},
		{
			Property: "smtp.username",
			Value:    req.SMTPUsername,
		},
	}
	err = h.GetDB().UpdateOrInsertPropertyValues(props)
	if err != nil {
		slog.Warn("Cannot update configuration", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	WriteSuccessResponse("config updated", nil, w)
}

func (h *Handler) SetupMyNautiqueCredentials(w http.ResponseWriter, r *http.Request) {
	// in case the configuration is already present, we do not
	// allow to edit it
	session := GetSessionFromContext(r)
	if viper.IsSet("mynautique.enabled") && AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	var req SetupMyNautiqueCredentialsRequest
	err := ReadBodyAndValidate(r, &req)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid request payload", w)
		return
	}

	viper.Set("mynautique.enabled", req.Enabled)
	viper.Set("mynautique.user", req.User)
	viper.Set("mynautique.password", req.Password)

	err = viper.WriteConfig()
	if err != nil {
		slog.Warn("Cannot store new configuration", slog.String("error", err.Error()))
		WriteFailureResponse("cannot write config", w)
		return
	}
	slog.Info("New mynautique configuration has been written to", slog.String("configfile", viper.GetString("configfile")))

	// Postprocess config change
	h.ReconnectDB()

	WriteSuccessResponse("config written", nil, w)
}

func (h *Handler) UploadLogoFile(w http.ResponseWriter, r *http.Request) {
	// only admins are supposed to upload a logo file
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}

	file, fileHeader, err := r.FormFile("logo")
	if err != nil {
		slog.Warn("cannot access the form file", slog.String("position", "logo"))
		WriteFailureResponse("Cannot read the uploaded file.", w)
		return
	}
	defer file.Close()

	// check for the max logo size to be 512kB
	if fileHeader.Size > 512*1024 {
		slog.Warn("attempt to upload a file larger than 512kB", slog.Int64("size", fileHeader.Size))
		WriteFailureResponse("File is larger than 512kB.", w)
		return
	}

	contentTypes := fileHeader.Header.Values("Content-Type")
	if len(contentTypes) == 0 {
		slog.Warn("Cannot read Content-Type of file")
		WriteFailureResponse("The uploaded file needs to be an image (e.g., Content-Type image/jpeg)", w)
		return
	}
	contentType := contentTypes[0]
	if contentType != "image/png" &&
		contentType != "image/jpeg" &&
		contentType != "image/jpg" &&
		contentType != "image/gif" {
		slog.Warn("Wrong Content-Type of file", slog.String("content_type", contentType))
		WriteFailureResponse("The uploaded file needs to be an image (e.g., Content-Type image/jpeg)", w)
		return
	}

	// store uploaded file into local path. We use a random filename
	// to store multiple files without collisions.
	nameTime := fmt.Sprintf("%s_%s", time.Now().Format(time.RFC3339), fileHeader.Filename)
	hash := sha256.Sum256([]byte(nameTime))
	fileHash := fmt.Sprintf("%x%s", hash[:16], filepath.Ext(fileHeader.Filename))
	localFileName := filepath.Join(viper.GetString("upload.path"), fileHash)
	if err := os.MkdirAll(filepath.Dir(localFileName), 0770); err != nil {
		slog.Warn("Cannot create directory for file", slog.String("file", localFileName), slog.String("error", err.Error()))
		WriteFailureResponse("File cannot get stored on server.", w)
		return
	}
	out, err := os.Create(localFileName)
	if err != nil {
		slog.Warn("Cannot create file on server for file", slog.String("file", localFileName), slog.String("error", err.Error()))
		WriteFailureResponse("File cannot get stored on server.", w)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		slog.Warn("Cannot write to file on server for file", slog.String("file", localFileName), slog.String("error", err.Error()))
		WriteFailureResponse("File cannot get stored on server.", w)
		return
	}

	err = h.GetDB().UpdateOrInsertPropertyValues([]database.Configuration{
		{
			Property: "logo.file",
			Value:    localFileName,
		},
	})
	if err != nil {
		slog.Warn("Cannot write file to configuration database", slog.String("file", localFileName), slog.String("error", err.Error()))
		WriteFailureResponse("File cannot get stored on server.", w)
		return
	}

	resp := &UploadLogoFileResponse{
		URI: localFileName,
	}
	WriteSuccessResponse("file uploaded", resp, w)
}

func (h *Handler) GetLogoPath(w http.ResponseWriter, r *http.Request) {
	conf, err := h.GetDB().GetPropertyValue("logo.file")
	if err != nil {
		slog.Warn("Cannot get logo file", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get logo path from server.", w)
		return
	}
	resp := &GetLogoPathResponse{}
	resp.URI = conf.Value
	WriteSuccessResponse("logo path", resp, w)
}

// TODO implement file removal
