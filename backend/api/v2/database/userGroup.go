package database

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (d *Mysql) CreateUserGroup(us UserStatus, p Pricing) error {
	return d.orm.Transaction(func(tx *gorm.DB) error {
		// check that group does not exist already
		var groups []UserStatus
		err := tx.Where("id = ? OR name = ?", us.ID, us.Name).Find(&groups).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if len(groups) > 0 {
			return fmt.Errorf("user group already exists with that ID or name")
		}

		// check that the user role exists
		var roles []UserRole
		err = tx.Debug().Where("id = ?", us.UserRoleID).Find(&roles).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if len(roles) != 1 {
			return fmt.Errorf("user role with id %d does not exist", us.UserRoleID)
		}

		// add the user status (aka user group)
		err = tx.Create(&us).Error
		if err != nil {
			return err
		}

		// add the pricing for the user group
		p.UserStatusID = us.ID
		err = tx.Create(&p).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (d *Mysql) ChangeUserGroup(us UserStatus, p Pricing) error {
	return d.orm.Transaction(func(tx *gorm.DB) error {
		// check that group does exist
		var group UserStatus
		err := tx.Where("id = ?", us.ID).First(&group).Error
		if err != nil {
			return err
		}
		// check that new name does not collide with other group
		if group.Name != us.Name {
			var collidingGroup UserStatus
			err = tx.Where("name = ?", us.Name).First(&collidingGroup).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if collidingGroup.ID != 0 {
				return fmt.Errorf("a group with this name already exists")
			}
		}
		// check that the price ID already exists
		var price Pricing
		err = tx.Where("id = ?", p.ID).First(&price).Error
		if err != nil {
			return err
		}
		// check that the user role exists
		var role UserRole
		err = tx.Where("id = ?", us.UserRoleID).First(&role).Error
		if err != nil {
			return err
		}

		// change the user status (aka user group)
		err = tx.Save(&us).Error
		if err != nil {
			return err
		}

		// change the pricing for the user group
		err = tx.Save(&p).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (d *Mysql) DeleteUserGroup(id uint) error {
	return d.orm.Transaction(func(tx *gorm.DB) error {
		// check that group does exist
		var group UserStatus
		err := tx.Where("id = ?", id).First(&group).Error
		if err != nil {
			return err
		}
		// get price ID
		var price Pricing
		err = tx.Where("user_status_id", id).First(&price).Error
		if err != nil {
			return err
		}
		// get assigned users
		var user User
		err = tx.Where("status = ?", id).First(&user).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if user.ID != 0 {
			return fmt.Errorf("user with ID %d is still member of this group", user.ID)
		}

		// delete entries
		err = tx.Delete(&price).Error
		if err != nil {
			return err
		}
		err = tx.Delete(&group).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (d *Mysql) SetUserGroup(userID uint, groupID uint) error {
	return d.orm.Model(&User{}).Where("id = ?", userID).Update("status", groupID).Error
}
