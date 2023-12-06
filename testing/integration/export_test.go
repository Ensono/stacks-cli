//go:build integration
// +build integration

package integration

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ExportSuite creates a suite of tests that will check the export command
// of the binary is being tested
type ExportSuite struct {
	BaseIntegration
}

// TestExportSuite runs the suite of tests for the export command
func TestExportSuite(t *testing.T) {

	s := new(ExportSuite)
	s.BinaryCmd = *binaryCmd
	s.Company = *company
	s.Project = *project
	s.ProjectDir = *projectDir

	s.Assert = assert.New(t)

	s.SetProjectDir()

	suite.Run(t, s)
}

func (suite *ExportSuite) TestExportInternalConfig() {

	// run the command and then check the output
	arguments := fmt.Sprintf("export -d %s", filepath.Join(suite.ProjectDir, "exported"))
	suite.BaseIntegration.RunCommand(suite.BinaryCmd, arguments, false)

	// check that the export directory exists
	suite.T().Run("Export directory exists", func(t *testing.T) {
		path := filepath.Join(suite.ProjectDir, "exported")
		exists := util.Exists(path)

		suite.Assert.Equal(true, exists, "Export directory should exist: %s", path)
	})

	// check that the files have been created in the specified directory
	suite.T().Run("Internal files have been written out", func(t *testing.T) {

		// create a slice of files to check for
		paths := []string{
			filepath.Join(suite.ProjectDir, "exported", "internal_config.yml"),
			filepath.Join(suite.ProjectDir, "exported", "azdo_variable_template.yml"),
		}

		// iterate around the files can check that they exist
		for _, path := range paths {
			exists := util.Exists(path)

			suite.Assert.Equal(true, exists, "File should exist: %s", path)
		}
	})

	// output should contain the files that have been exported to the filesystem
	suite.T().Run("Command output should contain the relevant information", func(t *testing.T) {

		// set the pattern to check in the output
		pattern := fmt.Sprintf(`(?m)config.*internal_config\.yml|azdo.*azdo_variable_template\.yml`)

		matched := suite.CheckCmdOutput(pattern)

		suite.Assert.Equal(true, matched, "Exported files should be displayed in the output")
	})
}
