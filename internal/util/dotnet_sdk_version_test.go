package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDotnetSDKVersionWithString(t *testing.T) {

	// create test table to iterate over to check the dotnet version
	tables := []struct {
		sdk  string
		test string
		msg  string
	}{
		{
			`{"sdk": {"version": "6.0.200"}}`,
			"6.0.200",
			"Version should be 6.0.200: %s",
		},
	}

	// iterate around the tests tables and perform the tests
	for _, table := range tables {

		version, _ := DotnetSDKVersion(table.sdk)

		if version != table.test {
			t.Errorf(table.msg, version)
		}
	}
}

func TestDotnetSDKVersionFromFile(t *testing.T) {

	// create test able with the name of the file and the content and
	// the expected result of the test
	tables := []struct {
		filename string
		content  string
		test     string
		msg      string
	}{
		{
			"global.json",
			`{"sdk": {"version": "6.0.200"}}`,
			"6.0.200",
			"Version should be 6.0.200: %s",
		},
	}

	// get the dir that the file should be written out to
	dir := t.TempDir()

	// iterate around the tables
	for _, table := range tables {

		// determine the path to the file
		file := filepath.Join(dir, table.filename)

		// create a file with the content in the table and the filename
		if err := os.WriteFile(file, []byte(table.content), 0666); err != nil {
			t.Fatalf("Unable to create '%s' file: %s", file, err.Error())
		}

		version, _ := DotnetSDKVersion(file)

		if version != table.test {
			t.Errorf(table.msg, version)
		}
	}
}
