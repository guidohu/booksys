package mynautique

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func responseMocker(statusCode int, response string) *http.Client {
	return &http.Client{Transport: RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: statusCode,
			Body:       ioutil.NopCloser(strings.NewReader(response)),
		}
	})}
}

func TestLogin(t *testing.T) {
	mockClient := responseMocker(http.StatusOK, `{
		"kind": "identitytoolkit#VerifyPasswordResponse",
		"localId": "foo.localId",
		"email": "foo@bar.com",
		"displayName": "",
		"idToken": "foo.idToken",
		"registered": true,
		"refreshToken": "foo.refreshToken",
		"expiresIn": "3600"
	}`)
	expectedAuth := loginResponse{
		IDToken:          "foo.idToken",
		ExpiresIn:        "3600",
		ExpiresInSeconds: 3600,
		Kind:             "identitytoolkit#VerifyPasswordResponse",
		LocalID:          "foo.localId",
		Email:            "foo@bar.com",
		DisplayName:      "",
		Registered:       true,
		RefreshToken:     "foo.refreshToken",
	}

	m := NewClient(&Options{
		User:       "foo",
		Password:   "bar",
		AuthAPIKey: "foobar",
		Client:     mockClient,
	})
	err := m.Login()
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(expectedAuth, m.auth) {
		t.Logf("got %v, want %v", m.auth, expectedAuth)
		t.Fail()
	}
}

func TestGetFleet(t *testing.T) {
	mockClient := responseMocker(http.StatusOK, `[
		{
			"owner":"foo@bar.com",
			"myear":"2022",
			"sDate":{"_seconds":1632700800,"_nanoseconds":0},
			"EngineSerial":"02-06L051FW-22-P999999",
			"dealer":1111,
			"fullHin":"USCTC99999H122",
			"model":"G23",
			"deliveryDate":"3/16/2022",
			"deliveryEntryDate":"3/16/2022 11:59:56 AM",
			"Engine":"ZZ6 450HP - 6.2L Direct Injection",
			"gpsEnabled":false,
			"hin":"99999"
		}
	]`)
	expectedFleet := []BoatInfo{
		{
			Owner:             "foo@bar.com",
			ModelYear:         2022,
			Date:              time.Date(2021, 9, 27, 0, 0, 0, 0, time.UTC),
			EngineSerial:      "02-06L051FW-22-P999999",
			DealerNumber:      1111,
			FullHin:           "USCTC99999H122",
			Model:             "G23",
			DeliveryDate:      time.Date(2022, 3, 16, 0, 0, 0, 0, time.UTC),
			DeliveryEntryDate: time.Date(2022, 3, 16, 11, 59, 56, 0, time.UTC),
			Engine:            "ZZ6 450HP - 6.2L Direct Injection",
			GPSEnabled:        false,
			Hin:               99999,
		},
	}
	m := NewClient(&Options{
		User:       "foo",
		Password:   "bar",
		AuthAPIKey: "foobar",
		Client:     mockClient,
	})
	m.auth.IDToken = "foo-bar-id-token"
	m.AuthUntil = time.Now().Add(time.Hour)
	err := m.GetFleet()
	if err != nil {
		t.Logf("got '%v', want '%v'", err.Error(), nil)
		t.Fail()
	}
	assert.Equal(t, expectedFleet, m.Fleet, "Fleets should be the same.")
}

