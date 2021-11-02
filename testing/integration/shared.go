// +build integration

package integration

import (
	"os"
	"path/filepath"
)

// ClearDir clears all of the files and folders within the specified
// directory. This is primarily used by the TearDown function of the test suites
// This is so that the parent directory does not get removed
func ClearDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))

	if err != nil {
		return err
	}

	// iterate around the files that have been found an remove them
	for _, file := range files {
		err := os.RemoveAll((file))
		if err != nil {
			return err
		}
	}

	return err
}

