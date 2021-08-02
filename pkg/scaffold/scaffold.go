package scaffold

import (
	"fmt"
	"os"
	"path"

	"github.com/amido/stacks-cli/internal/helper"
	"github.com/amido/stacks-cli/internal/util"
	"github.com/amido/stacks-cli/pkg/config"
)

// type Foo interface {
// 	Write(string, string) (config.Config, error)
// }

type Scaffold struct {
	Name   string // name of the processing template
	Config *config.Config
}

// New allocates a new ScaffoldPointer with the given config.
func New(conf *config.Config) *Scaffold {
	return &Scaffold{
		Config: conf,
	}
}

// Runs the scaffolding process
func (s *Scaffold) Run() error {
	if err := s.run(); err != nil {
		helper.ShowError(err)
		return err
	}
	return nil
}

// 1. determine action path based on input either API or config \n
// 2. get base source \n
// 3. generate replaceMap \n
// 4. replace placeholders in given files
// 5. copy to final output place
func (s *Scaffold) run() error {

	tmpPath, err := os.MkdirTemp(s.Config.Output.TmpPath, "source")
	if err != nil {
		return err
	}
	// invoke all helper functions from here so defer will be closed automatically on block exit
	defer os.RemoveAll(tmpPath)

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	helper.ShowInfo(fmt.Sprintf("Current Dir: %s\n", pwd))

	s.Config.Output.TmpPath = tmpPath
	s.Config.Output.ZipPath = fmt.Sprintf("%s/source.zip", tmpPath)
	s.Config.Output.UnzipPath = path.Join(tmpPath, "wip", s.Config.Self.Specific.Localpath)

	helper.ShowInfo(fmt.Sprintf("New Project Dir: %s\n", s.Config.Output.NewPath))

	if err := util.GitClone(s.Config.Self.Specific.Gitrepo, s.Config.Self.Specific.Gitref, s.Config.Output.TmpPath, s.Config.Output.ZipPath); err != nil {
		helper.ShowError(err)
		// cleanUpNewDirOnError(s.Config.Output.NewPath)
		return err
	}

	helper.ShowInfo(fmt.Sprintf("Cloned path %s\n\n", s.Config.Output.TmpPath))

	strs, e3 := s.sortFileOperations()
	if e3 != nil {
		helper.ShowError(err)
		cleanUpNewDirOnError(s.Config.Output.NewPath)
		return err
	}

	helper.ShowInfo(fmt.Sprintf("%s", strs))

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

// create replace map

// func walkDirFunc(path string, d fs.DirEntry, err error) error {
// 	fileListArr
// 	return nil
// }

// cleans up
func cleanUpNewDirOnError(newDir string) {
	os.RemoveAll(newDir)
	helper.ShowWarning("Removed would be New Directory")
}
