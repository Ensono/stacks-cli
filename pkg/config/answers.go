package config

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Ensono/stacks-cli/internal/util"
)

// Answers is the object that results from the command line, when run in interactive mode,
// are added to.
type Answers struct {
	CompanyName           string   `survey:"company_name"`
	CompanyDomain         string   `survey:"company_domain"`
	CompanyComponent      string   `survey:"company_component"`
	Pipeline              string   `survey:"pipeline"`
	CloudPlatform         string   `survey:"cloud_platform"`
	TerraformGroup        string   `survey:"terraform_group"`
	TerraformStorage      string   `survey:"terraform_storage"`
	TerraformContainer    string   `survey:"terraform_container"`
	NetworkDomainExternal string   `survey:"network_domain_external"`
	NetworkDomainInternal string   `survey:"network_domain_internal"`
	Options               []string `survey:"options"`
	WorkingDir            string   `survey:"working_dir"`
	ProjectCount          int      `survey:"project_count"`
	EnvironmentCount      int      `survey:"environment_count"`
}

// ProjectAnswers are the list of answers that are provided for each project defined
// on the command line
type ProjectAnswers struct {
	Name                    string `survey:"name"`
	FrameworkType           string `survey:"framework_type"`
	FrameworkOption         string `survey:"framework_option"`
	FrameworkVersion        string `survey:"framework_version"`
	FrameworkProperties     string `survey:"framework_properties"`
	FrameworkDeploymentMode string `survey:"framework_deployment_mode"`
	PlatformType            string `survey:"platform_type"`
	SourceControlType       string `survey:"source_control_type"`
	SourceControlUrl        string `survey:"source_control_url"`
	CloudRegion             string `survey:"cloud_region"`
	CloudGroup              string `survey:"cloud_group"`
}

type EnvironmentAnswers struct {
	Name      string `survey:"name"`
	Type      string `survey:"type"`
	DependsOn string `survey:"dependson"`
}

// getCoreQuestions returns the list of questions that need to be answered in interactive
// mode on the command line. This works in conjunction with the answers object
func (a *Answers) getCoreQuestions(config *Config) []*survey.Question {

	// get the default working directory
	workingDir := util.GetDefaultWorkingDir()

	var questions = []*survey.Question{
		{
			Name: "company_name",
			Prompt: &survey.Input{
				Message: "What is the name of your company?",
				Help:    "The name of your company that is used to help name resources in the chosen cloud platform",
				Default: config.Input.Business.Company,
			},
			Validate: survey.Required,
		},
		{
			Name: "company_domain",
			Prompt: &survey.Input{
				Message: "What is the scope or area of the company?",
				Help:    "As there can be many different projects being created within a team, this value is designed to be used to categorise the project",
				Default: config.Input.Business.Domain,
			},
			Validate: survey.Required,
		},
		{
			Name: "company_component",
			Prompt: &survey.Input{
				Message: "What component is being worked on?",
				Help:    "When using Ensono Stacks the scaffolded project is usually being created as part of a larger solution. This value is used to help name the project",
			},
			Validate: survey.Required,
		},
		{
			Name: "pipeline",
			Prompt: &survey.Select{
				Message: "What pipeline is being targeted?",
				Options: []string{"azdo", "gha"},
				Default: "azdo",
				Help:    "A number of different pipelines can be configured to work with Ensono Stacks. Select the pipeline that this project will use.",
			},
			Validate: survey.Required,
		},
		{
			Name: "cloud_platform",
			Prompt: &survey.Select{
				Message: "Which cloud is Stacks being setup in?",
				Options: []string{"aws", "azure"},
				Default: "azure",
				Help:    "Ensono Stacks supports being deployed into different cloud platforms. Select the cloud that will be used.",
			},
			Validate: survey.Required,
		},
		{
			Name: "terraform_group",
			Prompt: &survey.Input{
				Message: "Which group is the Terraform state being saved in?",
				Help:    "This is the group that contains the storage account for the Terraform state.",
				Default: config.Input.Terraform.Backend.Group,
			},
			Validate: survey.Required,
		},
		{
			Name: "terraform_storage",
			Prompt: &survey.Input{
				Message: "What is the name of the Terraform storage?",
				Help:    "The name of the Azure Storage account or S3 bucket being used to store the Terraform state",
				Default: config.Input.Terraform.Backend.Storage,
			},
			Validate: survey.Required,
		},
		{
			Name: "terraform_container",
			Prompt: &survey.Input{
				Message: "What is the name of the folder or container for Terraform state storage?",
				Help:    "For Azure storage accounts this is the name of the container to be used for the workspace. For AWS S3 Buckets this is the name of the folder to use.",
				Default: config.Input.Terraform.Backend.Container,
			},
			Validate: survey.Required,
		},
		{
			Name: "network_domain_external",
			Prompt: &survey.Input{
				Message: "What is the external domain of the solution?",
				Help:    "The domain that the projects will be accessible on. This needs to be a domain that is already registered and can be used.",
			},
			Validate: survey.Required,
		},
		{
			Name: "network_domain_internal",
			Prompt: &survey.Input{
				Message: "What is the internal domain of the solution?",
				Help:    "The internal domain of the application. If not specified this will be derived from the external domain. So if stacks.ensono.com is used then the internal will be stacks.ensono.internal",
			},
		},
		{
			Name: "options",
			Prompt: &survey.MultiSelect{
				Message: "What options would you like to enable, if any?",
				Options: []string{"Command Log", "Dry Run"},
				Help:    "A selection of options to apply to the CLI. Command Log - save all of the commands that are run to a log file; DryRun - run the command without performing any of the tasks.",
			},
		},
		{
			Name: "working_dir",
			Prompt: &survey.Input{
				Message: "Please specify the working directory for the projects?",
				Help:    "All of the projects that are generated by the CLI will be stored in this directory in a folder based on the project name",
				Default: workingDir,
			},
		},
		{
			Name: "project_count",
			Prompt: &survey.Input{
				Message: "How many projects would you like to configure?",
				Default: "1",
				Help:    "The CLI is designed to be able to scaffold multiple projects at once. By providing the number of projects, you will be asked a series of questions for each project",
			},
			Validate: survey.Required,
		},
		{
			Name: "environment_count",
			Prompt: &survey.Input{
				Message: "How many environments would you like to configure?",
				Default: "0",
			},
			Validate: survey.Required,
		},
	}

	return questions
}

