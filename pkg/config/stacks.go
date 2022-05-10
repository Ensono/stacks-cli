package config

import (
	"fmt"

	"github.com/amido/stacks-cli/internal/util"
)

type Stacks struct {
	Dotnet Dotnet `mapstructure:"dotnet"`
	Java   Java   `mapstructure:"java"`
	NodeJS NodeJS `mapstructure:"nodejs"`
	Infra  Infra  `mapstructure:"infra"`
}

type Dotnet struct {
	Webapi RepoInfo `mapstructure:"webapi"`
	CQRS   RepoInfo `mapstructure:"cqrs"`
	Events RepoInfo `mapstructure:"events"`
}

type Java struct {
	Webapi RepoInfo `mapstructure:"webapi"`
	CQRS   RepoInfo `mapstructure:"cqrs"`
	Events RepoInfo `mapstructure:"events"`
}

type NodeJS struct {
	CSR RepoInfo `mapstructure:"csr"`
	SSR RepoInfo `mapstructure:"ssr"`
}

type Infra struct {
	AKS RepoInfo `mapstructure:"aks"`
}

type RepoInfo struct {
	Options string `mapstructure:"options"`
	Version string `mapstructure:"version"`
	Type    string `mapstructure:"type"`
	Name    string `mapstructure:"name"`
	ID      string `mapstructure:"id"`

	// Allow support for previous version tags
	URL   string `mapstructure:"url"`
	Trunk string `mapstructure:"trunk"`
}

// GetSrcURLMap returns a map of the source control repositores
func (stacks *Stacks) GetSrcURLMap() map[string]RepoInfo {

	srcUrls := map[string]RepoInfo{
		"dotnet_webapi": stacks.Dotnet.Webapi,
		"dotnet_cqrs":   stacks.Dotnet.CQRS,
		"dotnet_events": stacks.Dotnet.Events,
		"java_webapi":   stacks.Java.Webapi,
		"java_cqrs":     stacks.Java.CQRS,
		"java_events":   stacks.Java.Events,
		"nodejs_csr":    stacks.NodeJS.CSR,
		"nodejs_ssr":    stacks.NodeJS.SSR,
		"infra_aks":     stacks.Infra.AKS,
	}

	return srcUrls
}

func (stacks *Stacks) GetSrcURL(key string) RepoInfo {
	srcUrls := stacks.GetSrcURLMap()
	return srcUrls[key]
}

// Normalize checks to see if the older API is being used and will take
// the values from that and populate the new structure
// It will return an error object to be used as a warning for people to update their structure
func (r *RepoInfo) Normalize() string {
	var msg string
	var usingOldApi bool

	if r.URL != "" && r.Name == "" {
		r.Name = r.URL
		usingOldApi = true
	}

	if r.Trunk != "" && r.Version == "" {
		r.Version = r.Trunk
		usingOldApi = true
	}

	if usingOldApi {
		msg = "Your configuration is using a deprecated version of the Stacks framework API, please update your configuration"
	}

	// if the type is empty, default to github
	if r.Type == "" {
		r.Type = "github"
	}

	// ensure that the type of the repo is correct
	validTypes := []string{"github", "nuget"}
	if !util.SliceContains(validTypes, r.Type) {
		msg = fmt.Sprintf("Specified type of '%s' is invalid, please check your configuration", r.Type)
	}

	return msg
}
