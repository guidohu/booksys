package database

import (
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

// AddBrowserSession adds a browser session to the table.
func (d *Mysql) AddBrowserSession(b BrowserSession) (string, error) {
	err := d.orm.Create(&b).Error
	if err != nil {
		slog.Error("Cannot add browser session", slog.Any("error", err))
		return "", err
	}

	return b.SessionSecret, nil
}

// GetBrowserSession returns the session or nil and an error in case it cannot
// be found.
func (d *Mysql) GetBrowserSession(id string) (*BrowserSession, error) {
	b := &BrowserSession{
		SessionSecret: id,
	}
	err := d.orm.First(b).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("No browser session found", slog.Any("error", err))
		} else {
			slog.Error("Cannot get browser session", slog.Any("error", err))
		}
		return nil, err
	}
	return b, nil
}

// UpdateBrowserSession writes back a modified login session.
func (d *Mysql) UpdateBrowserSession(b BrowserSession) error {
	return d.orm.Save(b).Error
}

// DeleteBrowserSession removes a login session.
func (d *Mysql) DeleteBrowserSession(b BrowserSession) error {
	return d.orm.Delete(&b).Error
}
