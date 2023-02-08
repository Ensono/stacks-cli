package scaffold

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/amido/stacks-cli/internal/config/static"
	"github.com/amido/stacks-cli/internal/models"
	"github.com/amido/stacks-cli/internal/util"
	"github.com/amido/stacks-cli/pkg/config"
	"github.com/amido/stacks-cli/pkg/downloaders"
	"github.com/amido/stacks-cli/pkg/interfaces"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Scaffold struct {
	Config *config.Config
	Logger *logrus.Logger
}

// New allocates a new ScaffoldPointer with the given config.
func New(conf *config.Config, logger *logrus.Logger) *Scaffold {
	return &Scaffold{
		Config: conf,
		Logger: logger,
	}
}

// Run performs the operations of the scaffolding sub command
// It will iterates around each of the projects that have been specified
// and performs all of the operations and that need to be done
func (s *Scaffold) Run() error {

	var err error

	// check the runtime configuration and set necessary defaults
	err = s.Config.Check()
	if err != nil {
		s.Logger.Fatalln(err.Error())
	}

	// determine if the configuration needs to be saved to a file
	savedConfigFile, err := s.Config.Save(viper.ConfigFileUsed())
	if savedConfigFile != "" {
		s.Logger.Infof("Configuration saved to file: %s", savedConfigFile)
	}
	if err != nil {
		s.Logger.Warnf("Issue saving configuration: %s", err.Error())
	}

	// Analyse the projects and the frameworks that have been chosen
	missing := s.Config.Input.CheckFrameworks(s.Config)
	errText := s.analyseMissing(missing)
	if errText != "" {
		s.Logger.Fatal(errText)
	}

	// validate the inputs
	validations := s.Config.Input.ValidateInput()
	if len(validations) > 0 {
		s.Logger.Infof("Some inputs have been modified:\n\t%s", strings.Join(validations, "\n\t"))
	}

	// create the temporary directory if it does not exist
	err = util.CreateIfNotExists(s.Config.Input.Directory.TempDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create temporary directory: %s", s.Config.Input.Directory.TempDir)
	}

	// Create the cache directory if it does not exist
	err = util.CreateIfNotExists(s.Config.Input.Directory.CacheDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create cache directory: %s", s.Config.Input.Directory.CacheDir)
	}

	// Cleanup the temporary dir after all the projects have been processed
	defer s.cleanup()

	// iterate around the projects that have been specified and
	// process each one in turn
	for _, project := range s.Config.Input.Project {
		s.processProject(project)
	}

	return err

}

// PerformOperation performs the operation as specified by the settings file for the project
// It is responsible for performing any template replacements using GoTemplate
//
// The method reads in the Action and determines what is requied
// The currently supported actions are
//
//	copy - copies data from the temporary dir to the working dir
//	cmd - run a command on the local machine
//		The command is set using the `command` parameter
func (s *Scaffold) PerformOperation(operation config.Operation, project *config.Project, path string, cloneDir string) error {

	var command string

	switch operation.Action {
	case "cmd":

		// define a replacements object so that all can be passed to the render function
		// the project is passed in as a seperate object as it is part of a slice
		replacements := config.Replacements{}
		replacements.Input = s.Config.Input
		replacements.Project = *project

		// create a string builder
		arguments := strings.Builder{}

		// get a list of the commands that are expected to be run by something that
		// uses the framework
		cmdList := s.Config.GetFrameworkCommands(project.Framework.Type)

		// check the operation command to see if has been specified
		// and that it is listed in the cmdList
		if operation.Command == "" {
			return fmt.Errorf("command has not been set for the operation")
		} else {
			if !util.SliceContains(cmdList, operation.Command) {
				return fmt.Errorf("command '%s' is not is the known list of commands for '%s'", operation.Command, project.Framework.Type)
			}
		}
		command = operation.Command

		// run the arguments that have been specified through the template engine
		args, err := s.Config.RenderTemplate("arguments", operation.Arguments, replacements)
		if err != nil {
			s.Logger.Errorf("Error resolving template: %s", err.Error())
			return err
		}

		// expand any OS based variables on the template
		arguments.WriteString(os.ExpandEnv(args))

		// if properties have nee supplied and the operation has the flag set to apply the properties
		// to the command, add them in here
		if operation.ApplyProperties && len(project.Framework.Properties) > 0 {
			arguments.WriteString(" ")
			arguments.WriteString(strings.Join(project.Framework.Properties, " "))
		}

		// Execute the command and check that it worked
		_, err = s.Config.ExecuteCommand(path, s.Logger, command, arguments.String(), false, false)
		if err != nil {
			s.Logger.Errorf("Issue running command: %s", err.Error())
		}
	case "copy":

		// copy the repository from the cloned directory to the project working directory
		// do not copy the git configuration folder
		util.CopyDirectory(cloneDir, path)
	}

	return nil
}

