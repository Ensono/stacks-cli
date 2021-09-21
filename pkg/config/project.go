package config

import (
	"strings"
)

type Project struct {
	Name          string        `mapstructure:"name"`
	Framework     Framework     `mapstructure:"framework"`
	Platform      Platform      `mapstructure:"platform"`
	SourceControl SourceControl `mapstructure:"sourcecontrol"`

	// define a map to hold the paths that are created for each project
	Paths map[string]string
}

// GetId returns a consistent identifier for the name of the project
// It will change all to lowercase and replace spaces with an "_"
func (project *Project) GetId() string {
	return strings.Replace(strings.ToLower(project.Name), " ", "_", -1)
}
