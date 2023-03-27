package config

type Package struct {
	Name    string `mapstructure:"name" yaml:"name"`
	Type    string `mapstructure:"type" yaml:"type"`
	URL     string `mapstructure:"url" yaml:"url"`
	Version string `mapstructure:"version" yaml:"version"`
	ID      string `mapstructure:"id" yaml:"id"`
}
