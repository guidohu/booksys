package database

// GetPricings returns the pricing of every user group.
func (d *Mysql) GetPricings() ([]Pricing, error) {
	var pricing []Pricing
	err := d.orm.Model(&Pricing{}).
		Preload("UserStatus").
		Preload("UserStatus.UserRole").
		Find(&pricing).Error
	return pricing, err
}

// GetUserStatusToPricingsMap returns the pricing of every user group, keyed by
// user group ID.
func (d *Mysql) GetUserStatusToPricingsMap() (map[uint]Pricing, error) {
	pMap := make(map[uint]Pricing)
	pricings, err := d.GetPricings()
	if err != nil {
		return pMap, err
	}
	for _, p := range pricings {
		pMap[p.UserStatusID] = p
	}
	return pMap, nil
}
