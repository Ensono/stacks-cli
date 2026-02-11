package downloaders

import (
	"os"
	"path/filepath"

	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/sirupsen/logrus"
)

var git = &Git{}

type Git struct {
	URL              string
	Version          string
	FrameworkVersion string
	TempDir          string
	Token            string

	logger     *logrus.Logger
	Filesystem billy.Filesystem
}

func NewGitDownloader(url string, version string, frameworkVersion string, tempDir string, token string) *Git {
	git.URL = url
	git.Version = version
	git.FrameworkVersion = frameworkVersion
	git.TempDir = tempDir
	git.Token = token

	return git
}

func (g *Git) fs() billy.Filesystem {
	if g.Filesystem != nil {
		return g.Filesystem
	}

	return osfs.New("/")
}

func (g *Git) Get() (string, error) {

	// declare method variables
	var dir string
	var err error
	tempDir := g.TempDir
	if tempDir != "" {
		var resolveErr error
		tempDir, resolveErr = g.resolveTempDir()
		if resolveErr != nil {
			return "", resolveErr
		}
		if err := util.RemoveAll(g.fs(), tempDir); err != nil {
			return "", err
		}
		if err := g.fs().MkdirAll(tempDir, os.ModePerm); err != nil {
			return "", err
		}
	}

	// call the GitClone method
	dir, err = util.GitClone(
		g.URL,
		g.FrameworkVersion,
		g.Version,
		tempDir,
		g.Token,
	)

	return dir, err

}

func (g *Git) PackageURL() string {
	return g.URL
}

func (g *Git) SetLogger(logger *logrus.Logger) {
	g.logger = logger
}

func (g *Git) resolveTempDir() (string, error) {
	if g.TempDir == "" || filepath.IsAbs(g.TempDir) {
		return g.TempDir, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, g.TempDir), nil
}
