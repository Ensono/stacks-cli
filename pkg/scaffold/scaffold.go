package scaffold

import (
	"fmt"
	"os"
	"path/filepath"

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

// 1. determine action path based on input either API or config \n
// 2. get base source \n
// TODO: still
// 3. generate replaceMap \n
// 4. replace placeholders in given files
// 5. copy to final output place
func (s *Scaffold) run() error {
	// invoke all helper functions from here so defer will be closed automatically on block exit
	// defer os.RemoveAll(s.Config.Input.Directory.TempDir)

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	s.Logger.Tracef("Current Dir: %s\n", pwd)

	// Iterate around the projects that have been configured
	for _, project := range s.Config.Input.Project {

		// Determine the project path
		project_path := s.setProjectPath(project.Name)
		s.Logger.Infof("Project path: %s\n", project_path)

		// create the directory
		err := os.MkdirAll(project_path, 0755)
		if err != nil {
			break
		}
	}

	return err

	s.Logger.Tracef("New Project Dir: %s\n", s.Config.Input.Directory.WorkingDir)

	if err := util.GitClone(s.Config.Self.Specific.Gitrepo, s.Config.Self.Specific.Gitref, s.Config.Output.TmpPath, s.Config.Output.ZipPath); err != nil {
		s.Logger.Trace(err.Error())
		// cleanUpNewDirOnError(s.Config.Output.NewPath)
		return err
	}

	// Add additional config values from Repos

	s.Logger.Tracef("Cloned path %s\n\n", s.Config.Output.TmpPath)

	strs, e3 := s.sortFileOperations()
	if e3 != nil {
		s.Logger.Trace(err.Error())
		cleanUpNewDirOnError(s.Config.Output.NewPath)
		return err
	}

	helper.TraceInfo(fmt.Sprintf("%s", strs))

	return nil
}

func (s *Scaffold) sortFileOperations() ([]string, error) {

	fileListArr, err := util.UnzipClone(s.Config.Output.ZipPath, s.Config.Output.UnzipPath)
	if err != nil {
		return nil, err
	}

	// create a map of replacements on each file

	return fileListArr, nil
}

func (s *Scaffold) setProjectPath(name string) string {
	project_path := filepath.Join(s.Config.Input.Directory.WorkingDir, name)

	return project_path
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
