// amidostacks scaffolding-cli command line tool
package main

import (
	"fmt"
	"log"
	"os"

	"amido.com/stacks-cli/internal/helper"
	"amido.com/stacks-cli/pkg/config"
	"amido.com/stacks-cli/pkg/scaffold"
	"github.com/urfave/cli/v2"
)

func main() {
	app := InitApp()
	err := app.Run(os.Args)
	// should exit here and only one possible route out
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: move this out to a helper AppInit
func InitApp() *cli.App {
	app := &cli.App{
		Name:     "scaffolding",
		HelpName: "scaffolding",
		Commands: []*cli.Command{
			{
				Flags: helper.StacksFlags(),
				Name:  "run",
				Usage: "runs the scaffolding process",
				Action: func(c *cli.Context) error {
					confFile := c.String("conf")
					replaceConf := c.String("replace-conf")
					sourceConf := c.String("source-conf")
					interActive := c.Bool("interactive")
					if interActive {
						initInterActiveFlow()
					}
					if confFile != "" {
						initConfigFileFlow(confFile, replaceConf, sourceConf)
					}
					return nil
				},
			},
			{
				Name:  "test",
				Usage: "runs test scaffolding only",
				Action: func(c *cli.Context) error {
					helper.ShowInfo(fmt.Sprintf("completed task: %s\n", c.Args().First()))
					return nil
				},
			},
		},
	}
	return app
}

func initConfigFileFlow(confFile, replaceConf, sourceConf string) error {

	helper.ShowInfo(fmt.Sprintf("Config file: %s\n", confFile))

	file, err := os.ReadFile(confFile)
	if err != nil {
		helper.ShowError(err)
		return err
	}

	conf, err := config.Create(file)
	if err != nil {
		helper.ShowError(err)
		return err
	}

	sc := scaffold.New(conf)

	if err = sc.Run(); err != nil {
		return err
	}

	helper.ShowInfo(fmt.Sprintf("Used Contents of: %s.\n", confFile))

	return nil
}

func initInterActiveFlow() {
	helper.ShowInfo("Started Interactive Flow\n")
}