// analyseMissing takes the list of missing commands and returns a string
// to be output to the console if there are any missing commands
func (s *Scaffold) analyseMissing(missing []models.Command) string {
	var list string

	// if there are no missing items, return
	if len(missing) == 0 {
		return ""
	}

	// iterate around the missing items
	for _, item := range missing {
		if item.Binary == "" {
			list += fmt.Sprintf("Framework '%s' may have been misspelled because the command for this framework cannot be determined\n", item.Framework)
		} else {
			list += fmt.Sprintf("Command '%s' for the '%s' framework cannot be located. Is '%s' installed and in your PATH?\n", item.Binary, item.Framework, item.Binary)
		}
	}

	// create the final message to return with
	message := fmt.Sprintf(`Some of the commands required by the specified frameworks do not exist on your machine or the framework has been specified incorrectly.

%s`, list)

	return message
}

// processProject configures the working directory according the project settings
func (s *Scaffold) processProject(project config.Project) {

	// output information about the project that is being setup
	s.Logger.Infof("Setting up project: %s", project.Name)

	// configure the directories for the project
	// if there is an error, e.g. the working directory already exists, display the
	// error and return to the calling function to move onto the next project
	err := s.setProjectDirs(&project)
	if err != nil {
		s.Logger.Error(err.Error())
		return
	}

	// Get the URL for the repository to download
	key := project.Framework.GetMapKey()
	repoInfo := s.Config.Input.Stacks.GetSrcURL(key)

	// if the URL is empty, emit error message and state why this might be the case
	if repoInfo == (config.RepoInfo{}) {
		s.Logger.Errorf(`The URL for the specified framework option, %s, is empty. Have you specified the correct framework option?`, project.Framework.Option)
		return
	}

	// ensure that the RepoInfo object is correctly configured
	msg := repoInfo.Normalize()
	if msg != "" {
		s.Logger.Warn(msg)
	}

	// download the template using the appropriate action
	var downloader interfaces.Downloader
	switch repoInfo.Type {
	case "github":

		// check that the URL is valid, if not skip this project and move onto the next one
		_, err = url.ParseRequestURI(repoInfo.Name)
		if err != nil {
			s.Logger.Errorf("Unable to download framework option as URL is invalid: %s", err.Error())
			return
		}

		downloader = downloaders.NewGitDownloader(
			repoInfo.Name,
			repoInfo.Version,
			project.Framework.Version,
			s.Config.Input.Directory.TempDir,
			s.Config.Input.Options.Token,
		)
	case "nuget":
		downloader = downloaders.NewNugetDownloader(
			repoInfo.Name,
			repoInfo.ID,
			repoInfo.Version,
			project.Framework.Version,
			s.Config.Input.Directory.CacheDir,
			s.Config.Input.Directory.TempDir,
		)
	}

	downloader.SetLogger(s.Logger)

	s.Logger.Infof("Retrieving framework option: %s", key)
	dir, err := downloader.Get()

	// set the clonedir as the project temporary directory
	project.Directory.TempDir = dir

	// if there was an error getting hold of the framework project display an error
	// and move onto the next project
	if err != nil {
		s.Logger.Errorf("Issue downloading the specified framework option\n\tURL: %s\n\tError: %s", downloader.PackageURL(), err.Error())
		return
	}

	// attempt to read in the settings for the framework option
	if project.SettingsFile != "" {
		s.Logger.Infof("Attempting to read project settings: %s", project.SettingsFile)
	}
	err = project.ReadSettings(dir, s.Config)
	if err != nil {
		s.Logger.Errorf("Error reading settings from project settings: %s", err.Error())
		s.Logger.Info("Please ensure you are running the latest version of the CLI")
		return
	}

	// check to see if any framework commands have been set and check the
	// version if they have
	incorrect := project.Settings.CheckCommandVersions(s.Config, s.Logger, project.Directory.WorkingDir, project.Directory.TempDir)
	if len(incorrect) > 0 {

		var parts []string

		for _, wrong := range incorrect {
			// iterate around the incorrect versions and create the body of the text to output
			parts = append(parts,
				fmt.Sprintf("\tVersion constraint for '%s' is '%s', but found '%s'", wrong.Binary, wrong.VersionRequired, wrong.VersionFound),
			)
		}

		s.Logger.Errorf("Unable to process project as framework versions are incorrect.\n\n%s", strings.Join(parts, "\n"))

		// there are issues with the versions of the commands, so move onto the next project
		// but only if Force is not set
		if s.Config.Force() {
			s.Logger.Warn("Continuing as the `force` option has been set. Your project may not configure properly with incorrect command versions")
		} else {
			return
		}
	}

	// iterate around the phases of the project and the operations contained therin
	for _, phase := range project.Phases {
		for _, op := range phase.Operations {

			// output information about the operation being performed
			s.Logger.Info(op.Description)

			// perform the operation
			err = s.PerformOperation(op, &project, phase.Directory, dir)

			if err != nil {
				s.Logger.Errorf("issue encountered performing '%s' operation: %s", phase.Name, err.Error())
				break
			}
		}
	}

	// configure the pipeline in the project
	s.configurePipeline(&project)

	// configure the git repository
	s.configureGitRepository(&project)

}

