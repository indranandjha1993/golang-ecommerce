package database

import (
	"golang-ecommerce/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	// Test
	db, err := Connect(config.GetTestConfig())

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Close()
}
