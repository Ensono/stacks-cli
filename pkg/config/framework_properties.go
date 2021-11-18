package config

type FrameworkProperties struct {
	Prop1 string `mapstructure:"prop1" yaml:",omitempty"`
	Prop2 string `mapstructure:"prop2" yaml:",omitempty"`
	Prop3 string `mapstructure:"prop3" yaml:",omitempty"`
	Prop4 string `mapstructure:"prop4" yaml:",omitempty"`
	Prop5 string `mapstructure:"prop5" yaml:",omitempty"`
}
