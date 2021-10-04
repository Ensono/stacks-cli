package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/amido/stacks-cli/internal/constants"
	"github.com/bobesa/go-domain-util/domainutil"
)

// import (
// 	"github.com/dnitsch/scaffold/internal/util"
// 	"github.com/dnitsch/scaffold/pkg/scaffold"
// )

// func ReadSelfConfigFile(input scaffold.InputConfig) (*scaffold.SelfConfig, error) {
// 	return readSelfConfigFile(input)
// }

// // readSelfConfigFile constructs self config for CLI based on bundle resources
// func readSelfConfigFile(input scaffold.InputConfig) (*scaffold.SelfConfig, error) {

// 	sharedT, err := ParseShared()
// 	util.CheckErrors(err)

// 	specificT, err := ParseSpecific(input.Platform, input.Deployment, input.ProjectType)

// 	s := scaffold.SelfConfig{
// 		Shared:   &sharedT,
// 		Specific: &specificT,
// 	}
// 	// TODO: feat request allow overwrite of self config from outside (as long as it can be parsed back to a SelfConfig)
// 	return &s, err
// }

type TypeDetail struct {
	Gitrepo                  string    `yaml:"git_repo"`
	Gitref                   string    `yaml:"git_ref"`
	Localpath                string    `yaml:"local_path"`
	FilenameReplacementPaths []string  `yaml:"filename_replacement_paths,omitempty"`
	Searchvalue              string    `yaml:"search_value,omitempty"`
	Foldermap                Foldermap `mapstructure:"folder_map"`
}

type Foldermap struct {
	Src  string `mapstructure:"src"`
	Dest string `mapstructure:"dest"`
}

type SelfConfig struct {
	Shared   *TypeDetail
	Specific *TypeDetail

	ProjectPaths map[string]string

	CmdLogPath string
}

type OutputConfig struct {
	TmpPath   string
	ZipPath   string
	UnzipPath string
	NewPath   string
}

type ReplaceConfig struct {
	Files  []string          `yaml:"files"`
	Values map[string]string `yaml:"values"`
}

type Config struct {
	Input   InputConfig
	Self    SelfConfig
	Output  OutputConfig
	Replace []ReplaceConfig
}

// GetVersion returns the current version of the application
// It will check to see uif the Version is empty, if it is, it will
// set and identifiable local build version
func (config *Config) GetVersion() string {
	var version string

	version = config.Input.Version

	fmt.Println(version)

	if version == "" {
		version = constants.DefaultVersion
	}

	return version
}

// SetPaths sets the current project path
func (selfConfig *SelfConfig) AddPath(project Project, path string) {
	if selfConfig.ProjectPaths == nil {
		selfConfig.ProjectPaths = make(map[string]string)
	}
	selfConfig.ProjectPaths[project.GetId()] = path
}

// GetPath returns the path for the current project
func (selfConfig *SelfConfig) GetPath(project Project) string {
	return selfConfig.ProjectPaths[project.GetId()]
}

// WriteVariablesFile writes out the variables template file for the build pipeline
func (config *Config) WriteVariablesFile(project *Project, pipelineSettings Pipeline, replacements Replacements) (string, error) {
	var err error
	var variableFile string
	var variableTemplate string

	variableFile = pipelineSettings.GetFilePath("file", project.Directory.WorkingDir, "variable")
	variableTemplate = pipelineSettings.GetVariableTemplate(project.Directory.WorkingDir)

	// render the variable file
	rendered, err := config.RenderTemplate(variableTemplate, replacements)

	if err != nil {
		return "Problem rendering variable template file", err
	}

	// get the dirname of the path and ensure it exists
	// this should not be needed in normal operation as the file structure should already exist
	dir := filepath.Dir(variableFile)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		fmt.Printf("%v", err)
	}

	err = os.WriteFile(variableFile, []byte(rendered), 0666)
	if err != nil {
		return "Problem writing out variable file", err
	}

	return "", err
}

// renderTemplate takes any string and attempts to replace items in it based
// on the values in the supplied Input object
func (config *Config) RenderTemplate(tmpl string, input Replacements) (string, error) {

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

// SetDefaultValues sets values in the config object that are based off other values in the
// config object
// For example, if the internal domain name has not been set then it will be based on the
// external domain name, with the TLD replaced with `internal`
func (config *Config) SetDefaultValues() {

	// Check that the internal domain name
	if config.Input.Network.Base.Domain.Internal == "" {

		// get the external domain and replace the suffix with internal
		internal := config.Input.Network.Base.Domain.External
		internal = strings.Replace(internal, domainutil.DomainSuffix(internal), "internal", -1)
		config.Input.Network.Base.Domain.Internal = internal
	}

	// Set the currentdirectory to the path that the CLI is currently running in
	cwd, _ := os.Getwd()
	config.Self.CmdLogPath = filepath.Join(cwd, "cmdlog.txt")
}

// WriteCmdLog writes the command out a log file in the directory that the CLI is being run
// The cmd is only written out if the option to do so has been set in the config
func (config *Config) WriteCmdLog(cmd string) error {

	var err error

	// return empty error if not logging commands
	if !config.Input.Options.CmdLog {
		return err
	}

	// get a reference to the file, either to create or append to the file
	f, err := os.OpenFile(config.Self.CmdLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// write out the cmd to the file
	if _, err := f.WriteString(fmt.Sprintf("%s\n", cmd)); err != nil {
		return err
	}

	return err
}

/*
// Create creates a config object based on parsed input config
func New(data InputConfig) (*Config, error) {
	tmpPath, err := os.MkdirTemp("", "source")
	if err != nil {
		return nil, err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	selfConf, err := readSelfConfigFile(data)
	if err != nil {
		helper.TraceInfo("Failed to read self config")
		return nil, err
	}

	conf := Config{
		Output: &OutputConfig{
			NewPath:   fmt.Sprintf("%s/%s", pwd, data.ProjectName),
			TmpPath:   tmpPath,
			ZipPath:   fmt.Sprintf("%s/source.zip", tmpPath),
			UnzipPath: path.Join(tmpPath, "wip", selfConf.Specific.Localpath),
		},
		Input: &data,
		Self:  selfConf,
	}

	helper.TraceInfo(fmt.Sprintf("New Project Dir: %s\n", conf.Output.NewPath))

	helper.TraceInfo(fmt.Sprintf("Temp Path: %s\n", conf.Output.TmpPath))

	return &conf, err
}

// Create creates a config object based on bytes stream read from a config file
func NewBytes(data []byte) (*Config, error) {

	t := InputConfig{}

	if err := yaml.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	conf, err := New(t)
	return conf, err
}

// readSelfConfigFile constructs self config for CLI based on bundle resources
func readSelfConfigFile(input InputConfig) (*SelfConfig, error) {
	// var err error

	sharedT, err := ParseLocalShared()
	// if errShared != nil {
	if err != nil {
		return nil, err
	}

	specificT, err := ParseLocalSpecific(input.Platform, input.Deployment, input.ProjectType)

	if err != nil {
		return nil, err
	}

	s := SelfConfig{
		Shared:   &sharedT,
		Specific: &specificT,
	}

	// TODO: feat request allow overwrite of self config from outside (as long as it can be parsed back to a SelfConfig)
	return &s, err
}
*/
