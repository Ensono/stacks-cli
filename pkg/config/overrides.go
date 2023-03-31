package config

type Overrides struct {
	InternalConfigPath string `mapstructure:"internal_config" yaml:"internal_config"`
	AdoVariablesPath   string `mapstructure:"ado_variables_path" yaml:"ado_variables_path"`
}
