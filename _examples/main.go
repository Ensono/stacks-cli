package main

import (
	"log"

	"github.com/amido/stacks-cli/pkg/config"
	"github.com/amido/stacks-cli/pkg/scaffold"
)

func main() {
	conf, err := config.New(config.InputConfig{ProjectName: "foo", ProjectType: "netcore", Deployment: "azdo", Platform: "aks"})

	if err != nil {
		log.Fatalf("%v\n\n", err.Error())
	}

	// Add replace config - simulating user overrides
	conf.Replace = &[]config.ReplaceConfig{
		config.ReplaceConfig{
			Files:  []string{"foo", "bar"},
			Values: map[string]string{"find": "replace"},
		}}
	sc := scaffold.New(conf)
	err = sc.Run()
	if err != nil {
		log.Fatalf("%v\n\n", err.Error())
	}
	log.Printf("%v\n\n", "Finished examples run")
}
