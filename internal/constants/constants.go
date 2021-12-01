package constants

const (
	// AppName states the name of the application
	AppName string = "Stacks CLI"

	// Set the timestamp format for logging
	LoggingTimestamp = "Mon, 02 Jan 2006 15:04:05 -0700"

	// Set the prefix that needs to be used when setting the configuration
	// using environment variables
	EnvVarPrefix = "amidostacks"

	// DefaultErrorMessage defines the default error message if one has not been set
	DefaultErrorMessage = "An error occurred in the application"

	// DefaultVersion sets a default version number if one is not set during the build
	// if this is seen when `stacks-cli -v` is run then it means it has been built
	// on a local machine
	DefaultVersion = "0.0.1-workstation"

	// SettingsFile is the default filename to be used when looking for the file in
	// a project that is to be used with stacks
	SettingsFile = "stackscli.yml"

	// GitHubRef is the org/name of the stacks-cli
	// This is used as on github api calls
	GitHubRef = "amido/stacks-cli"
)
