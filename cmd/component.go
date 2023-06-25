package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	componentCmd = &cobra.Command{
		Use:   "component",
		Short: "Perform operations on the named component",
		Long:  "",
	}

	foldersCmd = &cobra.Command{
		Use:   "folders",
		Short: "Display the folders for a command, usually used with a mono-repo",
		Long:  "",
		Run:   executeComponentFolders,
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

func executeComponentFolders(ccmd *cobra.Command, args []string) {

}
