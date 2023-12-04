package cmd

import (
	"github.com/amido/stacks-cli/pkg/setup"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Setup default configuration files",
		Long: `When running in interactive mode, the CLI will ask a series of questions,
		the answers to which maybe standard. This command allows values to be set for these
		types of questions`,
	}

	setupUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Add or update a setting in the configuration file",
		Long: `Stacks CLI has the ability to store parts of the configuration in multiple
		places. This command adds or updates a setting in the configuration file`,
		Run: executeSetupUpdate,
	}

	setupListCmd = &cobra.Command{
		Use:   "list",
		Short: "List locations where configuration files would be read from",
		Long: `Stacks CLI allows several locations to be set for configuration file, this command
		shows where those files would be read from`,
		Run: executeSetupList,
	}
)

func init() {

	// declare variable that will be populated from the command line
	var project_name string

	// add the sub commands for the setup command
	setupCmd.AddCommand(setupUpdateCmd)
	setupCmd.AddCommand(setupListCmd)

	// Add the run to the command root
	rootCmd.AddCommand(setupCmd)

	// Configure the flags for the update command
	setupUpdateCmd.Flags().StringVar(&project_name, "project", "", "The name of the project")

	// bind the flags to the configuration
	viper.BindPFlag("input.business.project", setupUpdateCmd.Flags().Lookup("project"))
}

func executeSetupUpdate(ccmd *cobra.Command, args []string) {

	// call the setup method
	setup := setup.New(&Config, App.Logger)
	err := setup.Upsert()
	if err != nil {
		App.Logger.Fatalf("Error running update: %s", err.Error())
	}
}

func executeSetupList(ccmd *cobra.Command, args []string) {

	// call the setup method
	setup := setup.New(&Config, App.Logger)
	err := setup.List()
	if err != nil {
		App.Logger.Fatalf("Error running list: %s", err.Error())
	}
}
