package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDontRunInteractive(t *testing.T) {

	// create a configuration object
	config := Config{}
	answers := Answers{}

	// set that the app should not run interactively
	config.Input.Interactive = false

	// call the RunInteractive method
	err := answers.RunInteractive(&config)

	// check that the err is nil
	assert.Equal(t, nil, err)
}
