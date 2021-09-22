package config

import "github.com/amido/stacks-cli/internal/models"

// Config is used to map the configuration onto the application models
type InputConfig struct {

	// State if running in Interactive mode
	Interactive bool `mapstructure:"interactive"`

	// Version of the application
	Version string

	// Define the logging parameters
	Log models.Log `mapstructure:"log"`

	Directory Directory `mapstructure:"directory"`

	Business     Business  `mapstructure:"business"`
	Cloud        Cloud     `mapstructure:"cloud"`
	Network      Network   `mapstructure:"network"`
	Pipeline     string    `mapstructure:"pipeline"`
	Project      []Project `mapstructure:"project"`
	Stacks       Stacks    `mapstructure:"stacks"` // Holds the information about the projects in stacks
	Terraform    Terraform `mapstructure:"terraform"`
	SettingsFile string    `mapstructure:"settingsfile"`
}
