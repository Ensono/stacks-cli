package setup

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Ensono/stacks-cli/internal/constants"
	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/Ensono/stacks-cli/pkg/config"
	"github.com/Ensono/stacks-cli/pkg/filter"
	"github.com/sirupsen/logrus"
)

type Setup struct {
	Config *config.Config
	Logger *logrus.Logger
}

func New(conf *config.Config, logger *logrus.Logger) *Setup {
	return &Setup{
		Config: conf,
		Logger: logger,
	}
}

func (s *Setup) Upsert() error {

	var err error
	var path string
	var perm uint32
	fs := s.Config.GetFilesystem()

	// configure variable to hold the path to the file after the basedir has been determined
	var slug []string = []string{constants.ConfigFileDir, fmt.Sprintf("%s.yml", constants.ConfigName)}

	// create a slice of dotted notation to see which values have been set that are not advised to be set globally
	dotted := []string{
		"business.project",
		"business.component",
		"terraform.backend.group",
		"terraform.backend.storage",
		"terraform.backend.container",
	}

	// determine the path to the configuration file
	// if global has been set then this will be the current users home directory.
	// otherwise it will be the current directory
	if s.Config.Input.Global {

		// create a slice to hold the values that have been set
		setvalues := []string{}

		// iterate around the dotted notation and check if the value has been set
		for _, path := range dotted {
			pathElements := strings.Split(path, ".")
			val, err := util.GetValueByDottedPath(s.Config.Input, path)
			if err != nil {
				s.Logger.Errorf("Unable to get value of field: %s", err.Error())
				continue
			}

			if val.(string) != "" {
				setvalues = append(setvalues, pathElements[len(pathElements)-1])
			}
		}

		// if the setvalues is not empty, display a warning
		if len(setvalues) > 0 {
			s.Logger.Warnf("It is not recommended to set the following values globally: %s", strings.Join(setvalues, ", "))
		}

		slug = append([]string{s.Config.Input.Directory.HomeDir}, slug...)
		perm = 0700
	} else {
		slug = append([]string{s.Config.Input.Directory.WorkingDir}, slug...)
		perm = 0755
	}

	// build up the full path to the file
	path = filepath.Join(slug...)

	s.Logger.Infof("Updating configuration file: %s", path)

	// Filter the configuration object and write out to the file
	filter := filter.New()
	filter.Filter(s.Config.Input, append(dotted, "business.company"))
	err = filter.WriteFile(fs, path, perm)

	return err

}

func (s *Setup) GetLatestInternalConfig() error {
	var err error

	// download the file from the specified endpoint to the home directory of the user
	s.Logger.Infof("Downloading latest configuration")

	// Get the data from the URL
	s.Logger.Debugf("Downloading file from: %s", s.Config.Input.Overrides.InternalConfigURL)
	resp, err := http.Get(s.Config.Input.Overrides.InternalConfigURL)
	if err != nil {
		s.Logger.Errorf("Error downloading file: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	// create the file
	// define the path to the file
	filepath := path.Join(util.GetStacksCLIDir(), "internal_config.yml")
	s.Logger.Debugf("Creating file: %s", filepath)
	out, err := os.Create(filepath)
	if err != nil {
		s.Logger.Errorf("Error creating file: %s", err.Error())
		return err
	}
	defer out.Close()

	// write out the body of the file
	_, err = io.Copy(out, resp.Body)

	return err
}

func (s *Setup) List() error {

	var err error

	return err
}
