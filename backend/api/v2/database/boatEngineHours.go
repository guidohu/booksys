package database

// GetEngineHourLatest returns the most recent engine hour reading.
func (d *Mysql) GetEngineHourLatest() (BoatEngineHour, error) {
	var b BoatEngineHour
	err := d.orm.Model(&BoatEngineHour{}).
		Preload("User").
		Preload("Type").
		Order("timestamp DESC").
		First(&b).Error
	return b, err
}

// GetEngineHoursEntry returns a single engine hour reading.
func (d *Mysql) GetEngineHoursEntry(id uint) (BoatEngineHour, error) {
	var b BoatEngineHour
	err := d.orm.Model(&BoatEngineHour{}).Where("id = ?", id).Preload("User").First(&b).Error
	return b, err
}

// GetEngineHours returns all engine hour readings, most recent first.
func (d *Mysql) GetEngineHours() ([]BoatEngineHour, error) {
	var b []BoatEngineHour
	err := d.orm.Model(&BoatEngineHour{}).
		Preload("User").
		Preload("Type").
		Order("timestamp DESC").
		Find(&b).Error
	return b, err
}

// AddEngineHours records an engine hour reading.
func (d *Mysql) AddEngineHours(b BoatEngineHour) error {
	return d.orm.Create(&b).Error
}

// UpdateEngineHours writes back a modified engine hour reading.
func (d *Mysql) UpdateEngineHours(b BoatEngineHour) error {
	return d.orm.Save(&b).Error
}
