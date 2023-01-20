package migrations

import (
	"golang-ecommerce/internal/models"

	"github.com/jinzhu/gorm"
)

func createUsersTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}).Error
}

func dropUsersTable(db *gorm.DB) error {
	return db.DropTable(&models.User{}).Error
}

func createUsersTokenTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.UserToken{}).Error
}

func dropUsersTokenTable(db *gorm.DB) error {
	return db.DropTable(&models.UserToken{}).Error
}
