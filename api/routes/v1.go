package routes

import (
	"golang-ecommerce/api/handlers/v1"
	"golang-ecommerce/api/middlewares"
	"golang-ecommerce/config"
	"golang-ecommerce/services/jwt"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func V1(e *echo.Group, db *gorm.DB, config *config.Config) *echo.Group {
	jwt := jwt.New(config.GetJWTSecret(), db)

	usersHandler := &handlers.UsersHandler{DB: db}
	authHandler := &handlers.AuthHandler{DB: db, Config: config, JWT: jwt}

	e.POST("/users", usersHandler.Create)
	e.GET("/users/:id", usersHandler.Get)
	e.GET("/users", usersHandler.List)
	e.PUT("/users/:id", usersHandler.Update)
	e.DELETE("/users/:id", usersHandler.Delete)

	e.POST("/login", authHandler.Login)
	e.POST("/logout", authHandler.Logout, middlewares.AuthMiddleware(config, db))
	e.GET("/profile", authHandler.Profile, middlewares.AuthMiddleware(config, db))

	return e
}
