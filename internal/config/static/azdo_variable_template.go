// This file contains the static template to be used for writing out the variables
// for the Azure DevOps pipeline

package static

import (
	_ "embed"
)

//go:embed azdo_variable_template.yml
var Azdo_Variable_Template_Tmpl string

func GetPipelineTemplate(name string) string {
	var template string

	switch name {
	case "azdo":
		template = Azdo_Variable_Template_Tmpl
	}

	return template
}
