package database

func (d *DBMysql) GetPricings() ([]Pricing, error) {
	var pricing []Pricing
	err := d.orm.Model(&Pricing{}).
		Preload("UserStatus").
		Preload("UserStatus.UserRole").
		Find(&pricing).Error
	return pricing, err
}

func (d *DBMysql) GetUserStatusToPricingsMap() (map[uint]Pricing, error) {
	pMap := map[uint]Pricing{}
	pricings, err := d.GetPricings()
	if err != nil {
		return pMap, err
	}
	for _, p := range pricings {
		pMap[p.UserStatusID] = p
	}
	return pMap, nil
}
