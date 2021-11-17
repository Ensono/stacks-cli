package static

import (
	_ "embed"
)

// The following static configuration sets the URLs to the repos for
// the location of the repositories
// This can be overriden by passing the configuration in as a configuration file
// but this will be the default
//go:embed stacks_frameworks.yml
var stacks_frameworks string

// Set the banner that is written out to the screen when stacks is run
//go:embed banner.txt
var Banner string

// Config byte parses static
func Config(key string) []byte {

	var result []byte

	switch key {
	case "stacks_frameworks":
		result = []byte(stacks_frameworks)
	}

	return result
}

// FrameworkCommand returns all of the commands that are associated with the specified
// framework and are expected to be run as part of the scaffolding
func FrameworkCommand(framework string) []string {
	commands := map[string][]string{
		"dotnet": {"dotnet", "git"},
		"java":   {"java", "git"},
	}

	return commands[framework]
}
