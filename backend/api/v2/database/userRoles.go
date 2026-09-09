package database

// GetUserRoles returns all access levels.
func (d *Mysql) GetUserRoles() ([]UserRole, error) {
	var roles []UserRole
	err := d.orm.Model(&UserRole{}).Find(&roles).Error
	return roles, err
}
