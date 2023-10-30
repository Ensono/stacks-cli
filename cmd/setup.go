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
		Run: executeSetup,
	}
)

func init() {

	// declare variable that will be populated from the command line
	var company_name string
	var project_name string
	var area string
	var tf_storage string
	var tf_group string
	var tf_container string
	var global bool

	// Add the run to the command root
	rootCmd.AddCommand(setupCmd)

	// Configure the flags
	setupCmd.Flags().StringVar(&company_name, "company", "", "The name of the company")
	setupCmd.Flags().StringVar(&project_name, "project", "", "The name of the project")
	setupCmd.Flags().StringVar(&area, "area", "", "The area of the project")
	setupCmd.Flags().StringVar(&tf_storage, "storage", "", "The name of the Azure Storage account or S3 bucket being used to store the Terraform state")
	setupCmd.Flags().StringVar(&tf_group, "group", "", "This is the group that contains the storage account for the Terraform state.")
	setupCmd.Flags().StringVar(&tf_container, "container", "", "For Azure storage accounts this is the name of the container to be used for the workspace. For AWS S3 Buckets this is the name of the folder to use.")
	setupCmd.Flags().BoolVarP(&global, "global", "g", false, "Set the values globally. These will be set in the user's home directory")

	// bind the flags to the configuration
	viper.BindPFlag("input.business.company", setupCmd.Flags().Lookup("company"))
	viper.BindPFlag("input.business.project", setupCmd.Flags().Lookup("project"))
	viper.BindPFlag("input.business.domain", setupCmd.Flags().Lookup("area"))
	viper.BindPFlag("input.terraform.backend.storage", setupCmd.Flags().Lookup("storage"))
	viper.BindPFlag("input.terraform.backend.group", setupCmd.Flags().Lookup("group"))
	viper.BindPFlag("input.terraform.backend.container", setupCmd.Flags().Lookup("container"))
	viper.BindPFlag("input.global", setupCmd.Flags().Lookup("global"))

}

func executeSetup(ccmd *cobra.Command, args []string) {

	// call the setup method
	setup := setup.New(&Config, App.Logger)
	err := setup.Run()
	if err != nil {
		App.Logger.Fatalf("Error running setup: %s", err.Error())
	}
}
