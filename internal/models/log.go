package models

type Log struct {
	// Level for the logging, ERROR, WARN, INFO, DEBUG
	Level string `mapstructure:"level" json:",omitempty"`

	// Format for the logging, TEXT or JSON
	Format string `mapstructure:"format" json:",omitempty"`

	// If the logging is in TEXT then should colour be used
	Colour bool `mapstructure:"colour"`

	// File that all logs should be saved to
	File string `mapstructure:"file" json:",omitempty"`
}
