package database

func (d *DBMysql) GetPropertyValue(key string) (Configuration, error) {
	var value Configuration
	err := d.orm.Where("property = ?", key).First(&value).Error
	return value, err
}
