package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Project struct {
	Name          string        `mapstructure:"name"`
	Framework     Framework     `mapstructure:"framework"`
	Platform      Platform      `mapstructure:"platform"`
	SourceControl SourceControl `mapstructure:"sourcecontrol"`
	SettingsFile  string        `mapstructure:"settingsfile"`
	Cloud         Cloud         `mapstructure:"cloud"`

	Directory Directory // Holds the workingdir and tempdir for the project

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
	err := project.setSettingsFilePath(path, config)
	if err != nil {
		return err
	}

	// read in the settings file
	// using Viper as this supports multiple file formats
	v := viper.New()
	v.SetConfigFile(project.SettingsFile)
	err = v.ReadInConfig()
	if err != nil {
		return err
	}

	// unmarshal the data into the settings for the project
	err = v.Unmarshal(&project.Settings)
	if err != nil {
		return err
	}

	return nil
}

func (project *Project) setSettingsFilePath(path string, config *Config) error {

	// set the default value of the filename
	settingsFilePath := config.Input.SettingsFile

	// if the SettingsFile has been set on the project then use that
	if project.SettingsFile != "" {
		settingsFilePath = project.SettingsFile
	}

	// determine if the filename exists as it is set, if not then prepend
	// the path to the filename if the filename is not ab absolute path
	_, err := os.Stat(settingsFilePath)
	if os.IsNotExist(err) && !filepath.IsAbs(settingsFilePath) {
		settingsFilePath = filepath.Join(path, settingsFilePath)
	}

	// determine if the file exists, if it does then set the SettingsFile
	// otherwise return an error
	if _, err := os.Stat(settingsFilePath); os.IsNotExist(err) {
		return err
	}

	project.SettingsFile = settingsFilePath

	return nil
}
