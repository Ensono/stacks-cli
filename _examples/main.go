package main

import (
	"log"

	"amido.com/scaffold/pkg/scaffold"
)

var sample = `
project_name: "test-app-infra-only"
project_type: "netcore"
platform: "aks"
deployment: "azdo"
advanced_config: true
cloud_region: "uksouth"
cloud_resource_group: "string"
business_company: "company"
business_project: "project"
business_domain: "core"
business_component: "infra"
sourcecontrol_repoType: "github"
sourcecontrol_repoName: "Name-Of-Repo"
sourceControl_repoUrl: "https://github/sample.git"
terraform_backend_storage: "replace_terraform_backend_storage"
terraform_backend_storagerg: "replace_terraform_backend_storage_rg"
terraform_backend_storagecontainer: "replace_terraform_backend_storage_container"
terraform_backend: 
networking_base_domain: domain.replace.me
`

func main() {
	conf, err := scaffold.Create([]byte(sample))
	if err != nil {
		log.Printf("%v\n\n")
	}

	sc := scaffold.New(&conf)
	err = sc.Run()
	if err != nil {
		log.Fatalf("%s\n\n", err.Error())
	}
	log.Printf("%v\n\n")
}
