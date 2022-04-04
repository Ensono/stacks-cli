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

// ConfigFileSuite performs a bunch of tests that check the output
// of the CLI and makes sure that all the replacements etc that are as expected
type ConfigFileSuite struct {
	BaseIntegration

	ProjectPath1 string
	ProjectPath2 string
}

// SetupSuite creates a new configuration file in the project directory and then
// runs the scaffold command with that configuration
func (suite *ConfigFileSuite) SetupSuite() {

	configFile := suite.BaseIntegration.WriteConfigFile("")

	// build up the command to run and use the configuration file to do so
	// this is only run once so that the files can be checked as required
	arguments := fmt.Sprintf("scaffold -c %s", configFile)

	suite.BaseIntegration.RunCommand(suite.BinaryCmd, arguments, false)
}

// TearDownSuite removes all of the files that have been generated in this suite
func (suite *ConfigFileSuite) TearDownSuite() {
	return
	err := suite.ClearDir(suite.ProjectDir)
	if err != nil {
		fmt.Printf("Error tearing down the ConfigFileSuite: %v", err)
	}
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

	// define the paths for the suite
	s.ProjectPath1 = filepath.Join(s.ProjectDir, fmt.Sprintf("%s-1", s.Project))
	s.ProjectPath2 = filepath.Join(s.ProjectDir, fmt.Sprintf("%s-2", s.Project))

	s.Assert = assert.New(t)

	s.SetProjectDir()

	suite.Run(t, s)
}

// TestProject1 checks that all files have been copied and created. It also checks that
// the directory has been setup as a git repository
func (suite *ConfigFileSuite) TestProject1() {

	// check that the project path exists
	suite.T().Run(fmt.Sprintf("%s-1 project directory exists", suite.Project), func(t *testing.T) {

		exists := util.Exists(suite.ProjectPath1)

		suite.Assert.Equal(true, exists, "Project directory should exist: %s", suite.ProjectPath1)
	})

	// ensure that the devops variable template exists
	suite.T().Run("Azure DevOps variable template file exist", func(t *testing.T) {
		path := filepath.Join(suite.ProjectPath1, "build", "azDevOps", "azure", "air-api-vars.yml")
		exists := util.Exists(path)

		suite.Assert.Equal(true, exists, "Azure DevOps variable template file should exist: %s", path)
	})

	// ensure that a git repo exists
	suite.T().Run("Git repo has been configured", func(t *testing.T) {
		path := filepath.Join(suite.ProjectPath1, ".git")
		exists := util.Exists(path)

		suite.Assert.Equal(true, exists, "Directory should have been configured as a git repository")
	})

	// check that the remote has been set in the git repo
	suite.T().Run("Check remote repo has been configured", func(t *testing.T) {

		// get the contents of the configuration file
		gitConfig := filepath.Join(suite.ProjectPath1, ".git", "config")
		config, err := ioutil.ReadFile(gitConfig)

		if err != nil {
			t.Errorf("Unable to read git config file: %s", err.Error())
		}

		// define pattern to check that the url has been set correctly
		pattern := fmt.Sprintf(`(?m)url\s+=\s+https://github\.com/dummy/%s-1`, suite.Project)

		t.Logf("Looking for pattern: %s", pattern)

		re := regexp.MustCompile(pattern)
		matched := re.MatchString(string(config))

		suite.Assert.Equal(true, matched, "Git should have been configured with remote repo")
	})

	// check that the project files have been namespaced with the companu name properly
	suite.T().Run("Ensure project files have been named correctly", func(t *testing.T) {
		var firstDir string

		basedir := filepath.Join(suite.ProjectPath1, "src", "api")
		files, _ := os.ReadDir(basedir)

		// iterate around the files and get the first directory
		for _, file := range files {
			filePath := filepath.Join(basedir, file.Name())
			info, err := os.Stat(filePath)
			if err != nil {
				suite.T().Fatalf("Problem analysing file: %v", err)
			}

			suite.T().Log(fmt.Sprintf("%s", file.Name()))

			if info.IsDir() {
				firstDir = file.Name()
				break
			}
		}

		// Check that the dirname begins with %company%
		pattern := fmt.Sprintf("^%s.*$", suite.Company)
		re := regexp.MustCompile(pattern)
		matched := re.MatchString(firstDir)

		suite.Assert.Equal(true, matched, "Project files should be namespaced with the company name '%s': %s", suite.Company, firstDir)
	})
}

// TestProject2 checks that the second project has been configured, but it has not been configured
// as a git repository
func (suite *ConfigFileSuite) TestProject2() {
	// check that the project path exists
	suite.T().Run(fmt.Sprintf("%s-2 project directory exists", suite.Project), func(t *testing.T) {

		exists := util.Exists(suite.ProjectPath2)

		suite.Assert.Equal(true, exists, "Project should exist: %s", suite.ProjectPath1)
	})

	// ensure that the devops variable template exists
	suite.T().Run("Azure DevOps variable template file exist", func(t *testing.T) {
		path := filepath.Join(suite.ProjectPath2, "build", "azDevOps", "azure", "air-api-vars.yml")
		exists := util.Exists(path)

		suite.Assert.Equal(true, exists, "Azure DevOps variable template file should exist: %s", path)
	})

	suite.T().Run("Git repo has not been configured", func(t *testing.T) {
		path := filepath.Join(suite.ProjectPath2, fmt.Sprintf("%s-2", suite.Project), ".git")
		exists := util.Exists(path)

		suite.Assert.Equal(false, exists, "Project should not have been configured as a Git repository")
	})
}

// TestCmdLogDoesNotExist checks that the cmdlog text file has not been created as the option do
// do so was not specified when the CLI was run
func (suite *ConfigFileSuite) TestCmdLogDoesNotExist() {
	path := filepath.Join(suite.ProjectDir, "cmdlog")
	exists := util.Exists(path)

	suite.Assert.Equal(false, exists, "Cmdlog file should not exist: %s", path)
}
