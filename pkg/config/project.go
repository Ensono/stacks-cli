package config

type Project struct {
	Name          string        `mapstructure:"name"`
	Framework     Framework     `mapstructure:"framework"`
	Platform      Platform      `mapstructure:"platform"`
	SourceControl SourceControl `mapstructure:"sourcecontrol"`
}
