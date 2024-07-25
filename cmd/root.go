package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/Ensono/stacks-cli/internal/config/staticFiles"
	"github.com/Ensono/stacks-cli/internal/constants"
	"github.com/Ensono/stacks-cli/internal/models"
	"github.com/Ensono/stacks-cli/internal/util"
	"github.com/Ensono/stacks-cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	cfgFile string
	tmpDir  string

	// App holds the configured objects, such as the logger and help messages
	App models.App

	// Config variable to hold the model after parsing
	Config config.Config

	// define a slice to hold a list of all the configuration files that have been read in
	ConfigFiles []string

	// Set a variable to hold the version number of the application
	version string

	// set shared variables that can be used as parameters across multiple commands
	business_company            string
	business_domain             string
	terraform_backend_storage   string
	terraform_backend_group     string
	terraform_backend_container string
	global                      bool
)

var rootCmd = &cobra.Command{
	Use:     "stacks-cli",
	Short:   "Build up a new project based on the Ensono Stacks system",
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
	var homeDir string

	var noBanner bool
	var noCLIVersionCheck bool
	var dryrun bool

	var githubToken string

	var override_internal_config string

	var folders []string

	cobra.OnInitialize(initConfig)

	// get the default directories
	defaultTempDir := util.GetDefaultTempDir()
	defaultWorkingDir := util.GetDefaultWorkingDir()
	defaultUserHomeDir := util.GetUserHomeDir()

	// Add flags that can be used in every command
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Path to the configuration file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "info", "Logging Level")
	rootCmd.PersistentFlags().StringVarP(&logFormat, "logformat", "f", "text", "Logging format, text or json")
	rootCmd.PersistentFlags().BoolVarP(&logColour, "logcolour", "", true, "State if colours should be used in the text output")
	rootCmd.PersistentFlags().StringVar(&logFile, "logfile", "", "File to write logs to")

	rootCmd.PersistentFlags().StringVarP(&workingDir, "workingdir", "w", defaultWorkingDir, "Directory to be used to create the new projects in")
	rootCmd.PersistentFlags().StringVar(&tmpDir, "tempdir", defaultTempDir, "Temporary directory to be used by the CLI")
	rootCmd.PersistentFlags().StringVar(&homeDir, "homedir", defaultUserHomeDir, "Users home directory")

	rootCmd.PersistentFlags().BoolVar(&dryrun, "dryrun", false, "Shows what actions would be taken but does not perform them")
	rootCmd.PersistentFlags().BoolVar(&noBanner, "nobanner", false, "Do not display the Stacks banner when running the command")
	rootCmd.PersistentFlags().BoolVar(&noCLIVersionCheck, "nocliversion", false, "Do not check for latest version of the CLI")
	rootCmd.PersistentFlags().BoolVarP(&onlineHelp, "onlinehelp", "H", false, "Open web browser with help for the command")

	rootCmd.PersistentFlags().StringVar(&githubToken, "token", "", "GitHub token to perform authenticated requests against the GitHub API")

	rootCmd.PersistentFlags().StringVar(&override_internal_config, "internalconfig", "", "Path to the configuration override file")

	rootCmd.PersistentFlags().StringSliceVar(&folders, "folders", []string{}, "List of additional folders to be used when running setup")

	rootCmd.PersistentFlags().StringVar(&business_company, "company", "", "The name of the company")
	rootCmd.PersistentFlags().StringVarP(&business_domain, "area", "A", "", "Area within the company that this project will belong to, e.g. core")
	rootCmd.PersistentFlags().StringVar(&terraform_backend_storage, "tfstorage", "", "Name of the storage to be used for Terraform state")
	rootCmd.PersistentFlags().StringVar(&terraform_backend_group, "tfgroup", "", "Name of the group that the storage account is in")
	rootCmd.PersistentFlags().StringVar(&terraform_backend_container, "tfcontainer", "", "Name of the container within the storage to use")
	rootCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "Set the values globally. These will be set in the user's home directory")

	// Bind command line arguments
	viper.BindPFlags(rootCmd.Flags())

	// Configure the logging options
	viper.BindPFlag("input.log.format", rootCmd.PersistentFlags().Lookup("logformat"))
	viper.BindPFlag("input.log.colour", rootCmd.PersistentFlags().Lookup("logcolour"))
	viper.BindPFlag("input.log.level", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("input.log.file", rootCmd.PersistentFlags().Lookup("logfile"))

	viper.BindPFlag("input.directory.working", rootCmd.PersistentFlags().Lookup("workingdir"))
	viper.BindPFlag("input.directory.temp", rootCmd.PersistentFlags().Lookup("tempdir"))
	viper.BindPFlag("input.directory.home", rootCmd.PersistentFlags().Lookup("homedir"))

	viper.BindPFlag("input.options.nobanner", rootCmd.PersistentFlags().Lookup("nobanner"))
	viper.BindPFlag("input.options.nocliversion", rootCmd.PersistentFlags().Lookup("nocliversion"))
	viper.BindPFlag("input.options.onlinehelp", rootCmd.PersistentFlags().Lookup("onlinehelp"))
	viper.BindPFlag("input.options.dryrun", rootCmd.PersistentFlags().Lookup("dryrun"))

	viper.BindPFlag("input.overrides.internal_config", rootCmd.PersistentFlags().Lookup("internalconfig"))
	viper.BindPFlag("input.folders", rootCmd.PersistentFlags().Lookup("folders"))

	viper.BindPFlag("input.business.company", rootCmd.PersistentFlags().Lookup("company"))
	viper.BindPFlag("input.business.domain", rootCmd.PersistentFlags().Lookup("area"))

	viper.BindPFlag("input.terraform.backend.storage", rootCmd.PersistentFlags().Lookup("tfstorage"))
	viper.BindPFlag("input.terraform.backend.group", rootCmd.PersistentFlags().Lookup("tfgroup"))
	viper.BindPFlag("input.terraform.backend.container", rootCmd.PersistentFlags().Lookup("tfcontainer"))

	viper.BindPFlag("input.global", rootCmd.PersistentFlags().Lookup("global"))
}

