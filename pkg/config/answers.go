package config

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/amido/stacks-cli/internal/config/static"
	"github.com/amido/stacks-cli/internal/util"
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
}

// ProjectAnswers are the list of answers that are provided for each project defined
// on the command line
type ProjectAnswers struct {
	Name              string `survey:"name"`
	FrameworkType     string `survey:"framework_type"`
	FrameworkOption   string `survey:"framework_option"`
	FrameworkVersion  string `survey:"framework_version"`
	PlatformType      string `survey:"platform_type"`
	SourceControlType string `survey:"source_control_type"`
	SourceControlUrl  string `survey:"source_control_url"`
	CloudRegion       string `survey:"cloud_region"`
	CloudGroup        string `survey:"cloud_group"`
}

// getCoreQuestions returns the list of questions that need to be answered in interactive
// mode on the command line. This works in conjunction with the answers object
func (a *Answers) getCoreQuestions() []*survey.Question {

	// get the default working directory
	workingDir := util.GetDefaultWorkingDir()

	var questions = []*survey.Question{
		{
			Name: "company_name",
			Prompt: &survey.Input{
				Message: "What is the name of your company?",
			},
			Validate: survey.Required,
		},
		{
			Name: "company_domain",
			Prompt: &survey.Input{
				Message: "What is the scope or area of the company?",
			},
			Validate: survey.Required,
		},
		{
			Name: "company_component",
			Prompt: &survey.Input{
				Message: "What component is being worked on?",
			},
			Validate: survey.Required,
		},
		{
			Name: "pipeline",
			Prompt: &survey.Select{
				Message: "What pipeline is being targeted?",
				Options: []string{"azdo"},
				Default: "azdo",
			},
			Validate: survey.Required,
		},
		{
			Name: "cloud_platform",
			Prompt: &survey.Select{
				Message: "Which cloud is Stacks being setup in?",
				Options: []string{"azure"},
				Default: "azure",
			},
			Validate: survey.Required,
		},
		{
			Name: "terraform_group",
			Prompt: &survey.Input{
				Message: "Which group is the Terraform state being saved in?",
			},
			Validate: survey.Required,
		},
		{
			Name: "terraform_storage",
			Prompt: &survey.Input{
				Message: "What is the name of the Terraform storage?",
			},
			Validate: survey.Required,
		},
		{
			Name: "terraform_container",
			Prompt: &survey.Input{
				Message: "What is the name of the Terraform storage container?",
			},
			Validate: survey.Required,
		},
		{
			Name: "network_domain_external",
			Prompt: &survey.Input{
				Message: "What is the external domain of the solution?",
			},
			Validate: survey.Required,
		},
		{
			Name: "network_domain_internal",
			Prompt: &survey.Input{
				Message: "What is the internal domain of the solution?",
				Help:    "If the internal domain is not set it will be derived from the specified external domain",
			},
		},
		{
			Name: "options",
			Prompt: &survey.MultiSelect{
				Message: "What options would you like to enable, if any?",
				Options: []string{"Command Log", "Dry Run"},
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
			},
			Validate: survey.Required,
		},
	}

	return questions
}

func (a *Answers) getProjectQuestions() []*survey.Question {

	var questions = []*survey.Question{
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
				Options: []string{"dotnet", "java"},
			},
			Validate: survey.Required,
		},
		{
			Name: "framework_option",
			Prompt: &survey.Select{
				Message: "Which style of the framework do you require?",
				Options: []string{"webapi", "cqrs", "events"},
			},
			Validate: survey.Required,
		},
		{
			Name: "framework_version",
			Prompt: &survey.Input{
				Message: "Which version of the framework option do you require?",
				Default: "latest",
			},
			Validate: survey.Required,
		},
		{
			Name: "platform_type",
			Prompt: &survey.Select{
				Message: "What platform is being used",
				Options: []string{"aks"},
				Default: "aks",
			},
			Validate: survey.Required,
		},
		{
			Name: "source_control_type",
			Prompt: &survey.Select{
				Message: "Please select the source control system being used",
				Options: []string{"github"},
				Default: "github",
			},
			Validate: survey.Required,
		},
		{
			Name: "source_control_url",
			Prompt: &survey.Input{
				Message: "What is the URL of the remote repository?",
			},
			Validate: survey.Required,
		},
		{
			Name: "cloud_region",
			Prompt: &survey.Input{
				Message: "Which cloud region should be used?",
			},
			Validate: survey.Required,
		},
		{
			Name: "cloud_group",
			Prompt: &survey.Input{
				Message: "What is the name of the group for all the resources?",
			},
			Validate: survey.Required,
		},
	}

	return questions
}

func (a *Answers) RunInteractive(config *Config) error {

	var err error

	// return without performing any tasks if interactive mode is not
	// enabled
	if !config.Input.Interactive {
		return err
	}

	// Output the banner to the screen
	fmt.Printf(static.Banner)

	// ask the questions
	err = survey.Ask(a.getCoreQuestions(), a)
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
		pa := ProjectAnswers{}
		err = survey.Ask(a.getProjectQuestions(), &pa)
		if err != nil {
			continue
		}

		// create a struct for the project
		project := Project{
			Name: pa.Name,
			Framework: Framework{
				Type:    pa.FrameworkType,
				Option:  pa.FrameworkOption,
				Version: pa.FrameworkVersion,
			},
			SourceControl: SourceControl{
				Type: pa.SourceControlType,
				URL:  pa.SourceControlUrl,
			},
			Cloud: Cloud{
				Region:        pa.CloudRegion,
				ResourceGroup: pa.CloudGroup,
			},
			SettingsFile: settingsFile,
		}

		// append this to the project list on the config object
		// if i == 0 {
		// 	config.Input.Project[0] = project
		// } else {
		projectList = append(projectList, project)
		// }
	}

	config.Input.Project = projectList

	return err
}
