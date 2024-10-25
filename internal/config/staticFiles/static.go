package staticFiles

import (
	_ "embed"
)

// Set the banner that is written out to the screen when stacks is run
//
//go:embed banner.txt
var IntFile_Banner string

//go:embed config.yml
var IntFile_config string

//go:embed ado_variable_template.yml
var Ado_Variable_Template_Tmpl string

//go:embed help.yml
var Help_Messages string

func GetPipelineTemplate(name string) string {
	var template string

	switch name {
	case "azdo":
		template = Ado_Variable_Template_Tmpl
	}

	return template
}
