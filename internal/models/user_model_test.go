package models

import (
	"golang-ecommerce/config"
	"golang-ecommerce/internal/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserModel(t *testing.T) {
	// Connect to the test database
	db, err := database.Connect(config.GetTestConfig())
	assert.NoError(t, err)
	defer db.Close()

	// Create a user
	user := &User{
		Name:     "John Doe",
		Email:    "johndoe@example.com",
		Password: "password",
		Role:     "admin",
	}
	err = db.Create(user).Error
	assert.NoError(t, err)

	// Test to check if user is created
	var result User
	err = db.First(&result, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.Password, result.Password)
	assert.Equal(t, user.Role, result.Role)

	// Test to check if user can be updated
	user.Name = "Jane Doe"
	err = db.Save(user).Error
	assert.NoError(t, err)
	err = db.First(&result, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Name, result.Name)

	// Test to check if user can be deleted
	err = db.Delete(user).Error
	assert.NoError(t, err)
	err = db.First(&result, user.ID).Error
	assert.Error(t, err)
}
