package database

func (d *DBMysql) GetExpenseTypes() ([]ExpenseType, error) {
	et := []ExpenseType{}
	err := d.orm.Model(&ExpenseType{}).Find(&et).Error
	return et, err
}
