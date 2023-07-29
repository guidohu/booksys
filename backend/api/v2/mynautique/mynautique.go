package mynautique

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

var api_url = "https://mynautique.azurewebsites.net/api/v2"
var auth_url = "https://identitytoolkit.googleapis.com/v1"

type Client struct {
	auth      loginResponse
	AuthUntil time.Time
	Fleet     []BoatInfo
	Options   *Options
}

type Options struct {
	User       string
	Password   string
	AuthAPIKey string
	Client     *http.Client
}

type loginResponse struct {
	IDToken          string `json:"idToken"`
	ExpiresIn        string `json:"expiresIn"`
	ExpiresInSeconds int64
	Kind             string `json:"kind"`
	LocalId          string `json:"localId"`
	Email            string `json:"email"`
	DisplayName      string `json:"displayName"`
	Registered       bool   `json:"registered"`
	RefreshToken     string `json:"refreshToken"`
}

type BoatInfo struct {
	Owner             string
	ModelYear         int64
	Date              time.Time
	EngineSerial      string
	DealerNumber      int64
	FullHin           string
	Model             string
	DeliveryDate      time.Time
	DeliveryEntryDate time.Time
	Engine            string
	GPSEnabled        bool
	Hin               int64
}

type boatInfoRaw struct {
	Owner                   string `json:"owner"`
	ModelYear               string `json:"myear"`
	Date                    sDate  `json:"sDate"`
	EngineSerial            string `json:"EngineSerial"`
	DealerNumber            int64  `json:"dealer"`
	FullHin                 string `json:"fullHin"`
	Model                   string `json:"model"`
	DeliveryDateString      string `json:"deliveryDate"`
	DeliveryEntryDateString string `json:"deliveryEntryDate"`
	Engine                  string `json:"Engine"`
	GPSEnabled              bool   `json:"gpsEnabled"`
	Hin                     string `json:"hin"`
}

type Telemetry struct {
	DeviceSerial                  int64
	ServiceMessage                string
	ServiceDetailsMessage         string
	GPSLongitude                  decimal.Decimal
	GPSLatitude                   decimal.Decimal
	IsWebSocketConnected          bool
	AccelXYZMagnitude             decimal.Decimal
	BallastBelly                  int64
	BallastPortSupplementary      int64
	BallastStarboardSupplementary int64
	BallastPortRear               int64
	BallastStarboardRear          int64
	BoatStatus                    int64
	DigitalInput4                 int64
	EngineHoursLinc               decimal.Decimal
	FuelLevelLinc                 int64
	GPSAltitude                   int64
	GPSCourse                     decimal.Decimal
	GSMSignalQualityBars          int64
	TemperatureCelsius            decimal.Decimal
	GPSQuality                    int64
	GPSSpeed                      decimal.Decimal
	In1Ana                        decimal.Decimal
	In2Ana                        decimal.Decimal
	In3Ana                        decimal.Decimal
	IsSleeping                    bool
	LastSensorUpdate              time.Time
	EngineSpeed                   decimal.Decimal
	GPSSpeedKmh                   decimal.Decimal
	GPSSpeedMph                   decimal.Decimal
	Timestamp                     time.Time
	GPSNumSatellites              int64
	AirTemperatureLinc            decimal.Decimal
	EngineTotalHoursOfOperation   decimal.Decimal
	LincDate                      time.Time
	WaterDepth                    decimal.Decimal
	WaterTemperature              decimal.Decimal
	FuelLevel1                    decimal.Decimal
	GSMConnected                  int64
	GSMSignalQualityDBm           int64
	ServiceDateCapture            int64
	ServiceHoursCapture           decimal.Decimal
	ServiceReminderEnabled        bool
}

