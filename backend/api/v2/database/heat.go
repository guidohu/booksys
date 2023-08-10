package database

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

func (d *DBMysql) GetUserHeats(userID uint, size int) ([]Heat, error) {
	var h []Heat
	err := d.orm.Where("user_id = ?", userID).
		Order("timestamp desc").
		Limit(size).
		Preload("User").
		Preload("Session").
		Find(&h).Error
	return h, err
}

// GetUserHeatStats returns total duration seconds and total cost
func (d *DBMysql) GetUserHeatStats(userID uint, start time.Time, end time.Time) (int64, decimal.Decimal, error) {
	type Stats struct {
		Duration int64
		Cost     decimal.Decimal
	}
	s := &Stats{}
	fmt.Println("start", start.Unix())
	fmt.Println("end", end.Unix())
	err := d.orm.Raw(`
	    SELECT sum(duration_s) as duration, sum(cost_chf) as cost
	    FROM heat
		WHERE UNIX_TIMESTAMP(heat.timestamp) >= ?
		  AND UNIX_TIMESTAMP(heat.timestamp) < ?`, start.Unix(), end.Unix()).
		Scan(s).Error
	return s.Duration, s.Cost, err
}

func (d *DBMysql) GetHeatInSessionCount(sessionID uint) (int, error) {
	var count int
	err := d.orm.Raw("SELECT count(*) FROM heat WHERE session_id = ?", sessionID).Scan(&count).Error
	return count, err
}

func (d *DBMysql) GetHeatsInSession(sessionID uint) ([]Heat, error) {
	var heats []Heat
	err := d.orm.Model(&Heat{}).
		Where("session_id = ?", sessionID).
		Preload("User").
		Preload("Session").
		Find(&heats).Error
	return heats, err
}

func (d *DBMysql) AddHeat(h *Heat) error {
	return d.orm.Create(h).Error
}

func (d *DBMysql) DeleteHeat(heatID uint) error {
	return d.orm.Exec("DELETE FROM heat WHERE id = ?", heatID).Error
}

func (d *DBMysql) GetHeat(heatID uint) (Heat, error) {
	var h Heat
	err := d.orm.Where("id = ?", heatID).First(&h).Error
	return h, err
}

func (d *DBMysql) ChangeHeat(h *Heat) error {
	return d.orm.Save(h).Error
}
