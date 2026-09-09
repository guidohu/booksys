package database

// GetExpenseTypes returns all expense types.
func (d *Mysql) GetExpenseTypes() ([]ExpenseType, error) {
	var et []ExpenseType
	err := d.orm.Model(&ExpenseType{}).Find(&et).Error
	return et, err
}
