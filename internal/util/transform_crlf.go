package util

import "strings"

// TransformCRLF takes an input and ensures that all line endings are in LF
// This is to ensure that files read on a Windows machine behave as required with
// some libraries
func TransformCRLF(input string) string {
	return strings.ReplaceAll(input, "\r\n", "\n")
}
