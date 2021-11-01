package util

import (
	"runtime"
	"testing"
)

func TestBuildCommand(t *testing.T) {

	// create test table to iterate over
	tables := []struct {
		platform  string
		command   string
		arguments string
		test      string
		count     int
	}{
		{
			"windows",
			"dotnet",
			"new -i .",
			"cmd",
			5,
		},
		{
			"windows",
			"echo",
			`"Hello Golang"`,
			"cmd",
			3,
		},
		{
			"linux",
			"dotnet",
			"new -i .",
			"dotnet",
			3,
		},
	}

	// iterate around the test tables and perform the tests
	for _, table := range tables {

		// run if the OS is the same as the platform
		if runtime.GOOS == table.platform {

			// get the cmd and args from the build command
			cmd, args := BuildCommand(table.command, table.arguments)

			// check that the command
			if cmd != table.test {
				t.Error("Command has not been sect correctly")
			}

			if len(args) != table.count {
				t.Error("Number of arguments is incorrect")
			}
		}

	}
}
