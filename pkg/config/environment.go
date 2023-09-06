package config

type Environment struct {
	Name      string `mapstructure:"name" yaml:",omitempty"`
	StageName  string `mapstructure:"stagename" yaml:",omitempty"`
	ShortName  string `mapstructure:"shortname" yaml:",omitempty"`
	IsProduction bool `mapstructure:"isproduction" yaml:",omitempty"`
	TriggerFromMainBranch bool `mapstructure:"triggerfrommainbranch" yaml:",omitempty"`
	DependsOn []string `mapstructure:"depends_on" yaml:",omitempty"`
}
