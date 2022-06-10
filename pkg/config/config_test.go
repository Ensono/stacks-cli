package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/amido/stacks-cli/internal/config/static"
	"github.com/amido/stacks-cli/internal/util"
	yaml "github.com/goccy/go-yaml"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setupConfigTests(t *testing.T) (func(t *testing.T), string) {
	// create a temporary directory
	tempDir := t.TempDir()

	deferFunc := func(t *testing.T) {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Logf("[ERROR] Unable to remove dir: %v", err)
		}
	}

	return deferFunc, tempDir
}

func TestGetVersionWithDefaultVersionNumber(t *testing.T) {

	config := Config{}

	// state what is expected from the method
	expected := "0.0.1-workstation"

	// get the actual response
	actual := config.GetVersion()

	assert.Equal(t, actual, expected)
}

func TestGetVersion(t *testing.T) {
	config := Config{}
	config.Input.Version = "100.98.99"

	// get the actual version
	actual := config.GetVersion()

	assert.Equal(t, actual, config.Input.Version)
}

// TestNonTemplate string ensures that a string without any placeholders
// comes back from the function the same as it went in
func TestNonTemplateString(t *testing.T) {

	cfg := Config{}

	// set the template string
	tmpl := "Hello World!"

	replacements := Replacements{}
	replacements.Input = cfg.Input

	// attempt to render the template
	rendered, err := cfg.RenderTemplate("string", tmpl, replacements)

	assert.Equal(t, nil, err)
	assert.Equal(t, tmpl, rendered)
}

// TestTemplateString tests that a template is correctly resolved when an
// Inpout object is passed to the render function
func TestTemplateString(t *testing.T) {

	// declare the cfg object
	cfg := Config{}

	// set some values for the config that represent what a user might set
	cfg.Input.Business.Company = "my-company"
	cfg.Input.Business.Domain = "website"

	replacements := Replacements{}
	replacements.Input = cfg.Input

	// create the template string
	tmpl := "Company: {{ .Input.Business.Company }}; Domain: {{ .Input.Business.Domain }}"

	// attempt to render the template
	rendered, err := cfg.RenderTemplate("string", tmpl, replacements)

	// define the expected value
	expected := "Company: my-company; Domain: website"

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, rendered)
}

func TestWriteVariableTemplate(t *testing.T) {

	// setup the environment
	cleanup, dir := setupConfigTests(t)
	defer cleanup(t)

	// create the necesssary objects
	project := Project{
		Name: "config_test",
		Directory: Directory{
			WorkingDir: dir,
		},
	}

	files := make([]PipelineFile, 1)
	files[0] = PipelineFile{
		Name: "variable",
		Path: "build/azDevOps/azure/azuredevops-vars.yml",
	}
	pipeline := Pipeline{
		Type: "azdo",
		File: files,
	}
	replacements := Replacements{}
	config := Config{}

	// call the method to create the variable file
	msg, err1 := config.WriteVariablesFile(&project, pipeline, replacements)

	// check to see if the file exists
	path := filepath.Join(dir, "build/azDevOps/azure/azuredevops-vars.yml")
	_, err2 := os.Stat(path)

	assert.Equal(t, "", msg)
	assert.Equal(t, nil, err1)
	assert.Equal(t, false, os.IsNotExist(err2))
}

func TestSetDefaultValueForInternalDomain(t *testing.T) {

	// configure the network domain settings so that it can be tested
	config := Config{
		Input: InputConfig{
			Network: Network{
				Base: NetworkBase{
					Domain: DomainType{
						External: "myproject.co.uk",
					},
				},
			},
		},
	}

	// set the default values
	config.SetDefaultValues()

	// check that the internal
	assert.Equal(t, "myproject.internal", config.Input.Network.Base.Domain.Internal)
}

func TestSetDefaultValueDoesNotChangeSpecifiedValue(t *testing.T) {

	// configure the network domain settings so that it can be tested
	config := Config{
		Input: InputConfig{
			Network: Network{
				Base: NetworkBase{
					Domain: DomainType{
						External: "myproject.co.uk",
						Internal: "myproject.newsuffix",
					},
				},
			},
		},
	}

	// set the default values
	config.SetDefaultValues()

	// check that the internal
	assert.Equal(t, "myproject.newsuffix", config.Input.Network.Base.Domain.Internal)
}

