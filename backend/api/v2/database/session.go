package database

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

// GetSessionsBetween get sessions between start and end
func (d *DBMysql) GetSessionsBetween(start, end time.Time) ([]Session, error) {
	var sessions []Session
	// Find sessions that
	// 1. start after 'start' and start before 'end'
	// 2. end after 'start' and end before 'end'
	// 3. start after 'start' and end before 'end'
	err := d.orm.Model(&Session{}).
		Where("UNIX_TIMESTAMP(start_time) >= ? AND UNIX_TIMESTAMP(start_time) < ?", start.Unix(), end.Unix()).
		Or("UNIX_TIMESTAMP(end_time) >= ? AND UNIX_TIMESTAMP(end_time) < ?", start.Unix(), end.Unix()).
		Or("UNIX_TIMESTAMP(start_time) < ? AND UNIX_TIMESTAMP(end_time) >= ?", start.Unix(), end.Unix()).
		Preload("SessionType").
		Preload("Creator").
		Find(&sessions).Error
	d.orm.Debug().Model(&Session{}).
		Where("UNIX_TIMESTAMP(start_time) >= ? AND UNIX_TIMESTAMP(start_time) < ?", start.Unix(), end.Unix()).
		Or("UNIX_TIMESTAMP(end_time) >= ? AND UNIX_TIMESTAMP(end_time) < ?", start.Unix(), end.Unix()).
		Or("UNIX_TIMESTAMP(start_time) < ? AND UNIX_TIMESTAMP(end_time) >= ?", start.Unix(), end.Unix()).
		Preload("SessionType").
		Preload("Creator").
		Find(&sessions)
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
	var sessionsMap = make(map[uint]bool)
	for _, s := range creatorSessions {
		if _, ok := sessionsMap[s.ID]; !ok {
			sessionsMap[s.ID] = true
			sessions = append(sessions, s)
		}
	}
	for _, s := range memberSessions {
		if _, ok := sessionsMap[s.ID]; !ok {
			sessionsMap[s.ID] = true
			sessions = append(sessions, s)
		}
	}
	sort.Sort(SessionByTime(sessions))
	return sessions, nil
}

func (d *DBMysql) GetUsersForSession(id uint) ([]UserToSession, error) {
	var users []UserToSession
	err := d.orm.Model(&UserToSession{}).
		Where("session_id = ?", id).
		Preload("User").
		Find(&users).
		Error
	return users, err
}

func (d *DBMysql) CreateSession(s Session) (uint, error) {
	err := d.orm.Create(&s).Error
	return s.ID, err
}

func (d *DBMysql) GetSession(sessionID uint) (Session, error) {
	var session Session
	err := d.orm.First(&session, sessionID).Error
	return session, err
}

func (d *DBMysql) UpdateSession(s Session) error {
	return d.orm.Save(&s).Error
}

func (d *DBMysql) DeleteUsersFromSession(sessionID uint) error {
	return d.orm.Exec("DELETE FROM user_to_session WHERE session_id = ?", sessionID).Error
}

func (d *DBMysql) DeleteSession(sessionID uint) error {
	return d.orm.Exec("DELETE FROM session WHERE id = ?", sessionID).Error
}

func (d *DBMysql) AddSessionToUserEntry(u UserToSession) error {
	return d.orm.Transaction(func(tx *gorm.DB) error {
		err := tx.Where(UserToSession{UserID: u.UserID, SessionID: u.SessionID}).
			FirstOrCreate(&u).Error
		if err != nil {
			return err
		}
		session := &Session{}
		err = tx.First(session, u.SessionID).Error
		if err != nil {
			return err
		}
		if session.FreeSpaces < 1 {
			return fmt.Errorf("cannot add user %d to sessioon %d because there are no free spaces left", u.UserID, u.SessionID)
		}
		return tx.Exec("UPDATE session SET free = free - 1 WHERE id = ?", u.SessionID).Error
	})
}

func (d *DBMysql) DeleteSessionToUserEntry(userID uint, sessionID uint) error {
	return d.orm.Transaction(func(tx *gorm.DB) error {
		err := d.orm.Exec("DELETE FROM user_to_session WHERE user_id = ? AND session_id = ?", userID, sessionID).Error
		if err != nil {
			return err
		}
		return tx.Exec("UPDATE session SET free = free + 1 WHERE id = ?", sessionID).Error
	})
}
