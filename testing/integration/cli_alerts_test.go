// +build integration

package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	"github.com/amido/stacks-cli/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// CLIAlertSuite creates a suite of tests that check the output from the CLI
// when things are not configured properly
type CLIAlertSuite struct {
	BaseIntegration

	ConfigFile string
}

// SetupSuite creates the environment for the tests to run
// In this case a configuration file is used so that the scaffold command can be run multiple times
func (suite *CLIAlertSuite) SetupSuite() {
	suite.ConfigFile = suite.BaseIntegration.WriteConfigFile()
}

// TestProjectAlreadyExists checks to see if the CLI has detected that the project
// already exists. This is done my checking that the output of the command has
// the alert
func (suite *CLIAlertSuite) TestProjectAlreadyExists() {

	// create the project directory before running the cli
	testProjectPath := filepath.Join(suite.ProjectDir, fmt.Sprintf("%s-1", suite.Project))
	err := util.CreateIfNotExists(testProjectPath, 066)
	if err != nil {
		suite.T().Errorf("Unable to create project dir: %s", err.Error())
	}

	// run the command and then check the output
	arguments := fmt.Sprintf("scaffold -c %s --nobanner", suite.ConfigFile)
	suite.BaseIntegration.RunCommand(suite.BinaryCmd, arguments, false)

	suite.T().Run("CLI does not overwrite an existing project", func(t *testing.T) {

		// create the pattern to match the output with
		pattern := fmt.Sprintf("project directory already exists, skipping: %s", testProjectPath)
		t.Logf("Looking for pattern: '%s'", pattern)

		re := regexp.MustCompile(pattern)
		matched := re.MatchString(suite.CmdOutput)

		assert.Equal(t, true, matched)
	})

	suite.T().Run("Second project is created", func(t *testing.T) {
		dir := filepath.Join(suite.ProjectDir, fmt.Sprintf("%s-2", suite.Project))
		exists := util.Exists(dir)
		assert.Equal(t, true, exists)
	})
}

// TestAppsNotFoundInPathEnvVar changes the PATH environment variable so that dotnet and git cannot
// be found by the CLI. This test checks that this is reported properly in the output
func (suite *CLIAlertSuite) TestAppsNotFoundInPathEnvVar() {

	var path string

	// get the PATH env var so that it can be restored
	envPath := os.Getenv("PATH")

	// set the path according to the os
	if runtime.GOOS == "windows" {
		path = "C:/Windows/System32"
	} else {
		path = "/usr/sbin"
	}

	err := os.Setenv("PATH", path)
	if err != nil {
		suite.T().Errorf("Unable to change PATH environment variable: %s", err.Error())
	}

	// run the command and then check the output
	arguments := fmt.Sprintf("scaffold -c %s --nobanner", suite.ConfigFile)
	suite.BaseIntegration.RunCommand(suite.BinaryCmd, arguments, true)

	suite.T().Run("`dotnet` binary cannot be located", func(t *testing.T) {

		// create the pattern to match the output
		pattern := "Command 'dotnet' for the 'dotnet' framework cannot be located."
		t.Logf("Looking for pattern: %s", pattern)

		re := regexp.MustCompile(pattern)
		matched := re.MatchString(suite.CmdOutput)

		assert.Equal(t, true, matched)
	})

	suite.T().Run("`git` binary cannot be located", func(t *testing.T) {

		// create the pattern to match the output
		pattern := "Command 'git' for the 'dotnet' framework cannot be located."
		t.Logf("Looking for pattern: %s", pattern)

		re := regexp.MustCompile(pattern)
		matched := re.MatchString(suite.CmdOutput)

		assert.Equal(t, true, matched)
	})

	// Reset the path variable
	err = os.Setenv("PATH", envPath)
	if err != nil {
		suite.T().Errorf("Unable to revert PATH environment variable: %s", err.Error())
	}
}

// TestBadConfigFile writes out a malformed YAML file and checks that the application
// errors out properly
func (suite *CLIAlertSuite) TestBadConfigFile() {

	// write out bad configuration file
	badConfig := fmt.Sprintf(`directory:\n\tworking:%s`, suite.ProjectDir)
	badConfigFile := filepath.Join(suite.ProjectDir, "bad-stacks.yml")
	err := ioutil.WriteFile(badConfigFile, []byte(badConfig), 0666)

	if err != nil {
		suite.T().Fatalf("Error writing out malformed configuration file: %s", err.Error())
	}

	arguments := fmt.Sprintf("scaffold -c %s --nobanner", badConfigFile)
	suite.BaseIntegration.RunCommand(suite.BinaryCmd, arguments, true)

	suite.T().Run("CLI states config file is malformed", func(t *testing.T) {

		// create the pattern to match the output
		pattern := "Unable to read in configuration file"
		t.Logf("Looking for pattern: %s", pattern)

		re := regexp.MustCompile(pattern)
		matched := re.MatchString(suite.CmdOutput)

		assert.Equal(t, true, matched)
	})

}

// TestIncorrectFrameworkOption tests that the CLI copes properly if someone
// specifies the wrong framework option
func (suite *CLIAlertSuite) TestIncorrectFrameworkOption() {

	// Set the framework option to use, this will be incorrect
	oldFrameworkOption := framework_option
	framework_option = "bus"

	// write out a configuration file
	configFile := suite.BaseIntegration.WriteConfigFile()

	// build up the command to run and use the configuration file to do so
	arguments := fmt.Sprintf("scaffold -c %s --nobanner", configFile)
	suite.BaseIntegration.RunCommand(suite.BinaryCmd, arguments, false)

	suite.T().Run("Ensure CLI errors gracefully", func(t *testing.T) {

		pattern := fmt.Sprintf("The URL for the specified framework option, %s, is empty", framework_option)
		t.Logf("Looking for pattern: %s", pattern)

		re := regexp.MustCompile(pattern)
		matched := re.MatchString(suite.CmdOutput)

		assert.Equal(t, true, matched)
	})

	// reset the framework option
	framework_option = oldFrameworkOption
}

// TestCLIAlertSuite runs the suite of tests to check that the CLI responds in the
// correct way when things are not quite right
func TestCLIAlertSuite(t *testing.T) {

	s := new(CLIAlertSuite)
	s.BinaryCmd = *binaryCmd
	s.Company = *company
	s.Project = *project
	s.ProjectDir = *projectDir

	// if the projectDir is . then set to the current dir
	if s.ProjectDir == "." {
		s.ProjectDir, _ = os.Getwd()
	}

	suite.Run(t, s)
}
