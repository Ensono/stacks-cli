package config

import (
	"fmt"

	"github.com/amido/stacks-cli/internal/util"
)

type Stacks struct {
	Components []StacksComponent `mapstructure:"components" yaml:"components"`

	named map[string]StacksComponent
}

func (s *Stacks) GetComponent(ref string) (StacksComponent, error) {

	var err error
	stacks_component := StacksComponent{}

	// if the named map contains the component return it, otherwise return
	// empty and set the error
	if val, ok := s.named[ref]; ok {
		stacks_component = val
	} else {
		err = fmt.Errorf("unable to find component with reference: %s", ref)
	}

	return stacks_component, err
}

func (s *Stacks) GetComponentCount() int {
	return len(s.named)
}

// SetUniqueComponents rewrites the slice of Stacks so that it is a unique
// list. Later values in the original slice take precedence
func (s *Stacks) SetUniqueComponents() {

	// create a map to hold the unique values of the slice
	s.named = make(map[string]StacksComponent)

	// iterate around the components that have been set
	for _, component := range s.Components {

		// update the named map with the component
		s.named[component.GetName()] = component
	}
}

func (s *Stacks) GetComponentPackage(name string) Package {
	component, _ := s.GetComponent(name)
	return component.Package

}

func (s *Stacks) GetComponentPackageRef(name string) string {
	var result string
	pkg := s.GetComponentPackage(name)

	switch pkg.Type {
	case "git":
		result = pkg.URL
	case "nuget":
		result = pkg.Name
	}

	return result
}

// Normalize checks to see if the older API is being used and will take
// the values from that and populate the new structure
// It will return an error object to be used as a warning for people to update their structure
func (p *Package) Normalize() string {
	var msg string

	// if the type is empty, default to github
	if p.Type == "" {
		p.Type = "github"
	}

	// ensure that the type of the repo is correct
	validTypes := []string{"git", "nuget"}
	if !util.SliceContains(validTypes, p.Type) {
		msg = fmt.Sprintf("Specified type of '%s' is invalid, please check your configuration", p.Type)
	}

	return msg
}

/*
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
*/
