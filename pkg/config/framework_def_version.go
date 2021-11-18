package config

type FrameworkDefVersion struct {
	Command   string `mapstructure:"command"`
	Arguments string `mapstructure:"arguments"`
	Pattern   string `mapstructure:"pattern"`
}