// setProjectDirs sets the temporary and working directories, based on the name of the
// project.
// It will attempt to create the directories, after defining them, but will error if
// the project directory already exists
func (s *Scaffold) setProjectDirs(project *config.Project) error {
	var err error

	project.Directory.WorkingDir = filepath.Join(s.Config.Input.Directory.WorkingDir, project.Name)

	// check to see if the workingdir already exists, if it does return an error
	// otherwise create them
	if util.Exists(project.Directory.WorkingDir) {

		// if Force is enabled, remove the directory with a warning
		if s.Config.Force() {
			s.Logger.Warnf("Removing existing project directory: %s", project.Directory.WorkingDir)

			err = os.RemoveAll(project.Directory.WorkingDir)
			return err
		} else {

			// determine if the dir is empty, if it is then allow overwriting
			empty, _ := util.IsEmpty(project.Directory.WorkingDir)
			if empty {
				s.Logger.Warnf("Overwriting empty directory: %s", project.Directory.WorkingDir)
			} else {
				return fmt.Errorf("project directory already exists, skipping: %s", project.Directory.WorkingDir)
			}
		}
	}

	err = util.CreateIfNotExists(project.Directory.WorkingDir, 0755)

	return err
}

// configurePipeline is responsible for setting up the build pipeline and variables file
func (s *Scaffold) configurePipeline(project *config.Project) {

	if len(project.Settings.Pipeline) == 0 {
		s.Logger.Info("No pipelines settings have been defined in the project for the CLI to configure")
		return
	}

	// get the pipeline settings
	pipelineSettings := project.Settings.GetPipeline(s.Config.Input.Pipeline)

	// define the replacements object so that all can be passed to the render function
	// the project is passed in a separate project as it is part of a slice
	replacements := config.Replacements{}
	replacements.Input = s.Config.Input
	replacements.Project = *project

	// attempt to write out the configuration file, unless in DryRun mode
	if s.Config.Input.Options.DryRun {
		s.Logger.Warn("Not creating variables template as in DRYRUN mode")
	} else {
		msg, err := s.Config.WriteVariablesFile(project, pipelineSettings, replacements)

		if err == nil {
			if msg == "" {
				s.Logger.Info("Created pipeline variable file")
			} else {
				s.Logger.Warn(msg)
			}
		} else {
			s.Logger.Error(msg)
		}
	}

	// perform any addition regex replacements
	errs := pipelineSettings.ReplacePatterns(project.Directory.WorkingDir)
	if len(errs) > 0 {
		for _, err := range errs {
			s.Logger.Error(err.Error())
		}
	}
}

// configureGitRepository configures the newly generated project as a git repository
// based on the settings that have been provided
func (s *Scaffold) configureGitRepository(project *config.Project) {

	s.Logger.Info("Configuring source control for the project")

	// check that the URL specific for the remote repo is a valid URL
	_, err := url.ParseRequestURI(project.SourceControl.URL)
	if err != nil {
		s.Logger.Errorf("Unable to configure remote repo: %s", err.Error())
		return
	}

	// iterate around the static git commands
	for _, command := range static.GitCmds {

		// split the command string into a cmd and args so that the operation model
		// can be configured and the PerformOperation method used
		commandParts := strings.SplitN(command, " ", 2)

		// build up the model
		op := config.Operation{
			Action:    "cmd",
			Command:   commandParts[0],
			Arguments: commandParts[1],
		}

		// call the PerformOperation function
		err := s.PerformOperation(op, project, project.Directory.WorkingDir, "")

		if err != nil {
			s.Logger.Error(err.Error())
			return
		}
	}
}

// cleanup is responsible for outputting completion messages and removing
// temporary directories
func (s *Scaffold) cleanup() {

	// state that nothing has been configured if DRYRUN has been enabled
	if s.Config.IsDryRun() {
		s.Logger.Warnf("CLI was run the with --dryrun option, no projects have been configured")
	}

	// if a command log has been requested then state the location of the log file
	if s.Config.UseCmdLog() {
		s.Logger.Infof("Command log has been created: %s", s.Config.Self.CmdLogPath)
	}

	// remove the temporary directory if permitted
	if s.Config.NoCleanup() {
		s.Logger.Warnf("Cleanup has been disabled, please perform the cleanup manually: %s", s.Config.Input.Directory.TempDir)
	} else {
		s.Logger.Info("Performing cleanup")
		s.Logger.Infof(" - removing temporary directory: %s", s.Config.Input.Directory.TempDir)

		err := os.RemoveAll(s.Config.Input.Directory.TempDir)
		if err != nil {
			s.Logger.Fatalf("Unable to remove temporary directory: %s", err.Error())
		}
	}
}
