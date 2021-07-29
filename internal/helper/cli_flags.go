package helper

import "github.com/urfave/cli/v2"

func StacksFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "conf",
			Aliases: []string{"c", "config"},
			Value:   "",
			Usage:   "Config file location",
			EnvVars: []string{"AMIDOSTACKS_CONFIG"},
		},
		&cli.StringFlag{
			Name:    "source-conf",
			Aliases: []string{"sc", "source-config"},
			Value:   "",
			Usage:   "Source Config file location",
			EnvVars: []string{"AMIDOSTACKS_CONFIG_SOURCE"},
		},
		&cli.StringFlag{
			Name:    "replace-conf",
			Aliases: []string{"rc", "replace-config"},
			Value:   "",
			Usage:   "Replace Config file location",
			EnvVars: []string{"AMIDOSTACKS_CONFIG_REPLACE"},
		},
		&cli.BoolFlag{
			Name:    "interactive",
			Aliases: []string{"i"},
			Value:   false,
			Usage:   "Run in interactive mode. will be ignored if a configuration file is provided",
		},
	}
}
