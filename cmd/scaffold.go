package cmd

import (
	"log"
	"os"

	"github.com/Ensono/stacks-cli/internal/config/staticFiles"
	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/Ensono/stacks-cli/pkg/scaffold"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	scaffoldCmd = &cobra.Command{
		Use:   "scaffold",
		Short: "Create a new project using Ensono Stacks",
		Long:  "",
		Run:   executeScaffoldRun,
	}
)

func init() {

	// declare variables that will be populated from the command line
	// - options
	var cmdlog bool
	var saveConfig bool
	var nocleanup bool
	var force bool
	var norun bool

	// - scaffold directories
	var cacheDir string

	// - project settings
	var project_name string
	var project_vcs_type string
	var project_vcs_url string
	var project_settings_file string

	// - framework settings
	var framework_type string
	var framework_option string
	var framework_version string
	var framework_properties []string

	// - platform settings
	var platform_type string

	// - pipeline
	var pipeline string

	// - cloud settings
	var cloud_platform string
	var cloud_region string
	var cloud_group string

	// - business settings
	var business_component string

	// - network settings
	var network_base_domain_external string
	var network_base_domain_internal string

	// - overrides
	var override_ado_variables string

	// get the default directories
	defaultCacheDir := util.GetDefaultCacheDir()

	// Add the run command to the root
	rootCmd.AddCommand(scaffoldCmd)

	// Configure the flags
	scaffoldCmd.Flags().StringVar(&cacheDir, "cachedir", defaultCacheDir, "Cache directory to be used for all downloads")

	scaffoldCmd.Flags().StringVarP(&project_name, "name", "n", "", "Name of the project to create")
	scaffoldCmd.Flags().StringVar(&project_vcs_type, "sourcecontrol", "github", "Type of source control being used")
	scaffoldCmd.Flags().StringVarP(&project_vcs_url, "sourcecontrolurl", "u", "", "Url of the remote for source control")
	scaffoldCmd.Flags().StringVar(&project_settings_file, "projectsettingsfile", "", "Path to a settings file to use for the project")

	scaffoldCmd.Flags().StringVarP(&framework_type, "framework", "F", "", "Framework for the project")
	scaffoldCmd.Flags().StringVarP(&framework_option, "frameworkoption", "O", "", "Option of the chosen framework to use")
	scaffoldCmd.Flags().StringVarP(&framework_version, "frameworkversion", "V", "latest", "Version of the framework package to download")

	// get the properties from the command line
	scaffoldCmd.Flags().StringSliceVar(&framework_properties, "frameworkprops", []string{}, "Properties to pass to the project settings")

	scaffoldCmd.Flags().StringVarP(&platform_type, "platformtype", "P", "", "Type of platform being deployed to")

	scaffoldCmd.Flags().StringVarP(&pipeline, "pipeline", "p", "", "Pipeline to use for CI/CD")

	scaffoldCmd.Flags().StringVarP(&cloud_platform, "cloud", "C", "", "Cloud platform being targetted")
	scaffoldCmd.Flags().StringVarP(&cloud_region, "cloudregion", "R", "", "Region that the resources should be deployed to")
	scaffoldCmd.Flags().StringVarP(&cloud_group, "cloudgroup", "G", "", "Group that the resources should belong to")

	scaffoldCmd.Flags().StringVar(&business_component, "component", "", "Business component, e.g. infrastructure")

	scaffoldCmd.Flags().StringVarP(&network_base_domain_external, "domain", "d", "", "External domain for the app")
	scaffoldCmd.Flags().StringVar(&network_base_domain_internal, "internaldomain", "", "Internal domain for the app")

	scaffoldCmd.Flags().StringVar(&override_ado_variables, "adovariables", "", "Path to the ado variables override file")

	scaffoldCmd.Flags().BoolVar(&cmdlog, "cmdlog", false, "Specify if commands should be logged")
	scaffoldCmd.Flags().BoolVar(&saveConfig, "save", false, "Save the the configuration from interactive or command line settings. Has no effect when using a configuration file.")
	scaffoldCmd.Flags().BoolVar(&nocleanup, "nocleanup", false, "If set, do not perform cleanup at the end of the scaffolding")
	scaffoldCmd.Flags().BoolVar(&force, "force", false, "If set, remove existing project directories before attempting to create new ones")
	scaffoldCmd.Flags().BoolVar(&norun, "norun", false, "When used in conjunction with --save, will not attempt to scaffold the projects but will just create config file")

	// Bind the flags to the configuration

	// The project is a slice, so that multiple projects can be specified, however
	// only one can be specified on the command line and in Environment variables
	// Viper works out that this is a slice and will bind to the first element of the slice
	viper.BindPFlag("input.project.name", scaffoldCmd.Flags().Lookup("name"))
	viper.BindPFlag("input.project.platform.type", scaffoldCmd.Flags().Lookup("platformtype"))
	viper.BindPFlag("input.project.sourcecontrol.type", scaffoldCmd.Flags().Lookup("sourcecontrol"))
	viper.BindPFlag("input.project.sourcecontrol.url", scaffoldCmd.Flags().Lookup("sourcecontrolurl"))

	viper.BindPFlag("input.project.settingsfile", scaffoldCmd.Flags().Lookup("projectsettingsfile"))
	viper.BindPFlag("input.project.cloud.region", scaffoldCmd.Flags().Lookup("cloudregion"))
	viper.BindPFlag("input.project.cloud.group", scaffoldCmd.Flags().Lookup("cloudgroup"))

	// configure the project framework settings
	viper.BindPFlag("input.project.framework.type", scaffoldCmd.Flags().Lookup("framework"))
	viper.BindPFlag("input.project.framework.option", scaffoldCmd.Flags().Lookup("frameworkoption"))
	viper.BindPFlag("input.project.framework.version", scaffoldCmd.Flags().Lookup("frameworkversion"))

	// -- bind the framework properties to the project framework
	viper.BindPFlag("input.project.framework.properties", scaffoldCmd.Flags().Lookup("frameworkprops"))

	viper.BindPFlag("input.settingsfile", scaffoldCmd.Flags().Lookup("settingsfile"))

	viper.BindPFlag("input.pipeline", scaffoldCmd.Flags().Lookup("pipeline"))

	viper.BindPFlag("input.cloud.platform", scaffoldCmd.Flags().Lookup("cloud"))

	viper.BindPFlag("input.business.component", scaffoldCmd.Flags().Lookup("component"))

	viper.BindPFlag("input.network.base.domain.external", scaffoldCmd.Flags().Lookup("domain"))
	viper.BindPFlag("input.network.base.domain.internal", scaffoldCmd.Flags().Lookup("internaldomain"))

	viper.BindPFlag("input.directory.cache", scaffoldCmd.Flags().Lookup("cachedir"))

	viper.BindPFlag("input.overrides.ado_variables_path", scaffoldCmd.Flags().Lookup("adovariables"))

	viper.BindPFlag("input.options.cmdlog", scaffoldCmd.Flags().Lookup("cmdlog"))
	viper.BindPFlag("input.options.save", scaffoldCmd.Flags().Lookup("save"))
	viper.BindPFlag("input.options.nocleanup", scaffoldCmd.Flags().Lookup("nocleanup"))
	viper.BindPFlag("input.options.force", scaffoldCmd.Flags().Lookup("force"))
	viper.BindPFlag("input.options.norun", scaffoldCmd.Flags().Lookup("norun"))
}

// ScaffoldOverrides updates the main configuration with any override files that have been specified on
// the command line.
// It is not set as a prerun function because it has to be read in at the appropriate point in the root command
func ScaffoldOverrides() {

	override_ado_variables := viper.GetString("input.overrides.ado_variables_path")
	if override_ado_variables != "" {
		data, err := os.ReadFile(override_ado_variables)
		if err != nil {
			log.Fatalf("Unable to read in ADO variables override file (%s): %s", err.Error(), override_ado_variables)
			App.Logger.Exit(3)
		}

		staticFiles.Ado_Variable_Template_Tmpl = string(data)
	}
}

func executeScaffoldRun(ccmd *cobra.Command, args []string) {

	// Call the scaffolding method
	scaff := scaffold.New(&Config, App.Logger)
	err := scaff.Run()
	if err != nil {
		App.Logger.Fatalf("Error running scaffold: %s", err.Error())
	}
}
