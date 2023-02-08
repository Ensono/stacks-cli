package config

import (
	"fmt"

	"github.com/amido/stacks-cli/internal/util"
)

type Stacks struct {
	Dotnet Dotnet `mapstructure:"dotnet" yaml:"dotnet"`
	Java   Java   `mapstructure:"java" yaml:"java"`
	Nx     Nx     `mapstructure:"nx" yaml:"nx"`
	Infra  Infra  `mapstructure:"infra" yaml:"infra"`
}

type Dotnet struct {
	Webapi RepoInfo `mapstructure:"webapi" yaml:"webapi"`
	CQRS   RepoInfo `mapstructure:"cqrs" yaml:"cqrs"`
}

type Java struct {
	Webapi RepoInfo `mapstructure:"webapi"  yaml:"webapi"`
	CQRS   RepoInfo `mapstructure:"cqrs"  yaml:"cqrs"`
	Events RepoInfo `mapstructure:"events" yaml:"events"`
}

type Nx struct {
	NextJs RepoInfo `mapstructure:"next" yaml:"next"`
	Apps   RepoInfo `mapstructure:"apps" yaml:"apps"`
}

type Infra struct {
	AKS RepoInfo `mapstructure:"aks" yaml:"aks"`
}

type RepoInfo struct {
	Options string `mapstructure:"options" yaml:"options,omitempty"`
	Version string `mapstructure:"version" yaml:"version,omitempty"`
	Type    string `mapstructure:"type" yaml:"type,omitempty"`
	Name    string `mapstructure:"name" yaml:"name,omitempty"`
	ID      string `mapstructure:"id" yaml:"id,omitempty"`
}

// GetSrcURLMap returns a map of the source control repositores
func (stacks *Stacks) GetSrcURLMap() map[string]RepoInfo {

	srcUrls := map[string]RepoInfo{
		"dotnet_webapi": stacks.Dotnet.Webapi,
		"dotnet_cqrs":   stacks.Dotnet.CQRS,
		"java_webapi":   stacks.Java.Webapi,
		"java_cqrs":     stacks.Java.CQRS,
		"java_events":   stacks.Java.Events,
		"nx_next":       stacks.Nx.NextJs,
		"nx_apps":       stacks.Nx.Apps,
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
