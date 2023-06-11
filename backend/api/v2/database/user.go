package database

import "fmt"

// GetUserByUsername find the user that has either the given username
// or the given email address.
func (d *DBMysql) GetUserByName(name string) (User, error) {
	var users []User
	d.orm.Where("username = ? or email = ?", name, name).Find(&users)
	if len(users) > 1 {
		return User{}, fmt.Errorf("multiple users match this username or email")
	}
	if len(users) == 0 {
		return User{}, fmt.Errorf("user not found")
	}
	return users[0], nil
}

// GetUserById returns a user with all its associations preloaded.
func (d *DBMysql) GetUserById(id uint) (User, error) {
	var user User
	err := d.orm.Model(&User{}).Preload("UserStatus").Preload("UserStatus.UserRole").First(&user, id).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}
