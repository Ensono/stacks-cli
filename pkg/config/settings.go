package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Ensono/stacks-cli/internal/models"
	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/Masterminds/semver"
	"github.com/sirupsen/logrus"
)

// Settings holds the settings for each project as read from
// the `stackscli.yml` file in the project
type Settings struct {
	Framework SettingsFramework `mapstructure:"framework"`
	Pipeline  []Pipeline        `mapstructure:"pipeline"`
	Init      Init              `mapstructure:"init"`
	Setup     Setup             `mapstructure:"setup"`
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
	Action          string   `mapstructure:"action"`
	Command         string   `mapstructure:"cmd"`
	Arguments       string   `mapstructure:"args"`
	Description     string   `mapstructure:"desc"`
	ApplyProperties bool     `mapstructure:"applyProperties"`
	Tags            []string `mapstructure:"tags"`
}

type SettingsFramework struct {
	Name     string                      `mapstructure:"name"`
	Commands []SettingsFrameworkCommands `mapstructure:"commands"`
}

type SettingsFrameworkCommands struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

// GetPipelines attempts to return all pipelines settings for the named pipeline
func (s *Settings) GetPipelines(name string) []Pipeline {
	pipeline := []Pipeline{}

	// iterate around the pipeline slice and find the ones with the type that matches
	// the specified name.
	for _, p := range s.Pipeline {
		if p.Type == name {
			pipeline = append(pipeline, p)
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

// CheckCmdVersions checks that all of the commands that have been specified for the
// component exist and that they are the correct version.
//
// It takes the following parameters:
// - config: A pointer to the Config struct containing the framework commands.
// - logger: A logger instance from the logrus package for logging errors and information.
// - path: The file path where the commands should be executed.
// - tmpPath: A temporary file path used for specific command checks.
//
// It returns a slice of models.Command containing the commands that do not meet the specified version constraints,
// and an info string providing additional information.
func (s *Settings) CheckCmdVersions(config *Config, logger *logrus.Logger, path string, tmpPath string) ([]models.Command, string) {

	var constraint string
	var incorrect []models.Command
	var info string
	var met bool
	var resultErrors []string
	var versionFound string

	// iterate around the framework commands and if a version constraint has been specified, run
	// the command and check against the pattern
	for _, cmd := range s.Framework.Commands {

		// Get the commands for the framework
		fCommands := config.GetFrameworkCommands(cmd.Name)

		// perform a check on the version of the command, if a version pattern has been specified
		// the version check is slightly different for .NET as it can be more specific
		for _, fCmd := range fCommands.Commands {

			// if the current cmd does not have a version specified then skip
			if fCmd.Version == (FrameworkDefVersion{}) {
				continue
			}

			// Run the command to get the the version of the command
			result, err := config.ExecuteCommand(path, logger, fCmd.Name, fCmd.Version.Arguments, false, true)

			// raise an error if the command failed for any reason
			if err != nil {
				logger.Errorf("Issue running command: %s", err.Error())
				continue
			}

			// if there are no results, add to the error slice so that all the errors can be
			// displayed at the end
			if result == "" {
				resultErrors = append(resultErrors, fmt.Sprintf("No versions for '%s' SDKs found", fCmd.Name))
				continue
			}

			// create an object the version struct
			version := models.Version{}
			version.Init(result, fCmd.Version.Pattern)

			// use the pattern to check for the version of the command
			// however as .NET has a different ruleset for this,
			switch cmd.Name {
			case "dotnet":

				// put together a path for the global.json file
				globalData := filepath.Join(tmpPath, "global.json")

				// determine if the file exists, if not see if there is a constraint set in the
				// framework project settings
				if !util.Exists(globalData) {
					if cmd.Version != "" {
						globalData = cmd.Version
					}
				}

				// set the global value
				version.DotNetGlobal(globalData)

				met, err = version.DotNet()

				// version_segments, err := util.VersionSegments(re, matches)
				if err != nil {
					logger.Errorf("version check error: %s", err.Error())
				}

				if err != nil {
					resultErrors = append(resultErrors, err.Error())
				}

			default:

				// for all other commands this is a semantic version check
				// so use the "version" named group to compare against

				// get the comparator to use and override the default if
				// one has been provided

				// If a version constraint has been specified, update the version object
				// to use it
				if cmd.Version != "" {
					version.SetSemverConstraint(cmd.Version)
				}

				logger.Debugf("Version of '%s' found: %s", fCmd.Name, versionFound)

				met, err = version.Semver()
				if err != nil {
					resultErrors = append(resultErrors, err.Error())
				}

			}

			// if not matched then create a command object and set in the array
			if !met {
				incorrect = append(incorrect, models.Command{
					Binary:          cmd.Name,
					VersionFound:    versionFound,
					VersionRequired: constraint,
				})
			}
		}
	}

	// if there any errors, log them and then exit
	if len(resultErrors) > 0 {
		for _, err := range resultErrors {
			logger.Errorf("error: %s", err)
		}
		os.Exit(7)
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

	// check if the version meets the constraint
	result = c.Check(v)

	return result
}
