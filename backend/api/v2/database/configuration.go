package database

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

const HiddenSecret = "hidden"

var AllowedProperties = map[string]interface{}{
	"currency":                    nil,
	"engine.hour.format":          nil,
	"fuel.payment.type":           nil,
	"location.address":            nil,
	"location.latitude":           nil,
	"location.longitude":          nil,
	"location.map":                nil,
	"location.timezone":           nil,
	"logo.file":                   nil,
	"mynautique.boat.id":          nil,
	"mynautique.enabled":          nil,
	"mynautique.fuel.capacity":    nil,
	"mynautique.password":         nil,
	"mynautique.user":             nil,
	"mynautique.api.key":          nil,
	"payment.account.bic":         nil,
	"payment.account.comment":     nil,
	"payment.account.iban":        nil,
	"payment.account.owner":       nil,
	"recaptcha.privatekey":        nil,
	"recaptcha.publickey":         nil,
	"smtp.password":               nil,
	"smtp.sender":                 nil,
	"smtp.server":                 nil,
	"smtp.username":               nil,
	"business.day.start":          nil,
	"business.day.end":            nil,
	"business.day.startatsunrise": nil,
	"business.day.endatsunset":    nil,
}

type MyNautiqueConfiguration struct {
	Enabled      bool
	BoatID       int64
	FuelCapacity decimal.Decimal
	Password     string
	User         string
	APIKey       string
}

type EmailConfiguration struct {
	Sender     string
	ServerPort string
	Server     string
	Port       string
	Username   string
	Password   string
}

func (e *EmailConfiguration) Empty() bool {
	return reflect.DeepEqual(e, &EmailConfiguration{})
}

func (d *Mysql) DeleteProperty(key string) error {
	return d.orm.Delete(&Configuration{}, "property = ?", key).Error
}

func (d *Mysql) GetPropertyValue(key string) (Configuration, error) {
	var value Configuration
	err := d.orm.Where("property = ?", key).First(&value).Error
	return value, err
}

func (d *Mysql) GetAllPropertyValues() ([]Configuration, error) {
	var values []Configuration
	err := d.orm.Find(&values).Error
	return values, err
}

func (d *Mysql) GetAllPropertyValuesMap() (map[string]Configuration, error) {
	values, err := d.GetAllPropertyValues()
	if err != nil {
		return nil, err
	}
	return GetPropertyValuesMapFromConfiguration(values), nil
}

func GetPropertyValuesMapFromConfiguration(c []Configuration) map[string]Configuration {
	configMap := map[string]Configuration{}
	for _, p := range c {
		configMap[p.Property] = p
	}
	return configMap
}

func (d *Mysql) UpdateOrInsertPropertyValues(conf []Configuration) error {
	err := d.orm.Transaction(func(tx *gorm.DB) error {
		for _, c := range conf {
			if _, exists := AllowedProperties[c.Property]; !exists {
				slog.Error("Property is not allowed and was not added/updated", slog.String("property", c.Property), slog.String("value", c.Value))
				continue
			}
			var prop Configuration
			err := tx.First(&prop, "property = ?", c.Property).Error
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					slog.Error("Property cannot be found", slog.String("property", c.Property), slog.String("value", c.Value))
					continue
				}
			}

			// The value "hidden" is used for sensitive data that should not be returned by the API.
			// If the value is hidden, we do not update it.
			if c.Value == HiddenSecret {
				continue
			}

			prop.Property = c.Property
			prop.Value = c.Value
			err = tx.Save(&prop).Error
			if err != nil {
				slog.Error("Property cannot be updated/set", slog.String("property", c.Property), slog.String("value", c.Value))
				continue
			}
		}
		if res := tx.Model(&ConfigurationVersion{ID: 1}).Update("version", gorm.Expr("version + 1")); res.Error != nil {
			slog.Error("Cannot update/set configuration version", slog.String("error", res.Error.Error()))
			return res.Error
		}
		if res := tx.Model(&ConfigurationVersion{ID: 1}).Update("time", time.Now()); res.Error != nil {
			slog.Error("Cannot update/set configuration time", slog.String("error", res.Error.Error()))
			return res.Error
		}
		return nil
	})

	if err != nil {
		slog.Error("Cannot update/set properties. Change has not been commited:", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (d *Mysql) GetTimezoneLocation() (*time.Location, error) {
	s, err := d.GetPropertyValue("location.timezone")
	if err != nil {
		slog.Error("Cannot get location.timezone", slog.String("error", err.Error()))
		return nil, err
	}
	return time.LoadLocation(s.Value)
}

func (d *Mysql) GetEmailConfiguration() (EmailConfiguration, error) {
	config := EmailConfiguration{}
	properties, err := d.GetAllPropertyValues()
	if err != nil {
		return EmailConfiguration{}, err
	}
	for _, p := range properties {
		switch p.Property {
		case "smtp.sender":
			config.Sender = p.Value
		case "smtp.server":
			config.ServerPort = p.Value
			host, port, err := net.SplitHostPort(p.Value)
			if err != nil {
				slog.Error("Cannot parse smtp.server into host and port", slog.String("error", err.Error()))
				return EmailConfiguration{}, err
			}
			config.Server = host
			config.Port = port
		case "smtp.password":
			config.Password = p.Value
		case "smtp.username":
			config.Username = p.Value
		}
	}
	return config, nil
}

func (d *Mysql) GetMyNautiqueConfiguration() (MyNautiqueConfiguration, error) {
	config := MyNautiqueConfiguration{}
	valid := true
	enabled, _ := d.GetPropertyValue("mynautique.enabled")
	if enabled.Value == "true" {
		config.Enabled = true
	}
	fuelCapacity, _ := d.GetPropertyValue("mynautique.fuel.capacity")
	if fuelCapacity.Value != "" {
		fc, err := decimal.NewFromString(fuelCapacity.Value)
		if err == nil {
			config.FuelCapacity = fc
		} else {
			slog.Error("Cannot parse mynautique.fuel.capacity")
			valid = false
		}
	}
	boatID, _ := d.GetPropertyValue("mynautique.boat.id")
	if boatID.Value != "" {
		bi, err := strconv.Atoi(boatID.Value)
		if err != nil {
			slog.Error("Cannot parse mynautique.boat.id")
			valid = false
		} else {
			config.BoatID = int64(bi)
		}
	}
	user, _ := d.GetPropertyValue("mynautique.user")
	if user.Value != "" {
		config.User = user.Value
	}
	password, _ := d.GetPropertyValue("mynautique.password")
	if password.Value != "" {
		config.Password = password.Value
	}
	apiKey, _ := d.GetPropertyValue("mynautique.api.key")
	if apiKey.Value != "" {
		config.APIKey = apiKey.Value
	}

	if valid {
		return config, nil
	}
	slog.Error("Return empty mynautique configuration")
	return MyNautiqueConfiguration{}, fmt.Errorf("invalid mynautique configuration")
}
