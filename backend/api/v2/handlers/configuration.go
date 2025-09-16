package handlers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"server/database"
	"strconv"
	"time"

	"golang.org/x/exp/slog"
)

// type GetDBConfigResponse struct {
// 	IsConfigured bool   `json:"is_configured"`
// 	DBServer     string `json:"db_server"`
// 	DBName       string `json:"db_name"`
// 	DBUser       string `json:"db_user"`
// 	DBPassword   string `json:"db_password,omitempty"`
// }

// PublicConfigurationMessage represents the public configuration of the server.
// It is available to regular users that are logged in.
type PublicConfigurationMessage struct {
	Currency               string  `json:"currency" validate:"required,excludesall={} []!()<>"`
	EngineHourFormat       string  `json:"engine_hour_format" validate:"required,oneof=hh.h hh:mm"`    // check if really needed
	FuelPaymentType        string  `json:"fuel_payment_type" validate:"required,oneof=billed instant"` // check if really needed
	LocationAddress        string  `json:"location_address" validate:"excludesall={}[]!><"`
	LocationLatitude       float32 `json:"location_latitude" validate:"required,latitude"`
	LocationLongitude      float32 `json:"location_longitude" validate:"required,longitude"`
	LocationMap            string  `json:"location_map" validate:"omitempty,googlemapsurl"`
	LocationTimeZone       string  `json:"location_time_zone" validate:"required"`
	LogoFilePath           string  `json:"logo_file" validate:"omitempty,filepath"`                                                      // TODO: validator for uploaded file | check if really needed
	MyNautiqueEnabled      bool    `json:"mynautique_enabled" validate:"omitempty,boolean"`                                              // check if really needed
	MyNautiqueFuelCapacity int     `json:"mynautique_fuel_capacity" validate:"required_if=MyNautiqueEnabled true,omitempty,number,gt=1"` // check if really needed
	MyNautiqueBoatID       int     `json:"mynautique_boat_id" validate:"required_if=MyNautiqueEnabled true,omitempty,number,gt=1"`       // required to query mynautique telemetry
	PaymentAccountBIC      string  `json:"payment_account_bic" validate:"omitempty,printascii"`
	PaymentAccountComment  string  `json:"payment_account_comment"`
	PaymentAccountIBAN     string  `json:"payment_account_iban" validate:"omitempty,printascii"`
	PaymentAccountOwner    string  `json:"payment_account_owner"`
	RecaptchaPublicKey     string  `json:"recaptcha_publickey" validate:"omitempty,recaptchakey,required_with=RecaptchaPrivateKey"` // check if really needed
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
	LogoFilePath           string  `json:"logo_file" validate:"omitempty,filepath"` // TODO: validator for uploaded file
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
	SMTPPassword           string  `json:"smtp_password" validate:"omitempty"`
	SMTPSender             string  `json:"smtp_sender" validate:"omitempty,required_with=SMTPSender,email"`
	SMTPServer             string  `json:"smtp_server" validate:"required_with=SMTPSender"`
	SMTPUsername           string  `json:"smtp_username" validate:"required_with=SMTPSender"`
}

type ConfigSource int

const (
	SourceCLI ConfigSource = iota
	SourceDatabase
)

type StringConfigValue struct {
	Value  string
	Source ConfigSource
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
	"MyNautiqueAPIKey":       "The API key for the my Nautique App needs to be provided.",
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
	URI      string `json:"uri"`
	FileName string `json:"filename"`
}

type GetLogoPathResponse struct {
	URI string `json:"uri"`
}

type GetRecaptchaKeyResponse struct {
	Key string `json:"key"`
}

