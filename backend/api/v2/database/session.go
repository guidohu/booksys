package database

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

// GetSessionsBetween returns the sessions that overlap the given time range.
func (d *Mysql) GetSessionsBetween(start, end time.Time) ([]Session, error) {
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
	return sessions, err
}

// GetSessionsByUser returns the sessions a user created or takes part in,
// deduplicated and sorted by start time.
func (d *Mysql) GetSessionsByUser(userID uint) ([]Session, error) {
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

// GetUsersForSession returns the participants of a session.
func (d *Mysql) GetUsersForSession(id uint) ([]UserToSession, error) {
	var users []UserToSession
	err := d.orm.Model(&UserToSession{}).
		Where("session_id = ?", id).
		Preload("User").
		Find(&users).
		Error
	return users, err
}

// CreateSession adds a session and returns the ID it was given.
func (d *Mysql) CreateSession(s Session) (uint, error) {
	err := d.orm.Create(&s).Error
	return s.ID, err
}

// GetSession returns a single session.
func (d *Mysql) GetSession(sessionID uint) (Session, error) {
	var session Session
	err := d.orm.First(&session, sessionID).Error
	return session, err
}

// UpdateSession writes back a modified session.
func (d *Mysql) UpdateSession(s Session) error {
	return d.orm.Save(&s).Error
}

// DeleteUsersFromSession removes every participant from a session.
func (d *Mysql) DeleteUsersFromSession(sessionID uint) error {
	return d.orm.Exec("DELETE FROM user_to_session WHERE session_id = ?", sessionID).Error
}

// DeleteSession removes a session.
func (d *Mysql) DeleteSession(sessionID uint) error {
	return d.orm.Exec("DELETE FROM session WHERE id = ?", sessionID).Error
}

// AddSessionToUserEntry adds a user to a session and takes one of its free
// spaces. It fails if the session is full.
func (d *Mysql) AddSessionToUserEntry(u UserToSession) error {
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

// DeleteSessionToUserEntry removes a user from a session and gives its free
// space back.
func (d *Mysql) DeleteSessionToUserEntry(userID uint, sessionID uint) error {
	return d.orm.Transaction(func(tx *gorm.DB) error {
		err := d.orm.Exec("DELETE FROM user_to_session WHERE user_id = ? AND session_id = ?", userID, sessionID).Error
		if err != nil {
			return err
		}
		return tx.Exec("UPDATE session SET free = free + 1 WHERE id = ?", sessionID).Error
	})
}
