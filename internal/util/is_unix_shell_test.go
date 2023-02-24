package util

import (
	"fmt"
	"testing"

	"github.com/amido/stacks-cli/internal/interfaces"
)

// Mock the exec calls to the undelying system
type tcShellCommand struct {
	OutputterFunc func() ([]byte, error)
}

func (sc tcShellCommand) Output() ([]byte, error) {
	return sc.OutputterFunc()
}

type execCommandFunc func(name string, arg ...string) interfaces.IShellCommand

func newMockShellCommanderForOutput(output string, err error, t *testing.T) execCommandFunc {
	testName := t.Name()
	return func(name string, arg ...string) interfaces.IShellCommand {
		fmt.Printf("exec.Command() for %v called with %v and %v", testName, name, arg)
		outputterFunc := func() ([]byte, error) {
			if err == nil {
				fmt.Printf("Output obtained for %v\n", testName)
			} else {
				fmt.Printf("Failed to get output for %v\n", testName)
			}
			return []byte(output), err
		}

		return tcShellCommand{
			OutputterFunc: outputterFunc,
		}
	}
}

func TestIsUnixShell(t *testing.T) {

	// swap out the shell commander
	curShellCommander := interfaces.ShellCommander
	defer func() { interfaces.ShellCommander = curShellCommander }()

	// create a list of the tests that need to be performed
	tables := []struct {
		output   string
		expected bool
		err      error
	}{
		{
			"",
			false,
			nil,
		},
		{
			"-bash",
			true,
			nil,
		},
	}

	// iterate around the test tables
	for _, table := range tables {

		interfaces.ShellCommander = newMockShellCommanderForOutput(table.output, table.err, t)
		result := IsUnixShell()

		if result != table.expected {
			t.Errorf("Unexpected output [\"%s\"], should be %v, got %v", table.output, table.expected, result)
		}
	}

}
