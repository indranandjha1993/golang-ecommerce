package handlers

import (
	"fmt"
	"golang-ecommerce/config"
	"golang-ecommerce/constants"
	"golang-ecommerce/internal/models"
	"golang-ecommerce/services/jwt"
	"golang-ecommerce/utils"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type AuthHandler struct {
	DB     *gorm.DB
	Config *config.Config
	JWT    *jwt.JWT
}

func (a *AuthHandler) Login(c echo.Context) error {
	user := new(models.User)
	// Get the request data
	req := new(struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	})

	if err := c.Bind(req); err != nil {
		return utils.WriteErrorResponse(c, constants.BadRequest, err, http.StatusBadRequest)
	}

	errs := utils.ValidateRequestData(req)
	if len(errs.Errors) > 0 {
		return utils.WriteErrorsResponse(c, constants.BadRequest, errs, http.StatusBadRequest)
	}

	// Find the user
	if err := a.DB.Where("email = ?", req.Email).First(user).Error; err != nil {
		return utils.WriteErrorResponse(c, constants.Unauthorized, err, http.StatusUnauthorized)
	}

	// Check the password
	if err := utils.CompareHashAndPassword(user.Password, req.Password); err != nil {
		return utils.WriteErrorResponse(c, constants.Unauthorized, err, http.StatusUnauthorized)
	}

	// Generate the token
	timeDur, err := a.Config.GetJWTExpiresIn()
	if err != nil {
		return utils.WriteErrorResponse(c, constants.ErrorMessage, err, http.StatusInternalServerError)
	}
	fmt.Printf("timeDur: %v", timeDur)

	token, err := a.JWT.GenerateToken(user.ID, timeDur)
	if err != nil {
		return utils.WriteErrorResponse(c, constants.ErrorMessage, err, http.StatusInternalServerError)
	}

	return utils.WriteSuccessResponse(c, constants.SuccessMessage, token, http.StatusOK)
}

func (a *AuthHandler) Logout(c echo.Context) error {
	user := c.Get("user").(*models.User)
	authToken := c.Get("authToken")

	// Revoke the token
	var userToken models.UserToken
	if err := a.DB.Model(&userToken).Where("user_id = ?", user.ID).Where("token = ?", authToken).Update("revoked", true).Error; err != nil {
		return utils.WriteErrorResponse(c, constants.ErrorMessage, err, http.StatusInternalServerError)
	}

	return utils.WriteSuccessResponse(c, constants.SuccessMessage, "", http.StatusOK)
}

func (a *AuthHandler) Profile(c echo.Context) error {
	user := c.Get("user").(*models.User)

	return utils.WriteSuccessResponse(c, constants.SuccessMessage, user, http.StatusOK)
}
