package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Create a new project using Amido Stacks",
		Long:  "",
		Run:   executeRun,
	}
)

func init() {

	// declare variables that will be populated from the command line
	var interactive bool

	// Add the run command to the root
	rootCmd.AddCommand(runCmd)

	// Configure the flags
	runCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Run in interactive mode")

	// Bind the flags to the configuration
	viper.BindPFlag("interactive", runCmd.Flags().Lookup("interactive"))
}

func executeRun(ccmd *cobra.Command, args []string) {

}
