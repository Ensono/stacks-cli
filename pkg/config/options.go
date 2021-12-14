package config

// Options holds the options for the CLI, such as turning on cmd logging
type Options struct {
	CmdLog       bool   `mapstructure:"cmdlog"`
	DryRun       bool   `mapstructure:"dryrun"`
	SaveConfig   bool   `mapstructure:"save" yaml:"-"`
	NoCleanup    bool   `mapstructure:"nocleanup" yaml:"-"`
	Force        bool   `mapstructure:"force" yaml:"-"`
	NoBanner     bool   `mapstructure:"nobanner"`
	NoCLIVersion bool   `mapstructure:"nocliversion"`
	Token        string `mapstructure:"token" json:"-"`
}
