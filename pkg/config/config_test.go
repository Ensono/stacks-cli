package config

import (
	"os"
	"path/filepath"
	"testing"

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
	rendered, err := cfg.RenderTemplate(tmpl, replacements)

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
	rendered, err := cfg.RenderTemplate(tmpl, replacements)

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
	pipeline := Pipeline{
		Type:         "azdo",
		VariableFile: "build/azDevOps/azure/azuredevops-vars.yml",
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
