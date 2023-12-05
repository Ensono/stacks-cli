package config

import (
	"fmt"
	"sort"

	"github.com/Ensono/stacks-cli/internal/util"
)

type Stacks struct {
	Components map[string]StacksComponent `mapstructure:"components" yaml:"components"`
}

func (s *Stacks) GetComponentCount() int {
	return len(s.Components)
}

func (s *Stacks) GetComponentPackage(name string) Package {
	// component, _ := s.GetComponent(name)
	component := s.Components[name]
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

// GetComponentNames returns a sorted slice of the components that are defined
// and removes any duplicates
func (s *Stacks) GetComponentNames() []string {

	// create a map to state if the name has already been encountered
	var exists map[string]bool = make(map[string]bool)

	// create the return slice
	var result []string

	for _, component := range s.Components {
		if _, ok := exists[component.Group]; !ok {
			// item has not been found in slice so add to the result and update the exists map
			exists[component.Group] = true
			result = append(result, component.Group)
		}
	}

	// sort the array
	sort.Strings(result)

	// return the result
	return result
}

// GetComponentOptions analyses the StacksComponent slice and returns all of the options that
// are associated with the specified framework
func (s *Stacks) GetComponentOptions(framework string) []string {

	var options []string

	// iterate around the component slice looking for the framework and append each option to the options slice
	for _, component := range s.Components {
		if component.Group == framework {
			options = append(options, component.Name)
		}
	}

	// sort the options
	sort.Strings(options)

	return options
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