func (a *Answers) getProjectQuestions(qType string, config *Config) []*survey.Question {

	var questions []*survey.Question

	// use the qType to determine which questions need to be asked
	switch qType {
	case "pre":
		questions = []*survey.Question{
			{
				Name: "name",
				Prompt: &survey.Input{
					Message: "What is the project name?",
				},
				Validate: survey.Required,
			},
			{
				Name: "framework_type",
				Prompt: &survey.Select{
					Message: "What framework should be used for the project?",
					Options: config.Stacks.GetComponentNames(), // []string{"dotnet", "java", "nx", "infra"},
					Help:    "Ensono Stacks supports a number of different frameworks. Select the framework that you would like to use. ",
				},
				Validate: survey.Required,
			},
		}
	case "java":

		// define the list of options for each of the different languages
		// options := []string{"webapi", "cqrs"}
		// if qType == "java" {
		// 	options = append(options, "events")
		// }

		questions = []*survey.Question{
			{
				Name: "framework_option",
				Prompt: &survey.Select{
					Message: "Which option of the framework do you require?",
					Options: config.Stacks.GetComponentOptions(qType),
					Help:    "Ensono Stacks has a number of options that are available for specific frameworks. Select the one that is appropriate for the desired workload.",
				},
				Validate: survey.Required,
			},
			{
				Name: "framework_properties",
				Prompt: &survey.Input{
					Message: "Specify any additional framework properties. (Use a comma to separate each one).",
					Default: "",
					Help:    "Additional properties that need to applied to the project when it is built. This is dependent on the framework that has been chosen. Multiple options can be specified by separating the options with a comma",
				},
			},
		}
	case "dotnet":

		// define the list of options for each of the different languages
		// options := []string{"webapi", "cqrs"}
		// if qType == "java" {
		// 	options = append(options, "events")
		// }

		questions = []*survey.Question{
			{
				Name: "framework_option",
				Prompt: &survey.Select{
					Message: "Which option of the framework do you require?",
					Options: config.Stacks.GetComponentOptions(qType),
					Help:    "Ensono Stacks has a number of options that are available for specific frameworks. Select the one that is appropriate for the desired workload.",
				},
				Validate: survey.Required,
			},
			{
				Name: "framework_properties",
				Prompt: &survey.Input{
					Message: "Specify any additional framework properties. (Use a comma to separate each one).",
					Default: "",
					Help:    "Additional properties that need to applied to the project when it is built. This is dependent on the framework that has been chosen. Multiple options can be specified by separating the options with a comma",
				},
			},
			{
				Name: "framework_deployment",
				Prompt: &survey.Input{
					Message: "What type of environment will this be deployed to?",
					Default: "AKS",
					Help:    "Containers can be deployed to an AKS cluster or ACA environment",
				},
			},
		}
	case "nx":
		// options := []string{"next", "apps"}

		questions = []*survey.Question{
			{
				Name: "framework_option",
				Prompt: &survey.Select{
					Message: "Which option of the framework do you require?",
					Options: config.Stacks.GetComponentOptions(qType),
					Description: func(value string, index int) string {
						if value == "next" {
							return "Stacks Workspace with NextJs"
						}
						if value == "apps" {
							return "Empty Stacks Nx Workspace"
						}
						return ""
					},
				},
				Validate: survey.Required,
			},
			{
				Name: "framework_properties",
				Prompt: &survey.Input{
					Message: "Specify any additional framework properties. (Use a comma to separate each one).",
					Default: "",
				},
			},
		}
	case "infra":
		questions = []*survey.Question{
			{
				Name: "framework_option",
				Prompt: &survey.Select{
					Message: "Which type of infrastructure is required?",
					Options: config.Stacks.GetComponentOptions(qType),
					Default: "aks",
					Help:    "A number of projects support different infrastructure. By answering this question, the CLI will prepare the project, if applicable, to the chosen cloud.",
				},
				Validate: survey.Required,
			},
		}
	case "post":
		questions = []*survey.Question{
			{
				Name: "framework_version",
				Prompt: &survey.Input{
					Message: "Which version of the framework option do you require?",
					Default: "latest",
					Help:    "There are a number of different versions of the frameworks that can be used. Specify the one that is required",
				},
				Validate: survey.Required,
			},
			{
				Name: "source_control_type",
				Prompt: &survey.Select{
					Message: "Please select the source control system being used",
					Options: []string{"github"},
					Default: "github",
					Help:    "This is the centralised source control that should be used",
				},
				Validate: survey.Required,
			},
			{
				Name: "source_control_url",
				Prompt: &survey.Input{
					Message: "What is the URL of the remote repository?",
					Help:    "When the project is scaffolded and configured as a Git repo, it will add in the origin to this URL.",
				},
				Validate: survey.Required,
			},
			{
				Name: "cloud_region",
				Prompt: &survey.Input{
					Message: "Which cloud region should be used?",
					Help:    "The region, of the chosen cloud, that the resources will be deployed to.",
				},
				Validate: survey.Required,
			},
			{
				Name: "cloud_group",
				Prompt: &survey.Input{
					Message: "What is the name of the group for all the resources?",
					Help:    "The name of the group into which all the resources for this project will be deployed into",
				},
				Validate: survey.Required,
			},
		}
	case "data":
		questions = []*survey.Question{
			{
				Name: "framework_option",
				Prompt: &survey.Select{
					Message: "Which option of the framework do you require?",
					Options: config.Stacks.GetComponentOptions(qType),
					Default: "data",
				},
			},
		}
	}

	return questions
}

