package models

type Log struct {
	// Level for the logging, ERROR, WARN, INFO, DEBUG
	Level string `mapstructure:"level"`

	// Format for the logging, TEXT or JSON
	Format string `mapstructure:"format"`

	// If the logging is in TEXT then should colour be used
	Colour bool `mapstructure:"colour"`
}
