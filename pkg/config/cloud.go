package config

type Cloud struct {
	Platform      string `mapstructure:"platform"`
	Region        string `mapstructure:"region"`
	ResourceGroup string `mapstructure:"group"`
}
