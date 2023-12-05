package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Ensono/stacks-cli/internal/constants"
	"github.com/spf13/viper"
)

type Project struct {
	Name          string        `mapstructure:"name"`
	Framework     Framework     `mapstructure:"framework"`
	Platform      Platform      `mapstructure:"platform" json:",omitempty"`
	SourceControl SourceControl `mapstructure:"sourcecontrol"`
	SettingsFile  string        `mapstructure:"settingsfile" json:",omitempty"`
	Cloud         Cloud         `mapstructure:"cloud"`

	Directory Directory `yaml:"-"` // Holds the workingdir and tempdir for the project

	Settings Settings `yaml:"-"` // Hold the settings for the current project

	Phases []Phase `yaml:"-"` // Holds the phases for the operations
}

// GetId returns a consistent identifier for the name of the project
// It will change all to lowercase and replace spaces with an "_"
func (project *Project) GetId() string {
	return strings.Replace(strings.ToLower(project.Name), " ", "_", -1)
}

// ReadSettings reads in the settings file for the current project
// Returns the path to the file that was read for the project settings and any errors
// that were raised
// The operations are also read in and the phases object is created
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

	// create the phases of the project
	project.Phases = []Phase{
		{
			Name:       "init",
			Directory:  path,
			Operations: project.Settings.Init.Operations,
		},
		{
			Name:       "setup",
			Directory:  project.Directory.WorkingDir,
			Operations: project.Settings.Setup.Operations,
		},
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

	// determine the path for the settings file
	// if it is empty then set as the path + `stackscli.yml`
	// if it is not empty and it is not absolute prepend the path to it
	// check that the path exists
	if settingsFilePath == "" {
		settingsFilePath = filepath.Join(path, constants.SettingsFile)
	} else if !filepath.IsAbs(settingsFilePath) {
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
