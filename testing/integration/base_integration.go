//go:build integration
// +build integration

package integration

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/Ensono/stacks-cli/pkg/config"
	yaml "github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var version string

// define the values for running the commands in the integration tests
// - the following properties can be set on the command line when running he
var company = flag.String("company", "MyCompany", "Name of the company")
var project = flag.String("project", "my-webapi", "Name of the project")
var framework_version = flag.String("version", "latest", "Version, branch or tag of the framework to use")
var projectDir = flag.String("projectdir", ".", "Project Directory")
var binaryCmd = flag.String("binarycmd", "stacks-cli", "Name and path of the binary to use to run the tests")

// - the propertues set here are used as is
var area = "core"
var component = "backend"
var cloud = "azure"
var cloud_region = "ukwest"
var cloud_group = "mywebapi-resources"
var domain = "stacks-example.com"
var framework = "dotnet"
var framework_option = "webapi"
var pipeline = "azdo"
var platform = "aks"
var tf_container = "mywebapi"
var tf_group = "supporting_group"
var tf_storage = "kjh56sdfnjnkjn"
var environment = "dev"
var environmentType = "Development"

// BaseIntegration declares the base struct that all integration tests will use
type BaseIntegration struct {
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

	ConfigFilename string

	Assert *assert.Assertions
}

// ClearDir clears all of the files and folders within the specified
// directory. This is primarily used by the TearDown function of the test suites
// This is so that the parent directory does not get removed
func (suite *BaseIntegration) ClearDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))

	if err != nil {
		return err
	}

	// iterate around the files that have been found an remove them
	for _, file := range files {
		err := os.RemoveAll((file))
		if err != nil {
			return err
		}
	}

	return err
}

func (suite *BaseIntegration) WriteConfigFile(filename string) string {

	if suite.ConfigFilename == "" {
		suite.ConfigFilename = "stacks.yml"
	}

	if filename != "" {
		suite.ConfigFilename = filename
	}

	cfg := config.Config{}
	cfg.Internal.AddFiles()

	// create the configuration
	input := config.Config{
		Input: config.InputConfig{
			Directory: config.Directory{
				WorkingDir: suite.ProjectDir,
			},
			Business: config.Business{
				Company:   suite.Company,
				Domain:    area,
				Component: component,
			},
			Cloud: config.Cloud{
				Platform: platform,
			},
			Network: config.Network{
				Base: config.NetworkBase{
					Domain: config.DomainType{
						External: domain,
					},
				},
			},
			Pipeline: pipeline,
			Project: []config.Project{
				{
					Name: fmt.Sprintf("%s-1", suite.Project),
					Framework: config.Framework{
						Type:    framework,
						Option:  framework_option,
						Version: *framework_version,
					},
					Platform: config.Platform{
						Type: platform,
					},
					SourceControl: config.SourceControl{
						Type: "github",
						URL:  fmt.Sprintf("https://github.com/dummy/%s-1", suite.Project),
					},
					Cloud: config.Cloud{
						Region:        cloud_region,
						ResourceGroup: cloud_group,
					},
				},
				{
					Name: fmt.Sprintf("%s-2", suite.Project),
					Framework: config.Framework{
						Type:    framework,
						Option:  framework_option,
						Version: *framework_version,
					},
					Platform: config.Platform{
						Type: platform,
					},
					SourceControl: config.SourceControl{
						Type: "github",
						URL:  "",
					},
					Cloud: config.Cloud{
						Region:        cloud_region,
						ResourceGroup: cloud_group,
					},
				},
				{
					Name: fmt.Sprintf("%s-3", suite.Project),
					Framework: config.Framework{
						Type:    framework,
						Option:  framework_option,
						Version: *framework_version,
					},
					Platform: config.Platform{
						Type: platform,
					},
					SourceControl: config.SourceControl{
						Type: "github",
						URL:  "",
					},
					Cloud: config.Cloud{
						Region:        cloud_region,
						ResourceGroup: cloud_group,
					},
				},
			},
			Terraform: config.Terraform{
				Backend: config.TerraformBackend{
					Storage:   tf_storage,
					Container: tf_container,
					Group:     tf_group,
				},
			},
		},
	}

	// serialize the input object and save to a file in the project directory called "stacks.yml"
	data, err := yaml.Marshal(&input)

	if err != nil {
		suite.T().Fatalf("Error serializing configuration: %s", err.Error())
	}

	configFile := filepath.Join(suite.ProjectDir, suite.ConfigFilename)

	err = os.WriteFile(configFile, data, 0666)

	if err != nil {
		suite.T().Fatalf("Error writing out configuration file: %s", err.Error())
	}

	return configFile
}

