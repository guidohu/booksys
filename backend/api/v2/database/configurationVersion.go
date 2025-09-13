package database

func (d *Mysql) GetConfigurationVersion() (ConfigurationVersion, error) {
	c := &ConfigurationVersion{}
	err := d.orm.Last(c).Error
	return *c, err
}