func (a *Answers) getEnvironmentQuestions(qType string, config *Config) []*survey.Question {

	var questions []*survey.Question

	// use the qType to determine which questions need to be asked
	switch qType {
	case "pre":
		questions = []*survey.Question{
			{
				Name: "name",
				Prompt: &survey.Input{
					Message: "What is the environment name?",
				},
				Validate: survey.Required,
			},
			{
				Name: "type",
				Prompt: &survey.Select{
					Message: "What is the type of environment?",
					Options: []string{"Development", "Testing", "Production"},
					Help: "Development: Deployed from a feature branch and is the first environment we deploy to.\n" +
						"Testing: Deployed from the main/master branch. \n" +
						"Production: Deployed from the main/master branch after it has been deployed to the 'Testing' environment(s)",
				},
				Validate: survey.Required,
			},
			{
				Name: "dependson",
				Prompt: &survey.Input{
					Message: "What environment does this environment depend on?",
					Help:    "This will determine the deployment order of the environments. You can provide multiple environments by seperating each environment with a comma e.g. Dev,Test",
				},
			},
		}
	}
	return questions
}

func (a *Answers) RunInteractive(config *Config) error {

	var err error

	// ask the questions
	err = survey.Ask(a.getCoreQuestions(config), a)
	if err != nil {
		return err
	}

	// add the information that has been gleaned from the survey to the
	// config object
	config.Input.Business.Company = a.CompanyName
	config.Input.Business.Domain = a.CompanyDomain
	config.Input.Business.Component = a.CompanyComponent
	config.Input.Pipeline = a.Pipeline
	config.Input.Cloud.Platform = a.CloudPlatform
	config.Input.Directory.WorkingDir = a.WorkingDir
	config.Input.Terraform.Backend.Group = a.TerraformGroup
	config.Input.Terraform.Backend.Storage = a.TerraformStorage
	config.Input.Terraform.Backend.Container = a.TerraformContainer
	config.Input.Network.Base.Domain.External = a.NetworkDomainExternal
	config.Input.Network.Base.Domain.Internal = a.NetworkDomainInternal

	if util.SliceContains(a.Options, "Command Log") {
		config.Input.Options.CmdLog = true
	}

	if util.SliceContains(a.Options, "Dry Run") {
		config.Input.Options.DryRun = true
	}

	projectList := []Project{}

	// Cobra and Viper will populate the projects slice with a single project based on
	// the input from the config or command line. If a project settings file has been
	// specified, it will be added to this first item
	// get the settings file path from the first project in the slice as the slice will
	// be overwritten
	settingsFile := config.Input.Project[0].SettingsFile

	// as a number of projects can be configured, the project questions need
	// to be asked project_count times
	for i := 0; i < a.ProjectCount; i++ {

		fmt.Printf("\nConfiguring project: %d\n", i+1)

		// ask the project questions
		// this is done in 3 stages so that the different options of the framework can be
		// modified based on the main framework that has been specified
		pa := ProjectAnswers{}

		// - pre questions
		err = survey.Ask(a.getProjectQuestions("pre", config), &pa)
		if err != nil {
			continue
		}

		// - specific questions based on the framework selected
		err = survey.Ask(a.getProjectQuestions(pa.FrameworkType, config), &pa)
		if err != nil {
			continue
		}

		// - post questions
		err = survey.Ask(a.getProjectQuestions("post", config), &pa)
		if err != nil {
			continue
		}

		// check to see if any properties have been specified, and if they have
		// create a properties object to work with the
		properties := []string{}
		if pa.FrameworkProperties != "" {

			// split the properties based comma and then iterate around setting the framework properties
			properties = strings.Split(pa.FrameworkProperties, ",")

		}

		// create a struct for the project
		project := Project{
			Name: pa.Name,
			Framework: Framework{
				Type:           pa.FrameworkType,
				Option:         pa.FrameworkOption,
				Version:        pa.FrameworkVersion,
				Properties:     properties,
				DeploymentMode: pa.FrameworkDeploymentMode,
			},
			SourceControl: SourceControl{
				Type: pa.SourceControlType,
				URL:  pa.SourceControlUrl,
			},
			Cloud: Cloud{
				Region:        pa.CloudRegion,
				ResourceGroup: pa.CloudGroup,
			},
			Platform: Platform{
				Type: pa.PlatformType,
			},
			SettingsFile: settingsFile,
		}

		// append this to the project list on the config object
		projectList = append(projectList, project)

	}

	config.Input.Project = projectList

	environmentList := []Environment{}

	// as a number of environments can be configured, the environments questions need
	// to be asked environments_count times
	for i := 0; i < a.EnvironmentCount; i++ {

		fmt.Printf("\nConfiguring environments: %d\n", i+1)

		// ask the environment questions
		// this is done in 3 stages so that the different options of the framework can be
		// modified based on the main framework that has been specified
		pa := EnvironmentAnswers{}

		// - pre questions
		err = survey.Ask(a.getEnvironmentQuestions("pre", config), &pa)
		if err != nil {
			continue
		}

		dependsOn := []string{}
		if len(pa.DependsOn) > 0 {
			dependsOn = strings.Split(pa.DependsOn, ",")
		}

		// create a struct for the environment
		environment := Environment{
			Name:      pa.Name,
			Type:      pa.Type,
			DependsOn: dependsOn,
		}

		// append this to the environment list on the config object
		environmentList = append(environmentList, environment)
	}

	config.Input.Environment = environmentList

	return err
}
