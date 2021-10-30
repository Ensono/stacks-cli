package util

import (
	"regexp"
	"runtime"
)

// BuildCommand builds up the command to be used, depending on the OS in use
// This is required so that on Windows the command is prepended with "cmd /C" and
// the arguments need to be split up into a slice so that they are passed to the exec.Command
// method
func BuildCommand(command string, arguments string) (string, []string) {

	var cmd string
	var args []string

	// split the argument string into a slice
	// this uses a regular expression to split up the arguments using space as a delimeter
	// However it will not split on a space that is contained within quotes (double or single)
	re := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
	args = re.FindAllString(arguments, -1)

	// if running on Windows then the cmd needs to be set to "cmd" and /C and the command prepended
	// to the args slice, otherwise set cmd to command
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = append([]string{"/C", command}, args...)
	} else {
		cmd = command
	}

	return cmd, args
}
