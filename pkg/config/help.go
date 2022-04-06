package config

// Help struct holds the urls for the help pages for the different commands
type Help struct {
	Root        string `mapstructure:"root"`
	Scaffold    string `mapstructure:"scaffold"`
	Interactive string `mapstructure:"interactive"`
	Version     string `mapstructure:"version"`
}

func (help *Help) GetUrl(cmd string) string {

	var url string

	switch cmd {
	case "root":
		url = help.Root
	case "scaffold":
		url = help.Scaffold
	case "interactive":
		url = help.Interactive
	case "version":
		url = help.Version
	}

	return url
}
