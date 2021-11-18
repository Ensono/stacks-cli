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

// Set the framework definitions which includes the commands that are associated with
// a framework as well as the commands that are required to get the version of that
// command on the machine that it is being run on
//go:embed framework_def.yml
var framework_defs string

// Config byte parses static
func Config(key string) []byte {

	var result []byte

	switch key {
	case "stacks_frameworks":
		result = []byte(stacks_frameworks)
	case "framework_defs":
		result = []byte(framework_defs)
	}

	return result
}

// FrameworkCommands returns all of the commands that are associated with the specified
// framework and are expected to be run as part of the scaffolding
func FrameworkCommands(framework string) []string {
	commands := map[string][]string{
		"dotnet": {"dotnet", "git"},
		"java":   {"java", "git"},
	}

	return commands[framework]
}
