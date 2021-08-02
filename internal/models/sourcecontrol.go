package models

type SourceControl struct {
	Type string `mapstructure:"type"`
	URL  string `mapstructure:"url"`
	Ref  string `mapstructure:"ref"`
}
