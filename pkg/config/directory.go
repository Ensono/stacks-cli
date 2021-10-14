package config

type Directory struct {
	WorkingDir string `mapstructure:"working"`
	TempDir    string `mapstructure:"temp" yaml:"-"`
}
