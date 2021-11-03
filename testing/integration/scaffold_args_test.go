// +build integration

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/amido/stacks-cli/internal/util"
	"github.com/stretchr/testify/suite"
)

// ArgumentTestSuite sets up the suite for checking that the command line args run properly
type ArgumentSuite struct {
	BaseIntegration
}

// SetupSuite sets up the tests by running the command to test the result from
func (suite *ArgumentSuite) SetupSuite() {

	// create the command and arguments required for the tests
	replacements := []interface{}{
		area,
		suite.Company,
		component,
		domain,
		framework,
		suite.Project,
		pipeline,
		platform,
		tf_container,
		tf_group,
		tf_storage,
		framework_option,
		suite.ProjectDir,
	}
	arguments := fmt.Sprintf(`scaffold -A %s --company %s --component %s --domain %s
-F %s -n %s -p %s -P %s --tfcontainer %s --tfgroup %s --tfstorage %s
-O %s --cmdlog -w %s`, replacements...)

	suite.BaseIntegration.RunCommand(suite.BinaryCmd, arguments, false)
}

// TearDownSuite removes all of the files that have been generated in this suite
func (suite *ArgumentSuite) TearDownSuite() {
	err := ClearDir(suite.ProjectDir)
	if err != nil {
		suite.T().Errorf("Error tearing down the Argument Suite: %v", err)
	}
}

// TestProjectDirExists tests that the project directory has been created properly
func (suite *ArgumentSuite) TestProjectDirExists() {
	exists := util.Exists(filepath.Join(suite.ProjectDir, suite.Project))
	suite.Equal(true, exists)
}

// TestCmdlogExists ensures that the command log has been created
func (suite *ArgumentSuite) TestCmdlogExists() {
	exists := util.Exists("cmdlog.txt")
	suite.Equal(true, exists)
}

// TestVariablesFileExists checks that the variable file required for the build pipeline exists
func (suite *ArgumentSuite) TestVariablesFileExists() {
	variablesFile := filepath.Join(suite.ProjectDir, suite.Project, "build", "azDevOps", "azure", "azuredevops-vars.yml")
	exists := util.Exists(variablesFile)
	suite.Equal(true, exists)
}

// TestNamespace checks that the source files for the project have been corrected set with the
// name of the company as supplied on the command line
// It does this by taking the first directory of the src/api folder and checks that it begins
// with the specified company name
func (suite *ArgumentSuite) TestNamespace() {
	var firstDir string

	basedir := filepath.Join(suite.ProjectDir, suite.Project, "src", "api")
	files, _ := os.ReadDir(basedir)

	// iterate around the files and get the first directory
	for _, file := range files {
		filePath := filepath.Join(basedir, file.Name())
		info, err := os.Stat(filePath)
		if err != nil {
			suite.T().Errorf("Problem analysing file: %v", err)
		}

		if info.IsDir() {
			firstDir = file.Name()
			break
		}
	}

	// Check that the dirname begins with %company%
	pattern := fmt.Sprintf("^%s.*$", suite.Company)
	re := regexp.MustCompile(pattern)
	matched := re.MatchString(firstDir)

	suite.Equal(true, matched)
}

// TestNoGitRepo checks that a .git directory has not been created as no remote URL
// has been set on the CLI
func (suite *ArgumentSuite) TestNoGitRepo() {

	gitDir := filepath.Join(suite.ProjectDir, suite.Project, ".git")

	exists := util.Exists(gitDir)

	suite.Equal(false, exists)
}

// TestRunSuite runs the testify test suite
// This is being used so that it possible to have a SetupSuite and a TearDown suite so that
// the necessary commands can be run and then the result checked
func TestArgumentSuite(t *testing.T) {

	s := new(ArgumentSuite)
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
