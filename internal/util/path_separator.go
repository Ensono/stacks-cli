package util

import (
	"strings"
)

// Normalise paths ensures that the path separator is correct for the specified platform
func NormalisePath(path string, separator string) string {

	var normalised string
	char := map[string]string{
		"/":  "\\",
		"\\": "/",
	}

	normalised = strings.Replace(path, char[separator], separator, -1)

	return normalised
}
