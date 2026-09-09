package database

import (
	"time"

	"gorm.io/gorm"
)

// AddPasswordResetToken stores a password reset token for a user.
func (d *Mysql) AddPasswordResetToken(entry PasswordReset) error {
	return d.orm.Create(&entry).Error
}

// GetActivePasswordResetEntry returns the live reset token of a user, meaning
// the most recently issued one that is still flagged valid and has not
// expired. It returns gorm.ErrRecordNotFound when there is none.
//
// It deliberately does not take the submitted token as an argument. Looking the
// row up by user alone is what makes the attempt counter work: a wrong guess
// still resolves to a row, so it can be counted and the token burned.
func (d *Mysql) GetActivePasswordResetEntry(userID uint) (PasswordReset, error) {
	ret := PasswordReset{}
	err := d.orm.
		Where("user_id = ? AND valid = true AND valid_until > ?", userID, time.Now()).
		Order("issued_at desc").
		First(&ret).Error
	return ret, err
}

// RegisterFailedPasswordResetAttempt counts a wrong token against an entry and
// invalidates it once maxAttempts guesses have been used up. Both writes happen
// in one transaction so that concurrent guesses cannot slip past the limit
// between the increment and the check.
func (d *Mysql) RegisterFailedPasswordResetAttempt(id uint, maxAttempts uint) error {
	return d.orm.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&PasswordReset{}).
			Where("id = ?", id).
			Update("attempts", gorm.Expr("attempts + 1")).Error
		if err != nil {
			return err
		}
		return tx.Model(&PasswordReset{}).
			Where("id = ? AND attempts >= ?", id, maxAttempts).
			Update("valid", false).Error
	})
}

// InvalidatePasswordResetEntries invalidates every outstanding reset token of
// a user.
func (d *Mysql) InvalidatePasswordResetEntries(userID uint) error {
	return d.orm.Exec("UPDATE password_reset SET valid = false WHERE user_id = ?", userID).Error
}
