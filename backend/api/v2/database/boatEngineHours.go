package database

func (d *DBMysql) GetEngineHourLatest() (BoatEngineHour, error) {
	var b BoatEngineHour
	err := d.orm.Model(&BoatEngineHour{}).Preload("User").Order("timestamp DESC").First(&b).Error
	return b, err
}

func (d *DBMysql) GetEngineHoursEntry(id uint) (BoatEngineHour, error) {
	var b BoatEngineHour
	err := d.orm.Model(&BoatEngineHour{}).Where("id = ?", id).Preload("User").First(&b).Error
	return b, err
}

func (d *DBMysql) GetEngineHours() ([]BoatEngineHour, error) {
	var b []BoatEngineHour
	err := d.orm.Model(&BoatEngineHour{}).Preload("User").Order("timestamp DESC").Find(&b).Error
	return b, err
}

func (d *DBMysql) AddEngineHours(b BoatEngineHour) error {
	return d.orm.Create(&b).Error
}

func (d *DBMysql) UpdateEngineHours(b BoatEngineHour) error {
	return d.orm.Save(&b).Error
}
