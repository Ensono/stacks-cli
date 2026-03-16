package downloaders

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/Ensono/stacks-cli/internal/models"
	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/sirupsen/logrus"
)

type Zip struct {
	URL      string
	Version  string
	TempDir  string
	CacheDir string
	Token    string

	// define private properties
	logger     *logrus.Logger
	Filesystem billy.Filesystem
}

func NewZipDownloader(zipURL string, version string, cacheDir string, tempDir string, token string) *Zip {
	return &Zip{
		URL:      zipURL,
		Version:  version,
		TempDir:  tempDir,
		CacheDir: cacheDir,
		Token:    token,
	}
}

func (z *Zip) fs() billy.Filesystem {
	if z.Filesystem != nil {
		return z.Filesystem
	}

	return osfs.New("/")
}

func (z *Zip) Get() (string, error) {

	if z.logger != nil {
		z.logger.Infof("Downloading zip file from: %s", z.URL)
	}

	// resolve the temp directory to an absolute path
	absoluteTempDir := z.TempDir
	if absoluteTempDir != "" {
		var err error
		absoluteTempDir, err = z.resolvePath(z.TempDir)
		if err != nil {
			return "", err
		}
		if err := util.RemoveAll(z.fs(), absoluteTempDir); err != nil {
			return "", err
		}
		if err := z.fs().MkdirAll(absoluteTempDir, os.ModePerm); err != nil {
			return "", err
		}
	}

	// ensure cache directory exists
	if z.CacheDir != "" {
		if err := z.fs().MkdirAll(z.CacheDir, os.ModePerm); err != nil {
			return "", err
		}
	}

	// determine the filename from the URL
	u, err := url.Parse(z.URL)
	if err != nil {
		return "", fmt.Errorf("unable to parse zip URL: %s", err.Error())
	}
	_, filename := path.Split(u.Path)
	if filename == "" {
		filename = "download.zip"
	}

	// determine the download path in the cache directory
	downloadPath := filepath.Join(z.CacheDir, filename)

	// download the zip file if it is not already cached
	ac := models.NewAPICall(z.URL, z.Token)
	if util.Exists(downloadPath) {
		if z.logger != nil {
			z.logger.Infof("Using zip file from local cache: %s", downloadPath)
		}
	} else {
		if z.logger != nil {
			z.logger.Info("Downloading zip file")
		}
		err = ac.Download(downloadPath)
		if err != nil {
			return "", fmt.Errorf("problem downloading zip file: %s", err.Error())
		}
	}

	// extract the zip file into the temp directory
	// util.Unzip already handles zip-slip protection
	dir, err := util.Unzip(downloadPath, absoluteTempDir)
	if err != nil {
		return "", fmt.Errorf("problem extracting zip file: %s", err.Error())
	}

	return dir, nil
}

func (z *Zip) PackageURL() string {
	return z.URL
}

func (z *Zip) SetLogger(logger *logrus.Logger) {
	z.logger = logger
}

func (z *Zip) resolvePath(dir string) (string, error) {
	if dir == "" || filepath.IsAbs(dir) {
		return dir, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, dir), nil
}
