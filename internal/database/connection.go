package database

import (
	"golang-ecommerce/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //import mysql dialect
)

// DB struct
type DB struct {
	*gorm.DB
}

// Connect function to connect to the database
func Connect(config *config.Config) (*DB, error) {
	// create a connection string
	db, err := gorm.Open("mysql", config.DSN())
	if err != nil {
		return nil, err
	}

	// return a new DB struct
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
