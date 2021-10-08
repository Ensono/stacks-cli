package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/amido/stacks-cli/internal/config/static"
	"github.com/amido/stacks-cli/internal/util"
	"github.com/amido/stacks-cli/pkg/config"
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

	// - options
	var cmdlog bool
	var dryrun bool

	// - project settings
	var project_name string
	var project_vcs_type string
	var project_vcs_url string
	var project_vcs_ref string
	var project_settings_file string

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
	var network_base_domain_external string
	var network_base_domain_internal string

	// Add the run command to the root
	rootCmd.AddCommand(scaffoldCmd)

	// Configure the flags
	// - run interactively
	scaffoldCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Run in interactive mode")

	scaffoldCmd.Flags().StringVarP(&project_name, "name", "n", "", "Name of the project to create")
	scaffoldCmd.Flags().StringVar(&project_vcs_type, "sourcecontrol", "github", "Type of source control being used")
	scaffoldCmd.Flags().StringVarP(&project_vcs_url, "sourcecontrolurl", "u", "", "Url of the remote for source control")
	scaffoldCmd.Flags().StringVar(&project_vcs_ref, "sourcecontrolref", "", "SHA reference or Tag to use to clone repo from")
	scaffoldCmd.Flags().StringVar(&project_settings_file, "projectsettingsfile", "", "Path to a settings file to use for the project")

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

	scaffoldCmd.Flags().StringVarP(&network_base_domain_external, "domain", "d", "", "External domain for the app")
	scaffoldCmd.Flags().StringVar(&network_base_domain_internal, "internaldomain", "", "Internal domain for the app")

	scaffoldCmd.Flags().BoolVar(&cmdlog, "cmdlog", false, "Specify if commands should be logged")
	scaffoldCmd.Flags().BoolVar(&dryrun, "dryrun", false, "Perform a dryrun of the CLI. No changes will be made on disk")

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
	viper.BindPFlag("project.settingsfile", scaffoldCmd.Flags().Lookup("projectsettingsfile"))
	viper.BindPFlag("project.cloud.region", scaffoldCmd.Flags().Lookup("cloudregion"))
	viper.BindPFlag("project.cloud.group", scaffoldCmd.Flags().Lookup("cloudgroup"))

	viper.BindPFlag("settingsfile", scaffoldCmd.Flags().Lookup("settingsfile"))

	viper.BindPFlag("pipeline", scaffoldCmd.Flags().Lookup("pipeline"))

	viper.BindPFlag("cloud.platform", scaffoldCmd.Flags().Lookup("cloud"))

	viper.BindPFlag("business.company", scaffoldCmd.Flags().Lookup("company"))
	viper.BindPFlag("business.domain", scaffoldCmd.Flags().Lookup("area"))
	viper.BindPFlag("business.component", scaffoldCmd.Flags().Lookup("component"))

	viper.BindPFlag("terraform.backend.storage", scaffoldCmd.Flags().Lookup("tfstorage"))
	viper.BindPFlag("terraform.backend.group", scaffoldCmd.Flags().Lookup("tfgroup"))
	viper.BindPFlag("terraform.backend.container", scaffoldCmd.Flags().Lookup("tfcontainer"))

	viper.BindPFlag("network.base.domain.external", scaffoldCmd.Flags().Lookup("domain"))
	viper.BindPFlag("network.base.domain.internal", scaffoldCmd.Flags().Lookup("internaldomain"))

	viper.BindPFlag("options.cmdlog", scaffoldCmd.Flags().Lookup("cmdlog"))
	viper.BindPFlag("options.dryrun", scaffoldCmd.Flags().Lookup("dryrun"))

	viper.BindPFlag("interactive", scaffoldCmd.Flags().Lookup("interactive"))
}

