package config

// import (
// 	"amido.com/scaffold/internal/util"
// 	"amido.com/scaffold/pkg/scaffold"
// )

// func ReadSelfConfigFile(input scaffold.InputConfig) (*scaffold.SelfConfig, error) {
// 	return readSelfConfigFile(input)
// }

// // readSelfConfigFile constructs self config for CLI based on bundle resources
// func readSelfConfigFile(input scaffold.InputConfig) (*scaffold.SelfConfig, error) {

// 	sharedT, err := ParseShared()
// 	util.CheckErrors(err)

// 	specificT, err := ParseSpecific(input.Platform, input.Deployment, input.ProjectType)

// 	s := scaffold.SelfConfig{
// 		Shared:   &sharedT,
// 		Specific: &specificT,
// 	}
// 	// TODO: feat request allow overwrite of self config from outside (as long as it can be parsed back to a SelfConfig)
// 	return &s, err
// }

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type TypeDetail struct {
	Gitrepo                  string   `yaml:"git_repo"`
	Gitref                   string   `yaml:"git_ref"`
	Localpath                string   `yaml:"local_path"`
	FilenameReplacementPaths []string `yaml:"filename_replacement_paths,omitempty"`
	Searchvalue              string   `yaml:"search_value,omitempty"`
	Foldermap                []struct {
		Src  string `yaml:"src"`
		Dest string `yaml:"dest"`
	} `yaml:"folder_map"`
}

type SelfConfig struct {
	Shared   *TypeDetail
	Specific *TypeDetail
}

type OutputConfig struct {
	TmpPath   string
	ZipPath   string
	UnzipPath string
	NewPath   string
}

type ReplaceConfig struct {
	Files  []string          `yaml:"files"`
	Values map[string]string `yaml:"values"`
}

type InputConfig struct {
	ProjectName                      string `yaml:"project_name"`
	ProjectType                      string `yaml:"project_type"`
	Platform                         string `yaml:"platform"`
	Deployment                       string `yaml:"deployment"`
	AdvancedcCnfig                   bool   `yaml:"advanced_config"`
	CreateConfig                     bool   `yaml:"create_config"`
	CloudRegion                      string `yaml:"cloud_egion"`
	CloudResourcegroup               string `yaml:"cloud_resourceGroup"`
	BusinessCompany                  string `yaml:"business_company"`
	BusinessProject                  string `yaml:"business_project"`
	BusinessDomain                   string `yaml:"business_domain"`
	BusinessComponent                string `yaml:"business_component"`
	SourcecontrolRepotype            string `yaml:"sourcecontrol_repo_type"`
	SourcecontrolReponame            string `yaml:"sourcecontrol_repo_name"`
	SourcecontrolRepourl             string `yaml:"sourcecontrol_repo_url"`
	TerraformBackendStorage          string `yaml:"terraformBackendStorage"`
	TerraformBackendStoragerg        string `yaml:"terraformBackendStorageRg"`
	TerraformBackendStoragecontainer string `yaml:"terraformBackendStorageContainer"`
	NetworkingBaseDomain             string `yaml:"networking_base_domain"`
}

type Config struct {
	Input   *InputConfig
	Self    *SelfConfig
	Output  *OutputConfig
	Replace *[]ReplaceConfig
}

// Create creates a config object based on user input or a config file
func Create(data []byte) (*Config, error) {

	tmpPath, err := os.MkdirTemp("", "source")
	if err != nil {
		return nil, err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// helper.ShowInfo(fmt.Sprintf("Current Dir: %s\n", pwd))

	// helper.ShowInfo(fmt.Sprintf("New Project Dir: %s\n", s.Config.Output.NewPath))

	t := InputConfig{}

	if err := yaml.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	selfConf, err := readSelfConfigFile(t)
	if err != nil {
		return nil, err
	}

	conf := Config{
		Output: &OutputConfig{
			NewPath:   fmt.Sprintf("%s/%s", pwd, t.ProjectName),
			TmpPath:   tmpPath,
			ZipPath:   fmt.Sprintf("%s/source.zip", tmpPath),
			UnzipPath: path.Join(tmpPath, "wip", selfConf.Specific.Localpath),
		},
		Input: &t,
		Self:  selfConf,
	}

	return &conf, err
}

// readSelfConfigFile constructs self config for CLI based on bundle resources
func readSelfConfigFile(input InputConfig) (*SelfConfig, error) {
	// var err error

	sharedT, err := ParseLocalShared()
	// if errShared != nil {
	if err != nil {
		return nil, err
	}

	specificT, err := ParseLocalSpecific(input.Platform, input.Deployment, input.ProjectType)

	if err != nil {
		return nil, err
	}

	s := SelfConfig{
		Shared:   &sharedT,
		Specific: &specificT,
	}

	// TODO: feat request allow overwrite of self config from outside (as long as it can be parsed back to a SelfConfig)
	return &s, err
}
