package config

type Cloud struct {
	Platform      string `mapstructure:"platform" yaml:",omitempty"`
	Region        string `mapstructure:"region" yaml:",omitempty"`
	ResourceGroup string `mapstructure:"group" yaml:"group,omitempty"`
}