type telemetryRaw struct {
	DeviceSerial                  int64   `json:"deviceSerial"`
	ServiceMessage                string  `json:"serviceMessage"`
	ServiceDetailsMessage         string  `json:"serviceDetailsMessage"`
	GPSLongitude                  float64 `json:"gps_long"`
	GPSLatitude                   float64 `json:"gps_lat"`
	IsWebSocketConnected          bool    `json:"isWebSocketConnected"`
	AccelXYZMagnitude             string  `json:"Accel_XYZ_Magnitude"`
	BallastBelly                  string  `json:"BALLAST_BELLY"`
	BallastPortSupplementary      string  `json:"BALLAST_PORT_SUPP"`
	BallastStarboardSupplementary string  `json:"BALLAST_STBD_SUPP"`
	BallastPortRear               string  `json:"BALLAST_REAR_PORT"`
	BallastStarboardRear          string  `json:"BALLAST_REAR_STBD"`
	BoatStatus                    string  `json:"BOAT_STATUS"`
	DigitalInput4                 string  `json:"Digital_Input_4"`
	EngineHoursLinc               string  `json:"ENGINE_HOURS_LINC"`
	FuelLevelLinc                 string  `json:"FUEL_LEVEL_LINC"`
	GPSAltitude                   string  `json:"GPS_Altitude_m"`
	GPSCourse                     string  `json:"GPS_Course"`
	GSMSignalQualityBars          string  `json:"GSM_Signal_Quality_bars"`
	TemperatureCelsius            string  `json:"Temp_C"`
	GPSQuality                    string  `json:"gps_qual"`
	GPSSpeed                      string  `json:"gps_speed"`
	In1Ana                        string  `json:"in1_ana"`
	In2Ana                        string  `json:"in2_ana"`
	In3Ana                        string  `json:"in3_ana"`
	IsSleeping                    bool    `json:"isSleeping"`
	LastSensorUpdate              int64   `json:"lastSensorUpdate"`
	EngineSpeed                   string  `json:"EngineSpeed"`
	GPSSpeedKmh                   string  `json:"GPS_Speed_kmh"`
	GPSSpeedMph                   string  `json:"GPS_Speed_mph"`
	Timestamp                     int64   `json:"ts"`
	GPSNumSatellites              string  `json:"GPS_Num_Satellites"`
	AirTemperatureLinc            string  `json:"AIR_TEMP_LINC"`
	EngineTotalHoursOfOperation   string  `json:"EngineTotalHoursOfOperation"`
	LincDate                      string  `json:"LINC_DATE"`
	WaterDepth                    string  `json:"WATER_DEPTH"`
	WaterTemperature              string  `json:"WATER_TEMP"`
	FuelLevel1                    string  `json:"FuelLevel1"`
	GSMConnected                  string  `json:"GSM_Connected"`
	GSMSignalQualityDBm           string  `json:"GSM_Signal_Quality_dBm"`
	ServiceDateCapture            string  `json:"SERVICE_DATE_CAPTURE"`
	ServiceHoursCapture           string  `json:"SERVICE_HOURS_CAPTURE"`
	ServiceReminderEnabled        string  `json:"SERVICE_REMINDER_ENABLE"`
}

type sDate struct {
	Seconds     int64 `json:"_seconds"`
	Nanoseconds int64 `json:"_nanoseconds"`
}

