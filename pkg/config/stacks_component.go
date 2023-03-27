package config

import "fmt"

type StacksComponent struct {
	Group   string  `mapstructure:"group" yaml:"group"`
	Name    string  `mapstructure:"name" yaml:"name"`
	Package Package `mapstructure:"package" yaml:"package"`
}

// GetName returns the name of the component by combining the group and the name
func (sc *StacksComponent) GetName() string {
	return fmt.Sprintf("%s_%s", sc.Group, sc.Name)
}
