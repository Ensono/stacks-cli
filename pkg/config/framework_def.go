package config

type FrameworkDef struct {
	Name     string              `mapstructure:"name"`
	Commands []string            `mapstructure:"commands"`
	Version  FrameworkDefVersion `mapstructure:"version"`
}
