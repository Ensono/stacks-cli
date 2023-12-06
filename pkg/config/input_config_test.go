package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Ensono/stacks-cli/internal/util"
	yaml "github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

func setupInputConfigTests(t *testing.T, create bool) (func(t *testing.T), string) {
	t.Log("Setting up InputConfig test environment")

	// create a temporary directory
	tempDir := t.TempDir()

	// create a new file in the tempDir to represent the framework binary and git
	if create {
		for _, filename := range []string{"dotnet", "git"} {

			if util.GetPlatformOS() == "windows" {
				filename += ".exe"
			}

			path := filepath.Join(tempDir, filename)
			file, _ := os.Create(path)
			defer file.Close()

			// if not running on windows set the executable bit
			if util.GetPlatformOS() != "windows" {
				_ = os.Chmod(path, 0555)
			}
		}
	}

	deferFunc := func(t *testing.T) {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Logf("[ERROR] Unable to remove dir: %v", err)
		}
	}

	return deferFunc, tempDir
}

func TestNonExistentFrameworkBinary(t *testing.T) {

	// setup the environment
	cleanup, tempDir := setupInputConfigTests(t, false)
	defer cleanup(t)

	// set the path to the tempDir so that no files can be found
	err := os.Setenv("PATH", tempDir)
	if err != nil {
		t.Errorf("Unable to set PATH environment variable")
	}

	// create a project
	project := Project{
		Framework: Framework{
			Type: "dotnet",
		},
	}

	config := Config{
		Input: InputConfig{},
	}
	config.Input.Project = append(config.Input.Project, project)
	config.Init()

	// get the static data and unmarshal into a config object
	err = yaml.Unmarshal(config.Internal.GetFileContent("config"), &config)
	if err != nil {
		t.Errorf("Error parsing the framework definitions: %s", err.Error())
	}

	// get a list of the commands from the CheckFramework
	missing := config.Input.CheckFrameworks(&config)

	assert.Equal(t, 2, len(missing))
}

func TestIncorrectFrameworkSet(t *testing.T) {

	// setup the environment
	cleanup, tempDir := setupInputConfigTests(t, true)
	defer cleanup(t)

	// set the path to the tempDir so that no files can be found
	err := os.Setenv("PATH", tempDir)
	if err != nil {
		t.Errorf("Unable to set PATH environment variable")
	}

	// create a project
	project := Project{
		Framework: Framework{
			Type: "unknown",
		},
	}

	config := Config{
		Input: InputConfig{},
	}
	config.Input.Project = append(config.Input.Project, project)

	// get a list of the commands from the CheckFramework
	missing := config.Input.CheckFrameworks(&config)

	assert.Equal(t, 0, len(missing))
}

func TestMultipleFrameworks(t *testing.T) {

	// setup the environment
	cleanup, tempDir := setupInputConfigTests(t, true)
	defer cleanup(t)

	// set the path to the tempDir so that no files can be found
	err := os.Setenv("PATH", tempDir)
	if err != nil {
		t.Errorf("Unable to set PATH environment variable")
	}

	// create a project
	projects := []Project{
		{
			Framework: Framework{
				Type: "dotnet",
			},
		},
		{
			Framework: Framework{
				Type: "java",
			},
		},
	}

	config := Config{
		Input: InputConfig{},
	}
	config.Input.Project = projects
	config.Init()

	// get the static data and unmarshal into a config object
	err = yaml.Unmarshal(config.Internal.GetFileContent("config"), &config)
	if err != nil {
		t.Errorf("Error parsing the framework definitions: %s", err.Error())
	}

	// get a list of the commands from the CheckFramework
	missing := config.Input.CheckFrameworks(&config)

	assert.Equal(t, 2, len(missing))
	assert.Equal(t, "java", missing[0].Framework)
}

func TestValidateInput(t *testing.T) {

	// create test tables
	tables := []struct {
		config Config
		test   string
		msg    string
		count  int
	}{
		{
			Config{
				Input: InputConfig{
					Business: Business{
						Company: "MyCompany",
					},
				},
			},
			"MyCompany",
			"String should not be modified as no spaces are present",
			0,
		},
		{
			Config{
				Input: InputConfig{
					Business: Business{
						Company: "My Company",
					},
				},
			},
			"My_Company",
			"Spaces in string should be replaced with an underscore",
			1,
		},
		{
			Config{
				Input: InputConfig{
					Business: Business{
						Company: "My Fantastic Company",
					},
				},
			},
			"My_Fantastic_Company",
			"All spaces in string should be replaced with an underscore",
			1,
		},
		{
			Config{
				Input: InputConfig{
					Business: Business{
						Company: "My  Company",
					},
				},
			},
			"My_Company",
			"Consecutive spaces in string should be replaced with an underscore",
			1,
		},
	}

	// iterate around the test tables
	for _, table := range tables {

		// call ValidateInput to check the results of the test
		validations := table.config.Input.ValidateInput()

		if table.config.Input.Business.Company != table.test {
			t.Error(table.msg)
		}

		if len(validations) != table.count {
			t.Errorf("There should be %d validation errors, %d found", table.count, len(validations))
		}
	}
}
