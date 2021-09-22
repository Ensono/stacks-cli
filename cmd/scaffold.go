package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/amido/stacks-cli/internal/constants"
	"github.com/amido/stacks-cli/internal/util"
	"github.com/amido/stacks-cli/pkg/scaffold"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	scaffoldCmd = &cobra.Command{
		Use:   "scaffold",
		Short: "Create a new project using Amido Stacks",
		Long:  "",
		Run:   executeRun,
	}
)

func init() {

	// declare variables that will be populated from the command line
	var interactive bool

	// - project settings
	var project_name string
	var project_vcs_type string
	var project_vcs_url string
	var project_vcs_ref string
	var settings_file string

	// - framework settings
	var framework_type string
	var framework_option string
	var framework_version string

	// - platform settings
	var platform_type string

	// - pipeline
	var pipeline string

	// - cloud settings
	var cloud_platform string
	var cloud_region string
	var cloud_group string

	// - business settings
	var business_company string
	var business_domain string
	var business_component string

	// - terraform settings
	var terraform_backend_storage string
	var terraform_backend_group string
	var terraform_backend_container string

	// - network settings
	var network_base_domain string

	// Add the run command to the root
	rootCmd.AddCommand(scaffoldCmd)

	// Configure the flags
	// - run interactively
	scaffoldCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Run in interactive mode")

	scaffoldCmd.Flags().StringVarP(&project_name, "name", "n", "", "Name of the project to create")
	scaffoldCmd.Flags().StringVar(&project_vcs_type, "sourcecontrol", "github", "Type of source control being used")
	scaffoldCmd.Flags().StringVarP(&project_vcs_url, "sourcecontrolurl", "u", "", "Url of the remote for source control")
	scaffoldCmd.Flags().StringVar(&project_vcs_ref, "sourcecontrolref", "", "SHA reference or Tag to use to clone repo from")

	scaffoldCmd.Flags().StringVarP(&framework_type, "framework", "F", "", "Framework for the project")
	scaffoldCmd.Flags().StringVarP(&framework_option, "frameworkoption", "O", "", "Option of the chosen framework to use")
	scaffoldCmd.Flags().StringVarP(&framework_version, "frameworkversion", "V", "latest", "Version of the framework package to download")

	scaffoldCmd.Flags().StringVarP(&platform_type, "platformtype", "P", "", "Type of platform being deployed to")

	scaffoldCmd.Flags().StringVarP(&pipeline, "pipeline", "p", "", "Pipeline to use for CI/CD")

	scaffoldCmd.Flags().StringVarP(&cloud_platform, "cloud", "C", "", "Cloud platform being targetted")
	scaffoldCmd.Flags().StringVarP(&cloud_region, "cloudregion", "R", "", "Region that the resources should be deployed to")
	scaffoldCmd.Flags().StringVarP(&cloud_group, "cloudgroup", "G", "", "Group that the resources should belong to")

	scaffoldCmd.Flags().StringVar(&business_company, "company", "", "Name of your company")
	scaffoldCmd.Flags().StringVarP(&business_domain, "area", "A", "", "Area within the company that this project will belong to, e.g. core")
	scaffoldCmd.Flags().StringVar(&business_component, "component", "", "Business component, e.g. infrastructure")

	scaffoldCmd.Flags().StringVar(&terraform_backend_storage, "tfstorage", "", "Name of the storage to be used for Terraform state")
	scaffoldCmd.Flags().StringVar(&terraform_backend_group, "tfgroup", "", "Name of the group that the storage account is in")
	scaffoldCmd.Flags().StringVar(&terraform_backend_container, "tfcontainer", "", "Name of the container within the storage to use")

	scaffoldCmd.Flags().StringVarP(&network_base_domain, "domain", "d", "", "Domain for the app")

	scaffoldCmd.Flags().StringVar(&settings_file, "settingsfile", constants.SettingsFile, "Name of the settings file to look for in a project")

	// Bind the flags to the configuration

	// The project is a slice, so that multiple projects can be specified, however
	// only one can be specified on the command line and in Environment variables
	// Viper works out that this is a slice and will bind to the first element of the slice
	viper.BindPFlag("project.name", scaffoldCmd.Flags().Lookup("name"))
	viper.BindPFlag("project.framework.type", scaffoldCmd.Flags().Lookup("framework"))
	viper.BindPFlag("project.framework.option", scaffoldCmd.Flags().Lookup("frameworkoption"))
	viper.BindPFlag("project.framework.version", scaffoldCmd.Flags().Lookup("frameworkversion"))
	viper.BindPFlag("project.platform.type", scaffoldCmd.Flags().Lookup("platformtype"))
	viper.BindPFlag("project.sourcecontrol.type", scaffoldCmd.Flags().Lookup("sourcecontrol"))
	viper.BindPFlag("project.sourcecontrol.url", scaffoldCmd.Flags().Lookup("sourcecontrolurl"))
	viper.BindPFlag("project.sourcecontrol.ref", scaffoldCmd.Flags().Lookup("sourcecontrolref"))

	viper.BindPFlag("settingsfile", scaffoldCmd.Flags().Lookup("settingsfile"))

	viper.BindPFlag("pipeline", scaffoldCmd.Flags().Lookup("pipeline"))

	viper.BindPFlag("platform.type", scaffoldCmd.Flags().Lookup("platformtype"))

	viper.BindPFlag("cloud.platform", scaffoldCmd.Flags().Lookup("cloud"))
	viper.BindPFlag("cloud.region", scaffoldCmd.Flags().Lookup("cloudregion"))
	viper.BindPFlag("cloud.group", scaffoldCmd.Flags().Lookup("cloudgroup"))

	viper.BindPFlag("business.company", scaffoldCmd.Flags().Lookup("company"))
	viper.BindPFlag("business.domain", scaffoldCmd.Flags().Lookup("area"))
	viper.BindPFlag("business.component", scaffoldCmd.Flags().Lookup("component"))

	viper.BindPFlag("terraform.backend.storage", scaffoldCmd.Flags().Lookup("tfstorage"))
	viper.BindPFlag("terraform.backend.group", scaffoldCmd.Flags().Lookup("tfgroup"))
	viper.BindPFlag("terraform.backend.container", scaffoldCmd.Flags().Lookup("tfcontainer"))

	viper.BindPFlag("network.base.domain", scaffoldCmd.Flags().Lookup("domain"))

	viper.BindPFlag("interactive", scaffoldCmd.Flags().Lookup("interactive"))
}

