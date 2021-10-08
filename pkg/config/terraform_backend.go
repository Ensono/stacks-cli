package config

type TerraformBackend struct {
	Storage   string `mapstructure:"storage"`
	Group     string `mapstructure:"group"`
	Container string `mapstructure:"container"`
}
