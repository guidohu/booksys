package database

func (d *Mysql) GetEngineHourLatest() (BoatEngineHour, error) {
	var b BoatEngineHour
	err := d.orm.Model(&BoatEngineHour{}).
		Preload("User").
		Preload("Type").
		Order("timestamp DESC").
		First(&b).Error
	return b, err
}

func (d *Mysql) GetEngineHoursEntry(id uint) (BoatEngineHour, error) {
	var b BoatEngineHour
	err := d.orm.Model(&BoatEngineHour{}).Where("id = ?", id).Preload("User").First(&b).Error
	return b, err
}

func (d *Mysql) GetEngineHours() ([]BoatEngineHour, error) {
	var b []BoatEngineHour
	err := d.orm.Model(&BoatEngineHour{}).
		Preload("User").
		Preload("Type").
		Order("timestamp DESC").
		Find(&b).Error
	return b, err
}

func (d *Mysql) AddEngineHours(b BoatEngineHour) error {
	return d.orm.Create(&b).Error
}

func (d *Mysql) UpdateEngineHours(b BoatEngineHour) error {
	return d.orm.Save(&b).Error
}