// RunCommand provides a way for all the Integration tests to run the CLI scaffold command
// in the same way
// The command and arguments are passed as strings, and the func will split up the arguments
// and then run accordingly. The output of the command is set on the suite struct so that
// the tests can analyse it
// A third option can be provided which is the ignore parameter. If set to true then the
// function will not err on a non 0 exit code. This is so that the output of the command
// can be check to make sure that the user has been informed as to why things have not worked
func (suite *BaseIntegration) RunCommand(command string, arguments string, ignore bool) {

	// use the util function to split the arguments
	cmd, args := util.BuildCommand(command, arguments)

	suite.T().Logf("Command: %s %s", cmd, args)

	// write out the command thst ius being run
	cmdlogFile := filepath.Join(suite.ProjectDir, "cmdlog.txt")

	err := os.WriteFile(cmdlogFile, []byte(fmt.Sprintf("%s %s", command, arguments)), 0666)

	if err != nil {
		suite.T().Fatalf("Error writing command to log file: %s", err.Error())
	}

	// configure the exec command to execute the command
	out, err := exec.Command(cmd, args...).Output()
	suite.T().Logf("Command Output:\n%s", string(out))
	if err != nil && !ignore {
		suite.T().Errorf("Error running command: %v", err)
	}
	suite.CmdOutput = string(out)
}

// SetProjectDir sets the path to the project directory
// If it has been set as "." then it will use the current directory
// If it is a relative path then it will prepend the current directory to it
func (suite *BaseIntegration) SetProjectDir() {

	// get the current directory
	cwd, _ := os.Getwd()

	// if the project dir is just a "." then set to the current dir
	if suite.ProjectDir == "." {
		suite.ProjectDir = cwd
	}

	// if hte project dir is not an absolute path, prepend the cwd to it
	if !filepath.IsAbs(suite.ProjectDir) {
		suite.ProjectDir = filepath.Join(cwd, suite.ProjectDir)
	}
}

// CreateDirectories creates any number of directories in the project directory
// This is used by test suites to create directories and test that the CLI does
// not clobber non-empty dirs
func (suite *BaseIntegration) CreateDirs(dirs []string) {

	var errors []string

	// iterate around the dirs and create each one
	for _, dir := range dirs {
		err := util.CreateIfNotExists(dir, os.ModePerm)
		if err != nil {
			errors = append(errors, fmt.Sprintf("unable to create project directory '%s': %s", dir, err.Error()))
		}
	}

	if len(errors) > 0 {
		suite.T().Fatalf("%s", strings.Join(errors, "\n"))
	}
}

// SetDotnetVersion sets the version of dotnet to use in the tests
// It will create the necessary global.json file
func (suite *BaseIntegration) SetDotnetVersion(version string) {

	// create the JSON string to declare the version to use
	sdkVersion := fmt.Sprintf(`{"sdk": {"version": "%s"}}`, version)

	// define the globalJsonPath
	globalJsonPath := filepath.Join(suite.ProjectDir, "global.json")

	// delete the file if it exists
	if util.Exists(globalJsonPath) {
		err := os.Remove(globalJsonPath)
		if err != nil {
			fmt.Errorf("unable to remove global.json file: %s", err.Error())
		}
	}

	os.WriteFile(globalJsonPath, []byte(sdkVersion), os.ModePerm)

}

// CheckCmdOutput checks the output of the command against the supplied pattern and
// returns a boolean stating if that pattern has been found
func (suite *BaseIntegration) CheckCmdOutput(pattern string) bool {
	var result bool

	suite.T().Logf("Looking for pattern: %s", pattern)

	re := regexp.MustCompile(pattern)
	result = re.MatchString(suite.CmdOutput)

	return result
}