func (h *Handler) SetupDBConfig(w http.ResponseWriter, r *http.Request) {
	configFile, _ := h.config.GetString("config")
	if configFile == "" {
		slog.Warn("No config file path provided to store configuration. Database setup not possible.")
		WriteFailureResponse("No config file path provided to store configuration. Database setup not possible, please provide a configuration file (--config flag).", w)
		return
	}
	// We only want to allow setup for the DB, if DB settings are not
	// already present in the configuration.
	if h.config.IsDBConfigured() {
		slog.Warn("SetupDB called but database settings are already configured.")
		WriteFailureResponse("Invalid request. Database config was created already and cannot be overwritten that way.", w)
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

	// Connect to database and initialize.
	err = h.Database.ConnectAndReplace(&database.Settings{
		User:     req.DBUser,
		Password: req.DBPassword,
		Protocol: "tcp",
		Host:     host,
		Port:     port,
		DBName:   req.DBName,
	})
	if err != nil {
		slog.Warn(fmt.Sprintf("New database parameters are not valid. Error returned from Connet(): %s", err))
		WriteFailureResponse("Cannot connect to database. Please make sure that the credentials are correct and the database is accepting connections.", w)
		return
	}

	// If database access was successful, we store the configuration.
	h.config.SetConfigFileValue("database.user", req.DBUser)
	h.config.SetConfigFileValue("database.password", req.DBPassword)
	h.config.SetConfigFileValue("database.protocol", "tcp")
	h.config.SetConfigFileValue("database.host", host)
	h.config.SetConfigFileValue("database.port", port)
	h.config.SetConfigFileValue("database.dbname", req.DBName)
	err = h.config.WriteConfigFile()
	if err != nil {
		slog.Warn("Cannot store new configuration", slog.String("error", err.Error()))
		WriteFailureResponse("Database was setup. No config file path provided to store database configuration, please start application by providing a config file.", w)
		return
	}
	slog.Info("New database configuration has been written to", slog.String("config", configFile))
	// TODO: The database manager should just be calles with a ConnectAndReplace like above
	// and we should be good and not need this part here.
	h.config.SetDB(h.Database)
	WriteSuccessResponse("config written", nil, w)
}

func (h *Handler) GetPublicConfiguration(w http.ResponseWriter, r *http.Request) {
	currency, _ := h.config.GetString("currency")
	engineHourFormat, _ := h.config.GetString("engine.hour.format")
	fuelPaymentType, _ := h.config.GetString("fuel.payment.type")
	locationAddress, _ := h.config.GetString("location.address")
	lat, _ := h.config.GetString("location.latitude")
	lon, _ := h.config.GetString("location.longitude")
	lat32, _ := strconv.ParseFloat(lat, 32)
	lon32, _ := strconv.ParseFloat(lon, 32)
	locationMap, _ := h.config.GetString("location.map")
	locationTimeZone, _ := h.config.GetString("location.timezone")
	logoFilePath, _ := h.config.GetString("logo.file")
	myNautiqueEnabled := h.config.GetBool("mynautique.enabled")
	myNautiqueBoatID := h.config.GetInt64("mynautique.boat.id")
	myNautiqueFuelCapacity := h.config.GetInt64("mynautique.fuel.capacity")
	paymentAccountBIC, _ := h.config.GetString("payment.account.bic")
	paymentAccountComment, _ := h.config.GetString("payment.account.comment")
	paymentAccountIBAN, _ := h.config.GetString("payment.account.iban")
	paymentAccountOwner, _ := h.config.GetString("payment.account.owner")
	recaptchePublicKey, _ := h.config.GetString("recaptcha.publickey")
	resp := &PublicConfigurationMessage{
		Currency:               currency,
		EngineHourFormat:       engineHourFormat,
		FuelPaymentType:        fuelPaymentType,
		LocationAddress:        locationAddress,
		LocationLatitude:       float32(lat32),
		LocationLongitude:      float32(lon32),
		LocationMap:            locationMap,
		LocationTimeZone:       locationTimeZone,
		LogoFilePath:           logoFilePath,
		MyNautiqueEnabled:      myNautiqueEnabled,
		MyNautiqueBoatID:       int(myNautiqueBoatID),
		MyNautiqueFuelCapacity: int(myNautiqueFuelCapacity),
		PaymentAccountBIC:      paymentAccountBIC,
		PaymentAccountComment:  paymentAccountComment,
		PaymentAccountIBAN:     paymentAccountIBAN,
		PaymentAccountOwner:    paymentAccountOwner,
		RecaptchaPublicKey:     recaptchePublicKey,
	}
	WriteSuccessResponse("configuration", resp, w)
}

// TODO create public and non public version of this
// E.g. Public version should not contain private info like mynautique, db, ... things.
func (h *Handler) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	currency, _ := h.config.GetString("currency")
	engineHourFormat, _ := h.config.GetString("engine.hour.format")
	fuelPaymentType, _ := h.config.GetString("fuel.payment.type")
	locationAddress, _ := h.config.GetString("location.address")
	lat, _ := h.config.GetString("location.latitude")
	lon, _ := h.config.GetString("location.longitude")
	lat32, _ := strconv.ParseFloat(lat, 32)
	lon32, _ := strconv.ParseFloat(lon, 32)
	locationMap, _ := h.config.GetString("location.map")
	locationTimeZone, _ := h.config.GetString("location.timezone")
	logoFilePath, _ := h.config.GetString("logo.file")
	myNautiqueEnabled := h.config.GetBool("mynautique.enabled")
	myNautiqueBoatID := h.config.GetInt64("mynautique.boat.id")
	myNautiqueFuelCapacity := h.config.GetInt64("mynautique.fuel.capacity")
	myNautiqueUser, _ := h.config.GetString("mynautique.user")
	paymentAccountBIC, _ := h.config.GetString("payment.account.bic")
	paymentAccountComment, _ := h.config.GetString("payment.account.comment")
	paymentAccountIBAN, _ := h.config.GetString("payment.account.iban")
	paymentAccountOwner, _ := h.config.GetString("payment.account.owner")
	recaptchePublicKey, _ := h.config.GetString("recaptcha.publickey")
	recaptchePrivateKey, _ := h.config.GetString("recaptcha.privatekey")
	smtpSender, _ := h.config.GetString("smtp.sender")
	smtpServer, _ := h.config.GetString("smtp.server")
	smtpUsername, _ := h.config.GetString("smtp.username")

	resp := &ConfigurationMessage{
		Currency:               currency,
		EngineHourFormat:       engineHourFormat,
		FuelPaymentType:        fuelPaymentType,
		LocationAddress:        locationAddress,
		LocationLatitude:       float32(lat32),
		LocationLongitude:      float32(lon32),
		LocationMap:            locationMap,
		LocationTimeZone:       locationTimeZone,
		LogoFilePath:           logoFilePath,
		MyNautiqueBoatID:       int(myNautiqueBoatID),
		MyNautiqueEnabled:      myNautiqueEnabled,
		MyNautiqueFuelCapacity: int(myNautiqueFuelCapacity),
		MyNautiquePassword:     "hidden",
		MyNautiqueUser:         myNautiqueUser,
		PaymentAccountBIC:      paymentAccountBIC,
		PaymentAccountComment:  paymentAccountComment,
		PaymentAccountIBAN:     paymentAccountIBAN,
		PaymentAccountOwner:    paymentAccountOwner,
		RecaptchaPrivateKey:    recaptchePrivateKey,
		RecaptchaPublicKey:     recaptchePublicKey,
		SMTPPassword:           "hidden",
		SMTPSender:             smtpSender,
		SMTPServer:             smtpServer,
		SMTPUsername:           smtpUsername,
	}
	// TODO return if we have a myNautique API key. This is to
	// decide whether to show the myNautique settings section in the
	// UI.
	WriteSuccessResponse("configuration", resp, w)
}

