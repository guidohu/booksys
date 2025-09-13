package database

func (d *Mysql) GetFuelEntries() ([]BoatFuel, error) {
	var b []BoatFuel
	err := d.orm.Model(&BoatFuel{}).Preload("User").Order("timestamp DESC").Find(&b).Error
	return b, err
}

func (d *Mysql) GetFuelEntry(id uint) (BoatFuel, error) {
	var b BoatFuel
	err := d.orm.Model(&BoatFuel{}).
		Where("id = ?", id).
		Preload("User").
		First(&b).
		Error
	return b, err
}

func (d *Mysql) AddFuelEntry(e BoatFuel) error {
	return d.orm.Create(&e).Error
}

func (d *Mysql) ChangeFuelEntry(e BoatFuel) error {
	return d.orm.Save(&e).Error
}

func (d *Mysql) RemoveFuelEntry(id uint) error {
	b := &BoatFuel{
		ID: id,
	}
	return d.orm.Delete(b).Error
}
