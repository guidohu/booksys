package database

func (d *DBMysql) GetMaintenance() ([]BoatMaintenance, error) {
	r := []BoatMaintenance{}
	err := d.orm.Model(&BoatMaintenance{}).Preload("User").Order("timestamp DESC").Find(&r).Error
	return r, err
}

func (d *DBMysql) AddMaintenanceEntry(m BoatMaintenance) error {
	return d.orm.Create(&m).Error
}
