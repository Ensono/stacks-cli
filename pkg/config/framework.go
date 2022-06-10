package config

import "fmt"

type Framework struct {
	Type       string   `mapstructure:"type"`
	Option     string   `mapstructure:"option"`
	Version    string   `mapstructure:"version" yaml:",omitempty"`    // Version of the project to download
	Properties []string `mapstructure:"properties" yaml:",omitempty"` // additional properties to be specified that need to be passed to project commands
}

// GetMapKey returns the key to be used in the srcUrl map to
// get the URL for cloning the repository
func (framework *Framework) GetMapKey() string {
	return fmt.Sprintf("%s_%s", framework.Type, framework.Option)
}
