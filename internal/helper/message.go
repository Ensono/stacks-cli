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

const (
	DEBUG_MESSAGE = "\n\n\nIf you want to see a full log trace - set env var AMIDOSTACKS_LOG to TRACE\n"
)

// Displays any Errors to end user
func ShowError(e error) {
	fmt.Fprintf(os.Stderr, ErrorColor, fmt.Sprintf("%s%s", e.Error(), DEBUG_MESSAGE))
	TraceError(e)
}

// Shows Info to end user
func ShowInfo(msg string) {
	fmt.Fprintf(os.Stdout, InfoColor, msg)

}

// Displays any Info to end user WHEN LOGGING enabled
func TraceInfo(msg string) {
	if show := showOutput(); show {
		fmt.Fprintf(os.Stdout, InfoColor, msg)
	}
}

// Displays any Warnings to end user WHEN LOGGING enabled
func TraceWarning(msg string) {
	if show := showOutput(); show {
		fmt.Fprintf(os.Stdout, WarningColor, msg)
	}
}

// Displays any Info to end user WHEN LOGGING enabled
func TraceError(e error) {
	if show := showOutput(); show {
		fmt.Fprintf(os.Stdout, WarningColor, e.Error())
	}
}

func showOutput() bool {
	val, present := os.LookupEnv("AMIDOSTACKS_LOG")
	return present && val == TRACE
}
