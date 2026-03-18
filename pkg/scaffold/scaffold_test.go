package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/Ensono/stacks-cli/internal/models"
	"github.com/Ensono/stacks-cli/pkg/config"

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

// createTestFileTree creates a directory tree with the given file paths relative to root.
// Each file is created with a short content string.
func createTestFileTree(t *testing.T, root string, files []string) {
	t.Helper()
	for _, f := range files {
		fullPath := filepath.Join(root, f)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory for %s: %v", f, err)
		}
		if err := os.WriteFile(fullPath, []byte("# "+f), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", f, err)
		}
	}
}

func TestPerformOperationCopyWithExclude(t *testing.T) {

	tables := []struct {
		name            string
		exclude         []string
		expectedFiles   []string
		unexpectedFiles []string
	}{
		{
			"no exclude copies everything",
			[]string{},
			[]string{
				"deploy/terraform/main.tf",
				"build/azdo/azure-pipelines.yml",
				".github/workflows/deploy.yml",
			},
			[]string{},
		},
		{
			"exclude github actions",
			[]string{".github/**"},
			[]string{
				"deploy/terraform/main.tf",
				"build/azdo/azure-pipelines.yml",
			},
			[]string{
				".github/workflows/deploy.yml",
			},
		},
		{
			"exclude azdo pipelines",
			[]string{"build/azdo/**"},
			[]string{
				"deploy/terraform/main.tf",
				".github/workflows/deploy.yml",
			},
			[]string{
				"build/azdo/azure-pipelines.yml",
			},
		},
		{
			"exclude multiple patterns",
			[]string{".github/**", "build/azdo/**"},
			[]string{
				"deploy/terraform/main.tf",
			},
			[]string{
				"build/azdo/azure-pipelines.yml",
				".github/workflows/deploy.yml",
			},
		},
	}

	sourceFiles := []string{
		"deploy/terraform/main.tf",
		"deploy/terraform/variables.tf",
		"build/azdo/azure-pipelines.yml",
		"build/azdo/pipeline-vars.template.yml",
		".github/workflows/deploy.yml",
		"stackscli.yml",
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {

			cloneDir := t.TempDir()
			destDir := t.TempDir()

			createTestFileTree(t, cloneDir, sourceFiles)

			cfg := config.Config{}
			logger := log.New()
			logger.SetLevel(log.FatalLevel)
			s := New(&cfg, logger)

			op := config.Operation{
				Action:  "copy",
				Exclude: table.exclude,
			}

			project := &config.Project{}
			err := s.PerformOperation(op, project, destDir, cloneDir)
			if err != nil {
				t.Fatalf("PerformOperation returned error: %v", err)
			}

			for _, f := range table.expectedFiles {
				path := filepath.Join(destDir, f)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("Expected file %s to exist but it was not found", f)
				}
			}

			for _, f := range table.unexpectedFiles {
				path := filepath.Join(destDir, f)
				if _, err := os.Stat(path); err == nil {
					t.Errorf("Expected file %s to be excluded but it exists", f)
				}
			}
		})
	}
}

func TestPipelineFiltering(t *testing.T) {

	tables := []struct {
		name             string
		opPipeline       string
		selectedPipeline string
		shouldSkip       bool
	}{
		{
			"empty pipeline runs for any selection",
			"",
			"azdo",
			false,
		},
		{
			"matching pipeline runs",
			"azdo",
			"azdo",
			false,
		},
		{
			"non-matching pipeline is skipped",
			"gha",
			"azdo",
			true,
		},
		{
			"case insensitive match",
			"AZDO",
			"azdo",
			false,
		},
	}

	sourceFiles := []string{
		"deploy/terraform/main.tf",
		"build/azdo/azure-pipelines.yml",
		".github/workflows/deploy.yml",
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {

			cloneDir := t.TempDir()
			destDir := t.TempDir()

			createTestFileTree(t, cloneDir, sourceFiles)

			cfg := config.Config{
				Input: config.InputConfig{
					Pipeline: table.selectedPipeline,
				},
			}
			logger := log.New()
			logger.SetLevel(log.FatalLevel)
			s := New(&cfg, logger)

			op := config.Operation{
				Action:   "copy",
				Pipeline: table.opPipeline,
			}

			project := config.Project{
				Directory: config.Directory{
					WorkingDir: destDir,
				},
				Phases: []config.Phase{
					{
						Name:       "setup",
						Directory:  destDir,
						Operations: []config.Operation{op},
					},
				},
				Framework: config.Framework{
					Option: "alz_management",
				},
			}

			// Run through the phase loop logic to test pipeline filtering
			for _, phase := range project.Phases {
				for _, phaseOp := range phase.Operations {
					if phaseOp.Pipeline != "" && !strings.EqualFold(phaseOp.Pipeline, s.Config.Input.Pipeline) {
						continue
					}
					if s.shouldRun(phaseOp.Tags, project.Framework.Option) {
						err := s.PerformOperation(phaseOp, &project, phase.Directory, cloneDir)
						if err != nil {
							t.Fatalf("PerformOperation returned error: %v", err)
						}
					}
				}
			}

			mainTf := filepath.Join(destDir, "deploy/terraform/main.tf")
			if table.shouldSkip {
				if _, err := os.Stat(mainTf); err == nil {
					t.Errorf("Expected operation to be skipped but files were copied")
				}
			} else {
				if _, err := os.Stat(mainTf); os.IsNotExist(err) {
					t.Errorf("Expected operation to run but files were not copied")
				}
			}
		})
	}
}

// equalFoldStr is no longer needed — using strings.EqualFold directly
