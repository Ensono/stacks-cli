package config

import (
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	Name          string        `mapstructure:"name"`
	Framework     Framework     `mapstructure:"framework"`
	Platform      Platform      `mapstructure:"platform"`
	SourceControl SourceControl `mapstructure:"sourcecontrol"`
	SettingsFile  string        `mapstructure:"settings_file"`

	Paths map[string]string // define a map to hold the paths that are created for each project

	Settings Settings // Hold the settings for the current project
}

// GetId returns a consistent identifier for the name of the project
// It will change all to lowercase and replace spaces with an "_"
func (project *Project) GetId() string {
	return strings.Replace(strings.ToLower(project.Name), " ", "_", -1)
}

// ReadSettings reads in the settings file for the current project
func (project *Project) ReadSettings(path string, config *Config) error {

	// get the path to the settings file
	err := project.getSettingsFilePath(path, config)
	if err != nil {
		return err
	}

	return nil
}

func (project *Project) getSettingsFilePath(path string, config *Config) error {

	// set the default value of the filename
	filename := config.Input.SettingsFile

	// if the SettingsFile has been set on the project then use that
	if project.SettingsFile != "" {
		filename = project.SettingsFile
	}

	// set the path to the file
	settingsFilePath := filepath.Join(path, filename)

	// determine if the file exists, if it does then set the SettingsFile
	// otherwise return an error
	if _, err := os.Stat(settingsFilePath); os.IsNotExist(err) {
		return err
	}

	project.SettingsFile = settingsFilePath

	return nil
}
