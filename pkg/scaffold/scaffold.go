package scaffold

import (
	"fmt"
	"os"

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
		helper.TraceError(err)
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
	// invoke all helper functions from here so defer will be closed automatically on block exit
	defer os.RemoveAll(s.Config.Output.TmpPath)

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	helper.TraceInfo(fmt.Sprintf("Current Dir: %s\n", pwd))

	helper.TraceInfo(fmt.Sprintf("New Project Dir: %s\n", s.Config.Output.NewPath))

	if err := util.GitClone(s.Config.Self.Specific.Gitrepo, s.Config.Self.Specific.Gitref, s.Config.Output.TmpPath, s.Config.Output.ZipPath); err != nil {
		helper.TraceError(err)
		// cleanUpNewDirOnError(s.Config.Output.NewPath)
		return err
	}

	helper.TraceInfo(fmt.Sprintf("Cloned path %s\n\n", s.Config.Output.TmpPath))

	strs, e3 := s.sortFileOperations()
	if e3 != nil {
		helper.TraceError(err)
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
