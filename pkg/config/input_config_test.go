package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/amido/stacks-cli/internal/config/static"
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

			if runtime.GOOS == "windows" {
				filename += ".exe"
			}

			path := filepath.Join(tempDir, filename)
			file, _ := os.Create(path)
			defer file.Close()

			// if not running on windows set the executable bit
			if runtime.GOOS != "windows" {
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

	// get the static data and unmarshal into a config object
	framework_defs := static.Config("framework_defs")
	err = yaml.Unmarshal(framework_defs, &config.FrameworkDefs)
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

	// get the static data and unmarshal into a config object
	framework_defs := static.Config("framework_defs")
	err = yaml.Unmarshal(framework_defs, &config.FrameworkDefs)
	if err != nil {
		t.Errorf("Error parsing the framework definitions: %s", err.Error())
	}

	// get a list of the commands from the CheckFramework
	missing := config.Input.CheckFrameworks(&config)

	assert.Equal(t, 2, len(missing))
	assert.Equal(t, "java", missing[0].Framework)
}
