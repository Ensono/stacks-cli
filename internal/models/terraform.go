package models

type Terraform struct {
	Backend TerraformBackend `mapstructure:"backend"`
}
