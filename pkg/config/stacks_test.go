package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	yaml "github.com/goccy/go-yaml"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var internalConfig = []byte(`
stacks:
  components:
    - group: dotnet
      name: webapi
      package:
        url: https://github.com/amido/stacks-dotnet-newfeature
        version: main
        type: git

    - group: infra
      name: keyvault
      package:
        url: https://github.com/amido/stacks-infrastructure-kv
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

// TestComponents tests the built in components and their values
func TestComponents(t *testing.T) {

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

	config.Stacks.SetUniqueComponents()

	// Ensure that the components have been set correctly, according to the built in config
	data := make(map[string]string)
	data["dotnet_webapi"] = "Amido.Stacks.Templates"
	data["dotnet_cqrs"] = "Amido.Stacks.CQRS.Templates"
	data["java_webapi"] = "https://github.com/amido/stacks-java"
	data["java_cqrs"] = "https://github.com/amido/stacks-java-cqrs"
	data["java_events"] = "https://github.com/amido/stacks-java-cqrs-events"
	data["nx_next"] = "https://github.com/amido/stacks-nx"
	data["nx_apps"] = "https://github.com/amido/stacks-nx"
	data["infra_aks"] = "https://github.com/amido/stacks-infrastructure-aks/"

	for key, value := range data {
		assert.Equal(t, value, config.Stacks.GetComponentPackageRef(key))
	}

	// check that there are the expected number of configurations
	if len(data) != config.Stacks.GetComponentCount() {
		t.Errorf("inconsistent number of components, expected %d but got %d", len(data), config.Stacks.GetComponentCount())
	}
}

// TestComponentOverrideAndAdd checks that an existing component can be overridden and that news
// ones can be supplied
func TestComponentOverrideAndAdd(t *testing.T) {

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
	viper.MergeConfig(strings.NewReader(string(data)))

	// unmarshal the data
	err = viper.Unmarshal(&config)
	if err != nil {
		t.Errorf("Unable to parse configuration data: %v", err)
	}

	config.Stacks.SetUniqueComponents()

	assert.Equal(t, "https://github.com/amido/stacks-dotnet-newfeature", config.Stacks.GetComponentPackageRef("dotnet_webapi"))
	assert.Equal(t, "https://github.com/amido/stacks-infrastructure-kv", config.Stacks.GetComponentPackageRef("infra_keyvault"))

}

func TestSetUniqueCommands(t *testing.T) {

	tables := []struct {
		cfg   Config
		count int
		msg   string
	}{
		{
			Config{
				Stacks: Stacks{
					Components: []StacksComponent{
						{
							Group: "dotnet",
							Name:  "webapi",
						},
					},
				},
			},
			1,
			"Slice should be unaffected",
		},
		{
			Config{
				Stacks: Stacks{
					Components: []StacksComponent{
						{
							Group: "dotnet",
							Name:  "webapi",
							Package: Package{
								URL: "https://github.com/amido/stacks-dotnet-cqrs",
							},
						},
						{
							Group: "dotnet",
							Name:  "webapi",
							Package: Package{
								URL: "https://github.com/amido/stacks-dotnet",
							},
						},
					},
				},
			},
			1,
			"There should only be one component in the final config",
		},
	}

	for _, table := range tables {

		// ensure that no duplicates exist
		table.cfg.Stacks.SetUniqueComponents()

		// check that the number of components is correct
		if table.cfg.Stacks.GetComponentCount() != table.count {
			t.Error(table.msg)
		}
	}
}
