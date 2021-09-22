package util

import (
	"testing"

	"github.com/amido/stacks-cli/pkg/config"
	"github.com/stretchr/testify/assert"
)

// TestNonTemplate string ensures that a string without any placeholders
// comes back from the function the same as it went in
func TestNonTemplateString(t *testing.T) {

	cfg := config.Config{}

	// set the template string
	tmpl := "Hello World!"

	replacements := config.Replacements{}
	replacements.Input = cfg.Input

	// attempt to render the template
	rendered, err := RenderTemplate(tmpl, replacements)

	assert.Equal(t, nil, err)
	assert.Equal(t, tmpl, rendered)
}

// TestTemplateString tests that a template is correctly resolved when an
// Inpout object is passed to the render function
func TestTemplateString(t *testing.T) {

	// declare the cfg object
	cfg := config.Config{}

	// set some values for the config that represent what a user might set
	cfg.Input.Business.Company = "my-company"
	cfg.Input.Business.Domain = "website"

	replacements := config.Replacements{}
	replacements.Input = cfg.Input

	// create the template string
	tmpl := "Company: {{ .Input.Business.Company }}; Domain: {{ .Input.Business.Domain }}"

	// attempt to render the template
	rendered, err := RenderTemplate(tmpl, replacements)

	// define the expected value
	expected := "Company: my-company; Domain: website"

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, rendered)
}
