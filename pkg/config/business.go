package config

type Business struct {
	Company   string `mapstructure:"company"`
	Project   string `mapstructure:"project" json:",omitempty"`
	Domain    string `mapstructure:"domain"`
	Component string `mapstructure:"component"`
}