func executeRun(ccmd *cobra.Command, args []string) {

	// determine if the interactive option has been set
	// if it has ask the user for input and then overwrite the configuration
	// that has been specified on the command line with the values as given
	// by the user
	answers := config.Answers{}
	err := answers.RunInteractive(&Config)
	if err != nil {
		App.Logger.Fatalf("Unable to perform interactive configuration: %s\n", err.Error())
	}

	// ensure that at least one project has been specified
	if len(Config.Input.Project) == 1 && Config.Input.Project[0].Name == "" {
		App.Logger.Fatalln("No projects have been defined")
	}

	// check that the specified pipeline is supported by the CLI
	pipeline := config.Pipeline{}
	if !pipeline.IsSupported(Config.Input.Pipeline) {
		App.Logger.Fatalf("Specified pipeline is not supported: %s %v\n", Config.Input.Pipeline, pipeline.GetSupported())
	}

	// Call the scaffolding method
	scaff := scaffold.New(&Config, App.Logger)
	err = scaff.Run()
	if err != nil {
		App.Logger.Fatalf("Error running scaffold: %s", err.Error())
	}

	// call the func to set default values in the object
	Config.SetDefaultValues()

	// iterate around the projects that have been specified
	for _, project := range Config.Input.Project {

		App.Logger.Infof("Setting up project: %s\n", project.Name)

		// Create the temporary and working directories for the current project
		project.Directory.TempDir = filepath.Join(Config.Input.Directory.TempDir, project.Name)
		project.Directory.WorkingDir = filepath.Join(Config.Input.Directory.WorkingDir, project.Name)

		// Clone the target repository into the temp directory
		key := project.Framework.GetMapKey()
		err := util.GitClone(
			Config.Input.Stacks.GetSrcURL(key),
			project.SourceControl.Ref,
			project.Directory.TempDir,
		)

		// if an error occured getting the code then show an error and move onto the next project
		if err != nil {
			App.Logger.Errorf("Error downloading code: %s", err.Error())
			continue
		}

		// Read in the configuration file for the project
		err = project.ReadSettings(project.Directory.TempDir, &Config)

		if err != nil {
			App.Logger.Errorf("Error reading settings from project settings: %s", err.Error())
			continue
		}

		// iterate around the init settings for the project
		// all of these operations will occur in the temporary directory of the project
		for _, op := range project.Settings.Init.Operations {

			// log information based on the description in the settings file
			App.Logger.Info(op.Description)

			err = scaff.PerformOperation(op, &Config, &project, project.Directory.TempDir)

			if err != nil {
				App.HandleError(err, "error", "Issue encountered performing init operation")
				break
			}
		}

		// iterate around the setup settings for the project
		// all of these operations will occur in the working directory of the project
		for _, op := range project.Settings.Setup.Operations {

			// log information based on the description in the settings file
			App.Logger.Info(op.Description)

			err = scaff.PerformOperation(op, &Config, &project, project.Directory.WorkingDir)

			if err != nil {
				App.HandleError(err, "error", "Issue encountered performing setup operation")
				break
			}
		}

		// replace the variable file in the working version of the project with the values
		// as specified in the CLI
		// use the path to the variable file from the stackscli based on the type of pipeline
		// being targeted

		// get the pipeline settings from the project settings file
		pipelineSettings := project.Settings.GetPipeline(Config.Input.Pipeline)

		// define a replacements object so that all can be passed to the render function
		// the project is passed in as a separate object as it is part of a slice
		replacements := config.Replacements{}
		replacements.Input = Config.Input
		replacements.Project = project

		if !Config.Input.Options.DryRun {
			msg, err := Config.WriteVariablesFile(&project, pipelineSettings, replacements)
			if err == nil {
				App.Logger.Info("Created pipeline variable file")
			} else {
				App.HandleError(err, "error", msg)
			}
		}

		// Replace patterns in the build file
		errs := pipelineSettings.ReplacePatterns(project.Directory.WorkingDir)
		if len(errs) > 0 {
			for _, err := range errs {
				App.HandleError(err, "error", "Issue performing replacements")
			}
		}

		// Initialise the working dir as a git repository
		// Iterate around the git commands and use the PerformOperation function so that the
		// commands get parsed by the template system
		App.Logger.Info("Configuring source control for project")
		for _, command := range static.GitCmds {

			// split the command string into cmd, args so that the operation model can be configured
			commandParts := strings.SplitN(command, " ", 2)

			op := config.Operation{
				Action:    "cmd",
				Command:   commandParts[0],
				Arguments: commandParts[1],
			}

			err = scaff.PerformOperation(op, &Config, &project, project.Directory.WorkingDir)

			if err != nil {
				App.HandleError(err, "error", "Issue encountered configuring project as git repository")
				break
			}

		}
	}

	// Output information about the run
	if Config.Input.Options.DryRun {
		App.Logger.Warnf("CLI was run with the --dryrun option, no projects have been configured")
	}

	if Config.Input.Options.CmdLog {
		App.Logger.Infof("Command log has been created: %s\n", Config.Self.CmdLogPath)
	}

	// Perform cleanup by removing the temporary directory
	App.Logger.Info("Performing cleanup")
	App.Logger.Infof(" - removing temporary directory: %s\n", Config.Input.Directory.TempDir)
	err = os.RemoveAll(Config.Input.Directory.TempDir)
	if err != nil {
		App.Logger.Fatalf("Unable to remove temporary directory: %s", Config.Input.Directory.TempDir)
	}

}
