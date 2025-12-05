package cmd

import (
	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/Ensono/stacks-cli/pkg/export"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export internal files to the filesystem",
		Long: `The CLI contains embedded configuration files that control
		how the application works. These files can be exported to the filesystem
		and modified and then set as arguments to the commands to change the apps
		behaviour`,
		Run: executeExportFiles,
	}
)

func init() {

	// declare variables that will be used to hold information from the supplied
	// arguments
	var directory string

	rootCmd.AddCommand(exportCmd)

	// Configure the flags
	exportCmd.Flags().StringVarP(&directory, "directory", "d", util.GetDefaultWorkingDir(), "Directory to be used to export the files to")

	viper.BindPFlag("input.directory.export", exportCmd.Flags().Lookup("directory"))
}

func executeExportFiles(ccmd *cobra.Command, args []string) {

	// Call the export method
	export := export.New(&Config, App.Logger)
	err := export.Run()
	if err != nil {
		msg := App.Help.GetMessage("GEN001", "export", err.Error())
		App.Logger.Fatalf("%s", msg)
	}

}
