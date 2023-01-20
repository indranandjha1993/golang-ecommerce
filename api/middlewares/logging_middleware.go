package middlewares

import (
	"fmt"

	"github.com/labstack/echo"
)

func LogRoutes(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Routes: ")
		for _, route := range c.Echo().Routes() {
			fmt.Printf("Path: %s, Method: %s\n", route.Path, route.Method)
		}
		return next(c)
	}
}