func TestGetBoatTelemetry(t *testing.T) {
	mockClient := responseMocker(http.StatusOK, `{
		"deviceSerial":21000099999,
		"serviceMessage":"",
		"serviceDetailsMessage":"",
		"gps_long":"12.749070779275361",
		"gps_lat":"23.38980444927536",
		"isWebSocketConnected":false,
		"Accel_XYZ_Magnitude":"0.124000",
		"BALLAST_BELLY":"0",
		"BALLAST_PORT_SUPP":"0",
		"BALLAST_REAR_PORT":"0",
		"BALLAST_REAR_STBD":"0",
		"BALLAST_STBD_SUPP":"0",
		"BOAT_STATUS":"0",
		"Digital_Input_4":"0",
		"ENGINE_HOURS_LINC":"187.300003",
		"FUEL_LEVEL_LINC":"71",
		"GPS_Altitude_m":"403",
		"GPS_Course":"0.000000",
		"GSM_Signal_Quality_bars":"0",
		"Temp_C":"51.750004",
		"gps_qual":"0",
		"gps_speed":"0.255600",
		"in1_ana":"13.032580",
		"in2_ana":"13.143538",
		"in3_ana":"0.032143",
		"isSleeping":false,
		"lastSensorUpdate":1690023626586,
		"EngineSpeed":"0",
		"GPS_Speed_kmh":"0",
		"GPS_Speed_mph":"0",
		"ts":1690023626329,
		"GPS_Num_Satellites":"18",
		"AIR_TEMP_LINC":"297",
		"EngineTotalHoursOfOperation":"187.250000",
		"LINC_DATE":"20221202",
		"WATER_DEPTH":"1.670000",
		"WATER_TEMP":"301",
		"FuelLevel1":"71.200005",
		"GSM_Connected":"1",
		"GSM_Signal_Quality_dBm":"-63",
		"SERVICE_DATE_CAPTURE":"0",
		"SERVICE_HOURS_CAPTURE":"197.600006",
		"SERVICE_REMINDER_ENABLE":"1"
	}`)
	AccelXYZMagnitude, _ := decimal.NewFromString("0.124000")
	GPSLongitude, _ := decimal.NewFromString("12.749070779275361")
	GPSLatitude, _ := decimal.NewFromString("23.38980444927536")
	EngineHoursLinc, _ := decimal.NewFromString("187.300003")
	GPSCourse, _ := decimal.NewFromString("0.000000")
	TemperatureCelsius, _ := decimal.NewFromString("51.750004")
	GPSSpeed, _ := decimal.NewFromString("0.255600")
	In1Ana, _ := decimal.NewFromString("13.032580")
	In2Ana, _ := decimal.NewFromString("13.143538")
	In3Ana, _ := decimal.NewFromString("0.032143")
	EngineSpeed, _ := decimal.NewFromString("0")
	GPSSpeedKmh, _ := decimal.NewFromString("0")
	GPSSpeedMph, _ := decimal.NewFromString("0")
	AirTemperatureLinc, _ := decimal.NewFromString("297")
	EngineTotalHoursOfOperation, _ := decimal.NewFromString("187.250000")
	WaterDepth, _ := decimal.NewFromString("1.670000")
	WaterTemperature, _ := decimal.NewFromString("301")
	FuelLevel1, _ := decimal.NewFromString("71.200005")
	ServiceHoursCapture, _ := decimal.NewFromString("197.600006")
	expectedTelemetry := Telemetry{
		DeviceSerial:                  21000099999,
		ServiceMessage:                "",
		ServiceDetailsMessage:         "",
		GPSLongitude:                  GPSLongitude,
		GPSLatitude:                   GPSLatitude,
		IsWebSocketConnected:          false,
		AccelXYZMagnitude:             AccelXYZMagnitude,
		BallastBelly:                  0,
		BallastPortSupplementary:      0,
		BallastStarboardSupplementary: 0,
		BallastPortRear:               0,
		BallastStarboardRear:          0,
		BoatStatus:                    0,
		DigitalInput4:                 0,
		EngineHoursLinc:               EngineHoursLinc,
		FuelLevelLinc:                 71,
		GPSAltitude:                   403,
		GPSCourse:                     GPSCourse,
		GSMSignalQualityBars:          0,
		TemperatureCelsius:            TemperatureCelsius,
		GPSQuality:                    0,
		GPSSpeed:                      GPSSpeed,
		In1Ana:                        In1Ana,
		In2Ana:                        In2Ana,
		In3Ana:                        In3Ana,
		IsSleeping:                    false,
		LastSensorUpdate:              time.UnixMilli(1690023626586).UTC(),
		EngineSpeed:                   EngineSpeed,
		GPSSpeedKmh:                   GPSSpeedKmh,
		GPSSpeedMph:                   GPSSpeedMph,
		Timestamp:                     time.UnixMilli(1690023626329).UTC(),
		GPSNumSatellites:              18,
		AirTemperatureLinc:            AirTemperatureLinc,
		EngineTotalHoursOfOperation:   EngineTotalHoursOfOperation,
		LincDate:                      time.Date(2022, 12, 2, 0, 0, 0, 0, time.UTC),
		WaterDepth:                    WaterDepth,
		WaterTemperature:              WaterTemperature,
		FuelLevel1:                    FuelLevel1,
		GSMConnected:                  1,
		GSMSignalQualityDBm:           -63,
		ServiceDateCapture:            0,
		ServiceHoursCapture:           ServiceHoursCapture,
		ServiceReminderEnabled:        true,
	}
	m := NewClient(&Options{
		User:       "foo",
		Password:   "bar",
		AuthAPIKey: "foobar",
		Client:     mockClient,
	})
	m.auth.IDToken = "foo-bar-id-token"
	m.AuthUntil = time.Now().Add(time.Hour)
	telemetry, err := m.GetBoatTelemetry(999)
	if err != nil {
		t.Logf("got %v, want %v", err.Error(), nil)
		t.Fail()
	}
	assert.Equal(t, expectedTelemetry, telemetry, "Telemetry should be the same.")
}
