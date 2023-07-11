package database

import "time"

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

func (d *DBMysql) GetUserHeatStats(userID uint, start time.Time, end time.Time) (int64, float64, error) {
	type stats struct {
		duration int64
		cost     float64
	}
	var s stats
	err := d.orm.Raw(`
	    SELECT sum(duration_s) as duration, sum(cost_chf) as cost
	    FROM heat
		WHERE UNIX_TIMESTAMP(heat.timestamp) >= ?
		  AND UNIX_TIMESTAMP(heat.timestamp) < ?`, start.Unix(), end.Unix()).
		Scan(&s).Error
	return s.duration, s.cost, err
}

func (d *DBMysql) GetHeatInSessionCount(sessionID uint) (int, error) {
	var count int
	err := d.orm.Raw("SELECT count(*) FROM heat WHERE session_id = ?", sessionID).Scan(&count).Error
	return count, err
}
