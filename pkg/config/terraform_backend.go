package config

type TerraformBackend struct {
	Storage   string `mapstructure:"storage" yaml:",omitempty"`
	Group     string `mapstructure:"group" yaml:",omitempty"`
	Container string `mapstructure:"container" yaml:",omitempty"`
}
