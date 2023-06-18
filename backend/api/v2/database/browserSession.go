package database

import (
	"errors"

	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

// AddBrowserSession adds a browser session to the table.
func (d *DBMysql) AddBrowserSession(b BrowserSession) (string, error) {
	err := d.orm.Create(&b).Error
	if err != nil {
		slog.Error("Cannot add browser session", slog.String("error", err.Error()))
		return "", err
	}

	return b.SessionSecret, nil
}

// GetBrowserSession returns the session or nil and an error in case it cannot
// be found.
func (d *DBMysql) GetBrowserSession(id string) (*BrowserSession, error) {
	b := &BrowserSession{
		SessionSecret: id,
	}
	err := d.orm.First(b).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("No browser session found", slog.String("error", err.Error()))
		} else {
			slog.Error("Cannot get browser session", slog.String("error", err.Error()))
		}
		return nil, err
	}
	return b, nil
}

func (d *DBMysql) UpdateBrowserSession(b BrowserSession) error {
	return d.orm.Save(b).Error
}

func (d *DBMysql) DeleteBrowserSession(b BrowserSession) error {
	return d.orm.Delete(&b).Error
}
