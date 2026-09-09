package database

// GetFuelEntries returns the refuelling log, most recent first.
func (d *Mysql) GetFuelEntries() ([]BoatFuel, error) {
	var b []BoatFuel
	err := d.orm.Model(&BoatFuel{}).Preload("User").Order("timestamp DESC").Find(&b).Error
	return b, err
}

// GetFuelEntry returns a single refuelling.
func (d *Mysql) GetFuelEntry(id uint) (BoatFuel, error) {
	var b BoatFuel
	err := d.orm.Model(&BoatFuel{}).
		Where("id = ?", id).
		Preload("User").
		First(&b).
		Error
	return b, err
}

// AddFuelEntry records a refuelling.
func (d *Mysql) AddFuelEntry(e BoatFuel) error {
	return d.orm.Create(&e).Error
}

// ChangeFuelEntry writes back a modified refuelling.
func (d *Mysql) ChangeFuelEntry(e BoatFuel) error {
	return d.orm.Save(&e).Error
}

// RemoveFuelEntry removes a refuelling.
func (d *Mysql) RemoveFuelEntry(id uint) error {
	b := &BoatFuel{
		ID: id,
	}
	return d.orm.Delete(b).Error
}
