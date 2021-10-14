package static

import (
	_ "embed"
)

var shared = `
git_repo: ""
git_ref: ""
local_path: ""
folder_map:
  - src: shared/_gitignore
    dest: "./.gitignore"
`

var aks_azdo_ssr = `
git_repo: https://github.com/amido/stacks-typescript-ssr.git
git_ref: 4a51f638d0a5597648c2f6a1b96c31460dced7ab
local_path: src/ssr
folder_map:
  - src: shared/README.md
    dest: ./README.md
  - src: src/ssr/build/azDevops/azure/templates
    dest: build/azDevops/azure/templates
  - src: src/ssr/build/azDevops/azure/azure-pipelines-ssr-aks.yml
    dest: build/azDevops/azure/azure-pipelines-ssr-aks.yml
  - src: src/ssr/deploy/azure/app/kube
    dest: deploy/azure/app/kube
  - src: src/ssr/deploy/k8s/app/base_app-deploy.yml
    dest: deploy/k8s/app/base_app-deploy.yml
  - src: src/ssr/test/testcafe
    dest: test/testcafe
  - src: src/ssr/src/ssr
    dest: src/ssr
  - src: src/ssr/.gitattributes
    dest: ".gitattributes"
  - src: src/ssr/.gitignore
    dest: ".gitignore"
  - src: src/ssr/yamllint.conf
    dest: yamllint.conf
`

var aks_azdo_netcore = `
git_repo: https://github.com/amido/stacks-dotnet.git
git_ref: 2a810ec620360bb6d69c220abdf409bc6af349be
local_path: src/netcore
filename_replacement_paths:
  - "/src"
search_value: xxAMIDOxx.xxSTACKSxx
folder_map:
  - src: shared/README.md
    dest: "./README.md"
  - src: src/netcore/build/azDevops/azure/templates/steps/build/build-netcore.yml
    dest: build/azDevops/azure/templates/steps/build/build-netcore.yml
  - src: src/netcore/build/azDevops/azure/azure-pipelines-netcore-k8s.yml
    dest: build/azDevops/azure/azure-pipelines-netcore-k8s.yml
  - src: src/netcore/deploy/azure/app/kube
    dest: deploy/azure/app/kube
  - src: src/netcore/deploy/k8s/app/base_api-deploy.yml
    dest: deploy/k8s/app/base_api-deploy.yml
  - src: src/netcore/src
    dest: src
  - src: src/netcore/.gitattributes
    dest: ".gitattributes"
  - src: src/netcore/.gitignore
    dest: ".gitignore"
  - src: src/netcore/yamllint.conf
    dest: yamllint.conf
`

// The following static configuration sets the URLs to the repos for
// the location of the repositories
// This can be overriden by passing the configuration in as a configuration file
// but this will be the default
//go:embed stacks_frameworks.yml
var stacks_frameworks string

func FrameworkCommand(framework string) []string {
	commands := map[string][]string{
		"dotnet": {"dotnet"},
		"java":   {"java"},
	}

	return commands[framework]
}

// Set the banner that is written out to the screen when stacks is run
//go:embed banner.txt
var Banner string

// Config byte parses static
func Config(key string) []byte {
	switch key {
	case "shared":
		return []byte(shared)
	case "aks_azdo_ssr":
		return []byte(aks_azdo_ssr)
	case "aks_azdo_netcore":
		return []byte(aks_azdo_netcore)
	// NOTE: add more mappings here
	case "stacks_frameworks":
		return []byte(stacks_frameworks)
	default:
		return []byte(shared)
	}
}
