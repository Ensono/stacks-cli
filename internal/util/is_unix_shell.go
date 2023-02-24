package util

import (
	"github.com/amido/stacks-cli/internal/interfaces"
)

// IsUnixShell determines if the program is running in a Unix
// shell. This is required for the situations where Bash is being
// run on Windows for example.
//
// This function is used to determine if the CLI should use forward
// slashes in any paths that are generated and shown to the user
func IsUnixShell() bool {

	var result = false

	// now determine if running in a Unix like shell
	// cmd := exec.Command("echo", "$0")
	cmd := interfaces.ShellCommander("echo", "$0")

	output, err := cmd.Output()

	if err == nil {

		// if the output is not null then it is running in a *Nix like shell
		// so change the delimiter
		if string(output) != "" && string(output) != "$0" {
			result = true
		}
	}

	return result
}
