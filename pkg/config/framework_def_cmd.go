package config

type FrameworkDefCmd struct {
	Name    string              `mapstructure:"name" yaml:"name"`
	Version FrameworkDefVersion `mapstructure:"version" yaml:"version"`
}
