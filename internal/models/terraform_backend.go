package models

type TerraformBackend struct {
	Storage       string `mapstructure:"storage"`
	ResourceGroup string `mapstructure:"group"`
	Container     string `mapstructure:"container"`
}
