package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Ensono/stacks-cli/internal/config/staticFiles"
	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/sirupsen/logrus"
)

type Pipeline struct {
	Type         string                `mapstructure:"type"`
	File         []PipelineFile        `mapstructure:"files"`
	Template     []PipelineFile        `mapstructure:"templates"`
	Items        []string              `mapstructure:"items"`
	Replacements []PipelineReplacement `mapstructure:"replacements"`
	Logger       *logrus.Logger        `mapstructure:"-"`
}

// SetLogger sets the logger for the Pipeline struct
func (p *Pipeline) SetLogger(logger *logrus.Logger) {
	p.Logger = logger
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

		// update the path if it is not an absolute path
		path := templateVarPath
		if !filepath.IsAbs(path) {
			path = filepath.Join(workingDir, templateVarPath)
		}

		if util.Exists(path) {
			content, _ := os.ReadFile(path)
			template = string(content)
		}
	} else {
		template = staticFiles.GetPipelineTemplate(p.Type)
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
		"gha",
	}

	return pipelines
}

// ReplacePatterns replaces the phrases that are found in the build file according to
// the regex pattern with the specified value
func (p *Pipeline) ReplacePatterns(config *Config, inputs Replacements, dir string) []error {

	var count int
	var errs []error
	var filelist []string
	errs = make([]error, 0)

	// Return if there are no replacements to perform
	if len(p.Replacements) == 0 {
		return errs
	}

	// iterate around the p.File and get a list of all the files
	// these can be treated as globs from the filesystem
	for _, item := range p.File {

		// if no replace has been set on the file then continue to the next iteration
		if item.NoReplace {
			continue
		}

		files, _ := util.GetFileList(item.Path, dir)

		// add the files to the filelist
		filelist = append(filelist, files...)
	}

	// now iterate over the items that have been set
	for _, item := range p.Items {
		files, _ := util.GetFileList(item, dir)
		filelist = append(filelist, files...)
	}

	// Ensure the filelist does not contain duplicates
	unique := make(map[string]struct{})
	result := make([]string, 0, len(filelist))
	for _, f := range filelist {
		if _, seen := unique[f]; !seen {
			unique[f] = struct{}{}
			result = append(result, f)
		}
	}
	filelist = result

	// iterate around all the files that have been set
	for _, item := range filelist {

		p.Logger.Infof("\t%s", item)

		// read the file into a variable
		contentBytes, err := os.ReadFile(item)
		if err != nil {
			errs = append(errs, err)
			return errs
		}
		content := string(contentBytes)

		// iterate around the replacements to get the pattern and the replacement value
		for _, replacement := range p.Replacements {

			// TODO Ensure that the template for the replacement is rendered. This will mean that the value
			// the regex is replacing can come from the inputs of the CLI
			// render the replacement value as a template
			replacement_value, err := config.RenderTemplate("regex", replacement.Value, inputs)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			// create the regex object from the pattern
			pattern := fmt.Sprintf(`(?m)%s`, replacement.Pattern)

			// The following replacement handles file paths from Windows machines
			// The path C:\Users\anon will be read in as C:\\Users\\anon, which fails any matching
			// By replacing these literal double backslashes with single backslashes, the regex can match correctly
			pattern = strings.ReplaceAll(pattern, `\\`, `\`)
			re, err := regexp.Compile(pattern)

			if err != nil {
				errs = append(errs, err)
				continue
			}

			p.Logger.Debugf("\t\tPattern: %q\n", pattern)
			p.Logger.Debugf("\t\tReplacement: %q\n", replacement_value)

			// Use FindAllStringSubmatchIndex to get all matches and submatches in one pass
			count = 0
			matches := re.FindAllStringSubmatchIndex(content, -1)
			if len(matches) == 0 {
				p.Logger.Infof("\t\t%d replacements made for pattern `%s`", count, pattern)
				continue
			}

			var resultBuilder strings.Builder
			lastIndex := 0
			for _, matchIdx := range matches {
				// matchIdx[0] and matchIdx[1] are the start and end of the full match
				start, end := matchIdx[0], matchIdx[1]
				// Write content before the match
				resultBuilder.WriteString(content[lastIndex:start])

				// Prepare replacement string
				result := replacement_value
				// For each submatch, replace ${i} with the submatch value
				for i := 1; i < len(matchIdx)/2; i++ {
					subStart, subEnd := matchIdx[2*i], matchIdx[2*i+1]
					submatch := ""
					if subStart != -1 && subEnd != -1 {
						submatch = content[subStart:subEnd]
					}
					placeholder := fmt.Sprintf("${%d}", i)
					result = strings.ReplaceAll(result, placeholder, submatch)
				}

				resultBuilder.WriteString(result)
				count++
				lastIndex = end
			}
			// Write the rest of the content after the last match
			resultBuilder.WriteString(content[lastIndex:])
			content = resultBuilder.String()
			p.Logger.Infof("\t\t%d replacements made for pattern `%s`", count, pattern)
		}

		// write out the file
		err = os.WriteFile(item, []byte(content), 0666)
		if err != nil {
			errs = append(errs, err)
		}

	}

	return errs

}
