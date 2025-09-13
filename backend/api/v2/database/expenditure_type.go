package database

func (d *Mysql) GetExpenseTypes() ([]ExpenseType, error) {
	et := []ExpenseType{}
	err := d.orm.Model(&ExpenseType{}).Find(&et).Error
	return et, err
}
