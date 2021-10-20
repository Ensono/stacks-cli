package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/amido/stacks-cli/internal/config/static"
	"github.com/amido/stacks-cli/internal/helper"
	"github.com/amido/stacks-cli/internal/util"
	"github.com/amido/stacks-cli/pkg/config"
	"github.com/sirupsen/logrus"
)

// type Foo interface {
// 	Write(string, string) (config.Config, error)
// }
type Scaffold struct {
	Name   string // name of the processing template
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

// Runs the scaffolding process
func (s *Scaffold) Run() error {
	if err := s.run(); err != nil {
		helper.TraceError(err)
		return err
	}
	return nil
}

// run iterates around all the projects that have been specified and sets up the
// working directory for each of them
func (s *Scaffold) run() error {

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	s.Logger.Tracef("Current Dir: %s\n", pwd)

	// Iterate around the projects that have been configured
	for _, project := range s.Config.Input.Project {

		// Determine the project path
		s.Config.Self.AddPath(project, s.setProjectPath(project.Name))
		s.Logger.Infof("Project path: %s\n", s.Config.Self.GetPath(project))
		s.Logger.Debugf("Project ID: %s", project.GetId())

		// create the directory
		err := os.MkdirAll(s.Config.Self.GetPath(project), 0755)
		if err != nil {
			break
		}
	}

	return err
}

func (s *Scaffold) setProjectPath(name string) string {
	project_path := filepath.Join(s.Config.Input.Directory.WorkingDir, name)

	return project_path
}

// PerformOperation performs the operation as specified by the settings file for the project
// It is responsible for performing any template replacements using GoTemplate
//
// The method reads in the Action and determines what is requied
// The currently supported actions are
//		copy - copies data from the temporary dir to the working dir
//		cmd - run a command on the local machine
//			The command is set using the `command` parameter
func (s *Scaffold) PerformOperation(operation config.Operation, cfg *config.Config, project *config.Project, path string) error {

	var command string

	switch operation.Action {
	case "cmd":

		// define a replacements object so that all can be passed to the render function
		// the project is passed in as a seperate object as it is part of a slice
		replacements := config.Replacements{}
		replacements.Input = cfg.Input
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
		args, err := cfg.RenderTemplate(operation.Arguments, replacements)
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
		err = cfg.WriteCmdLog(path, fmt.Sprintf("%s %s", command, args))
		if err != nil {
			s.Logger.Warnf("Unable to write command to log: %s", err.Error())
		}

		// set the command that needs to be executed
		cmd := exec.Command(command, args)
		// cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = path

		// only run the command if not in dryrun mode
		if !cfg.Input.Options.DryRun {
			if err = cmd.Run(); err != nil {
				s.Logger.Errorf("Error running command: %s", err.Error())
				return err
			}
		}
	}

	return nil
}

// create replace map

// func walkDirFunc(path string, d fs.DirEntry, err error) error {
// 	fileListArr
// 	return nil
// }

// cleans up
func cleanUpNewDirOnError(newDir string) {
	os.RemoveAll(newDir)
	helper.ShowInfo("Removed would be New Directory")
}
