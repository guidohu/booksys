package database

func (d *DBMysql) GetFuelEntries() ([]BoatFuel, error) {
	var b []BoatFuel
	err := d.orm.Model(&BoatFuel{}).Preload("User").Order("timestamp DESC").Find(&b).Error
	return b, err
}

func (d *DBMysql) GetFuelEntry(id uint) (BoatFuel, error) {
	var b BoatFuel
	err := d.orm.Model(&BoatFuel{}).
		Where("id = ?", id).
		Preload("User").
		First(&b).
		Error
	return b, err
}

func (d *DBMysql) AddFuelEntry(e BoatFuel) error {
	return d.orm.Create(&e).Error
}

func (d *DBMysql) ChangeFuelEntry(e BoatFuel) error {
	return d.orm.Save(&e).Error
}
