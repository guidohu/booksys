package database

import (
	"time"

	"github.com/shopspring/decimal"
)

// GetUserHeats returns up to size of the user's most recent rides.
func (d *Mysql) GetUserHeats(userID uint, size int) ([]Heat, error) {
	var h []Heat
	err := d.orm.Where("user_id = ?", userID).
		Order("timestamp desc").
		Limit(size).
		Preload("User").
		Preload("Session").
		Find(&h).Error
	return h, err
}

// GetUserHeatsBySession returns up to size of the user's most recent rides in
// a given session.
func (d *Mysql) GetUserHeatsBySession(sessionID uint, userID uint, size int) ([]Heat, error) {
	var h []Heat
	err := d.orm.Where("user_id = ?", userID).
		Where("session_id = ?", sessionID).
		Order("timestamp desc").
		Limit(size).
		Preload("User").
		Preload("Session").
		Find(&h).Error
	return h, err
}

// GetUserHeatStats returns the total duration in seconds and the total cost.
func (d *Mysql) GetUserHeatStats(userID uint, start time.Time, end time.Time) (int64, decimal.Decimal, error) {
	type Stats struct {
		Duration int64
		Cost     decimal.Decimal
	}
	s := &Stats{}
	err := d.orm.Raw(`
	    SELECT sum(duration_s) as duration, sum(cost_chf) as cost
	    FROM heat
		WHERE UNIX_TIMESTAMP(heat.timestamp) >= ?
		  AND UNIX_TIMESTAMP(heat.timestamp) < ?
		  AND user_id = ?`, start.Unix(), end.Unix(), userID).
		Scan(s).Error
	return s.Duration, s.Cost, err
}

// GetHeatInSessionCount returns how many rides a session has.
func (d *Mysql) GetHeatInSessionCount(sessionID uint) (int, error) {
	var count int
	err := d.orm.Raw("SELECT count(*) FROM heat WHERE session_id = ?", sessionID).Scan(&count).Error
	return count, err
}

// GetHeatsInSession returns all rides of a session.
func (d *Mysql) GetHeatsInSession(sessionID uint) ([]Heat, error) {
	var heats []Heat
	err := d.orm.Model(&Heat{}).
		Where("session_id = ?", sessionID).
		Preload("User").
		Preload("Session").
		Find(&heats).Error
	return heats, err
}

// AddHeat records a ride. The ID of h is filled in on success.
func (d *Mysql) AddHeat(h *Heat) error {
	return d.orm.Create(h).Error
}

// DeleteHeat removes a ride.
func (d *Mysql) DeleteHeat(heatID uint) error {
	return d.orm.Exec("DELETE FROM heat WHERE id = ?", heatID).Error
}

// GetHeat returns a single ride.
func (d *Mysql) GetHeat(heatID uint) (Heat, error) {
	var h Heat
	err := d.orm.Where("id = ?", heatID).First(&h).Error
	return h, err
}

// ChangeHeat writes back a modified ride.
func (d *Mysql) ChangeHeat(h *Heat) error {
	return d.orm.Save(h).Error
}