// TestSetDefaultValueWorkingDir checks that the working directory is prepended with the
// current directory if it is set as relative
// It will also check that if an absolute path is given the path is not modified
func TestSetDefaultValueWorkingDir(t *testing.T) {

	// get the current directory
	cwd, _ := os.Getwd()

	// set the absolute path to use based on the OS
	// this is required because the filepath.IsAbs function works out what an absolute path
	// based on the platform. Thus on windows "/" is not absolute
	var abs_path string
	if runtime.GOOS == "windows" {
		abs_path = "c:\\Users\\operator\\test"
	} else {
		abs_path = "/home/operator/test"
	}

	// configure a relative path for the working dir
	config_relative := Config{
		Input: InputConfig{
			Directory: Directory{
				WorkingDir: "test",
			},
		},
	}

	// configure an absolute path for the working dir
	config_absolute := Config{
		Input: InputConfig{
			Directory: Directory{
				WorkingDir: abs_path,
			},
		},
	}

	// set the default values for each of the objects
	config_relative.SetDefaultValues()
	config_absolute.SetDefaultValues()

	// check that the relative object has the correct path
	assert.Equal(t, filepath.Join(cwd, "test"), config_relative.Input.Directory.WorkingDir)

	// check that the absolute object has the correct path
	assert.Equal(t, abs_path, config_absolute.Input.Directory.WorkingDir)

}

func TestWriteCmdLog(t *testing.T) {

	// setup the environment
	cleanup, dir := setupConfigTests(t)
	defer cleanup(t)

	config := Config{}

	// check for empty error when the cmdlog is not enabled
	err := config.WriteCmdLog(dir, "test for no error")
	assert.Equal(t, nil, err)

	// update the config object to write out the cmdlog
	config.Input.Options.CmdLog = true

	// set the default values so that the cmdlog path is defined
	config.SetDefaultValues()

	// write out something to the cmdlog
	err = config.WriteCmdLog(dir, "my-command args")

	// determine if the file exists
	cmdlogExists := util.Exists(config.Self.CmdLogPath)

	// get the content of the cmdlog so it can be checked to be what is expected
	expected := fmt.Sprintf("[%s] my-command args\n", dir)
	actual, _ := ioutil.ReadFile(config.Self.CmdLogPath)

	assert.Equal(t, nil, err)                 // ensure no errors
	assert.Equal(t, true, cmdlogExists)       // file exists
	assert.Equal(t, expected, string(actual)) // check that the contents of the is as expected

	// remove the cmdlog file from the machine
	_ = os.Remove(config.Self.CmdLogPath)

}

func TestCheck(t *testing.T) {

	// create a test table to iterate around
	tables := []struct {
		conf Config
		test error
		msg  string
	}{
		{
			Config{
				Input: InputConfig{
					Project: []Project{
						{
							Name: "",
						},
					},
				},
			},
			nil,
			"An error should have been raised as no projects have been specified",
		},
		{
			Config{
				Input: InputConfig{
					Pipeline: "fred",
				},
			},
			nil,
			"'fred' is not a valid pipeline and an error should have been raised",
		},
		{
			Config{
				Input: InputConfig{
					Pipeline: "azdo",
					Project: []Project{
						{
							Name: "my-webapi",
						},
					},
				},
			},
			errors.New(""),
			"No error should be raised as a valid pipeline and project exist",
		},
	}

	for _, table := range tables {
		conf := table.conf
		res := conf.Check()

		if res == table.test {
			t.Error(table.msg)
		}
	}

}

func TestSave(t *testing.T) {

	// setup the environment
	// this creates a temporary directory into which the configuration
	// can be saved
	cleanup, dir := setupConfigTests(t)
	defer cleanup(t)

	// create the test table to work with
	tables := []struct {
		conf           Config
		usedConfigFile string
		savedFile      string
		test           error
		msg            string
	}{
		{
			Config{
				Input: InputConfig{
					Options: Options{
						SaveConfig: false,
					},
				},
			},
			"",
			"",
			nil,
			"Configuration should not be saved as the option to save is false",
		},
		{
			Config{
				Input: InputConfig{
					Options: Options{
						SaveConfig: true,
					},
				},
			},
			"config.yml",
			"",
			nil,
			"Configuration should not be saved as a configuration file has been used",
		},
		{
			Config{
				Input: InputConfig{
					Directory: Directory{
						WorkingDir: dir,
					},
					Options: Options{
						SaveConfig: true,
					},
				},
			},
			"",
			filepath.Join(dir, "stacks.yml"),
			nil,
			"Saved file is not saved in the expected location",
		},
	}

	for _, table := range tables {
		conf := table.conf
		path, res := conf.Save(table.usedConfigFile)

		if res == table.test && path != table.savedFile {
			t.Error(table.msg)
		}

		// check to see if the path exists
		if path != "" {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Error("Saved configuration file cannot be found")
			}
		}
	}
}

