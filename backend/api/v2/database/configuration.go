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

var AllowedProperties = map[string]interface{}{
	"currency":                 nil,
	"engine.hour.format":       nil,
	"fuel.payment.type":        nil,
	"location.address":         nil,
	"location.latitude":        nil,
	"location.longitude":       nil,
	"location.map":             nil,
	"location.time.zone":       nil,
	"logo.file":                nil,
	"mynautique.boat.id":       nil,
	"mynautique.enabled":       nil,
	"mynautique.fuel.capacity": nil,
	"mynautique.password":      nil,
	"mynautique.user":          nil,
	"payment.account.bic":      nil,
	"payment.account.comment":  nil,
	"payment.account.iban":     nil,
	"payment.account.owner":    nil,
	"recaptcha.privatekey":     nil,
	"recaptcha.publickey":      nil,
	"smtp.password":            nil,
	"smtp.sender":              nil,
	"smtp.server":              nil,
	"smtp.username":            nil,
}

type MyNautiqueConfiguration struct {
	Enabled      bool
	BoatID       int64
	FuelCapacity decimal.Decimal
	Password     string
	User         string
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
	if reflect.DeepEqual(e, &EmailConfiguration{}) {
		return true
	}
	return false
}

func (d *DBMysql) GetPropertyValue(key string) (Configuration, error) {
	var value Configuration
	err := d.orm.Where("property = ?", key).First(&value).Error
	return value, err
}

func (d *DBMysql) GetAllPropertyValues() ([]Configuration, error) {
	var values []Configuration
	err := d.orm.Find(&values).Error
	return values, err
}

func getPropertyValuesMap(c []Configuration) map[string]Configuration {
	configMap := map[string]Configuration{}
	for _, p := range c {
		configMap[p.Property] = p
	}
	return configMap
}

func (d *DBMysql) UpdateOrInsertPropertyValues(conf []Configuration) error {
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

			prop.Property = c.Property
			prop.Value = c.Value
			err = tx.Save(&prop).Error
			if err != nil {
				slog.Error("Property cannot be updated/set", slog.String("property", c.Property), slog.String("value", c.Value))
				continue
			}
		}
		return nil
	})

	if err != nil {
		slog.Error("Cannot update/set properties. Change has not been commited:", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (d *DBMysql) GetTimezoneLocation() (*time.Location, error) {
	s, err := d.GetPropertyValue("location.time.zone")
	if err != nil {
		slog.Error("Cannot get location.time.zone", slog.String("error", err.Error()))
		return nil, err
	}
	return time.LoadLocation(s.Value)
}

func (d *DBMysql) GetEmailConfiguration() (EmailConfiguration, error) {
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

func (d *DBMysql) GetMyNautiqueConfiguration() (MyNautiqueConfiguration, error) {
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

	if valid {
		return config, nil
	}
	slog.Error("Return empty mynautique configuration")
	return MyNautiqueConfiguration{}, fmt.Errorf("invalid mynautique configuration")
}
