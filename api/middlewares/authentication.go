package middlewares

import (
	"errors"
	"golang-ecommerce/config"
	"golang-ecommerce/constants"
	"golang-ecommerce/internal/models"
	"golang-ecommerce/services/jwt"
	"golang-ecommerce/utils"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func AuthMiddleware(config *config.Config, db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the token from the request
			authToken := c.Request().Header.Get("Authorization")
			if authToken == "" {
				return utils.WriteErrorResponse(c, constants.Unauthenticated, errors.New("authorization token is missing"), http.StatusUnauthorized)
			}
			// Parse the token
			jwt := jwt.New(config.GetJWTSecret(), db)
			claims, err := jwt.ParseToken(authToken)
			if err != nil {
				return utils.WriteErrorResponse(c, constants.Unauthorized, err, http.StatusUnauthorized)
			}

			// Get the user from the database
			user := new(models.User)
			if err := db.First(&user, claims.UserID).Error; err != nil {
				return utils.WriteErrorResponse(c, constants.Unauthorized, errors.New("invalid token"), http.StatusUnauthorized)
			}

			// Add the user to the context
			c.Set("user", user)
			c.Set("authToken", authToken)

			return next(c)
		}
	}
}
