package config

// Settings holds the settings for each project as read from
// the `stackscli.yml` file in the project
type Settings struct {
	Framework string `mapstructure:"framework"`
	Init      Init   `mapstructure:"init"`
	Setup     Setup  `mapstructure:"setup"`
}

// Init holds the operations that should be performed before any work
// is done on the working project
type Init struct {
	Operations []Operation `mapstructure:"operations"`
}

// Setup holds the operaions that should be performed after the projet
// has been added to the working directory
type Setup struct {
	Operations []Operation `mapstructure:"operations"`
}

type Operation struct {
	Action      string `mapstructure:"action"`
	Command     string `mapstructure:"cmd"`
	Arguments   string `mapstructure:"args"`
	Description string `mapstructure:"desc"`
}
