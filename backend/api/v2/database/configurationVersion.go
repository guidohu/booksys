package database

// GetConfigurationVersion returns the current configuration version.
func (d *Mysql) GetConfigurationVersion() (ConfigurationVersion, error) {
	c := &ConfigurationVersion{}
	err := d.orm.Last(c).Error
	return *c, err
}
