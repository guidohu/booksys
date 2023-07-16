package database

func (d *DBMysql) GetPricings() ([]Pricing, error) {
	var pricing []Pricing
	err := d.orm.Model(&Pricing{}).Preload("UserStatus").Find(&pricing).Error
	return pricing, err
}
