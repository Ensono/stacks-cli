package models

import (
	"errors"
	"os"
	"strings"
)

type Platform struct{}

// IsMingW returns a boolean value stating if the
// the application is running under the Mingw64 system
//
// This is essential when checking the environments in order
// to determine if the interactive component of the app will work
func (p Platform) isMingW() bool {
	var result bool

	// Get the environment variable MSYSTEM and check its value if it exists
	value, present := os.LookupEnv("MSYSTEM")

	// if the variable is present check the value
	if present && strings.ToLower(value) == "mingw64" {
		result = true
	}

	return result
}

func (p Platform) RunEnvironmentChecks() error {
	var err error = nil

	// determine if the shell is MINGW64
	if p.isMingW() {
		err = errors.New("Unable to run interactive session in your shell as it is unsupported [MINGW64]")
	}

	return err
}
