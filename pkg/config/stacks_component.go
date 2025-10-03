package config

import "fmt"

type StacksComponent struct {
	Group        string  `mapstructure:"group" yaml:"group"`
	Name         string  `mapstructure:"name" yaml:"name"`
	Package      Package `mapstructure:"package" yaml:"package"`
	TemplateMode *bool   `mapstructure:"template_mode" yaml:"template_mode"`
}

// GetName returns the name of the component by combining the group and the name
func (sc *StacksComponent) GetName() string {
	return fmt.Sprintf("%s_%s", sc.Group, sc.Name)
}

// IsTemplateModeEnabled returns the TemplateMode value, with a default of true if not set
func (sc *StacksComponent) IsTemplateModeEnabled() bool {
	if sc.TemplateMode == nil {
		return true // default value
	}
	return *sc.TemplateMode
}
