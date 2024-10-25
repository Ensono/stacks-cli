package config

type FrameworkDef struct {
	Name     string            `mapstructure:"name" yaml:"name"`
	Commands []FrameworkDefCmd `mapstructure:"commands" yaml:"commands"`
}

func (fd *FrameworkDef) GetCmdList() []string {
	cmds := []string{}
	for _, c := range fd.Commands {
		cmds = append(cmds, c.Name)
	}
	return cmds
}
