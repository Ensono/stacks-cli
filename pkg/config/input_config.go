package config

import (
	"os/exec"

	"github.com/amido/stacks-cli/internal/config/static"
	"github.com/amido/stacks-cli/internal/models"
	"github.com/amido/stacks-cli/internal/util"
)

// Config is used to map the configuration onto the application models
type InputConfig struct {

	// State if running in Interactive mode
	Interactive bool `mapstructure:"interactive"`

	// Version of the application
	Version string

	// Define the logging parameters
	Log models.Log `mapstructure:"log"`

	Directory Directory `mapstructure:"directory"`

	Business     Business  `mapstructure:"business"`
	Cloud        Cloud     `mapstructure:"cloud"`
	Network      Network   `mapstructure:"network"`
	Pipeline     string    `mapstructure:"pipeline"`
	Project      []Project `mapstructure:"project"`
	Stacks       Stacks    `mapstructure:"stacks"` // Holds the information about the projects in stacks
	Terraform    Terraform `mapstructure:"terraform"`
	SettingsFile string    `mapstructure:"settingsfile" json:",omitempty"`
	Options      Options   `mapstructure:"options"`
}

// CheckFrameworks iterates around each of the projects and builds up a list of the frameworks
// that have been specified. It will then check that each of the framework binaries
// are present in the path.
// If there are not then the ones that are not present are returned to the calling function
func (ic *InputConfig) CheckFrameworks() []models.Command {

	var frameworkTypes []string
	var missing []models.Command

	// iterate around the projects
	// if the framework does not already exist in the slice, check if it the executable
	// exists in the path
	for _, project := range ic.Project {

		// add the framework type to the frameworks if not already present
		if !util.SliceContains(frameworkTypes, project.Framework.Type) {
			frameworkTypes = append(frameworkTypes, project.Framework.Type)

			// get the binaries for this framework type
			binaries := static.FrameworkCommand(project.Framework.Type)

			for _, binary := range binaries {
				// create a command object
				command := models.Command{
					Framework: project.Framework.Type,
					Binary:    binary,
				}

				// if the binary is null then the framework has not been specified properly so
				// add the command to the missing list
				// otherwise check that the binary exists in the path
				if binary == "" {
					missing = append(missing, command)
				} else {

					// determine if the binary is in the path
					_, err := exec.LookPath(command.Binary)

					// if there is an error then the command cannot be found in the path, so
					// add it to the missing slice
					if err != nil {
						missing = append(missing, command)
					}
				}
			}
		}

	}

	return missing
}
