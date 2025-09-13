package database

func (d *Mysql) AddPasswordResetToken(entry PasswordReset) error {
	return d.orm.Create(&entry).Error
}

func (d *Mysql) GetPasswordResetEntry(userID uint, token string) (PasswordReset, error) {
	ret := PasswordReset{}
	err := d.orm.Where("user_id = ? AND token = ? AND valid = true AND timestamp > now()", userID, token).First(&ret).Error
	return ret, err
}

func (d *Mysql) InvalidatePasswordResetEntries(userID uint) error {
	return d.orm.Exec("UPDATE password_reset SET valid = false WHERE user_id = ?", userID).Error
}
