package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/amido/stacks-cli/internal/config/static"
	"github.com/amido/stacks-cli/internal/constants"
	"github.com/amido/stacks-cli/internal/models"
	"github.com/amido/stacks-cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	// App holds the configured objects, such as the logger
	App models.App

	// Config variable to hold the model after parsing
	Config config.Config

	// Set a variable to hold the version number of the application
	version string
)

var rootCmd = &cobra.Command{
	Use:     "stacks-cli",
	Short:   "Build up a new project based on the Amido Stacks system",
	Long:    "",
	Version: version,

	// Call prerun method to unmarshal the config into the app models
	PersistentPreRun: preRun,
}

// Execute is the entry point for the application
func Execute() {
	// Determine if there was an error in the application
	err := rootCmd.Execute()

	if err != nil {
		log.Fatalf("%v", err)
	}
}

func init() {
	// Declare variables to accept the flag values
	var logLevel string
	var logFormat string
	var logColour bool

	var workingDir string
	var tmpDir string

	cobra.OnInitialize(initConfig)

	// get the default directories
	defaultTempDir, defaultWorkingDir := getDefaultDirectories()

	// Add flags that are to be used in every command
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Path to the configuration file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "info", "Logging Level")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "logformat", "f", "text", "Logging format, text or json")
	rootCmd.PersistentFlags().BoolVarP(&logColour, "logcolour", "", true, "State if colours should be used in the text output")

	rootCmd.PersistentFlags().StringVarP(&workingDir, "workingdir", "w", defaultWorkingDir, "Directory to be used to create the new projects in")
	rootCmd.PersistentFlags().StringVar(&tmpDir, "tempdir", defaultTempDir, "Temporary directory to be used by the CLI")

	// Bind command line arguments
	viper.BindPFlags(rootCmd.Flags())

	// Configure the logging options
	viper.BindPFlag("log.format", rootCmd.PersistentFlags().Lookup("logformat"))
	viper.BindPFlag("log.colour", rootCmd.PersistentFlags().Lookup("logcolour"))
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("loglevel"))

	viper.BindPFlag("directory.working", rootCmd.PersistentFlags().Lookup("workingdir"))
	viper.BindPFlag("directory.temp", rootCmd.PersistentFlags().Lookup("tempdir"))
}

// initConfig reads in a confiig file and ENV vars if set
// Sets up logging with the specified log level
func initConfig() {

	if cfgFile != "" {
		// Use the config file from the cobra flag
		viper.SetConfigFile(cfgFile)
	}

	// Allow configuration options to be set using Environment variables
	viper.SetEnvPrefix(constants.EnvVarPrefix)

	// The configuration settings are nested
	// Change the `.` delimiter to a `_` when accessing from an Environment Variable
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv() // read in environment variables that match

	// Read  in the static configuration
	stacks_frameworks := string(static.Config("stacks_frameworks"))
	stacks_config := strings.NewReader(stacks_frameworks)
	viper.SetConfigType("yaml")
	viper.MergeConfig(stacks_config)

	// if a configuration file is found, read it in
	if err := viper.MergeInConfig(); err == nil {
		fmt.Println("Using configuration file:", viper.ConfigFileUsed())
	}
}

func preRun(ccmd *cobra.Command, args []string) {

	err := viper.Unmarshal(&Config.Input)
	if err != nil {
		log.Fatalf("Unable to read configuration into models: %v", err)
	}

	// Configure application logging
	// This is done after unmarshalling of the configuration so that the
	// model values can be used rather than the strings from viper
	App.ConfigureLogging(Config.Input.Log)

	// Set the version of the app in the configuration
	Config.Input.Version = version

}

// setDefaultDirectoies sets the workingdir to the current directory and the
// tempdir to the system temporary directory
func getDefaultDirectories() (string, string) {

	tmpPath, err := os.MkdirTemp("", "stackscli")
	if err != nil {
		log.Fatalf("Unable to create temporary directory")
	}

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to determine current directory")
	}

	return tmpPath, workingDir
}
