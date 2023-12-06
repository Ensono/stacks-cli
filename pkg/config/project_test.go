package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Ensono/stacks-cli/internal/constants"
	"github.com/stretchr/testify/assert"
)

func setupProjectTestCase(t *testing.T) (func(t *testing.T), string) {
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

// TestSettingsFileDoesNotExist tests to see of the default stackscli.yml file is missing
func TestSettingsFileDoesNotExist(t *testing.T) {

	config := Config{}

	// set the default value for the settings file using the constant, which is
	// what the CLI options will do when running
	config.Input.SettingsFile = constants.SettingsFile

	cleanup, tempDir := setupProjectTestCase(t)
	defer cleanup(t)

	// setup the project configuration
	config.Input.Project = make([]Project, 1)

	err := config.Input.Project[0].setSettingsFilePath(tempDir, &config)

	// TODO: Need to work out how to check for a specific type of error
	// In this case need to check for os.IsNotExist
	// This is so that the function can return more error types if required
	assert.NotEqual(t, err, nil)
}

// TestDefaultSettingsFileExists tests to see if the default settings file exists and
// that the full path to the file has been set
func TestDefaultSettingsFileExists(t *testing.T) {

	config := Config{}

	// set the default value for the settings file using the constant, which is
	// what the CLI options will do when running
	config.Input.SettingsFile = constants.SettingsFile

	cleanup, tempDir := setupProjectTestCase(t)
	defer cleanup(t)

	// create a settings file in the tempDir
	file, _ := os.Create(filepath.Join(tempDir, config.Input.SettingsFile))
	defer file.Close()

	// setup the project configuration
	config.Input.Project = make([]Project, 1)

	err := config.Input.Project[0].setSettingsFilePath(tempDir, &config)

	// define the expected result
	expected := filepath.Join(tempDir, "stackscli.yml")

	assert.Equal(t, err, nil)
	assert.Equal(t, expected, config.Input.Project[0].SettingsFile)
}

func TestSpecifiedSettingsFileExists(t *testing.T) {

	config := Config{}

	// specific a specific name for the settings file, this is the equivalent
	// of setting it on the command line or in the configuration file
	config.Input.SettingsFile = "altsettings.yml"

	cleanup, tempDir := setupProjectTestCase(t)
	defer cleanup(t)

	// create a settings file in the tempDir
	file, _ := os.Create(filepath.Join(tempDir, config.Input.SettingsFile))
	defer file.Close()

	// setup the project configuration
	config.Input.Project = make([]Project, 1)

	err := config.Input.Project[0].setSettingsFilePath(tempDir, &config)

	// define the expected result
	expected := filepath.Join(tempDir, "altsettings.yml")

	assert.Equal(t, err, nil)
	assert.Equal(t, expected, config.Input.Project[0].SettingsFile)
}
