package database

func (d *DBMysql) GetUserSessionPayments(userID uint) (float64, error) {
	var res float64
	err := d.orm.Raw("SELECT COALESCE(sum(amount_chf), 0) as total FROM payment WHERE user_id = ? AND type_id = ?", userID, ExpenseTypeSession).
		Scan(&res).Error
	return res, err
}
