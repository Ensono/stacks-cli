// +build integration

package integration

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ActiveState/termtest"
	"github.com/amido/stacks-cli/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// InterativeSuite contains the propertries for running the Stacks CLI in an interactive way
type InteractiveSuite struct {
	BaseIntegration
}

// TearDownSuite removes all of the files that have been generated in this suite
func (suite *InteractiveSuite) TearDownSuite() {
	err := ClearDir(suite.ProjectDir)
	if err != nil {
		suite.T().Errorf("Error tearing down the Argument Suite: %v", err)
	}
}

// SetupSuite sets up the test by running the interactive command of the Stacks CLI
// It will set the output to the cmdOutput
// The Go-Expect library is used to test the interactiveness
func (suite *InteractiveSuite) TestInteractiveMode() {

	// create the command and arguments required for the tests
	command := suite.BinaryCmd
	arguments := "interactive --nobanner"

	// use the util function to split the arguments up and run the command
	cmd, args := util.BuildCommand(command, arguments)

	// configure command and the arguments that Termtest needs to run
	opts := termtest.Options{
		CmdName: cmd,
		Args:    args,
	}

	cp, err := termtest.NewTest(suite.T(), opts)
	require.NoError(suite.T(), err, "create console process")
	defer cp.Close()

	cp.Expect("What is the name of your company?")
	cp.SendLine(suite.Company)
	cp.Expect("What is the scope or area of the company?")
	cp.SendLine(area)
	cp.Expect("What component is being worked on?")
	cp.SendLine(component)
	cp.Expect("What pipeline is being targeted?")
	cp.SendLine(pipeline)
	cp.Expect("Which cloud is Stacks being setup in?")
	cp.SendLine(cloud)
	cp.Expect("Which group is the Terraform state being saved in?")
	cp.SendLine(tf_group)
	cp.Expect("What is the name of the Terraform storage?")
	cp.SendLine(tf_storage)
	cp.Expect("What is the name of the Terraform storage container?")
	cp.SendLine(tf_container)
	cp.Expect("What is the external domain of the solution?")
	cp.SendLine(domain)
	cp.Expect("What is the internal domain of the solution?")
	cp.SendLine("")
	cp.Expect("What options would you like to enable, if any?")
	cp.SendLine("")
	cp.Expect("Please specify the working directory for the projects?")
	cp.SendLine(suite.ProjectDir)
	cp.Expect("How many projects would you like to configure?")
	cp.SendLine("1")

	cp.Expect("What is the project name?")
	cp.SendLine(suite.Project)
	cp.Expect("What framework should be used for the project?")
	cp.SendLine(framework)
	cp.Expect("Which style of the framework do you require?")
	cp.SendLine(framework_option)
	cp.Expect("Which version of the framework option do you require?")
	cp.SendLine("")
	cp.Expect("Specify any additional framework properties. (Use a comma to separate each one).")
	cp.SendLine("")
	cp.Expect("What platform is being used")
	cp.SendLine(platform)
	cp.Expect("Please select the source control system being used")
	cp.SendLine("")
	cp.Expect("What is the URL of the remote repository?")
	cp.SendLine(fmt.Sprintf("https://github.com/dummy/%s", suite.Project))
	cp.Expect("Which cloud region should be used?")
	cp.SendLine(cloud_region)
	cp.Expect("What is the name of the group for all the resources?")
	cp.SendLine(cloud_group)

	cp.ExpectExitCode(0)

	// run a test to ensure that the configFile has been created
	suite.T().Run("configuration file should exist", func(t *testing.T) {
		// check that the stacks configuration file exists
		configFile := filepath.Join(suite.ProjectDir, "stacks.yml")

		exists := util.Exists(configFile)

		assert.Equal(t, true, exists)
	})

}

// expect the prompts from the command line and then send the appropriate response
func TestInteractiveSuite(t *testing.T) {

	s := new(InteractiveSuite)
	s.BinaryCmd = *binaryCmd
	s.Company = *company
	s.Project = *project
	s.ProjectDir = *projectDir

	s.SetProjectDir()

	// only run the interactive test when not on windows
	// this is because the questions are repeated in the Windows console which causes
	// an issue with the expect
	if runtime.GOOS != "windows" {
		suite.Run(t, s)
	}
}
