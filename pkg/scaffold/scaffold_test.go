package scaffold

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/amido/stacks-cli/pkg/config"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func TestSetProjectPath(t *testing.T) {

	// create a configuration object representing the config that might
	// be injected
	project_name := "test_project"

	cfg := config.Config{}
	// config.Input.Project[0].Name = project_name

	project := make([]config.Project, 1)
	project[0].Name = project_name
	cfg.Input.Project = project

	cfg.Input.Directory.WorkingDir = filepath.Join(os.TempDir())

	logger := log.New()

	s := New(&cfg, logger)

	actual := filepath.Join(os.TempDir(), project_name)

	assert.Equal(t, actual, s.setProjectPath(project_name))
}
