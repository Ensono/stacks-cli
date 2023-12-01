package setup

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/amido/stacks-cli/internal/util"
	"github.com/amido/stacks-cli/pkg/config"

	log "github.com/sirupsen/logrus"
)

func setupSetupTestCase(t *testing.T) (func(t *testing.T), string) {
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

func TestCreateProjectSettings(t *testing.T) {

	// create temporary directory to work with
	cleanup, tempDir := setupSetupTestCase(t)
	defer cleanup(t)

	// create a list of tests that need to be performed
	tables := []struct {
		config  config.Config
		pattern string
		msg     string
	}{
		{
			config.Config{
				Input: config.InputConfig{
					Directory: config.Directory{
						WorkingDir: tempDir,
					},
					Business: config.Business{
						Company: "ensono",
					},
				},
			},
			`company:\s+%s`,
			"File content must match: %s",
		},
		{
			config.Config{
				Input: config.InputConfig{
					Global: true,
					Directory: config.Directory{
						HomeDir: tempDir,
					},
					Business: config.Business{
						Company:   "ensono",
						Component: "stacks",
					},
				},
			},
			`company:\s+%s`,
			"Global file content must match: %s",
		},
	}

	logger := log.New()
	for _, table := range tables {
		var dir string
		var buf bytes.Buffer
		setup := New(&table.config, logger)

		// if the configuration global, check that a warning has been raised
		if table.config.Input.Global {
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()
		}

		// call the Upsert function
		err := setup.Upsert()

		// ensure there are no errors
		if err != nil {
			t.Error(err)
		}

		if table.config.Input.Global {
			dir = table.config.Input.Directory.HomeDir

			// check that the output contains a message to say that it is not recommended to set values
			if strings.Contains(buf.String(), "not recommended to set") {
				t.Error("Setting properties in global file that are not recommended")
			}

		} else {
			dir = table.config.Input.Directory.WorkingDir
		}
		path := path.Join(dir, ".stackscli", "config.yml")

		// check that the a .stackscli/config.yml file exists in the tempdir
		if !util.Exists(path) {
			t.Errorf("Configuration file does not exist: %s", path)
		} else {
			// Ensure the contents of the file is correct
			// read the contents of the file
			content, err := os.ReadFile(path)
			if err != nil {
				t.Error(err)
			}

			// compare the content with the pattern
			pattern := fmt.Sprintf(table.pattern, table.config.Input.Business.Company)
			re := regexp.MustCompile(pattern)
			matched := re.MatchString(string(content))

			if !matched {
				t.Errorf(table.msg, pattern)
			}
		}
	}

}
