package config

type Package struct {
	ID      string `mapstructure:"id" yaml:"id"`
	Name    string `mapstructure:"name" yaml:"name"`
	Path    string `mapstructure:"path" yaml:"path"`
	Type    string `mapstructure:"type" yaml:"type"`
	URL     string `mapstructure:"url" yaml:"url"`
	Version string `mapstructure:"version" yaml:"version"`
}
