package handlers

import (
	"golang-ecommerce/constants"
	"golang-ecommerce/internal/models"
	"golang-ecommerce/utils"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type UsersHandler struct {
	DB *gorm.DB
}

func (u *UsersHandler) Create(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return utils.WriteErrorResponse(c, constants.BadRequest, err, http.StatusBadRequest)
	}

	var errs utils.ErrorWarper
	if err := utils.ValidateMinLength(user.Name, 3); err != nil {
		errs.Add("name", err.Error())
	}
	if err := utils.ValidateMaxLength(user.Name, 100); err != nil {
		errs.Add("name", err.Error())
	}
	if err := utils.ValidateNonEmpty(user.Name); err != nil {
		errs.Add("name", err.Error())
	}
	if !utils.ValidateEmail(user.Email) {
		errs.Add("email", "invalid email id")
	}
	if err := utils.ValidateNonEmpty(user.Email); err != nil {
		errs.Add("email", err.Error())
	}
	if err := utils.ValidatePassword(user.Password); err != nil {
		errs.Add("password", err.Error())
	}
	if err := utils.ValidateNonEmpty(user.Password); err != nil {
		errs.Add("password", err.Error())
	}
	if len(errs.Errors) > 0 {
		return utils.WriteErrorsResponse(c, constants.BadRequest, errs, http.StatusBadRequest)
	}

	password, _ := utils.HashPassword(user.Password)
	user.Password = password

	if err := u.DB.Create(user).Error; err != nil {
		return utils.WriteErrorResponse(c, constants.ErrorMessage, err, http.StatusInternalServerError)
	}

	return utils.WriteSuccessResponse(c, constants.CreatedMessage, user, http.StatusCreated)
}

func (u *UsersHandler) Get(c echo.Context) error {
	id := c.Param("id")

	var user models.User
	if err := u.DB.Find(&user, id).Error; err != nil {
		return utils.WriteErrorResponse(c, constants.NotFound, err, http.StatusNotFound)
	}
	return utils.WriteSuccessResponse(c, constants.SuccessMessage, user, http.StatusOK)
}

func (u *UsersHandler) List(c echo.Context) error {
	var users []models.User
	users_list, err := utils.Paginate(c, u.DB, &users)
	if err != nil {
		return utils.WriteErrorResponse(c, constants.ErrorMessage, err, http.StatusInternalServerError)
	}

	return utils.WriteSuccessResponse(c, constants.SuccessMessage, users_list, http.StatusOK)
}

func (u *UsersHandler) Update(c echo.Context) error {
	id := c.Param("id")
	user := new(models.User)
	if err := u.DB.Find(&user, id).Error; err != nil {
		return utils.WriteErrorResponse(c, constants.NotFound, err, http.StatusNotFound)
	}

	if err := c.Bind(user); err != nil {
		return utils.WriteErrorResponse(c, constants.NotFound, err, http.StatusNotFound)
	}

	if errs := utils.ValidateUpdateUser(*user); len(errs.Errors) > 0 {
		return utils.WriteErrorsResponse(c, constants.BadRequest, errs, http.StatusBadRequest)
	}

	if user.Password != "" {
		password, _ := utils.HashPassword(user.Password)
		user.Password = password
	}

	tx := u.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return utils.WriteErrorResponse(c, constants.ErrorMessage, err, http.StatusInternalServerError)
	}

	if user.Name != "" {
		tx.Model(&user).Update("name", user.Name)
	}
	if user.Email != "" {
		tx.Model(&user).Update("email", user.Email)
	}
	if user.Password != "" {
		tx.Model(&user).Update("password", user.Password)
	}

	if err := tx.Commit().Error; err != nil {
		return utils.WriteErrorResponse(c, constants.ErrorMessage, err, http.StatusInternalServerError)
	}

	return utils.WriteSuccessResponse(c, constants.UpdatedMessage, user, http.StatusOK)
}

func (u *UsersHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	user := new(models.User)
	if err := u.DB.Find(&user, id).Error; err != nil {
		return utils.WriteErrorResponse(c, constants.NotFound, err, http.StatusNotFound)
	}

	if err := u.DB.Delete(&user).Error; err != nil {
		return utils.WriteErrorResponse(c, constants.ErrorMessage, err, http.StatusInternalServerError)
	}

	return utils.WriteSuccessResponse(c, constants.DeletedMessage, user, http.StatusOK)
}
