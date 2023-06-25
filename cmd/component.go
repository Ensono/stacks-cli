package cmd

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/amido/stacks-cli/internal/constants"
	"github.com/amido/stacks-cli/internal/util"
	"github.com/amido/stacks-cli/pkg/config"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	componentCmd = &cobra.Command{
		Use:   "component",
		Short: "Perform operations on the named component",
		Long:  "",
	}

	foldersCmd = &cobra.Command{
		Use:   "folders [list of projects]",
		Short: "Display the folders for a command, usually used with a mono-repo",
		Long:  "",
		Run:   executeComponentFolders,
		Args:  cobra.MinimumNArgs(1),
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List the known components",
		Long: `The components that the CLI can interact with is built it at build time or can be overridden at runtime.
		This command displays the known components, from either source, along with the respective URL.`,
		Run: executeComponentList,
	}
)

func init() {

	// Add the command to the rootCmd
	rootCmd.AddCommand(componentCmd)

	// Add the subcommands for the component command
	componentCmd.AddCommand(foldersCmd, listCmd)

}

// executeComponentList outputs the list of stacks.components
func executeComponentList(ccmd *cobra.Command, args []string) {

	// create a table object
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Group", "Name", "Type", "Name / URL", "Version / ID"})

	// Iterate around the components and display the details
	for _, component := range Config.Stacks.Components {

		name_url := ""
		version_id := ""

		switch component.Package.Type {
		case "nuget":
			name_url = component.Package.Name
			version_id = component.Package.ID
		case "git":
			name_url = component.Package.URL
			version_id = component.Package.Version
		}

		// Append a row for each of the components
		t.AppendRow(
			table.Row{
				component.Group, component.Name, component.Package.Type, name_url, version_id,
			},
		)
	}

	// Sort the data in the table
	t.SortBy([]table.SortBy{
		{Name: "Group", Mode: table.Asc},
		{Name: "Name", Mode: table.Asc},
	})

	// Group titles
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
	})

	t.SetStyle(table.StyleLight)

	t.Render()
}

// execute Component folders retrieves the stackscli.yml file from the remote
// repository and analyses it to display a list of folders that are supported
func executeComponentFolders(ccmd *cobra.Command, args []string) {

	// define a slice to hold the valid entries
	var valid []string

	// Iterate around each of the args and ensure that they ar valid
	// this is done so that all the errors about any invalid projects are shown in one place
	for _, name := range args {
		_, ok := Config.Stacks.Components[name]
		if ok {
			valid = append(valid, name)
		} else {
			App.Logger.Warnf("Invalid component: %s", name)
		}
	}

	// Create a new table writer to display the output
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Component", "Folders"})

	// iterate around the valid projects and get the settings for the project
	for _, name := range valid {

		projectSettings := config.Settings{}

		// get the component details to work with
		component := Config.Stacks.Components[name]

		// select the download method based on the type
		switch component.Package.Type {
		case "git":

			// Not using the downloader here as only getting one file, will use the raw endpoint
			// Parse the URL to get the component parts
			fileURL, err := url.Parse(component.Package.URL)
			if err != nil {
				App.Logger.Errorf("Unable to parse package URL: %s", component.Package.URL)
				continue
			}

			App.Logger.Infof("[%s] Retrieving project settings", name)

			// Create the URL to the raw file
			downloadUrl := fmt.Sprintf("https://raw.githubusercontent.com%s/%s/%s", fileURL.Path, component.Package.Version, constants.SettingsFile)

			App.Logger.Debugf("Download URL: %s", downloadUrl)

			data, err := util.DownloadFile(downloadUrl)
			if err != nil {
				App.Logger.Errorf("Error downloading settings file: %s", err.Error())
				continue
			}

			v := viper.New()
			v.SetConfigType("yaml")
			err = v.ReadConfig(bytes.NewBuffer(data))
			if err != nil {
				App.Logger.Errorf("Error reading settings file: %s", err.Error())
				continue
			}

			// unmarshal the data into an object
			err = v.Unmarshal(&projectSettings)
			if err != nil {
				App.Logger.Errorf("Error in settings file: %s", err.Error())
				continue
			}
		}

		// specify the cell info for the folders
		cellText := ""
		if len(projectSettings.Folders) > 0 {
			cellText = strings.Join(projectSettings.Folders, "\n")
		} else {
			cellText = "<NONE_CONFIGURED>"
		}

		t.AppendRow(
			table.Row{
				name, cellText,
			},
		)

		t.AppendSeparator()
	}

	// Output the table
	t.SetStyle(table.StyleLight)

	t.Render()
}
