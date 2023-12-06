//go:build integration
// +build integration

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// CLIAlertSuite creates a suite of tests that check the output of the CLI
// when things are not configured properly
type CLIAlertSuite struct {
	BaseIntegration

	ConfigFile                     string
	BadConfigFile                  string
	WrongFrameworkOptionConfigFile string

	ProjectPath1 string
	ProjectPath2 string
	ProjectPath3 string
	Project3File string

	FrameworkOption string
}

// SetupSuite creates the necessary environment for the tests to run
// In this case a configuration file is used so that scaffold command can be run
// several times
func (suite *CLIAlertSuite) SetupSuite() {

	suite.ConfigFile = suite.BaseIntegration.WriteConfigFile("")

	// create two directories that are to be encountered when running the scaffold
	suite.CreateDirs([]string{
		suite.ProjectPath1,
		suite.ProjectPath3,
	})

	// create a file in one of the dirs so that it is not empty
	file, err := os.Create(suite.Project3File)
	if err != nil {
		suite.T().Fatalf("unable to create project file '%s': %s", suite.Project3File, err.Error())
	}
	file.Close()

	// create the badConfigFile to check that the CLI throws the correct error
	malformed := fmt.Sprintf(`directory:\n\tworkingDir:%s`, suite.ProjectDir)
	suite.BadConfigFile = filepath.Join(suite.ProjectDir, "malformed-stacks.yml")
	err = os.WriteFile(suite.BadConfigFile, []byte(malformed), os.ModePerm)
	if err != nil {
		suite.T().Fatalf("unable to create malformed configuration file: %s", err.Error())
	}

	// update the configuration so that it contains an invalid framework option
	// then write this file out so that the CLI can read it in
	suite.FrameworkOption = "bus"
	old := framework_option
	framework_option = suite.FrameworkOption
	reset := func() { framework_option = old }
	defer reset()

	// write out configuration file
	suite.WrongFrameworkOptionConfigFile = suite.BaseIntegration.WriteConfigFile("wrong-stacks.yml")

}

// TearDownSuite removes all of the files that have been generated
// in this suite
func (suite *CLIAlertSuite) TearDownSuite() {
	err := suite.ClearDir(suite.ProjectDir)
	if err != nil {
		fmt.Printf("error tearing down the CLIAlert suite: %s", err.Error())
	}
}

// TestCLIAlertSuite runs the the suite of tests to check that CLI alerts that can
// be raised during operation
// This uses the same configuration file as the configfile tests, but is only
// interested in the alerts that are genereted by the CLI
func TestCLIAlertSuite(t *testing.T) {

	s := new(CLIAlertSuite)
	s.BinaryCmd = *binaryCmd
	s.Company = *company
	s.Project = *project
	s.ProjectDir = *projectDir

	// define the paths for the suite
	s.ProjectPath1 = filepath.Join(s.ProjectDir, fmt.Sprintf("%s-1", s.Project))
	s.ProjectPath2 = filepath.Join(s.ProjectDir, fmt.Sprintf("%s-2", s.Project))
	s.ProjectPath3 = filepath.Join(s.ProjectDir, fmt.Sprintf("%s-3", s.Project))
	s.Project3File = filepath.Join(s.ProjectPath3, "project.json")

	s.BaseIntegration.Assert = assert.New(t)

	s.SetProjectDir()

	suite.Run(t, s)
}

// TestConfigFileExists checks that the configuration file has been written out
// properly. If not then none of the tests in this suite will work properly
func (suite *CLIAlertSuite) TestConfigFileExists() {

	exists := util.Exists(suite.ConfigFile)

	suite.Assert.Equal(true, exists, "Configuration file does not exist")
}

