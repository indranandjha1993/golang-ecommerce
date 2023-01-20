package main

import (
	"fmt"
	"golang-ecommerce/api/routes"
	"golang-ecommerce/config"
	"golang-ecommerce/internal/database"
	"golang-ecommerce/internal/migrations"
)

func main() {
	// Load config
	c, err := config.SetConfig("config.ini")
	if err != nil {
		fmt.Printf("Failed to get config: %v", err)
		return
	}

	// Connect to database
	db, err := database.Connect(c)
	if err != nil {
		fmt.Printf("Failed to connect DB: %v", err)
		return
	}
	defer db.Close()

	// Migrate database
	errs := migrations.Migrate(db.DB, "up")
	if len(errs.Errors) > 0 {
		fmt.Printf("Failed to migrate DB: %v", err)
		return
	}

	// Seed database
	// err = seeders.Seed(db.DB)
	// if err != nil {
	// 	fmt.Printf("Failed to seed DB: %v", err)
	// 	return
	// }

	e := routes.NewRouter("v1", db.DB, c)
	// e.Use(middlewares.LogRoutes)
	e.Logger.Fatal(e.Start(c.Server.Address + ":" + c.Server.Port))
}
