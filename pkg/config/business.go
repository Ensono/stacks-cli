package config

type Business struct {
	Company   string `mapstructure:"company" yaml:",omitempty"`
	Project   string `mapstructure:"project" json:",omitempty" yaml:",omitempty"`
	Domain    string `mapstructure:"domain" yaml:",omitempty"`
	Component string `mapstructure:"component" yaml:",omitempty"`
}
