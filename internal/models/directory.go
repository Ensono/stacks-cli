package models

type Directory struct {
	WorkingDir string `mapstructure:"working"`
	TempDir    string `mapstructure:"temp"`
}
