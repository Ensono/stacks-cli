package cmd

import (
	"github.com/Ensono/stacks-cli/pkg/interactive"
	"github.com/spf13/cobra"
)

var (
	interactiveCmd = &cobra.Command{
		Use:   "interactive",
		Short: "Generate a configuration though interactive questions",
		Long: `Setting up the configuration file for the Stacks CLI can seem a bit daunting.

By using the "interactive" sub command you will be asked a series of questions about the projects
you wish to configure. The answers to these questions will then be written out to a configuration
file tha can be read by the "scaffold" sub command.`,
		Run: executeInteractiveRun,
	}
)

func init() {
	// Add the command to the root
	rootCmd.AddCommand(interactiveCmd)
}

func executeInteractiveRun(ccmd *cobra.Command, args []string) {

	inter := interactive.New(&Config, App.Logger)
	err := inter.Run()
	if err != nil {
		App.Logger.Fatalf("Error running interactive configuration: %s", err.Error())
	}
}
