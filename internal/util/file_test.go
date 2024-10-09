package util

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

// setupFileTests creates a temporary directory and files that can be used for the tests
func setupFileTests(t *testing.T) (func(t *testing.T), string) {

	// create a temporary directory
	tempDir := t.TempDir()

	// create a slice of the files that need to be created
	files := []string{
		"build/taskctl/contexts.yaml",
		"build/azureDevops/pipeline.yaml",
	}

	// iterate around the files and create them in the tempDir
	for _, file := range files {
		path := filepath.Join(tempDir, file)

		// ensure that the directory for the file exists
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			log.Fatalf("failed to create directory: %v", err)
		}

		// create the file
		f, err := os.Create(path)
		if err != nil {
			log.Fatalf("failed to create file: %v", err)
		}
		defer f.Close()
	}

	return func(t *testing.T) {
		log.Println("Cleaning up")
	}, tempDir
}

// TestGetFileList tests that the GetFileList function returns the correct list of files
// based on a glob pattern
func TestGetFileList(t *testing.T) {

	// Call the setup function for the test
	teardownTest, dir := setupFileTests(t)
	defer teardownTest(t)

	// create the test table, which will have the different settings and check how many
	// files have been found by the GetFileList function
	tables := []struct {
		pattern  string
		expected int
	}{
		{
			"build/taskctl/contexts.yaml",
			1,
		},
		{
			"build/taskctl/*.yaml",
			1,
		},
		{
			"build/**/*.yaml",
			2,
		},
	}

	// iterate around the test tables
	for _, table := range tables {
		files, _ := GetFileList(table.pattern, dir)

		if len(files) != table.expected {
			t.Errorf("Expected %d files, but got %d", table.expected, len(files))
		}
	}

}
