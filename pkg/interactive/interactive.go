package interactive

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/amido/stacks-cli/internal/util"
	"github.com/amido/stacks-cli/pkg/config"
	yaml "github.com/goccy/go-yaml"
	"github.com/sirupsen/logrus"
)

type Interactive struct {
	Config *config.Config
	Logger *logrus.Logger
}

func New(conf *config.Config, logger *logrus.Logger) *Interactive {
	return &Interactive{
		Config: conf,
		Logger: logger,
	}
}

// Run starts the interactive questions and saves the file to the
// working directory that has been specified
func (i *Interactive) Run() error {

	var err error

	// create an answers config object so that the questions can be asked
	answers := config.Answers{}
	err = answers.RunInteractive(i.Config)

	if err != nil {
		return err
	}

	// determine the full path to the configuration file
	path := i.getPath()

	i.Logger.Infof("Saving configuration to file: %s", path)

	// marshal the InputObject and write to the specified path
	//
	// The github.com/goccy libray is being used here instead because it allows items
	// to be omitted when they are empty, muchlike the built in JSON parser
	// The built in YAML parser does not support this
	data, err := yaml.Marshal(&i.Config.Input)

	if err != nil {
		return err
	}

	// write the data out to the file
	err = ioutil.WriteFile(path, data, 0666)

	// output information about what to run next
	helpText := fmt.Sprintf(`To scaffold the new projects, run the following command
	
stacks-cli scaffold -c %s`, path)
	i.Logger.Info(helpText)

	return err
}

// getPath builds the path to where the configuration file should be saved
func (i *Interactive) getPath() string {
	path := filepath.Join(i.Config.Input.Directory.WorkingDir, "stacks.yml")

	// check to see if in Unix shell and ensure using forward `/` slash if it is
	// This is a edge-case where Bash could be running on Windows which requires that the
	// path delimiters need to be set to `/`
	if util.IsUnixShell() {
		path = filepath.ToSlash(path)
	}

	return path
}
