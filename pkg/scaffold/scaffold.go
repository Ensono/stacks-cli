package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/amido/stacks-cli/internal/config/static"
	"github.com/amido/stacks-cli/internal/models"
	"github.com/amido/stacks-cli/internal/util"
	"github.com/amido/stacks-cli/pkg/config"
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
	missing := s.Config.Input.CheckFrameworks()
	errText := s.analyseMissing(missing)
	if errText != "" {
		s.Logger.Fatal(errText)
	}

	// iterate around the projects that have been specified and
	// process each one in turn
	for _, project := range s.Config.Input.Project {
		s.processProject(project)
	}

	return err

}

// run iterates around all the projects that have been specified and sets up the
// working directory for each of them
// func (s *Scaffold) run() error {

// 	pwd, err := os.Getwd()
// 	if err != nil {
// 		return err
// 	}

// 	s.Logger.Tracef("Current Dir: %s\n", pwd)

// 	// Iterate around the projects that have been configured
// 	for _, project := range s.Config.Input.Project {

// 		// Determine the project path
// 		s.Config.Self.AddPath(project, s.setProjectPath(project.Name))
// 		s.Logger.Infof("Project path: %s\n", s.Config.Self.GetPath(project))
// 		s.Logger.Debugf("Project ID: %s", project.GetId())

// 		// create the directory
// 		err := os.MkdirAll(s.Config.Self.GetPath(project), 0755)
// 		if err != nil {
// 			break
// 		}
// 	}

// 	return err
// }

// PerformOperation performs the operation as specified by the settings file for the project
// It is responsible for performing any template replacements using GoTemplate
//
// The method reads in the Action and determines what is requied
// The currently supported actions are
//		copy - copies data from the temporary dir to the working dir
//		cmd - run a command on the local machine
//			The command is set using the `command` parameter
func (s *Scaffold) PerformOperation(operation config.Operation, project *config.Project, path string) error {

	var command string

	switch operation.Action {
	case "cmd":

		// define a replacements object so that all can be passed to the render function
		// the project is passed in as a seperate object as it is part of a slice
		replacements := config.Replacements{}
		replacements.Input = s.Config.Input
		replacements.Project = *project

		// get a list of the commands that are expected to be run by something that
		// uses the framework
		cmdList := static.FrameworkCommand(project.Framework.Type)

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

		// run the args that have been specified through the template engine
		args, err := s.Config.RenderTemplate(operation.Arguments, replacements)
		if err != nil {
			s.Logger.Errorf("Error resolving template: %s", err.Error())
			return err
		}

		// set the command to be run if the platform is windows
		if runtime.GOOS == "windows" {
			args = fmt.Sprintf("/C %s %s", command, args)
			command = "cmd"
		}

		// output the command being run if in debug mode
		s.Logger.Debugf("Command: %s %s", command, args)

		// Write out the command log
		err = s.Config.WriteCmdLog(path, fmt.Sprintf("%s %s", command, args))
		if err != nil {
			s.Logger.Warnf("Unable to write command to log: %s", err.Error())
		}

		// set the command that needs to be executed
		cmd := exec.Command(command, args)
		// cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = path

		// only run the command if not in dryrun mode
		if !s.Config.IsDryRun() {
			if err = cmd.Run(); err != nil {
				s.Logger.Errorf("Error running command: %s", err.Error())
				return err
			}
		}
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
			list += fmt.Sprintf("Framework '%s' may have been misspelled because the command for this framework cannot be determined", item.Framework)
		} else {
			list += fmt.Sprintf("Command '%s' for the '%s' framework cannot be located. Is '%s' installed and in your PATH?", item.Binary, item.Framework, item.Binary)
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

	// Get the specified framework project based on the information supplied in
	// the configuration
	key := project.Framework.GetMapKey()
	s.Logger.Infof("Retrieving framework option: %s", key)
	dir, err := util.GitClone(
		s.Config.Input.Stacks.GetSrcURL(key),
		project.SourceControl.Ref,
		s.Config.Input.Directory.TempDir,
	)

	// if there was an error getting hold of the framework project display an error
	// and move onto the next project
	if err != nil {
		s.Logger.Errorf("Error downloading the specific framework option: %s", err.Error())
		return
	}

	// attempt to read in the settings for the framework option
	err = project.ReadSettings(dir, s.Config)
	s.Logger.Infof("Attempting to read project settings: %s", project.SettingsFile)
	if err != nil {
		s.Logger.Errorf("Error reading settings from project settings: %s", err.Error())
		return
	}

	// iterate around the phases of the project and the operations contained therin
	for _, phase := range project.Phases {
		for _, op := range phase.Operations {

			// output information about the operation being performed
			s.Logger.Info(op.Description)

			// perform the operation
			err = s.PerformOperation(op, &project, phase.Directory)

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

	// cleanup
	s.cleanup()

}

// setProjectDirs sets the temporary and working directories, based on the name of the
// project.
// It will attempt to create the directories, after defining them, but will error if
// the project directory already exists
func (s *Scaffold) setProjectDirs(project *config.Project) error {
	var err error

	project.Directory.WorkingDir = filepath.Join(s.Config.Input.Directory.WorkingDir, project.Name)
	project.Directory.TempDir = filepath.Join(s.Config.Input.Directory.TempDir, project.Name)

	// check to see if the workingdir already exists, if it does return an error
	// otherwise create them
	if util.Exists(project.Directory.WorkingDir) {

		// if Clobber is turned on, remove the directory with a warning
		if s.Config.Clobber() {
			s.Logger.Warnf("Removing existing project directory: %s", project.Directory.WorkingDir)

			err = os.RemoveAll(project.Directory.WorkingDir)
			return err
		} else {
			s.Logger.Warnf("Project directory already exists, skipping: %s", project.Directory.WorkingDir)
			return nil
		}
	}

	err = util.CreateIfNotExists(project.Directory.WorkingDir, 0755)

	return err
}

// configurePipeline is responsible for setting up the build pipeline and variables file
func (s *Scaffold) configurePipeline(project *config.Project) {

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
			s.Logger.Info("Created pipeline variable file")
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
		err := s.PerformOperation(op, project, project.Directory.WorkingDir)

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
