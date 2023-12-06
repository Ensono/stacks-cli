package config

type Terraform struct {
	Backend TerraformBackend `mapstructure:"backend" yaml:",omitempty"`
}
