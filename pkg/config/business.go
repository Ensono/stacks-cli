package config

type Business struct {
	Company   string `mapstructure:"company"`
	Project   string `mapstructure:"project"`
	Domain    string `mapstructure:"domain"`
	Component string `mapstructure:"component"`
}
