package config

type Overrides struct {
	InternalConfigPath string `mapstructure:"internal_config" yaml:"internal_config"`
	InternalConfigURL  string `mapstructure:"internal_config_url" yaml:"internal_config_url"`
	AdoVariablesPath   string `mapstructure:"ado_variables_path" yaml:"ado_variables_path"`
}
