package database

func (d *Mysql) GetMaintenance() ([]BoatMaintenance, error) {
	r := []BoatMaintenance{}
	err := d.orm.Model(&BoatMaintenance{}).Preload("User").Order("timestamp DESC").Find(&r).Error
	return r, err
}

func (d *Mysql) AddMaintenanceEntry(m BoatMaintenance) error {
	return d.orm.Create(&m).Error
}
