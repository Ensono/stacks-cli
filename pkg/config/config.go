package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/amido/stacks-cli/internal/constants"
	"github.com/bobesa/go-domain-util/domainutil"
	yaml "github.com/goccy/go-yaml"
)

type SelfConfig struct {
	ProjectPaths map[string]string

	CmdLogPath string
}

type ReplaceConfig struct {
	Files  []string          `yaml:"files"`
	Values map[string]string `yaml:"values"`
}

type Config struct {
	Input   InputConfig
	Self    SelfConfig
	Replace []ReplaceConfig
}

// Check checks the configuration and ensures that there are some projects
// to work with and that the chosen pipeline is supported
// It also sets some defaults based on other settings in the configuration
func (c *Config) Check() error {
	var err error

	// determine if any projects have been specified
	if len(c.Input.Project) == 1 && c.Input.Project[0].Name == "" {
		return fmt.Errorf("no projects have been defined")
	}

	// check to see if the the specified pipeline is supported
	pipeline := Pipeline{}
	if !pipeline.IsSupported(c.Input.Pipeline) {
		return fmt.Errorf("specified pipeline is not supported - %s %v", c.Input.Pipeline, pipeline.GetSupported())
	}

	// set necessary default values
	c.SetDefaultValues()

	return err
}

// IsDryRun returns the boolean value of the dryrun option
func (c *Config) IsDryRun() bool {
	return c.Input.Options.DryRun
}

// UseCmdLog states of the command log should be used
func (c *Config) UseCmdLog() bool {
	return c.Input.Options.CmdLog
}

// NoCleanup returns a boolean stating if the app should perform cleanup functions
func (c *Config) NoCleanup() bool {
	return c.Input.Options.NoCleanup
}

// NoBanner returns the option to no display the Stacks banner
func (c *Config) NoBanner() bool {
	return c.Input.Options.NoBanner
}

// Force states if projects should be overwritten
func (c *Config) Force() bool {
	return c.Input.Options.Force
}

// Save saves the user's configuration to a file
// This is only applicable if a configuration file has not been used and the option
// to save the configuration has been set as an option
func (c *Config) Save(usedConfig string) (string, error) {
	var err error
	var savedConfigFile string

	// return with a nil error if a configuration file has been specified
	// or the SaveConfig has not been set
	if !c.Input.Options.SaveConfig {
		return "", nil
	} else if c.Input.Options.SaveConfig && usedConfig == "" {
		return "", nil
	}

	// determine the path to the saveConfigFile
	savedConfigFile = filepath.Join(c.Input.Directory.WorkingDir, "stacks.yml")

	// deserialise the data from the Config.Input object
	data, err := yaml.Marshal(&c.Input)
	if err != nil {
		return savedConfigFile, fmt.Errorf("problem converting configuration to YAML syntax")
	}

	err = ioutil.WriteFile(savedConfigFile, data, 0)
	if err != nil {
		return savedConfigFile, fmt.Errorf("problem writing configuration to file")
	}

	return savedConfigFile, err
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
		return fmt.Sprintf("Problem rendering variable template file: %s", err.Error()), err
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
		return fmt.Sprintf("Problem writing out variable file: %s", err.Error()), err
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

	// If the working directory that has been set for the projects is relative, prepend the
	// the current directory to it
	if !filepath.IsAbs(config.Input.Directory.WorkingDir) {
		config.Input.Directory.WorkingDir = filepath.Join(cwd, config.Input.Directory.WorkingDir)
	}
}

// WriteCmdLog writes the command out a log file in the directory that the CLI is being run
// The cmd is only written out if the option to do so has been set in the config
func (config *Config) WriteCmdLog(path string, cmd string) error {

	var err error

	// return empty error if not logging commands
	if !config.UseCmdLog() {
		return err
	}

	// get a reference to the file, either to create or append to the file
	f, err := os.OpenFile(config.Self.CmdLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// write out the cmd to the file
	if _, err := f.WriteString(fmt.Sprintf("[%s] %s\n", path, cmd)); err != nil {
		return err
	}

	return err
}
