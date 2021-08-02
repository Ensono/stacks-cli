package models

type Framework struct {
	Type    string `mapstructure:"type"`
	Option  string `mapstructure:"option"`
	Version string `mapstructure:"version"` // Version of the project to download
}
