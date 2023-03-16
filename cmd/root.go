package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/amido/stacks-cli/internal/config/staticFiles"
	"github.com/amido/stacks-cli/internal/constants"
	"github.com/amido/stacks-cli/internal/models"
	"github.com/amido/stacks-cli/internal/util"
	"github.com/amido/stacks-cli/pkg/config"
	yaml "github.com/goccy/go-yaml"
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
	var logFile string
	var onlineHelp bool

	var workingDir string
	var tmpDir string

	var noBanner bool
	var noCLIVersionCheck bool
	var dryrun bool

	var githubToken string

	cobra.OnInitialize(initConfig)

	// get the default directories
	defaultTempDir := util.GetDefaultTempDir()
	defaultWorkingDir := util.GetDefaultWorkingDir()

	// Add flags that are to be used in every command
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Path to the configuration file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "info", "Logging Level")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "logformat", "f", "text", "Logging format, text or json")
	rootCmd.PersistentFlags().BoolVarP(&logColour, "logcolour", "", true, "State if colours should be used in the text output")
	rootCmd.PersistentFlags().StringVar(&logFile, "logfile", "", "File to write logs to")

	rootCmd.PersistentFlags().StringVarP(&workingDir, "workingdir", "w", defaultWorkingDir, "Directory to be used to create the new projects in")
	rootCmd.PersistentFlags().StringVar(&tmpDir, "tempdir", defaultTempDir, "Temporary directory to be used by the CLI")

	rootCmd.PersistentFlags().BoolVar(&dryrun, "dryrun", false, "Shows what actions would be taken but does not perform them")
	rootCmd.PersistentFlags().BoolVar(&noBanner, "nobanner", false, "Do not display the Stacks banner when running the command")
	rootCmd.PersistentFlags().BoolVar(&noCLIVersionCheck, "nocliversion", false, "Do not check for latest version of the CLI")
	rootCmd.PersistentFlags().BoolVarP(&onlineHelp, "onlinehelp", "H", false, "Open web browser with help for the command")

	rootCmd.PersistentFlags().StringVar(&githubToken, "token", "", "GitHub token to perform authenticated requests against the GitHub API")

	// Bind command line arguments
	viper.BindPFlags(rootCmd.Flags())

	// Configure the logging options
	viper.BindPFlag("log.format", rootCmd.PersistentFlags().Lookup("logformat"))
	viper.BindPFlag("log.colour", rootCmd.PersistentFlags().Lookup("logcolour"))
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("log.file", rootCmd.PersistentFlags().Lookup("logfile"))

	viper.BindPFlag("directory.working", rootCmd.PersistentFlags().Lookup("workingdir"))
	viper.BindPFlag("directory.temp", rootCmd.PersistentFlags().Lookup("tempdir"))

	viper.BindPFlag("options.nobanner", rootCmd.PersistentFlags().Lookup("nobanner"))
	viper.BindPFlag("options.nocliversion", rootCmd.PersistentFlags().Lookup("nocliversion"))
	viper.BindPFlag("options.onlinehelp", rootCmd.PersistentFlags().Lookup("onlinehelp"))
	viper.BindPFlag("options.dryrun", rootCmd.PersistentFlags().Lookup("dryrun"))

}

// initConfig reads in a config file and ENV vars if set
// Sets up logging with the specified log level
func initConfig() {

	Config.Init()

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
	// viper.SetConfigType("yaml")
	// viper.MergeConfig(strings.NewReader(Config.Internal.GetFileContentString("stacks_frameworks")))

	err := viper.MergeInConfig()
	if err != nil && viper.ConfigFileUsed() != "" {
		fmt.Printf("Unable to read in configuration file: %s", err.Error())
		return
	}
}

func preRun(ccmd *cobra.Command, args []string) {

	// Unmarshal the data from the static config file into the config object
	err := yaml.Unmarshal(Config.Internal.GetFileContent("config"), &Config)
	if err != nil {
		log.Fatalf("Unable to read internal configuration: %v", err)
		App.Logger.Exit(1)
	}

	err = viper.Unmarshal(&Config.Input)
	if err != nil {
		log.Fatalf("Unable to read configuration into models: %v", err)
		App.Logger.Exit(2)
	}

	// Configure application logging
	// This is done after unmarshalling of the configuration so that the
	// model values can be used rather than the strings from viper
	App.ConfigureLogging(Config.Input.Log)

	// Set the version of the app in the configuration
	Config.Input.Version = version

	// output the banner, unless it has been disabled or the parent command is completion
	if !Config.NoBanner() && ccmd.Parent().Use != "completion" && !Config.OnlineHelp() {
		fmt.Println(staticFiles.IntFile_Banner)
	}

	// Check that the CLI is online
	// use a DNS lookup to check that github can be accessed
	// this is so that the check is not performed if the the environment is not
	// connected to the internet
	App.Logger.Info("Performing connectivity check")
	err = util.CheckConnectivity("github.com")
	if err != nil {
		App.Logger.Fatal(err.Error())
		return
	}

	// set the urls to use to open the web based help for a command
	// err = yaml.Unmarshal(Config.Internal.GetFileContent("help_urls"), &Config.Help)
	// if err != nil {
	//	App.Logger.Fatalf("Unable to parse help URL data: %s", err.Error())
	//}

	// Determine if the online help option has been specified, it is has
	// open up the webpage for the specified command and then exit
	if Config.OnlineHelp() {

		// call the command to open the webpage
		status := Config.OpenOnlineHelp(ccmd.Use, App.Logger)

		// if the onlinehelp has not worked display the normal command line based help
		if !status {
			ccmd.Help()
		}

		// exit the program
		os.Exit(0)
	}

	// Call method to determine if this version of the CLI is the latest one
	checkCLIVersion()

	// output the configuration file that is being used
	configFileUsed := viper.ConfigFileUsed()
	if configFileUsed != "" {
		App.Logger.Infof("Using configuration file: %s", configFileUsed)
	}

	// set the framework definitions on the config object
	//err = yaml.Unmarshal(Config.Internal.GetFileContent("framework_defs"), &Config.FrameworkDefs)
	//if err != nil {
	//	App.Logger.Fatalf("Unable to parse framework definition data: %s", err.Error())
	//}
}

func checkCLIVersion() {

	// do not perform version check if it has been turned off
	// or no token has been supplied
	if Config.Input.Options.NoCLIVersion || Config.Input.Options.Token == "" {
		return
	}

	App.Logger.Info("Checking for latest version of CLI")

	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", constants.GitHubRef)
	releaseMap, err := util.CallHTTPAPI(url, Config.Input.Options.Token)

	if err != nil {
		App.Logger.Errorf("Unable to get latest CLI version: %s", err.Error())
		return
	}

	// use semantic version checks to see if the current version of the app is the latest
	// version
	settings := config.Settings{}

	// get values from the map
	latestVersion := releaseMap["tag_name"].(string)
	releaseUrl := releaseMap["url"].(string)

	// create a version constraint
	constraint := fmt.Sprintf("= %s", strings.Replace(latestVersion, "v", "", -1))
	uptodate := settings.CompareVersion(constraint, Config.GetVersion(), App.Logger)

	if !uptodate {
		fmt.Printf("A newer release version of the Stacks CLI is available, %s, you are running version %s.\n\nMore information at %s\n\n", latestVersion, Config.GetVersion(), releaseUrl)
	}
}
