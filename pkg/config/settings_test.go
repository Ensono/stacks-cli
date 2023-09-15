package config

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setupPipelines(t *testing.T) Settings {
	t.Log("Setting up pipeline settings")

	// setup the pipelines
	pipelines := make([]Pipeline, 1)

	pipelines[0] = Pipeline{
		Type: "azdo",
	}

	// create some framework versions to work with
	versions := make([]SettingsFrameworkCommands, 2)
	versions[0] = SettingsFrameworkCommands{
		Name:    "dotnet",
		Version: "3.1",
	}
	versions[1] = SettingsFrameworkCommands{
		Name: "git",
	}

	settings := Settings{
		Framework: SettingsFramework{
			Name:     "dotnet",
			Commands: versions,
		},
	}
	settings.Pipeline = pipelines

	return settings
}

func TestPipelineExists(t *testing.T) {

	// get the settings from the setup
	settings := setupPipelines(t)

	assert.NotEqual(t, []Pipeline{}, settings.GetPipelines("azdo"))
}

func TestPipelineDoesNotExist(t *testing.T) {

	// get the settings from the setup
	settings := setupPipelines(t)

	assert.Equal(t, []Pipeline{}, settings.GetPipelines("jenkins"))
}

func TestPipelineReturnCount(t *testing.T) {

	// get the settings from the setup
	settings := setupPipelines(t)

	pipelineCount := len(settings.GetPipelines("azdo"))

	assert.Equal(t, 1, pipelineCount)
}

func TestFrameworkCommandVersions(t *testing.T) {

	// get the settings from the setup
	settings := setupPipelines(t)

	// create a table of tests
	tables := []struct {
		command string
		test    string
		msg     string
	}{
		{
			"dotnet",
			"3.1",
			"Dotnet should come back with a version number",
		},
		{
			"git",
			"",
			"Git should have an empty version number, even though it has been defined",
		},
		{
			"java",
			"",
			"No version number should be returned for java",
		},
	}

	// iterate around the table and perform each test
	for _, table := range tables {

		// get the result for the framework command
		res := settings.GetRequiredVersion(table.command)

		if res != table.test {
			t.Error(table.msg)
		}
	}
}

func TestCompareVersion(t *testing.T) {

	settings := Settings{}
	logger := log.New()

	// create a table of tests to check that all version checks and
	// constraints work as expected
	tables := []struct {
		constraint string
		version    string
		test       bool
		message    string
	}{
		{
			"",
			"5.0.38",
			true,
			"Version should pass version comparison as no constraint has been set",
		},
		{
			">= 3.1",
			"5.0.38",
			true,
			"Version should pass because 5.0.38 is greater than 3.1",
		},
		{
			">= 3.1, < 3.2",
			"5.0.38",
			false,
			"Version should not pass because 5.0.38 is not in the range 3.1 to 3.2",
		},
		{
			"> 6",
			"5.0.38",
			false,
			"Version should not pass because 5.0.38 is less than 6",
		},
		{
			"> 1.8",
			"1.8.0_301",
			true,
			"Version should pass because the point release is greater than 0",
		},
	}

	// iterate around the tables
	for _, table := range tables {

		// get the result of the compareVersion
		res := settings.CompareVersion(table.constraint, table.version, logger)

		if res != table.test {
			t.Error(table.message)
		}
	}
}
