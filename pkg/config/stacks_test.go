package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	yaml "github.com/goccy/go-yaml"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var basicConfiguration = []byte(`
input:
  project:
    - name: tests
      framework:
        type: dotnet
        option: webapi
`)

var srcUrlConfiguration = []byte(`
input:
  project:
    - name: tests
      framework:
        type: dotnet
        option: webapi

stacks:
  dotnet:
    webapi:
      name: https://github.com/amido/stacks-dotnet-newfeature
      version: main
`)

func setupTestCase(t *testing.T, configuration []byte) (func(t *testing.T), string) {
	t.Log("Setting up configuration test environment")

	// create a temporary directory
	tempDir := t.TempDir()

	// write out the configuration file to the directory
	configFilePath := filepath.Join(tempDir, "testconfig.yml")
	if err := os.WriteFile(configFilePath, configuration, 0666); err != nil {
		t.Logf("[ERROR] Unable to write out configuration file: %v", err)
	} else {
		t.Logf("Config file successfully written")
	}

	deferFunc := func(t *testing.T) {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Logf("[ERROR] Unable to remove dir: %v", err)
		}
	}

	return deferFunc, configFilePath
}

func TestDefaultSrcUrlMap(t *testing.T) {

	config := Config{}
	config.Init()

	// setup the enviornment
	cleanup, configFile := setupTestCase(t, basicConfiguration)
	defer cleanup(t)

	// Read in the configuration file
	viper.SetConfigFile(configFile)

	// Read in the configuration file
	if err := viper.MergeInConfig(); err != nil {
		fmt.Printf("[ERROR] Unable to read in configuration file: %v\n", err)
	}

	// Unmarshal the internal static config
	err := yaml.Unmarshal(config.Internal.GetFileContent("config"), &config)
	if err != nil {
		t.Errorf("Unable to parse internal config: %v", err)
	}

	// unmarshal the data
	err = viper.Unmarshal(&config.Input)
	if err != nil {
		t.Errorf("Unable to parse configuration data: %v", err)
	}

	// get the src URL map
	srcURLs := config.Stacks.GetSrcURLMap()

	assert.Equal(t, "Amido.Stacks.Templates", srcURLs["dotnet_webapi"].Name)
	assert.Equal(t, "Amido.Stacks.CQRS.Templates", srcURLs["dotnet_cqrs"].Name)
	assert.Equal(t, "https://github.com/amido/stacks-java", srcURLs["java_webapi"].Name)
	assert.Equal(t, "https://github.com/amido/stacks-java-cqrs", srcURLs["java_cqrs"].Name)
	assert.Equal(t, "https://github.com/amido/stacks-java-cqrs-events", srcURLs["java_events"].Name)
	assert.Equal(t, "https://github.com/amido/stacks-nx", srcURLs["nx_next"].Name)
	assert.Equal(t, "https://github.com/amido/stacks-nx", srcURLs["nx_apps"].Name)
}

func TestSrcUrlMap(t *testing.T) {

	config := Config{}
	config.Init()

	// Unmarshal the internal static config
	err := yaml.Unmarshal(config.Internal.GetFileContent("config"), &config)
	if err != nil {
		t.Errorf("Unable to parse internal config: %v", err)
	}

	// setup the enviornment
	cleanup, configFile := setupTestCase(t, srcUrlConfiguration)
	defer cleanup(t)

	// Read in the configuration file
	viper.Reset()
	viper.SetConfigFile(configFile)

	// Read in the configuration file
	if err := viper.MergeInConfig(); err != nil {
		fmt.Printf("[ERROR] Unable to read in configuration file: %v\n", err)
	}

	// unmarshal the data
	err = viper.Unmarshal(&config)
	if err != nil {
		t.Errorf("Unable to parse configuration data: %v", err)
	}

	// get the src URL map
	srcURLs := config.Stacks.GetSrcURLMap()

	assert.Equal(t, "https://github.com/amido/stacks-dotnet-newfeature", srcURLs["dotnet_webapi"].Name)
}
