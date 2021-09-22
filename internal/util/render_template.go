package util

import (
	"bytes"
	"text/template"

	"github.com/amido/stacks-cli/pkg/config"
)

// renderTemplate takes any string and attempts to replace items in it based
// on the values in the supplied Input object
func RenderTemplate(tmpl string, input config.Replacements) (string, error) {

	// declare var to hold the rendered string
	var rendered bytes.Buffer

	// create an object of the template
	t := template.Must(
		template.New("").Parse(tmpl),
	)

	// render the template into the variable
	err := t.Execute(&rendered, input)
	if err != nil {
		return "", err
	}

	return rendered.String(), nil
}
