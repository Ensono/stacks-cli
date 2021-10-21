package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/amido/stacks-cli/internal/config/static"
	"github.com/amido/stacks-cli/internal/util"
)

type Pipeline struct {
	Type         string                `mapstructure:"type"`
	File         []PipelineFile        `mapstructure:"files"`
	Template     []PipelineFile        `mapstructure:"templates"`
	Replacements []PipelineReplacement `mapstructure:"replacements"`
}

type PipelineFile struct {
	Name      string `mapstructure:"name"`
	Path      string `mapstructure:"path"`
	NoReplace bool   `mapstructure:"noreplace"`
}

type PipelineReplacement struct {
	Pattern string `mapstructure:"pattern"`
	Value   string `mapstructure:"value"`
}

// GetFilePath iterates around the either the File or Template slice
// looking for the specified name, if the name is found then it will
// return the path associated with the name
func (p *Pipeline) GetFilePath(filetype string, workingDir string, name string) string {

	var path string
	var list []PipelineFile

	switch filetype {
	case "file":
		list = p.File
	case "template":
		list = p.Template
	}

	// iterate around the list looking for the specified name
	for _, item := range list {
		if item.Name == name {
			path = item.Path
			break
		}
	}

	// only prepend the workingDir to the path if it is not ""
	if path != "" {
		path = filepath.Join(workingDir, path)
	}

	return path
}

func (p *Pipeline) GetVariableTemplate(workingDir string) string {
	var template string

	// determine if a template variable path has been set
	templateVarPath := p.GetFilePath("template", workingDir, "variable")

	// if the variableTemplate has been set attempt to find the file and read in its contents
	if templateVarPath != "" {
		path := filepath.Join(workingDir, templateVarPath)

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

// ReplacePatterns replaces the phrases that are found in the build file according to
// the regex pattern with the specified value
func (p *Pipeline) ReplacePatterns(dir string) []error {

	var err error
	var errs []error
	errs = make([]error, 0)

	// Return if there are no replacements to perform
	if len(p.Replacements) == 0 {
		return errs
	}

	// iterate around all the files that have been set
	for _, item := range p.File {

		// continue onto the next iteration if the file is has NoReplace set
		if item.NoReplace {
			continue
		}

		// determine the path to the build file
		buildFile := filepath.Join(dir, item.Path)
		if !util.Exists(buildFile) {
			err = fmt.Errorf("unable to find '%s' file: %s", item.Name, buildFile)
			errs = append(errs, err)
			return errs
		}

		// read the file into a variable
		content, err := ioutil.ReadFile(buildFile)
		if err != nil {
			errs = append(errs, err)
			return errs
		}

		// iterate around the replacements to get the pattern and the replacement value
		for _, replacement := range p.Replacements {

			// create the regex object from the pattern
			pattern := fmt.Sprintf(`(?m)%s`, replacement.Pattern)
			re, err := regexp.Compile(pattern)

			if err != nil {
				errs = append(errs, err)
				continue
			}
			content = re.ReplaceAll(content, []byte(replacement.Value))
		}

		// write out the file
		err = os.WriteFile(buildFile, content, 0666)
		if err != nil {
			errs = append(errs, err)
		}

	}

	return errs

}