// initConfig reads in a config file and ENV vars if set
// Sets up logging with the specified log level
func initConfig() {

	Config.Init()

	// Ensure that the home directory for stacks for the user exists
	stacksCliHome := util.GetStacksCLIDir()
	util.CreateIfNotExists(stacksCliHome, 0755)

	// determine the root path based on the operating system
	root_path := "/"
	if runtime.GOOS == "windows" {
		root_path = fmt.Sprintf("%s\\", os.Getenv("SystemDrive"))
	}

	// get a list of the paths from the cwd
	var directories []string
	for _path := util.GetDefaultWorkingDir(); _path != root_path; _path = filepath.Dir(_path) {
		directories = append(directories, util.NormalisePath(_path, string(os.PathSeparator)))
	}
	slices.Reverse(directories)

	// if additional folders have been specified on the commandline add them to the list
	folders := viper.GetStringSlice("input.folders")
	if len(folders) > 0 {
		directories = append(directories, folders...)
	}

	// set multiple paths that a configuration file can be read from
	// - home directory file
	// - all folders from the root to the current directory, this is so that the closest configuration
	//   is the one that is read in last
	// - any additional paths that have been specified on the command line
	homeDirConfigfile := fmt.Sprintf("%s.yml", path.Join(util.GetUserHomeDir(), constants.ConfigFileDir, constants.ConfigName))
	if util.Exists(homeDirConfigfile) {
		ConfigFiles = append(ConfigFiles, util.NormalisePath(homeDirConfigfile, string(os.PathSeparator)))
		viper.SetConfigName(constants.ConfigName)
		viper.AddConfigPath(path.Join(util.GetUserHomeDir(), constants.ConfigFileDir))
		viper.MergeInConfig()
	}

	// reverse the order of the directories so that the closest one is read in last
	for i := len(directories) - 1; i >= 0; i-- {
		// build up the path to a possible configuration file and check if it exists
		configfile := util.NormalisePath(fmt.Sprintf("%s.yml", path.Join(directories[i], constants.ConfigFileDir, constants.ConfigName)), string(os.PathSeparator))
		if util.Exists(configfile) {

			if !util.SliceContains(ConfigFiles, configfile) {
				ConfigFiles = append(ConfigFiles, configfile)
			}

			// add the config path to the viper instance
			viper.AddConfigPath(path.Join(directories[i], constants.ConfigFileDir))

			// merge the configuration file into the viper instance
			viper.MergeInConfig()
		}
	}

	// if a configuation file has been specified on the command line, copy it to the tempdir and
	// add to the viper instance. This is so that it can be named correctly
	if cfgFile != "" {
		cfgFile = util.NormalisePath(cfgFile, string(os.PathSeparator))
		// if the cfgfile can be found, copy it to the tempdir
		if util.Exists(cfgFile) {

			ConfigFiles = append(ConfigFiles, cfgFile)

			// get the name of the configuration file and add it to the configuration path
			filename := filepath.Base(cfgFile)
			viper.SetConfigName(strings.TrimSuffix(filename, filepath.Ext(filename)))
			viper.AddConfigPath(filepath.Dir(cfgFile))
			viper.MergeInConfig()

			// reset the name of the configuration file
			viper.SetConfigName(constants.ConfigName)
		}
	}

	// Determine if an environment variable has been set that states the string to be
	// used as the environment variable prefix
	// This has been done so that the it is possible for people to set a different value and be
	// compatible with older version of the Stacks CLI
	envvarprefix := constants.EnvVarPrefix
	if os.Getenv("STACKSCLI_ENVVARPREFIX") != "" {
		envvarprefix = os.Getenv("STACKSCLI_ENVVARPREFIX")
	}

	// Allow configuration options to be set using Environment variables
	viper.SetEnvPrefix(envvarprefix)

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
	var err error

	// Read in the internal configuration file
	// This can come from three locations, in this order:
	// - file that has been specified on the command line
	// - downloaded configuration file that exists in the user home directory
	// - the internal configuration file that is shipped with the application

	// set the content of the file so it can be read in
	var configContent []byte

	if viper.GetString("input.overrides.internal_config") != "" {
		configContent, err = util.GetFileContent(viper.GetString("input.overrides.internal_config"))
	} else if util.Exists(path.Join(util.GetStacksCLIDir(), "internal_config.yml")) {
		configContent, err = util.GetFileContent(path.Join(util.GetStacksCLIDir(), "internal_config.yml"))
	} else {
		configContent = Config.Internal.GetFileContent("config")
	}

	if err != nil {

		log.Fatalf("Unable to read internal configuration: %v", err)
		App.Logger.Exit(1)
	}

	err = yaml.Unmarshal(configContent, &Config)
	if err != nil {
		log.Fatalf("Unable to read internal configuration: %v", err)
		App.Logger.Exit(1)
	}

	ScaffoldOverrides()

	// Unmarshal the configuration into the models in the application
	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("Unable to read configuration into models: %v", err)
		App.Logger.Exit(4)
	}

	// Set the help messages in the App.Help object
	App.LoadHelp(Config.Internal.GetFileContent("help"))

	// Configure application logging
	// This is done after unmarshalling of the configuration so that the
	// model values can be used rather than the strings from viper
	App.ConfigureLogging(Config.Input.Log)

	if Config.Input.Overrides.InternalConfigPath != "" {
		msg := App.Help.GetMessage("INT001", Config.Input.Overrides.InternalConfigPath)
		App.Logger.Infof(msg)
	}

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

	// output a list of the configuration files that have been read in
	App.Logger.Infof("Configuration files read:\n\t%s", strings.Join(ConfigFiles, "\n\t"))

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
