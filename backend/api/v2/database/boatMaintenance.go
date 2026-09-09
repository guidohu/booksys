package database

// GetMaintenance returns the maintenance log, most recent first.
func (d *Mysql) GetMaintenance() ([]BoatMaintenance, error) {
	var r []BoatMaintenance
	err := d.orm.Model(&BoatMaintenance{}).Preload("User").Order("timestamp DESC").Find(&r).Error
	return r, err
}

// AddMaintenanceEntry records a maintenance job.
func (d *Mysql) AddMaintenanceEntry(m BoatMaintenance) error {
	return d.orm.Create(&m).Error
}
