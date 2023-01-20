package seeders

import "github.com/jinzhu/gorm"

func Seed(db *gorm.DB) error {
	err := seedUsers(db)
	if err != nil {
		return err
	}
	return nil
}
