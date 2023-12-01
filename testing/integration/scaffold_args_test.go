//go:build integration
// +build integration

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/amido/stacks-cli/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ArgsSuite sets up the suite for checking the CLI works as expected when
// running with a command line args
type ArgsSuite struct {
	BaseIntegration

	ProjectPath string
}

// SetupSuite sets up the tests by running the CLI command with the necessary arguments
func (suite *ArgsSuite) SetupSuite() {

	var arguments string

	// create a map of the args and the values associated with the args
	args := map[string]string{
		"-A":            area,
		"--company":     suite.Company,
		"--component":   component,
		"--domain":      domain,
		"-F":            framework,
		"-n":            suite.Project,
		"-p":            pipeline,
		"-P":            platform,
		"--tfcontainer": tf_container,
		"--tfgroup":     tf_group,
		"--tfstorage":   tf_storage,
		"-O":            framework_option,
		"-V":            *framework_version,
		"--cmdlog":      "",
		"-w":            suite.ProjectDir,
		"--save":        "",
		"--nobanner":    "",
	}

	// join the arguments up into a string to pass to the RunCommand
	for key, value := range args {

		// check the value, if it is "" then it is a switch arg
		if value == "" {
			arguments += fmt.Sprintf("%s ", key)
		} else {
			arguments += fmt.Sprintf("%s %s ", key, value)
		}
	}

	// execute the command
	trimmed := fmt.Sprintf("scaffold %s", strings.TrimSpace(arguments))
	suite.BaseIntegration.RunCommand(suite.BinaryCmd, trimmed, false)
}

// TearDownSuite removes all of the files that have been created in this test suite
func (suite *ArgsSuite) TearDownSuite() {
	err := suite.ClearDir(suite.ProjectDir)
	if err != nil {
		fmt.Printf("Error tearing down the Argument Suite: %v", err)
	}
}

// TestArgsSuite runs the testify test suite
// This is being used so that it possible to have a SetupSuite and a TearDown suite so that
// the necessary commands can be run and then the result checked
func TestArgsSuite(t *testing.T) {

	s := new(ArgsSuite)
	s.BinaryCmd = *binaryCmd
	s.Company = *company
	s.Project = *project
	s.ProjectDir = *projectDir

	// define the paths for the suite
	s.ProjectPath = filepath.Join(s.ProjectDir, s.Project)

	s.Assert = assert.New(t)

	s.SetProjectDir()

	suite.Run(t, s)
}

// TestProjectPath checks that the project has been setup correctly.
func (suite *ArgsSuite) TestProject() {

	// check that the project path exists
	suite.T().Run(fmt.Sprintf("%s project directory exists", suite.Project), func(t *testing.T) {

		exists := util.Exists(suite.ProjectPath)

		suite.Assert.Equal(true, exists, "Project should exist: %s", suite.ProjectPath)
	})

	// ensure that the devops variable template exists
	suite.T().Run("Azure DevOps variable template file exist", func(t *testing.T) {
		path := filepath.Join(suite.ProjectPath, "build", "azDevOps", "azure", "air-api-vars.yml")
		exists := util.Exists(path)

		suite.Assert.Equal(true, exists, "Project should exist: %s", "Azure DevOps variable template file should exist: %s", path)
	})

	// check that no git repo has been created as no remote URL has been supplied
	suite.T().Run("Git repo has not been configured", func(t *testing.T) {
		path := filepath.Join(suite.ProjectPath, ".git")
		exists := util.Exists(path)

		suite.Assert.Equal(false, exists, "Directory should not have been configured as a git repository")
	})

	// check that the project files have been namespaced with the company name properly
	suite.T().Run("Ensure project files have been named correctly", func(t *testing.T) {
		var firstDir string
		var list []string

		basedir := filepath.Join(suite.ProjectPath, "src", "api")
		files, _ := os.ReadDir(basedir)

		t.Logf("Reading dir: %s", basedir)
		t.Logf("Found %d files", len(files))

		// iterate around the files and get the first directory
		for _, file := range files {
			filePath := filepath.Join(basedir, file.Name())
			info, err := os.Stat(filePath)
			if err != nil {
				suite.T().Fatalf("Problem analysing file: %v", err)
			}

			t.Logf(file.Name())

			if info.IsDir() {
				list = append(list, file.Name())
			}
		}

		t.Logf("Files: %s", strings.Join(list, ", "))

		// Check that the dirname begins with %company%
		pattern := fmt.Sprintf("^%s.*$", suite.Company)
		re := regexp.MustCompile(pattern)
		matched := re.MatchString(list[0])

		suite.Assert.Equal(true, matched, fmt.Sprintf("Project files should be namespaced with the company nam. [%s !match %s]", list[0], pattern))
	})
}

// TestCmdLogExists checks that the cmdlog text file has not been created as the option do
// do so was not specified when the CLI was run
func (suite *ArgsSuite) TestCmdLogExists() {
	path := filepath.Join(suite.ProjectDir, "cmdlog.txt")
	exists := util.Exists(path)

	suite.Assert.Equal(true, exists, "cmdlog file should exist: %s", path)
}

// TestSavedConfig checks that the configuration has been saved to a file

func (suite *ArgsSuite) TestSavedConfig() {
	path := filepath.Join(suite.ProjectDir, "stacks.yml")
	exists := util.Exists(path)

	suite.Assert.Equal(true, exists, "saved config file should exist: %s", path)
}