func (t *Telemetry) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	var r telemetryRaw
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	// parse all the strings into meaningful types
	errs := []error{}
	accelXYZ, err := decimal.NewFromString(r.AccelXYZMagnitude)
	errs = append(errs, err)
	ballastBelly, err := strconv.Atoi(r.BallastBelly)
	errs = append(errs, err)
	ballastPortSupplementary, err := strconv.Atoi(r.BallastPortSupplementary)
	errs = append(errs, err)
	ballastStarboardSupplementary, err := strconv.Atoi(r.BallastStarboardSupplementary)
	errs = append(errs, err)
	ballastPortRear, err := strconv.Atoi(r.BallastPortRear)
	errs = append(errs, err)
	ballastStarboardRear, err := strconv.Atoi(r.BallastStarboardRear)
	errs = append(errs, err)
	boatStatus, err := strconv.Atoi(r.BoatStatus)
	errs = append(errs, err)
	digitalInput4, err := strconv.Atoi(r.DigitalInput4)
	errs = append(errs, err)
	engineHoursLinc, err := decimal.NewFromString(r.EngineHoursLinc)
	errs = append(errs, err)
	fuelLevelLinc, err := strconv.Atoi(r.FuelLevelLinc)
	errs = append(errs, err)
	GPSAltitude, err := strconv.Atoi(r.GPSAltitude)
	errs = append(errs, err)
	GPSCourse, err := decimal.NewFromString(r.GPSCourse)
	errs = append(errs, err)
	GSMSignalQualityBars, err := strconv.Atoi(r.GSMSignalQualityBars)
	errs = append(errs, err)
	temperatureCelsius, err := decimal.NewFromString(r.TemperatureCelsius)
	errs = append(errs, err)
	GPSQuality, err := strconv.Atoi(r.GPSQuality)
	errs = append(errs, err)
	GPSSpeed, err := decimal.NewFromString(r.GPSSpeed)
	errs = append(errs, err)
	in1Ana, err := decimal.NewFromString(r.In1Ana)
	errs = append(errs, err)
	in2Ana, err := decimal.NewFromString(r.In2Ana)
	errs = append(errs, err)
	in3Ana, err := decimal.NewFromString(r.In3Ana)
	errs = append(errs, err)
	lastSensorUpdate := time.UnixMilli(r.LastSensorUpdate).UTC()
	engineSpeed, err := decimal.NewFromString(r.EngineSpeed)
	errs = append(errs, err)
	GPSSpeedKmh, err := decimal.NewFromString(r.GPSSpeedKmh)
	errs = append(errs, err)
	GPSSpeedMph, err := decimal.NewFromString(r.GPSSpeedMph)
	errs = append(errs, err)
	timestamp := time.UnixMilli(r.Timestamp).UTC()
	GPSNumSatellites, err := strconv.Atoi(r.GPSNumSatellites)
	errs = append(errs, err)
	airTemperatureLinc, err := decimal.NewFromString(r.AirTemperatureLinc)
	errs = append(errs, err)
	engineTotalHoursOfOperation, err := decimal.NewFromString(r.EngineTotalHoursOfOperation)
	errs = append(errs, err)
	lincDate, err := time.Parse("20060102", r.LincDate)
	errs = append(errs, err)
	waterDepth, err := decimal.NewFromString(r.WaterDepth)
	errs = append(errs, err)
	waterTemperature, err := decimal.NewFromString(r.WaterTemperature)
	errs = append(errs, err)
	fuelLevel1, err := decimal.NewFromString(r.FuelLevel1)
	errs = append(errs, err)
	GSMConnected, err := strconv.Atoi(r.GSMConnected)
	errs = append(errs, err)
	GSMSignalQualityDBm, err := strconv.Atoi(r.GSMSignalQualityDBm)
	errs = append(errs, err)
	serviceDateCapture, err := strconv.Atoi(r.ServiceDateCapture)
	errs = append(errs, err)
	serviceHoursCapture, err := decimal.NewFromString(r.ServiceHoursCapture)
	errs = append(errs, err)
	serviceReminderEnabled := false
	if r.ServiceReminderEnabled == "1" {
		serviceReminderEnabled = true
	}

	// check all the errors
	for i, e := range errs {
		if e != nil {
			return fmt.Errorf("cannot parse raw telemetry (internal ref: %d): %s", i, e.Error())
		}
	}

	*t = Telemetry{
		DeviceSerial:                  r.DeviceSerial,
		ServiceMessage:                r.ServiceMessage,
		ServiceDetailsMessage:         r.ServiceDetailsMessage,
		GPSLongitude:                  decimal.NewFromFloat(r.GPSLongitude),
		GPSLatitude:                   decimal.NewFromFloat(r.GPSLatitude),
		IsWebSocketConnected:          false,
		AccelXYZMagnitude:             accelXYZ,
		BallastBelly:                  int64(ballastBelly),
		BallastPortSupplementary:      int64(ballastPortSupplementary),
		BallastStarboardSupplementary: int64(ballastStarboardSupplementary),
		BallastPortRear:               int64(ballastPortRear),
		BallastStarboardRear:          int64(ballastStarboardRear),
		BoatStatus:                    int64(boatStatus),
		DigitalInput4:                 int64(digitalInput4),
		EngineHoursLinc:               engineHoursLinc,
		FuelLevelLinc:                 int64(fuelLevelLinc),
		GPSAltitude:                   int64(GPSAltitude),
		GPSCourse:                     GPSCourse,
		GSMSignalQualityBars:          int64(GSMSignalQualityBars),
		TemperatureCelsius:            temperatureCelsius,
		GPSQuality:                    int64(GPSQuality),
		GPSSpeed:                      GPSSpeed,
		In1Ana:                        in1Ana,
		In2Ana:                        in2Ana,
		In3Ana:                        in3Ana,
		IsSleeping:                    false,
		LastSensorUpdate:              lastSensorUpdate,
		EngineSpeed:                   engineSpeed,
		GPSSpeedKmh:                   GPSSpeedKmh,
		GPSSpeedMph:                   GPSSpeedMph,
		Timestamp:                     timestamp,
		GPSNumSatellites:              int64(GPSNumSatellites),
		AirTemperatureLinc:            airTemperatureLinc,
		EngineTotalHoursOfOperation:   engineTotalHoursOfOperation,
		LincDate:                      lincDate,
		WaterDepth:                    waterDepth,
		WaterTemperature:              waterTemperature,
		FuelLevel1:                    fuelLevel1,
		GSMConnected:                  int64(GSMConnected),
		GSMSignalQualityDBm:           int64(GSMSignalQualityDBm),
		ServiceDateCapture:            int64(serviceDateCapture),
		ServiceHoursCapture:           serviceHoursCapture,
		ServiceReminderEnabled:        serviceReminderEnabled,
	}

	return nil
}

