package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	yaml "github.com/goccy/go-yaml"
	"github.com/spf13/viper"
)

var internalConfig = []byte(`
stacks:
  components:
    dotnet_webapi:
      group: dotnet
      name: webapi
      package:
        url: https://github.com/amido/stacks-dotnet-newfeature
        version: main
        type: git

    infra_keyvault:
      group: infra
      name: keyvault
      package:
        url: https://github.com/amido/stacks-infrastructure-kv
        version: main
        type: git

    golang_webapi:
      group: golang
      name: webapi
      package:
        url: https://github.com/amido/stacks-golang-webapi
        version: main
        type: git

`)

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
`)

func setupTestCase(t *testing.T, configuration []byte) (func(t *testing.T), string, string) {
	t.Log("Setting up configuration test environment")

	// create a temporary directory
	tempDir := t.TempDir()

	// write out the input configuration file to the directory
	inputConfigFilePath := filepath.Join(tempDir, "testconfig.yml")
	if err := os.WriteFile(inputConfigFilePath, configuration, 0666); err != nil {
		t.Logf("[ERROR] Unable to write out configuration file: %v", err)
	} else {
		t.Logf("Config file successfully written: %s", inputConfigFilePath)
	}

	// write out the internal configuration to the directory
	internalConfigFilePath := filepath.Join(tempDir, "internalconfig.yml")
	if err := os.WriteFile(internalConfigFilePath, internalConfig, 0666); err != nil {
		t.Logf("[ERROR] Unable to write out configuration file: %v", err)
	} else {
		t.Logf("Internal config file successfully written: %s", internalConfigFilePath)
	}

	deferFunc := func(t *testing.T) {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Logf("[ERROR] Unable to remove dir: %v", err)
		}
	}

	return deferFunc, inputConfigFilePath, internalConfigFilePath
}

func TestStacksComponents(t *testing.T) {

	var expected int = 13

	config := Config{}
	config.Init()

	// Unmarshal the internal static config
	static_config := config.Internal.GetFileContent("config")
	err := yaml.Unmarshal(static_config, &config)
	if err != nil {
		t.Errorf("Unable to parse internal config: %v", err)
	}

	// setup the enviornment
	cleanup, configFile, _ := setupTestCase(t, basicConfiguration)
	defer cleanup(t)

	// Read in the configuration file
	viper.SetConfigFile(configFile)

	// Read in the configuration file
	if err := viper.MergeInConfig(); err != nil {
		fmt.Printf("[ERROR] Unable to read in configuration file: %v\n", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		t.Errorf("Unable to parse configuration data: %v", err)
	}

	if len(config.Stacks.Components) != expected {
		t.Errorf("Unexpected number of components, %d instead of %d", len(config.Stacks.Components), expected)
	}
}

func TestOverriddenStacksComponents(t *testing.T) {
	var expected int = 15

	config := Config{}
	config.Init()

	// Unmarshal the internal static config
	static_config := config.Internal.GetFileContent("config")
	err := yaml.Unmarshal(static_config, &config)
	if err != nil {
		t.Errorf("Unable to parse internal config: %v", err)
	}

	// setup the enviornment
	cleanup, configFile, internalConfigFile := setupTestCase(t, basicConfiguration)
	defer cleanup(t)

	// Read in the configuration file
	viper.SetConfigFile(configFile)

	// Read in the configuration file
	if err := viper.MergeInConfig(); err != nil {
		fmt.Printf("[ERROR] Unable to read in configuration file: %v\n", err)
	}

	// add in the internal configuration
	data, err := os.ReadFile(internalConfigFile)
	if err != nil {
		t.Errorf("Unable to read in override for internal configuration: %v", err)
	}
	err = viper.MergeConfig(strings.NewReader(string(data)))
	if err != nil {
		t.Errorf("Unable to merge in override configuration: %v", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		t.Errorf("Unable to parse configuration data: %v", err)
	}

	if len(config.Stacks.Components) != expected {
		t.Errorf("Unexpected number of components, %d instead of %d", len(config.Stacks.Components), expected)
	}

	// check that the dotnet_webapi URL has been overrwritten
	expected_url := "https://github.com/amido/stacks-dotnet-newfeature"
	if config.Stacks.Components["dotnet_webapi"].Package.URL != expected_url {
		t.Errorf("'dotnet_webapi' URL should have been overridden, expected %s but got %s", expected_url, config.Stacks.Components["dotnet_webapi"].Package.URL)
	}

}

// TestGetComponentNames ensures that the expected names are returned from the system
func TestGetComponentNames(t *testing.T) {

	// create a set of tests
	tables := []struct {
		expectedFrameworks []string
		override           bool
	}{
		{
			[]string{"data", "dotnet", "infra", "java", "nx"},
			false,
		},
		{
			[]string{"data", "dotnet", "golang", "infra", "java", "nx"},
			true,
		},
	}

	for _, table := range tables {

		// initialise the configuration each time so that each time a new config is being tested
		config := Config{}
		config.Init()

		// Unmarshal the internal static config
		static_config := config.Internal.GetFileContent("config")
		err := yaml.Unmarshal(static_config, &config)
		if err != nil {
			t.Errorf("Unable to parse internal config: %v", err)
		}

		// if the override has been set then ensure the internalconfig os overwritten
		if table.override {
			// setup the enviornment
			cleanup, configFile, internalConfigFile := setupTestCase(t, basicConfiguration)
			defer cleanup(t)

			// Read in the configuration file
			viper.SetConfigFile(configFile)

			// Read in the configuration file
			if err := viper.MergeInConfig(); err != nil {
				fmt.Printf("[ERROR] Unable to read in configuration file: %v\n", err)
			}

			// add in the internal configuration
			data, err := os.ReadFile(internalConfigFile)
			if err != nil {
				t.Errorf("Unable to read in override for internal configuration: %v", err)
			}
			err = viper.MergeConfig(strings.NewReader(string(data)))
			if err != nil {
				t.Errorf("Unable to merge in override configuration: %v", err)
			}

			err = viper.Unmarshal(&config)
			if err != nil {
				t.Errorf("Unable to parse configuration data: %v", err)
			}
		}

		frameworks := config.Stacks.GetComponentNames()

		actualCount := len(frameworks)
		expectedCount := len(table.expectedFrameworks)

		if actualCount != expectedCount {
			t.Errorf("Unexpected number of options for %s. Expected %d got %d", table.expectedFrameworks, expectedCount, actualCount)
		}

		if !reflect.DeepEqual(frameworks, table.expectedFrameworks) {
			t.Errorf("Unexpected options returned. Expected %v got %v", table.expectedFrameworks, frameworks)
		}
	}
}

// TestGetComponentOptionsDefault checks that the options returned are expected for the internal config
// as well as ones that been added
func TestGetComponentOptionsDefault(t *testing.T) {

	tables := []struct {
		framework       string
		expectedOptions []string
		override        bool
	}{
		{
			"data",
			[]string{"fabric"},
			false,
		},
		{
			"dotnet",
			[]string{"cqrs", "webapi"},
			false,
		},
		{
			"java",
			[]string{"cqrs", "events", "webapi"},
			false,
		},
		{
			"nx",
			[]string{"apps", "next"},
			false,
		},
		{
			"infra",
			[]string{"aca", "aks", "data", "eks", "template"},
			false,
		},
		{
			"infra",
			[]string{"aca", "aks", "data", "eks", "keyvault", "template"},
			true,
		},
		{
			"golang",
			[]string{"webapi"},
			true,
		},
	}

	for _, table := range tables {

		// initialise the configuration each time so that each time a new config is being tested
		config := Config{}
		config.Init()

		// Unmarshal the internal static config
		static_config := config.Internal.GetFileContent("config")
		err := yaml.Unmarshal(static_config, &config)
		if err != nil {
			t.Errorf("Unable to parse internal config: %v", err)
		}

		// if the override has been set then ensure the internalconfig os overwritten
		if table.override {
			// setup the enviornment
			cleanup, configFile, internalConfigFile := setupTestCase(t, basicConfiguration)
			defer cleanup(t)

			// Read in the configuration file
			viper.SetConfigFile(configFile)

			// Read in the configuration file
			if err := viper.MergeInConfig(); err != nil {
				fmt.Printf("[ERROR] Unable to read in configuration file: %v\n", err)
			}

			// add in the internal configuration
			data, err := os.ReadFile(internalConfigFile)
			if err != nil {
				t.Errorf("Unable to read in override for internal configuration: %v", err)
			}
			err = viper.MergeConfig(strings.NewReader(string(data)))
			if err != nil {
				t.Errorf("Unable to merge in override configuration: %v", err)
			}

			err = viper.Unmarshal(&config)
			if err != nil {
				t.Errorf("Unable to parse configuration data: %v", err)
			}
		}

		// get the options associated with the framework
		options := config.Stacks.GetComponentOptions(table.framework)

		// get the expected number of options
		expectedCount := len(table.expectedOptions)

		if len(options) != expectedCount {
			t.Errorf("Unexpected number of options for %s. Expected %d got %d", table.framework, expectedCount, len(options))
		}

		if !reflect.DeepEqual(options, table.expectedOptions) {
			t.Errorf("Unexpected options returned. Expected %v got %v", table.expectedOptions, options)
		}
	}

}
