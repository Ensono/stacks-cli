package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the version of the app",
		Long:  "",
		Run:   showVersion,
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func showVersion(ccmd *cobra.Command, args []string) {
	fmt.Println("Version: ", Config.GetVersion())
}
