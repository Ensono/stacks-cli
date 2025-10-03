package downloaders

import (
	"os"
	"strings"

	cp "github.com/otiai10/copy"

	"github.com/sirupsen/logrus"
)

var filesystem = &Filesystem{}

type Filesystem struct {
	Path    string
	TempDir string

	// define private properties
	logger *logrus.Logger
}

func NewFilesystemDownloader(path string, tempDir string) *Filesystem {
	filesystem.Path = path
	filesystem.TempDir = tempDir

	return filesystem
}

func (f *Filesystem) Get() (string, error) {

	if f.logger != nil {
		f.logger.Infof("Copying files from: %s", f.Path)
	}

	// copy the repository from the cloned directory to the project working directory
	// do not copy the git configuration folder
	opt := cp.Options{
		Skip: func(info os.FileInfo, src, dest string) (bool, error) {
			return strings.HasSuffix(src, ".git"), nil
		},
	}
	err := cp.Copy(f.Path, f.TempDir, opt)
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
