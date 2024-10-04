package config

type FrameworkDefVersion struct {
	Arguments  string `mapstructure:"arguments" yaml:"arguments"`
	Pattern    string `mapstructure:"pattern" yaml:"pattern"`
	Comparator string `mapstructure:"comparator" yaml:"comparator"`
}
