package utils

import (
	"errors"
	"golang-ecommerce/internal/models"
	"image"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strconv"
	"unicode"

	"gopkg.in/go-playground/validator.v9"
)

func ValidateNonEmpty(s string) error {
	if s == "" {
		return errors.New("cannot be empty")
	}
	return nil
}

func ValidateMinLength(s string, min int) error {
	if len(s) < min {
		return errors.New("must be at least " + strconv.Itoa(min) + " characters")
	}
	return nil
}

func ValidateMaxLength(s string, max int) error {
	if len(s) > max {
		return errors.New("cannot be more than " + strconv.Itoa(max) + " characters")
	}
	return nil
}

func ValidateMobile(mobile string) error {
	re := regexp.MustCompile(`^(1\s?)?((\([0-9]{3}\))|[0-9]{3})[\s\-]?[\0-9]{3}[\s\-]?[0-9]{4}$`)
	if !re.MatchString(mobile) {
		return errors.New("invalid mobile number")
	}
	return nil
}

func ValidateURL(url string) error {
	re := regexp.MustCompile(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)?[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)
	if !re.MatchString(url) {
		return errors.New("invalid URL")
	}
	return nil
}

func ValidateEmail(email string) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return regex.MatchString(email)
}

func ValidateImage(file multipart.File, header *multipart.FileHeader) error {
	// check file size
	if header.Size > 5000000 {
		return errors.New("file is too large (max 5MB)")
	}

	// check file extension
	ext := filepath.Ext(header.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		return errors.New("invalid file type (must be jpg, jpeg, png or gif)")
	}

	// check image resolution
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return errors.New("could not decode image")
	}
	if img.Width < 800 || img.Height < 600 {
		return errors.New("image resolution too low (must be at least 800x600)")
	}

	return nil
}
func ValidatePassword(password string) error {
	// Check if password is at least 8 characters
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Check if password contains at least 1 uppercase letter
	hasUpper := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return errors.New("password must contain at least 1 uppercase letter")
	}

	// Check if password contains at least 1 lowercase letter
	hasLower := false
	for _, char := range password {
		if unicode.IsLower(char) {
			hasLower = true
			break
		}
	}
	if !hasLower {
		return errors.New("password must contain at least 1 lowercase letter")
	}

	// Check if password contains at least 1 number
	hasNumber := false
	for _, char := range password {
		if unicode.IsNumber(char) {
			hasNumber = true
			break
		}
	}
	if !hasNumber {
		return errors.New("password must contain at least 1 number")
	}

	// Check if password contains at least 1 special character
	re := regexp.MustCompile(`[!@#\$%^&*(),.?":{}|<>]`)
	hasSpecial := re.MatchString(password)
	if !hasSpecial {
		return errors.New("password must contain at least 1 special character")
	}

	// If all checks pass, password is considered secure
	return nil
}

// ValidateUpdateUser validates the fields of a user struct
func ValidateUpdateUser(user models.User) ErrorWarper {

	var e ErrorWarper

	// Validate name
	if user.Name != "" {
		if err := ValidateMinLength(user.Name, 3); err != nil {
			e.Add("name", err.Error())
		}
		if err := ValidateMaxLength(user.Name, 100); err != nil {
			e.Add("name", err.Error())
		}
	}
	// Validate email
	if user.Email != "" {
		if !ValidateEmail(user.Email) {
			e.Add("email", "invalid email id")
		}
	}
	// Validate password
	if user.Password != "" {
		if err := ValidatePassword(user.Password); err != nil {
			e.Add("password", err.Error())
		}
	}
	return e
}

func ValidateRequestData(req interface{}) ErrorWarper {
	var errs ErrorWarper
	v := validator.New()
	if err := v.Struct(req); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errs.Add(err.Field(), err.Tag())
		}
	}
	return errs
}
