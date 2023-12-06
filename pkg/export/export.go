package export

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/Ensono/stacks-cli/pkg/config"
	"github.com/sirupsen/logrus"
)

type Export struct {
	Config *config.Config
	Logger *logrus.Logger
}

// New allocates a new ExportPointer with the given config.
func New(conf *config.Config, logger *logrus.Logger) *Export {
	return &Export{
		Config: conf,
		Logger: logger,
	}
}

// Run performs the task of exporting all of the internal static
// configuration files to the specified directory
func (e *Export) Run() error {
	var err error

	// determine if the path has been set, is relative or absolute
	// if the specified path is relative prepend the cwd to the path
	if !filepath.IsAbs(e.Config.Input.Directory.Export) {
		e.Config.Input.Directory.Export = filepath.Join(e.Config.Input.Directory.WorkingDir, e.Config.Input.Directory.Export)
	}

	// Check that the path exists
	err = util.CreateIfNotExists(e.Config.Input.Directory.Export, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create export directory: %s", e.Config.Input.Directory.Export)
	}

	e.Logger.Info("Exporting static configuration")

	if e.Config.IsDryRun() {
		e.Logger.Warn("Running export in DryRun mode, files will not be written to disk")
	}

	// Create a slice of all the static files to export and then write them out
	statics := make([]string, 0)
	statics = append(statics, "config")
	statics = append(statics, "azdo")
	// statics = append(statics, "help_urls")

	// iterate around the slice and export each static file
	for _, name := range statics {

		// determine the filename of the exported file
		filename := filepath.Join(e.Config.Input.Directory.Export, e.Config.Internal.GetFilename(name))

		e.Logger.Infof("\t%s - %s", name, filename)

		if !e.Config.Input.Options.DryRun {
			err = os.WriteFile(filename, e.Config.Internal.GetFileContent(name), os.ModePerm)

			if err != nil {
				e.Logger.Error(err.Error())
			}
		}
	}

	return err
}
