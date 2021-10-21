package config

type Directory struct {
	WorkingDir string `mapstructure:"working" yaml:"working"`
	TempDir    string `mapstructure:"temp" yaml:"-"`
}
