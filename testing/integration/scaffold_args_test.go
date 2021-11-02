// +build integration

package integration

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/amido/stacks-cli/internal/util"
	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var company = flag.String("company", "MyCompany", "Name of the company")
var project = flag.String("project", "my-webapi", "Name of the project")
var projectDir = flag.String("projectdir", ".", "Project Directory")
var binaryCmd = flag.String("binarycmd", "stacks-cli", "Name and path of the binary to use to run the tests")

// ArgumentTestSuite sets up the suite for checking that the command line args run properly
type ArgumentSuite struct {
	suite.Suite

	// set the name of the project to create
	Project    string
	ProjectDir string

	// Set the name of the company for which the project is being setup for
	Company string

	// The name of the command to run
	BinaryCmd string

	// Cmdoutput to be used for analysis
	CmdOutput string
}

// SetupSuite sets up the tests by running the command to test the result from
func (suite *ArgumentSuite) SetupSuite() {

	// create the command and arguments required for the tests
	command := suite.BinaryCmd
	arguments := fmt.Sprintf(`scaffold -A core --company %s --component backend --domain stacks-example.com
-F dotnet -n %s -p azdo -P aks --tfcontainer mywebapi --tfgroup supporting-group --tfstorage kjh56sdfnjnkjn
-O webapi --cmdlog -w %s`, suite.Company, suite.Project, suite.ProjectDir)

	// use the util function to split the arguments up and run the command
	cmd, args := util.BuildCommand(command, arguments)

	// configure the exec command to execute the command
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		suite.T().Errorf("Error running command: %s", err.Error())
	}
	suite.CmdOutput = string(out)

	// cmdLine.Stdout = os.Stdout
	// cmdLine.Stderr = os.Stderr

	// execute the command
	//if err := cmdLine.Run(); err != nil {
	//	suite.T().Errorf("Error running command: %s", err.Error())
	//}
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
