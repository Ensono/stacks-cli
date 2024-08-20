package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDotnetSDKVersionWithString(t *testing.T) {

	// create test table to iterate over to check the dotnet version
	tables := []struct {
		sdk         string
		testVersion string
		msgVersion  string
		logMsg      string
	}{
		{
			`{"sdk": {"version": "6.0.200"}}`,
			"6.0.200",
			"Version should be 6.0.200: %s",
			"",
		},
		{
			`{"sdk": {"version": "6.0.200", "rollForward": "latestPatch"}}`,
			"6.0.x",
			"Version should be 6.0.x: %s",
			"Rolling forward to the latest patch/feature version of .NET SDK: 6.0.x",
		},
		{
			`{"sdk": {"version": "6.0.200", "rollForward": "latestFeature"}}`,
			"6.0.x",
			"Version should be 6.0.x: %s",
			"Rolling forward to the latest patch/feature version of .NET SDK: 6.0.x",
		},
	}

	// iterate around the tests tables and perform the tests
	for _, table := range tables {

		version, logMsg, _ := DotnetSDKVersion(table.sdk)

		if version != table.testVersion {
			t.Errorf(table.msgVersion, version)
		}

		if logMsg != table.logMsg {
			t.Errorf("Log message should be '%s' but was '%s'", table.logMsg, logMsg)
		}
	}
}

func TestDotnetSDKVersionFromFile(t *testing.T) {

	// create test able with the name of the file and the content and
	// the expected result of the test
	tables := []struct {
		filename    string
		content     string
		testVersion string
		msgVersion  string
		logMsg      string
	}{
		{
			"global.json",
			`{"sdk": {"version": "6.0.200"}}`,
			"6.0.200",
			"Version should be 6.0.200: %s",
			"",
		},
		{
			"global.json",
			`{"sdk": {"version": "6.0.200", "rollForward": "latestPatch"}}`,
			"6.0.x",
			"Version should be 6.0: %s",
			"Rolling forward to the latest patch/feature version of .NET SDK: 6.0.x",
		},
		{
			"global.json",
			`{"sdk": {"version": "6.0.200", "rollForward": "latestFeature"}}`,
			"6.0.x",
			"Version should be 6.0: %s",
			"Rolling forward to the latest patch/feature version of .NET SDK: 6.0.x",
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

		version, logMsg, _ := DotnetSDKVersion(file)

		if version != table.testVersion {
			t.Errorf(table.msgVersion, version)
		}

		if logMsg != table.logMsg {
			t.Errorf("Log message should be '%s' but was '%s'", table.logMsg, logMsg)
		}
	}
}
