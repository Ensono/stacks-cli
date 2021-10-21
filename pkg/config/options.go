package config

// Options holds the options for the CLI, such as turning on cmd logging
type Options struct {
	CmdLog     bool `mapstructure:"cmdlog"`
	DryRun     bool `mapstructure:"dryrun"`
	SaveConfig bool `mapstructure:"save" yaml:"-"`
	NoCleanup  bool `mapstructure:"nocleanup" yaml:"-"`
	Clobber    bool `mapstructure:"clobber" yaml:"-"`
}
