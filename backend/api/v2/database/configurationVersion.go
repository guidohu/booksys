package database

func (d *DBMysql) GetConfigurationVersion() (ConfigurationVersion, error) {
	c := &ConfigurationVersion{}
	err := d.orm.Last(c).Error
	return *c, err
}