func TestIsDryRun(t *testing.T) {

	tables := []struct {
		cfg  Config
		test bool
		msg  string
	}{
		{
			Config{
				Input: InputConfig{
					Options: Options{
						DryRun: false,
					},
				},
			},
			false,
			"DryRun should not be enabled",
		},
		{
			Config{
				Input: InputConfig{
					Options: Options{
						DryRun: true,
					},
				},
			},
			true,
			"DryRun should be enabled",
		},
	}

	for _, table := range tables {
		res := table.cfg.IsDryRun()

		if res != table.test {
			t.Error(table.msg)
		}
	}
}

func TestUseCmdLog(t *testing.T) {

	tables := []struct {
		cfg  Config
		test bool
		msg  string
	}{
		{
			Config{
				Input: InputConfig{
					Options: Options{
						CmdLog: false,
					},
				},
			},
			false,
			"CmdLog should not be enabled",
		},
		{
			Config{
				Input: InputConfig{
					Options: Options{
						CmdLog: true,
					},
				},
			},
			true,
			"CmdLog should be enabled",
		},
	}

	for _, table := range tables {
		res := table.cfg.UseCmdLog()

		if res != table.test {
			t.Error(table.msg)
		}
	}
}

func TestNoCleanup(t *testing.T) {

	tables := []struct {
		cfg  Config
		test bool
		msg  string
	}{
		{
			Config{
				Input: InputConfig{
					Options: Options{
						NoCleanup: false,
					},
				},
			},
			false,
			"NoCleanup should not be enabled",
		},
		{
			Config{
				Input: InputConfig{
					Options: Options{
						NoCleanup: true,
					},
				},
			},
			true,
			"NoCleanup should be enabled",
		},
	}

	for _, table := range tables {
		res := table.cfg.NoCleanup()

		if res != table.test {
			t.Error(table.msg)
		}
	}
}

func TestNoBanner(t *testing.T) {

	tables := []struct {
		cfg  Config
		test bool
		msg  string
	}{
		{
			Config{
				Input: InputConfig{
					Options: Options{
						NoBanner: false,
					},
				},
			},
			false,
			"NoBanner should not be enabled",
		},
		{
			Config{
				Input: InputConfig{
					Options: Options{
						NoBanner: true,
					},
				},
			},
			true,
			"NoBanner should be enabled",
		},
	}

	for _, table := range tables {
		res := table.cfg.NoBanner()

		if res != table.test {
			t.Error(table.msg)
		}
	}
}

func TestForce(t *testing.T) {

	tables := []struct {
		cfg  Config
		test bool
		msg  string
	}{
		{
			Config{
				Input: InputConfig{
					Options: Options{
						Force: false,
					},
				},
			},
			false,
			"Force should not be enabled",
		},
		{
			Config{
				Input: InputConfig{
					Options: Options{
						Force: true,
					},
				},
			},
			true,
			"Force should be enabled",
		},
	}

	for _, table := range tables {
		res := table.cfg.Force()

		if res != table.test {
			t.Error(table.msg)
		}
	}
}

func TestExecuteCommand(t *testing.T) {

	// Create a logger
	config := Config{}
	logger := log.New()

	// setup the environment
	cleanup, dir := setupConfigTests(t)
	defer cleanup(t)

	// create the command and arguments that need to be run
	command := "echo"
	arguments := "HelloWorld!"

	// call the ExecuteCommand to run the command and get the result
	result, err := config.ExecuteCommand(dir, logger, command, arguments, false, false)

	// perform the necessary tests
	assert.Equal(t, nil, err)
	assert.Equal(t, "HelloWorld!", result)

}

func TestGetFrameworkCommands(t *testing.T) {

	config := Config{}

	// get the static data and unmarshal into a config object
	framework_defs := static.Config("framework_defs")
	err := yaml.Unmarshal(framework_defs, &config.FrameworkDefs)
	if err != nil {
		t.Errorf("Error parsing the framework definitions: %s", err.Error())
	}

	// create the testing tables that are required
	tables := []struct {
		name    string
		count   int
		message string
	}{
		{
			"dotnet",
			2,
			"Commands should contain 2 elements",
		},
		{
			"java",
			3,
			"Commands should contain 2 elements",
		},
		{
			"notvalid",
			0,
			"Commands should contain no elements",
		},
	}

	for _, table := range tables {

		res := config.GetFrameworkCommands(table.name)

		if len(res) != table.count {
			t.Error(table.message)
		}
	}
}

func TestGetFrameworks(t *testing.T) {

	config := Config{}

	// attempt to get the stacks as a model
	stacks, err := config.GetFrameworks()

	assert.Equal(t, nil, err)
	assert.Equal(t, "Amido.Stacks.Templates", stacks.Dotnet.Webapi.Name)

}
