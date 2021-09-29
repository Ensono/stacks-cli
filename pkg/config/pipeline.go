package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/amido/stacks-cli/internal/config/static"
	"github.com/amido/stacks-cli/internal/util"
)

type Pipeline struct {
	Type         string `mapstructure:"type"`
	TemplateFile string `mapstructure:"templateFile"`
	VariableFile string `mapstructure:"variableFile"`
}

func (p *Pipeline) GetVariableFilePath(workingDir string) string {
	path := filepath.Join(workingDir, p.VariableFile)

	return path
}

func (p *Pipeline) GetVariableTemplate(workingDir string) string {
	var template string

	// if the variableTemplate has been set attempt to find the file and read in its contents
	if p.TemplateFile != "" {
		path := filepath.Join(workingDir, p.TemplateFile)

		if util.Exists(path) {
			content, _ := ioutil.ReadFile(path)
			template = string(content)
		}
	} else {
		template = static.GetPipelineTemplate(p.Type)
	}

	return template
}

// IsSupported states if the specified pipeline is supported by Stacks
// This is only used to state which overall pipelines are possible, each
// project can define the pipelines that it supports within this overall group
func (p *Pipeline) IsSupported(name string) bool {
	var result bool

	// check to see if the lowercase version of the pipeline name is supported
	for _, v := range p.GetSupported() {
		if v == strings.ToLower(name) {
			result = true
		}
	}

	return result
}

// GetSupported returns a slice of all the currently supported pipelines
// This is determined using reflection on the current object
func (p *Pipeline) GetSupported() []string {
	// create slice of the pipelines that are supported
	// do this by getting all the fields of the pipeline object
	pipelines := []string{
		"azdo",
	}

	return pipelines
}
