package seeders

import (
	"golang-ecommerce/internal/models"

	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{
		Name:     "John Doe",
		Email:    "johndoe@example.com",
		Password: "password",
		Role:     "admin",
	},
	{
		Name:     "Jane Smith",
		Email:    "janesmith@example.com",
		Password: "password",
		Role:     "customer",
	},
}

func seedUsers(db *gorm.DB) error {
	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}
	return nil
}
