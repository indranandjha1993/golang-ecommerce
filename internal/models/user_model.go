package models

import (
	"github.com/jinzhu/gorm"
)

// User struct
type User struct {
	gorm.Model
	Name       string      `gorm:"type:varchar(100);not null"`
	Email      string      `gorm:"type:varchar(100);not null;unique"`
	Password   string      `gorm:"type:varchar(255);not null"`
	Role       string      `gorm:"type:varchar(20);not null"`
	UserTokens []UserToken `gorm:"ForeignKey:UserID"`
}

type UserToken struct {
	gorm.Model
	UserID    uint `gorm:"not null"`
	User      User
	Token     string `gorm:"type:varchar(255);not null"`
	ExpiresAt int64  `gorm:"not null"`
	Revoked   bool   `gorm:"not null"`
	LoginInfo string `gorm:"type:text;"`
}
