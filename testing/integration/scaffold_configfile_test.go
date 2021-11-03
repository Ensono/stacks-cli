// +build integration

package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/amido/stacks-cli/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ConfigFileSuite contains the properties for running the Stacks CLI with a
// configuration file
type ConfigFileSuite struct {
	BaseIntegration
}

// SetupSuite creates a new configuration file in the project directory and then
// runs the scaffold command with that configuration
func (suite *ConfigFileSuite) SetupSuite() {

	configFile := suite.BaseIntegration.WriteConfigFile()

	// build up the command to run and use the configuration file to do so
	arguments := fmt.Sprintf("scaffold -c %s", configFile)

	suite.BaseIntegration.RunCommand(suite.BinaryCmd, arguments, false)

}

// TearDownSuite removes all of the files that have been generated in this suite
func (suite *ConfigFileSuite) TearDownSuite() {
	err := ClearDir(suite.ProjectDir)
	if err != nil {
		suite.T().Errorf("Error tearing down the Argument Suite: %v", err)
	}
}

func (suite *ConfigFileSuite) TestProject1() {

	// check that the project path exists
	suite.T().Run(fmt.Sprintf("%s-1 project directory exists", suite.Project), func(t *testing.T) {
		path := filepath.Join(suite.ProjectDir, fmt.Sprintf("%s-1", suite.Project))
		exists := util.Exists(path)
		assert.Equal(t, true, exists)
	})

	// ensure that the devops variable template exists
	suite.T().Run("Azure DevOps variable template file exist", func(t *testing.T) {
		path := filepath.Join(suite.ProjectDir, fmt.Sprintf("%s-1", suite.Project), "build", "azDevOps", "azure", "azuredevops-vars.yml")
		exists := util.Exists(path)
		assert.Equal(t, true, exists)
	})

	// ensure that a git repo exists
	suite.T().Run("Git repo has been configured", func(t *testing.T) {
		path := filepath.Join(suite.ProjectDir, fmt.Sprintf("%s-1", suite.Project), ".git")
		exists := util.Exists(path)
		assert.Equal(t, true, exists)
	})

	// check that the remote has been set in the git repo
	suite.T().Run("Check remote repo has been configured", func(t *testing.T) {

		// get the contents of the configuration file
		gitConfig := filepath.Join(suite.ProjectDir, fmt.Sprintf("%s-1", suite.Project), ".git", "config")
		config, err := ioutil.ReadFile(gitConfig)

		if err != nil {
			t.Errorf("Unable to read git config file: %s", err.Error())
		}

		// define pattern to check that the url has been set correctly
		pattern := fmt.Sprintf(`(?m)url\s+=\s+https://github\.com/dummy/%s-1`, suite.Project)

		t.Logf("Looking for pattern: %s", pattern)

		re := regexp.MustCompile(pattern)
		matched := re.MatchString(string(config))

		assert.Equal(t, true, matched)
	})

}

func (suite *ConfigFileSuite) TestProject2() {
	// check that the project path exists
	suite.T().Run(fmt.Sprintf("%s-2 project directory exists", suite.Project), func(t *testing.T) {
		path := filepath.Join(suite.ProjectDir, fmt.Sprintf("%s-2", suite.Project))
		exists := util.Exists(path)
		assert.Equal(t, true, exists)
	})

	// ensure that the devops variable template exists
	suite.T().Run("Azure DevOps variable template file exist", func(t *testing.T) {
		path := filepath.Join(suite.ProjectDir, fmt.Sprintf("%s-2", suite.Project), "build", "azDevOps", "azure", "azuredevops-vars.yml")
		exists := util.Exists(path)
		assert.Equal(t, true, exists)
	})

	suite.T().Run("Git repo has not been configured", func(t *testing.T) {
		path := filepath.Join(suite.ProjectDir, fmt.Sprintf("%s-2", suite.Project), ".git")
		exists := util.Exists(path)
		assert.Equal(t, false, exists)
	})
}

func (suite *ConfigFileSuite) TestCmdLogDoesNotExist() {
	path := filepath.Join(suite.ProjectDir, "cmdlog")
	exists := util.Exists(path)
	assert.Equal(suite.T(), false, exists)
}

// TestRunSuite runs the test suite for running stacks cli with a configuration file
// This is being used so that it possible to have a SetupSuite and a TearDown suite so that
// the necessary commands can be run and then the result checked
func TestConfigFileSuite(t *testing.T) {

	s := new(ConfigFileSuite)
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