func executeRun(ccmd *cobra.Command, args []string) {

	// ensure that at least one project has been specified
	if len(Config.Input.Project) == 1 && Config.Input.Project[0].Name == "" {
		App.Logger.Fatalln("No projects have been defined")
	}

	// Call the scaffolding method
	err := scaffold.New(&Config, App.Logger).Run()
	if err != nil {
		App.Logger.Fatalf("Error running scaffold: %s", err.Error())
	}

	// Ensure that the temp directory is removed
	App.Logger.Info("Performing cleanup")
	err = os.RemoveAll(Config.Input.Directory.TempDir)
	if err != nil {
		App.Logger.Fatalf("Unable to remove temporary directory: %s", Config.Input.Directory.TempDir)
	}

	// iterate around the projects that have been specified
	for _, project := range Config.Input.Project {

		App.Logger.Infof("Setting up project: %s\n", project.Name)

		// Create the temporary and working directories for the current project
		projectTempDir := filepath.Join(Config.Input.Directory.TempDir, project.Name)
		// projectDir := filepath.Join(Config.Input.Directory.WorkingDir, project.Name)

		// Clone the target repository into the temp directory
		err := util.GitClone(
			Config.Input.Stacks.GetSrcURL(project.Framework.GetMapKey()),
			project.SourceControl.Ref,
			projectTempDir,
		)

		if err == nil {
			fmt.Println("clones")
		}

		// Read in the configuration file for the project

		// copy the contents of the temporary project dir to the working directory
		// util.CopyDirectory(filepath.Join(projectTempDir, "*"), Config.Self.GetPath(project))
	}

}
