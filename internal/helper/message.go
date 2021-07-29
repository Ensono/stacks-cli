package helper

import (
	"fmt"
	"os"
)

// Only set trace to display internal messages - like
const (
	TRACE = "TRACE"
	// INFO = "INFO"
	// WARN = "WARN"
	// ERROR = "ERROR"
)

const (
	InfoColor    = "\033[1;34m%v\033[0m\n\n"
	NoticeColor  = "\033[1;36m%v\033[0m\n\n"
	WarningColor = "\033[1;33m%v\033[0m\n\n"
	ErrorColor   = "\033[1;31m%v\033[0m\n\n"
	DebugColor   = "\033[0;36m%v\033[0m\n\n"
)

func ShowError(e error) {
	if show := showOutput(); show {
		fmt.Fprintf(os.Stderr, ErrorColor, e.Error())
	}
}

func ShowInfo(msg string) {
	if show := showOutput(); show {
		fmt.Fprintf(os.Stdout, InfoColor, msg)
	}
}

func ShowWarning(msg string) {
	if show := showOutput(); show {
		fmt.Fprintf(os.Stdout, WarningColor, msg)
	}
}

func showOutput() bool {
	val, present := os.LookupEnv("AMIDOSTACKS_CONFIG")
	return present && val == TRACE
}
