package database

import (
	"errors"
	"log/slog"
)

// GetUserByName finds the user that has either the given username or the
// given email address.
func (d *Mysql) GetUserByName(name string) (User, error) {
	if d.orm == nil {
		return User{}, errors.New("database handler not available")
	}
	var users []User
	if err := d.orm.Where("username = ? or email = ?", name, name).Find(&users).Error; err != nil {
		return User{}, err
	}
	if len(users) > 1 {
		return User{}, errors.New("multiple users match this username or email")
	}
	if len(users) == 0 {
		return User{}, errors.New("user not found")
	}
	return users[0], nil
}

// GetUserByID returns a user with all its associations preloaded.
func (d *Mysql) GetUserByID(id uint) (User, error) {
	var user User
	err := d.orm.Model(&User{}).Preload("UserStatus").Preload("UserStatus.UserRole").First(&user, id).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// AddUser adds a user and returns the ID it was given.
func (d *Mysql) AddUser(u User) (uint, error) {
	tx := d.orm.Create(&u)
	return u.ID, tx.Error
}

// DeleteUserByID flags a user as deleted and clears their personal data. The
// row is kept so that the transactions referencing it stay intact.
func (d *Mysql) DeleteUserByID(id uint) error {
	user := User{
		ID: id,
	}
	return d.orm.Model(&user).Updates(map[string]any{
		"username":      "",
		"password_salt": 0,
		"password_hash": "none",
		"first_name":    "",
		"last_name":     "",
		"address":       "",
		"city":          "",
		"plz":           0,
		"mobile":        "",
		"email":         "",
		"license":       0,
		"status":        UserStatusGuest,
		"locked":        1,
		"comment":       "deleted user",
		"deleted":       1,
	}).Error
}

// UpdateUser updates user values that are supposed to be updated
// by the users themselves. It does not update PasswordSalt, PasswordHash
// Locked, UserStatusID, Comment and IsDeleted
func (d *Mysql) UpdateUser(userID uint, user User) error {
	// update the ID as we want it to be specified explicitly
	user.ID = userID

	return d.orm.Model(&user).
		Select("*").
		Omit("PasswordSalt", "PasswordHash", "Locked", "UserStatusID", "Comment", "IsDeleted").
		Updates(user).Error
}

// UpdatePassword updates salt and password of a user. It does not
// update anything else, although an entire user can be provided. It
// ignores all values except the hash and the salt.
func (d *Mysql) UpdatePassword(userID uint, user User) error {
	// update the ID as we want it to be specified explicitly
	user.ID = userID

	return d.orm.Model(&user).
		Select("PasswordSalt", "PasswordHash").
		Updates(user).Error
}

// CountAdminUsers returns how many users are administrators.
func (d *Mysql) CountAdminUsers() int64 {
	var count int64
	d.orm.Model(&User{}).Where("status = ?", UserStatusAdmin).Count(&count)
	return count
}

// GetAdminUsers returns all administrators.
func (d *Mysql) GetAdminUsers() ([]User, error) {
	var users []User
	err := d.orm.Where("status = ?", UserStatusAdmin).Find(&users).Error
	return users, err
}

// IsAdminUser reports whether a user is an administrator.
func (d *Mysql) IsAdminUser(userID uint) (bool, error) {
	var count int64
	err := d.orm.Model(&User{}).Where("id = ? AND status = ?", userID, UserStatusAdmin).Count(&count).Error
	return count == 1, err
}

// UsersExist reports whether the user table holds any user at all. It is used
// to detect a database that has not been set up yet.
func (d *Mysql) UsersExist() (bool, error) {
	var exists bool
	err := d.db.QueryRow("SELECT IF(COUNT(*), 'true', 'false') AS users FROM user;").Scan(&exists)
	if err != nil {
		slog.Error("Cannot check user table for existing users", slog.Any("error", err))
		return false, err
	}
	return exists, nil
}

// UserExists reports whether a user with the given ID exists.
func (d *Mysql) UserExists(id uint) bool {
	var user User
	err := d.orm.Model(&User{}).Where("id = ?", id).First(&user).Error
	return err == nil
}

// ChangeUserStatus moves a user into a different user group.
func (d *Mysql) ChangeUserStatus(id uint, userStatusID uint) error {
	return d.orm.Model(&User{}).Where("id = ?", id).Updates(User{UserStatusID: userStatusID}).Error
}

// ChangeLock locks or unlocks a user account.
func (d *Mysql) ChangeLock(id uint, locked bool) error {
	return d.orm.Model(&User{}).Where("id = ?", id).Update("locked", locked).Error
}

// GetUsers returns all users, optionally including the ones flagged as
// deleted.
func (d *Mysql) GetUsers(includeDeleted bool) ([]User, error) {
	var users []User
	err := d.orm.Where("deleted = 0 OR deleted = ?", includeDeleted).Find(&users).Error
	return users, err
}
