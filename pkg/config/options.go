package config

// Options holds the options for the CLI, such as turning on cmd logging
type Options struct {
	CmdLog bool `mapstructure:"cmdlog"`
	DryRun bool `mapstructure:"dryrun"`
}
