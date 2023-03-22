package scaffold

import (
	"os"
	"regexp"
	"testing"

	"github.com/amido/stacks-cli/internal/models"
	"github.com/amido/stacks-cli/pkg/config"

	log "github.com/sirupsen/logrus"
)

func setupScaffoldTestCase(t *testing.T) (func(t *testing.T), string) {
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

func TestAnalyseMissing(t *testing.T) {

	// create test tables
	tables := []struct {
		missing []models.Command
		pattern string
		msg     string
	}{
		{
			[]models.Command{},
			"",
			"There should be no error message as there are no missing comamnds",
		},
		{
			[]models.Command{
				{
					Binary:    "",
					Framework: "dotnet",
				},
			},
			`(?m)Framework 'dotnet' may have been misspelled because the command for this framework cannot be determined`,
			"An error message should be returned as there is 1 missing command",
		},
	}

	// create the necessary objects
	cfg := config.Config{}
	logger := log.New()
	scaffold := New(&cfg, logger)

	// iterate around the test tables and perform the tests
	for _, table := range tables {
		res := scaffold.analyseMissing(table.missing)

		// compare the result with the pattern
		re := regexp.MustCompile(table.pattern)
		matched := re.MatchString(res)

		if !matched {
			t.Error(table.msg)
		}
	}

}

func TestConfigurePipeline(t *testing.T) {

	cleanup, tempDir := setupScaffoldTestCase(t)
	defer cleanup(t)

	// create the test tables for the different configurations
	tables := []struct {
		cfg  config.Config
		test bool
	}{
		{
			config.Config{
				Input: config.InputConfig{
					Pipeline: "azdo",
					Options: config.Options{
						DryRun: true,
					},
					Project: []config.Project{
						{
							Directory: config.Directory{
								WorkingDir: tempDir,
							},
						},
					},
				},
			},
			false,
		},
	}

	// iterate around the tables
	for _, table := range tables {

		// create the necessary objects
		logger := log.New()
		scaffold := New(&table.cfg, logger)

		scaffold.configurePipeline(&table.cfg.Input.Project[0])

	}
}

func TestShouldRun(t *testing.T) {

	// create the test tables
	tables := []struct {
		tags     []string
		keyword  string
		expected bool
		msg      string
	}{
		{
			[]string{},
			"webapi",
			true,
			"Operation should run as there are no tags",
		},
		{
			[]string{"cqrs"},
			"webapi",
			false,
			"Operation should not run as the keyword is not in the tags",
		},
		{
			[]string{"nextjs", "apps"},
			"apps",
			true,
			"Operation should run because the keyword exists in tags",
		},
	}

	// iterate around the tables
	for _, table := range tables {
		s := Scaffold{}

		result := s.shouldRun(table.tags, table.keyword)

		if result != table.expected {
			t.Error(table.msg)
		}
	}
}
