package config

type Directory struct {
	WorkingDir string `mapstructure:"working" yaml:"working"`
	TempDir    string `mapstructure:"temp" yaml:"-"`
	CacheDir   string `mapstructure:"cache" yaml:"-"`
	Export     string `mapstructure:"export"`
}
