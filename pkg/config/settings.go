package config

// Settings holds the settings for each project as read from
// the `stackscli.yml` file in the project
type Settings struct {
	Framework string     `mapstructure:"framework"`
	Pipeline  []Pipeline `mapstructure:"pipeline"`
	Init      Init       `mapstructure:"init"`
	Setup     Setup      `mapstructure:"setup"`
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

// GetPipeline attempts to return the pipeline settings for the named pipeline
func (s *Settings) GetPipeline(name string) Pipeline {
	pipeline := Pipeline{}

	// iterate around the pipeline slice and find the one with the type that matches
	// the specified name
	for _, p := range s.Pipeline {
		if p.Type == name {
			pipeline = p
			break
		}
	}

	return pipeline
}
