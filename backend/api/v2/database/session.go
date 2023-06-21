package database

import "time"

// GetSessionsBetween get sessions between start and end
func (d *DBMysql) GetSessionsBetween(start, end time.Time) ([]Session, error) {
	var sessions []Session
	err := d.orm.Model(&Session{}).
		Where("UNIX_TIMESTAMP(start_time) >= ? AND UNIX_TIMESTAMP(start_time) < ?", start.Unix(), end.Unix()).
		Or("UNIX_TIMESTAMP(end_time) >= ? AND UNIX_TIMESTAMP(end_time) < ?", start.Unix(), end.Unix()).
		Or("UNIX_TIMESTAMP(start_time) < ? AND UNIX_TIMESTAMP(end_time) >= ?", start.Unix(), end.Unix()).
		Preload("SessionType").
		Preload("Creator").
		Find(&sessions).Error
	return sessions, err
}

func (d *DBMysql) GetUsersForSession(id uint) ([]UserToSession, error) {
	var users []UserToSession
	err := d.orm.Model(&UserToSession{}).
		Where("session_id = ?", id).
		Preload("User").
		Error
	return users, err
}
