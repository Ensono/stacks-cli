package config

type Environment struct {
	Name      string `mapstructure:"name" yaml:",omitempty"`
	StageName  string `mapstructure:"stagename" yaml:",omitempty"`
	ShortName  string `mapstructure:"shortname" yaml:",omitempty"`
	ProductionEquivalent bool `mapstructure:"production_equivalent" yaml:",omitempty"`
	TriggerFromMainBranch bool `mapstructure:"trigger_from_main_branch" yaml:",omitempty"`
	DependsOn []string `mapstructure:"depends_on" yaml:",omitempty"`
}
