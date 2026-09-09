package database

import "github.com/shopspring/decimal"

// GetUserSessionPayments returns the total a user paid for sessions.
func (d *Mysql) GetUserSessionPayments(userID uint) (decimal.Decimal, error) {
	var res decimal.Decimal
	err := d.orm.Raw("SELECT COALESCE(sum(amount_chf), 0) as total FROM payment WHERE user_id = ? AND type_id = ?", userID, ExpenseTypeSession).
		Scan(&res).Error
	return res, err
}
