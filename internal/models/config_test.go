package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersionWithDefaultVersionNumber(t *testing.T) {

	config := Config{}

	// state what is expected from the method
	expected := "0.0.1-workstation"

	// get the actual response
	actual := config.GetVersion()

	assert.Equal(t, actual, expected)
}

func TestGetVersion(t *testing.T) {
	config := Config{
		Version: "100.98.99",
	}

	// get the actual version
	actual := config.GetVersion()

	assert.Equal(t, actual, config.Version)
}