func (b *BoatInfo) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	var r boatInfoRaw
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	// parse all the strings into meaningful types
	errs := []error{}
	modelYear, err := strconv.Atoi(r.ModelYear)
	errs = append(errs, err)
	date := time.Unix(r.Date.Seconds, r.Date.Nanoseconds).UTC()
	deliveryDate, err := time.Parse("1/2/2006", r.DeliveryDateString)
	errs = append(errs, err)
	deliveryEntryDate, err := time.Parse("1/2/2006 03:04:05 PM", r.DeliveryEntryDateString)
	errs = append(errs, err)
	hin, err := strconv.Atoi(r.Hin)
	errs = append(errs, err)

	// check all the errors
	for i, e := range errs {
		if e != nil {
			return fmt.Errorf("cannot parse raw boat info (internal ref: %d): %s", i, e.Error())
		}
	}

	*b = BoatInfo{
		Owner:             r.Owner,
		ModelYear:         int64(modelYear),
		Date:              date,
		EngineSerial:      r.EngineSerial,
		DealerNumber:      r.DealerNumber,
		FullHin:           r.FullHin,
		Model:             r.Model,
		DeliveryDate:      deliveryDate,
		DeliveryEntryDate: deliveryEntryDate,
		Engine:            r.Engine,
		GPSEnabled:        r.GPSEnabled,
		Hin:               int64(hin),
	}
	return nil
}

func NewMyNautiqueClient(opts *Options) *Client {
	c := &Client{
		Options: &Options{
			Client: &http.Client{},
		},
		AuthUntil: time.Now().Add(-1 * time.Hour),
	}
	if opts == nil {
		return c
	}
	c.Options.User = opts.User
	c.Options.Password = opts.Password
	c.Options.AuthAPIKey = opts.AuthAPIKey
	if opts.Client != nil {
		c.Options.Client = opts.Client
	}
	return c
}

func (m *Client) Login() error {
	now := time.Now()
	postBody, _ := json.Marshal(map[string]any{
		"email":             m.Options.User,
		"password":          m.Options.Password,
		"returnSecureToken": true,
	})
	requestBody := bytes.NewBuffer(postBody)
	url := fmt.Sprintf("%s/accounts:signInWithPassword?key=%s", auth_url, m.Options.AuthAPIKey)
	req, err := http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		return fmt.Errorf("cannot build request: %s", err.Error())
	}
	resp, err := m.Options.Client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot login: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read login response: %s", err.Error())
	}

	var response loginResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("cannot parse login response: %s", err.Error())
	}
	expIn, err := strconv.Atoi(response.ExpiresIn)
	if err == nil {
		response.ExpiresInSeconds = int64(expIn)
	}
	m.AuthUntil = now.Add(time.Second * time.Duration(response.ExpiresInSeconds-10))
	m.auth = response
	return nil
}

func (m *Client) GetFleet() error {
	if m.isAuthExpired() {
		m.Login()
	}

	url := fmt.Sprintf("%s/fleet/get-fleet", api_url)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("token", m.auth.IDToken)
	resp, err := m.Options.Client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot get fleet: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response from get fleet: %s", err.Error())
	}

	boats := []BoatInfo{}
	err = json.Unmarshal(body, &boats)
	if err != nil {
		return fmt.Errorf("cannot parse fleet response: %s", err.Error())
	}
	m.Fleet = boats
	return nil
}

func (m *Client) GetBoatTelemetry(id int64) (Telemetry, error) {
	if m.isAuthExpired() {
		m.Login()
	}

	t := Telemetry{}
	url := fmt.Sprintf("%s/boat/get-boat-telemetry/%d", api_url, id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("token", m.auth.IDToken)
	resp, err := m.Options.Client.Do(req)
	if err != nil {
		return t, fmt.Errorf("cannot get telemetry: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return t, fmt.Errorf("cannot read response from get telemetry: %s", err.Error())
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return t, fmt.Errorf("cannot parse telemetry response: %s", err.Error())
	}
	return t, nil
}

func (m *Client) isAuthExpired() bool {
	return time.Now().After(m.AuthUntil)
}