// TestProjectDir already exists checks that the CLI has detected a project directory
// already exists. It will overwrite the project directory if it is empty, but it will
// not do so if files or directories exist below it
func (suite *CLIAlertSuite) TestProjectDirAlreadyExists() {

	// run the scaffold command
	arguments := fmt.Sprintf("scaffold -c %s --nobanner", suite.ConfigFile)
	suite.BaseIntegration.RunCommand(suite.BinaryCmd, arguments, false)

	suite.T().Run("CLI overwrites existing empty directory", func(t *testing.T) {

		escapedProjectPath := strings.Replace(suite.ProjectPath1, "\\", "\\\\\\\\", -1)
		pattern := fmt.Sprintf(`(?i)overwriting empty directory: %s`, escapedProjectPath)

		matched := suite.CheckCmdOutput(pattern)

		suite.Assert.Equal(true, matched, "CLI should overwrite empty directory")
	})

	suite.T().Run("Second project is created", func(t *testing.T) {

		exists := util.Exists(suite.ProjectPath2)

		suite.Assert.Equal(true, exists, "CLI should create the directory for the second project")
	})

	suite.T().Run("CLI does not overwrite and existing project directory", func(t *testing.T) {

		escapedProjectPath := strings.Replace(suite.ProjectPath3, "\\", "\\\\\\\\", -1)
		pattern := fmt.Sprintf(`(?i)project directory already exists, skipping: %s`, escapedProjectPath)

		matched := suite.CheckCmdOutput(pattern)

		suite.Assert.Equal(true, matched, "CLI should not overwrite an existing directory, with data in it")
	})
}

// TestFrameworkAppsNotFound changes the PATH environment variable for the machine so that
// framework commands cannot be found.
// The test then checks that the expected messages are in the output
// The PATH is restored on the machine after being run
// The stacks-cli is passed as an absolute to the integration test so that does not need to be
// found in the path
func (suite *CLIAlertSuite) TestFrameworkAppsNotFound() {

	var err error

	// declare path variable that will be declared
	var path string

	// get the current envPath from the machine
	var envPath string = os.Getenv("PATH")

	// set the path according to the OS
	switch util.GetPlatformOS() {
	case "windows":
		path = "C:/Windows/System32"
	default:
		path = "/usr/sbin"
	}

	// set the OS path to the path that has been defined
	err = os.Setenv("PATH", path)
	if err != nil {
		suite.Assert.Fail("Unable to set temporary PATH environment variable: %s", err.Error())
	}

	// Reset thePATH env var at the end of the test
	reset := func() {
		err = os.Setenv("PATH", envPath)
		if err != nil {
			suite.Assert.Fail("Unable to reset PATH environment variable: %s", err.Error())
		}
	}
	defer reset()

	// run the scaffold command
	// the exit code is ignored here so that the output of the command can be seen
	// otherwise the tests just stop
	arguments := fmt.Sprintf("scaffold -c %s --nobanner", suite.ConfigFile)
	suite.BaseIntegration.RunCommand(suite.BinaryCmd, arguments, true)

	// create the test tables to use to perform the necessary tests
	tables := []struct {
		binary  string
		pattern string
		msg     string
	}{
		{
			"dotnet",
			"(?i)command 'dotnet' for the 'dotnet' framework cannot be located",
			"CLI should state that the dotnet command cannot be found",
		},
		{
			"git",
			"(?i)command 'git' for the 'dotnet' framework cannot be located",
			"CLI should state that the git command cannot be found",
		},
	}

	// iterate around the test tables
	for _, table := range tables {

		suite.T().Run(fmt.Sprintf("`%s` command cannot be located", table.binary), func(t *testing.T) {
			matched := suite.CheckCmdOutput(table.pattern)

			suite.Assert.Equal(true, matched, table.msg)
		})
	}
}

// TestIncorrectFrameworkOption tests that the CLI copes property if someone specifies
// an invalid framework option, e.g. one that the CLI does not know about
func (suite *CLIAlertSuite) TestIncorrectFrameworkOption() {

	// run the scaffold command against the incorrect framrowkr version configuration file
	arguments := fmt.Sprintf("scaffold -c %s --nobanner", suite.WrongFrameworkOptionConfigFile)
	suite.BaseIntegration.RunCommand(suite.BinaryCmd, arguments, false)

	suite.T().Run("Ensure CLI errors gracefully", func(t *testing.T) {

		pattern := fmt.Sprintf("(?i)the url for the specified framework option, %s, is empty", suite.FrameworkOption)
		matched := suite.CheckCmdOutput(pattern)

		suite.Assert.Equal(true, matched, "CLI should error because the framework option is invalid")
	})
}
