package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amido/stacks-cli/internal/config/static"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var basicConfiguration = []byte(`
project:
- name: tests
  framework:
    type: dotnet
    option: webapi
`)

var srcUrlConfiguration = []byte(`
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

	// setup the enviornment
	cleanup, configFile := setupTestCase(t, basicConfiguration)
	defer cleanup(t)

	// Read in the configuration file
	viper.SetConfigFile(configFile)

	// read in the static configuration of the src repo urls
	stacks_config := strings.NewReader(string(static.Config("stacks_frameworks")))
	viper.MergeConfig(stacks_config)

	// Read in the configuration file
	if err := viper.MergeInConfig(); err != nil {
		fmt.Printf("[ERROR] Unable to read in configuration file: %v\n", err)
	}

	// unmarshal the data
	err := viper.Unmarshal(&config.Input)
	if err != nil {
		t.Errorf("Unable to parse configuration data: %v", err)
	}

	// get the src URL map
	srcURLs := config.Input.Stacks.GetSrcURLMap()

	assert.Equal(t, "https://github.com/amido/stacks-dotnet", srcURLs["dotnet_webapi"].Name)
	assert.Equal(t, "https://github.com/amido/stacks-dotnet-cqrs", srcURLs["dotnet_cqrs"].Name)
	assert.Equal(t, "https://github.com/amido/stacks-dotnet-cqrs-events", srcURLs["dotnet_events"].Name)
	assert.Equal(t, "https://github.com/amido/stacks-java", srcURLs["java_webapi"].Name)
	assert.Equal(t, "https://github.com/amido/stacks-java-cqrs", srcURLs["java_cqrs"].Name)
	assert.Equal(t, "https://github.com/amido/stacks-java-cqrs-events", srcURLs["java_events"].Name)
	assert.Equal(t, "https://github.com/amido/stacks-typescript-csr", srcURLs["nodejs_csr"].Name)
	assert.Equal(t, "https://github.com/amido/stacks-typescript-ssr", srcURLs["nodejs_ssr"].Name)
}

func TestSrcUrlMap(t *testing.T) {

	config := Config{}

	// setup the enviornment
	cleanup, configFile := setupTestCase(t, srcUrlConfiguration)
	defer cleanup(t)

	// Read in the configuration file
	viper.Reset()
	viper.SetConfigFile(configFile)

	// read in the static configuration of the src repo urls
	stacks_config := strings.NewReader(string(static.Config("stacks_frameworks")))
	viper.MergeConfig(stacks_config)

	// Read in the configuration file
	if err := viper.MergeInConfig(); err != nil {
		fmt.Printf("[ERROR] Unable to read in configuration file: %v\n", err)
	}

	// unmarshal the data
	err := viper.Unmarshal(&config.Input)
	if err != nil {
		t.Errorf("Unable to parse configuration data: %v", err)
	}

	// get the src URL map
	srcURLs := config.Input.Stacks.GetSrcURLMap()

	assert.Equal(t, "https://github.com/amido/stacks-dotnet-newfeature", srcURLs["dotnet_webapi"].Name)
}
