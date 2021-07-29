package models

import "github.com/amido/stacks-cli/internal/constants"

// Config is used to map the configuration onto the application models
type Config struct {

	// State if running in Interactive mode
	Interactive bool `mapstructure:interactive`

	// Version of the application
	Version string

	// Define the logging parameters
	Log Log `mapstructure:log`
}

// GetVersion returns the current version of the application
// It will check to see uif the Version is empty, if it is, it will
// set and identifiable local build version
func (config *Config) GetVersion() string {
	var version string

	version = config.Version

	if version == "" {
		version = constants.DefaultVersion
	}

	return version
}
