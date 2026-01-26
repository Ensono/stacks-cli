package downloaders

import (
	"os"
	"path/filepath"
	"strings"

	cp "github.com/otiai10/copy"

	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/sirupsen/logrus"
)

var filesystem = &Filesystem{}

type Filesystem struct {
	Path    string
	TempDir string

	// define private properties
	logger     *logrus.Logger
	Filesystem billy.Filesystem
}

func NewFilesystemDownloader(path string, tempDir string) *Filesystem {
	filesystem.Path = path
	filesystem.TempDir = tempDir

	return filesystem
}

func (f *Filesystem) fs() billy.Filesystem {
	if f.Filesystem != nil {
		return f.Filesystem
	}

	return osfs.New("/")
}

func (f *Filesystem) Get() (string, error) {

	if f.logger != nil {
		f.logger.Infof("Copying files from: %s", f.Path)
	}

	absoluteTempDir := f.TempDir
	if absoluteTempDir != "" {
		var err error
		absoluteTempDir, err = f.resolvePath(f.TempDir)
		if err != nil {
			return "", err
		}
		if err := util.RemoveAll(f.fs(), absoluteTempDir); err != nil {
			return "", err
		}
		if err := f.fs().MkdirAll(absoluteTempDir, os.ModePerm); err != nil {
			return "", err
		}
	}

	// copy the repository from the cloned directory to the project working directory
	// do not copy the git configuration folder
	opt := cp.Options{
		Skip: func(info os.FileInfo, src, dest string) (bool, error) {
			return strings.HasSuffix(src, ".git"), nil
		},
	}
	err := cp.Copy(f.Path, absoluteTempDir, opt)
	if err != nil {
		if f.logger != nil {
			f.logger.Errorf("Issue copying files: %s", err.Error())
		}
	}

	return f.TempDir, err
}

func (f *Filesystem) SetLogger(logger *logrus.Logger) {
	f.logger = logger
}

func (f *Filesystem) PackageURL() string {
	return f.Path
}

func (f *Filesystem) resolvePath(dir string) (string, error) {
	if dir == "" || filepath.IsAbs(dir) {
		return dir, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, dir), nil
}
