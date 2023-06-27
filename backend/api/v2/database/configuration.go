package database

import (
	"errors"

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
