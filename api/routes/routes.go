package routes

import (
	"fmt"
	"golang-ecommerce/config"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func NewRouter(version string, db *gorm.DB, config *config.Config) *echo.Echo {
	// Create an Echo instance
	e := echo.New()

	api := e.Group("/api")
	api = api.Group("/" + version)

	if version == "v1" {
		fmt.Println("Its working")
		V1(api, db, config)
	} else {
		api.Any("/*", func(c echo.Context) error {
			fmt.Println("Its not working")
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Not Found!",
			})
		})
	}

	return e
}
