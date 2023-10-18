package config

type Environment struct {
	Name      string `mapstructure:"name" yaml:",omitempty"`
	Type 	  string `mapstructure:"type" yaml:",omitempty"`
	DependsOn []string `mapstructure:"dependson" yaml:",omitempty"`
}
