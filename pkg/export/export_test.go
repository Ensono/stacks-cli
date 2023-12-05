package export

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/Ensono/stacks-cli/pkg/config"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func setupExportTestCase(t *testing.T) (func(t *testing.T), string) {

	// create a temporary directory
	tempDir := t.TempDir()

	deferFunc := func(t *testing.T) {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Logf("[ERROR] Unable to remove dir: %v", err)
		}
	}

	return deferFunc, tempDir
}

func TestExportedFiles(t *testing.T) {

	cleanup, tempDir := setupExportTestCase(t)
	defer cleanup(t)

	// create test tables
	tables := []struct {
		dir    string
		dryrun bool
		files  []string
	}{
		{
			filepath.Join(tempDir, "exported"),
			false,
			[]string{
				"azdo_variable_template.yml",
				"internal_config.yml",
			},
		},
	}

	// iterate around the test tables and ensure the files are written out to the correct folder
	for _, table := range tables {

		// create the configuration object
		cfg := config.Config{
			Input: config.InputConfig{
				Directory: config.Directory{
					Export: table.dir,
				},
			},
		}
		cfg.Internal.AddFiles()

		// create the necessary objects
		logger := log.New()
		export := New(&cfg, logger)

		// run the export
		export.Run()

		// check that the expected directory exists
		if !util.Exists(table.dir) {
			t.Errorf("Expected directory does not exist,: %s", table.dir)
		}

		// get the files in the exported directory and compare against the specified list
		files, _ := os.ReadDir(table.dir)

		filenames := []string{}
		for _, info := range files {
			filenames = append(filenames, info.Name())
		}

		if len(files) != len(table.files) {
			t.Errorf("Expected number of files not found, %d: %d", len(table.files), len(files))
		}

		// perform a deep comparison of the files, to ensure what is retrieved is what is expected
		assert.Equal(t, filenames, table.files)

	}
}
