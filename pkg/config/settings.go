package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/amido/stacks-cli/internal/models"
	"github.com/amido/stacks-cli/internal/util"
	"github.com/sirupsen/logrus"
)

// Settings holds the settings for each project as read from
// the `stackscli.yml` file in the project
type Settings struct {
	Framework SettingsFramework `mapstructure:"framework"`
	Pipeline  []Pipeline        `mapstructure:"pipeline"`
	Init      Init              `mapstructure:"init"`
	Setup     Setup             `mapstructure:"setup"`
	Folders   []string          `mapstructure:"folders"`
}

// Init holds the operations that should be performed before any work
// is done on the working project
type Init struct {
	Operations []Operation `mapstructure:"operations"`
}

// Setup holds the operaions that should be performed after the projet
// has been added to the working directory
type Setup struct {
	Operations []Operation `mapstructure:"operations"`
}

type Operation struct {
	Action          string      `mapstructure:"action"`
	Command         string      `mapstructure:"cmd"`
	Arguments       string      `mapstructure:"args"`
	Description     string      `mapstructure:"desc"`
	ApplyProperties bool        `mapstructure:"applyProperties"`
	Tags            []string    `mapstructure:"tags"`
	Items           interface{} `mapstructure:"items"`
}

type SettingsFramework struct {
	Name     string                      `mapstructure:"name"`
	Commands []SettingsFrameworkCommands `mapstructure:"commands"`
}

type SettingsFrameworkCommands struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

// GetPipeline attempts to return the pipeline settings for the named pipeline
func (s *Settings) GetPipeline(name string) Pipeline {
	pipeline := Pipeline{}

	// iterate around the pipeline slice and find the one with the type that matches
	// the specified name
	for _, p := range s.Pipeline {
		if p.Type == name {
			pipeline = p
			break
		}
	}

	return pipeline
}

// GetRequiredVersion interrogates the list of framework versions to determine if
// any specific version is required for building the application
func (s *Settings) GetRequiredVersion(name string) string {

	version := ""

	// iterate around the list of framework versions and determine if the command
	// that has been specified has a required version number
	for _, v := range s.Framework.Commands {
		if v.Name == name {
			version = v.Version
		}
	}

	return version

}

// CheckCommandVersions checks all of the framework commands that have been specified
// and ensures that they are the correct version
//
// path - pa
func (s *Settings) CheckCommandVersions(config *Config, logger *logrus.Logger, path string, tmpPath string) ([]models.Command, string) {

	var err error
	var incorrect []models.Command
	var versionCmd string
	var versionArgs string
	var specificVersion string
	var re regexp.Regexp
	var info string

	// iterate around the framework commands
	for _, cmd := range s.Framework.Commands {

		// define the command to get the version of it
		switch cmd.Name {
		case "dotnet":

			versionCmd = "dotnet"
			versionArgs = "--info"
			re = *regexp.MustCompile(`\.NET.*SDK.*:\r?\n\sVersion:\s+(?P<version>.*?)\r?\n`)

			// check to see if a global.json file exists in the project dir, if it is read it in
			// so that the version of dotnet can be matched against it as a more specific check
			globalJsonPath := filepath.Join(tmpPath, "global.json")
			specificVersion, err = util.DotnetSDKVersion(globalJsonPath)

			if err != nil {
				logger.Warnf("Issue retrieving global .NET version: %s", err.Error())
			}

			if specificVersion != "" {
				info = "Specific version constraint has been found in project using the 'global.json' file"
			}

		case "java":

			versionCmd = "java"
			versionArgs = "-version"
			re = *regexp.MustCompile(`"(?P<version>.*)"`)

		case "nx":

			versionCmd = "node"
			versionArgs = "--version"
			re = *regexp.MustCompile(`"v(?P<version>.*)"`)

		default:
			versionCmd = ""
		}

		// execute the command if one has been set
		if versionCmd == "" {
			continue
		}

		result, err := config.ExecuteCommand(path, logger, versionCmd, versionArgs, false, true)

		// check for errors
		if err != nil {
			logger.Errorf("Issue running command: %s", err.Error())
			continue
		}

		// get the version from the result so it can be tested using semver
		matches := re.FindStringSubmatch(result)
		idx := re.SubexpIndex(("version"))
		versionFound := matches[idx]

		logger.Debugf("Tool version found: %s", versionFound)

		// get the constraint that should be used to check for
		// if a specific version has been specified modify this constraint
		constraint := cmd.Version
		if specificVersion != "" {
			constraint = fmt.Sprintf("= %s", specificVersion)
		}

		met := s.CompareVersion(constraint, versionFound, logger)

		// if not matched then create a command object and set in the array
		if !met {
			incorrect = append(incorrect, models.Command{
				Binary:          cmd.Name,
				VersionFound:    versionFound,
				VersionRequired: constraint,
			})
		}

	}

	return incorrect, info
}

// CompareVersion compares the specified version against the contsraint
func (s *Settings) CompareVersion(constraint string, version string, logger *logrus.Logger) bool {
	var result bool

	// check that the version string can be turned into a semantic version
	// this is done by removing characters that should not be there
	pattern := "_"
	re := regexp.MustCompile(pattern)
	matched := re.MatchString(version)
	if matched {
		old := version
		version = strings.ReplaceAll(old, "_", "")
		logger.Warnf("Version has been modified so it can be parsed as a semver, from '%s' to '%s'", old, version)
	}

	// create a semver constraint to compare the version number against
	// if the constraint is null set this to a constraint that matches everything
	// greater than 0
	if constraint == "" {
		constraint = ">= 0"
	}
	c, err := semver.NewConstraint(constraint)

	if err != nil {
		logger.Errorf("Unable to parse version constraint '%s': %s", constraint, err.Error())
		return false
	}

	// set the version that has been returned from the command
	v, err := semver.NewVersion(version)
	if err != nil {
		logger.Errorf("Unable to parse found version number '%s': %s", version, err.Error())
		return false
	}

	// check if the version meets the contraint
	result = c.Check(v)

	return result
}
