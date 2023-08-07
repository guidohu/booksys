package database

func (d *DBMysql) GetUserRoles() ([]UserRole, error) {
	var roles []UserRole
	err := d.orm.Model(&UserRole{}).Find(&roles).Error
	return roles, err
}
