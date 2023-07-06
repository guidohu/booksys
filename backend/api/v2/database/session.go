package database

import (
	"server/generics"
	"sort"
	"time"
)

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

func (d *DBMysql) GetSessionsByUser(userID uint) ([]Session, error) {
	var creatorSessions []Session
	// Get sessions where the user is creator
	err := d.orm.Model(&Session{}).Where("creator_id = ?", userID).Find(&creatorSessions).Error
	if err != nil {
		return []Session{}, err
	}
	// Get sessions where the user is part of
	var memberSessions []Session
	err = d.orm.Table("session").
		Select("session.*").
		Joins("JOIN user_to_session ON session.id = user_to_session.session_id").
		Where("user_to_session.user_id", userID).
		Preload("SessionType").
		Find(&memberSessions).Error
	if err != nil {
		return []Session{}, err
	}

	var sessions []Session
	sessions = append(creatorSessions, memberSessions...)
	sessions = generics.Unique(sessions)
	sort.Sort(SessionByTime(sessions))
	return sessions, nil
}

func (d *DBMysql) GetUsersForSession(id uint) ([]UserToSession, error) {
	var users []UserToSession
	err := d.orm.Model(&UserToSession{}).
		Where("session_id = ?", id).
		Preload("User").
		Error
	return users, err
}