func (h *Handler) SetConfiguration(w http.ResponseWriter, r *http.Request, req ConfigurationMessage, hCtx *HandlerCtx) {
	// TODO check smtp.password
	// - needs to be set in case there is other smtp configuration
	// - can be empty in case there is a password in the db already
	// - empty will not change the password in the db

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
			Property: "location.timezone",
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
	dbh := hCtx.Database
	err := dbh.UpdateOrInsertPropertyValues(props)
	if err != nil {
		slog.Warn("Cannot update configuration", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	WriteSuccessResponse("config updated", nil, w)
}

func (h *Handler) SetupMyNautiqueCredentials(w http.ResponseWriter, r *http.Request, req SetupMyNautiqueCredentialsRequest, _ *HandlerCtx) {
	// in case the configuration is already present, we do not
	// allow to edit it
	if h.config.IsSet("mynautique.enabled") {
		slog.Warn("mynautique.enabled is already set, not allowing call to setup my nautique credentials.")
		return
	}

	err := h.config.SetPropertyValue("mynautique.enabled", fmt.Sprintf("%t", req.Enabled))
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid request payload", w)
		return
	}
	err = h.config.SetPropertyValue("mynautique.user", req.User)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid request payload", w)
		return
	}
	err = h.config.SetPropertyValue("mynautique.password", req.Password)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse("Invalid request payload", w)
		return
	}

	WriteSuccessResponse("config written", nil, w)
}

func (h *Handler) UploadLogoFile(w http.ResponseWriter, r *http.Request) {
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

	// store uploaded file into local path. We use a content based filename
	// to store multiple files without collisions.
	nameTime := fmt.Sprintf("%s_%s", time.Now().Format(time.RFC3339), fileHeader.Filename)
	hash := sha256.Sum256([]byte(nameTime))
	fileHash := fmt.Sprintf("%x%s", hash[:16], filepath.Ext(fileHeader.Filename))
	storageDir, _ := h.config.GetString("http.uploadpath")
	localFileName := filepath.Join(storageDir, fileHash)
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

	resp := &UploadLogoFileResponse{
		URI:      localFileName,
		FileName: fileHash,
	}
	WriteSuccessResponse("file uploaded", resp, w)
}

func (h *Handler) GetLogoPath(w http.ResponseWriter, r *http.Request) {
	conf, _ := h.config.GetString("logo.file")
	uploadDir, _ := h.config.GetString("http.uploadpath")
	slog.Warn("DEBUG: uploadDir", slog.String("dir", uploadDir))
	resp := &GetLogoPathResponse{}
	if conf != "" {
		resp.URI = filepath.Join(uploadDir, conf)
	}
	WriteSuccessResponse("logo path", resp, w)
}

// TODO implement file removal

func (h *Handler) GetRecaptchaKey(w http.ResponseWriter, r *http.Request) {
	key, _ := h.config.GetString("recaptcha.publickey")
	resp := &GetRecaptchaKeyResponse{}
	resp.Key = key
	WriteSuccessResponse("recaptcha key", resp, w)
}
