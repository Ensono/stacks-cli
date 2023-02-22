package interfaces

import (
	"os/exec"
)

type IShellCommand interface {
	Output() ([]byte, error)
}

type execShellCommand struct {
	*exec.Cmd
}

func newExecShellCommander(name string, arg ...string) IShellCommand {
	execCmd := exec.Command(name, arg...)
	return execShellCommand{Cmd: execCmd}
}

var ShellCommander = newExecShellCommander
