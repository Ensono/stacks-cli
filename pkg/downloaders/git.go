package downloaders

import (
	"github.com/amido/stacks-cli/internal/util"
	"github.com/sirupsen/logrus"
)

var git = &Git{}

type Git struct {
	URL              string
	Version          string
	FrameworkVersion string
	TempDir          string
	Token            string

	logger *logrus.Logger
}

func NewGitDownloader(url string, version string, frameworkVersion string, tempDir string, token string) *Git {
	git.URL = url
	git.Version = version
	git.FrameworkVersion = frameworkVersion
	git.TempDir = tempDir
	git.Token = token

	return git
}

func (g *Git) Get() (string, error) {

	// declare method variables
	var dir string
	var err error

	// call the GitClone method
	dir, err = util.GitClone(
		g.URL,
		g.FrameworkVersion,
		g.Version,
		g.TempDir,
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
