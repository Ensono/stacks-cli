package filter

import (
	"log"
	"regexp"
	"testing"

	"github.com/Ensono/stacks-cli/pkg/config"
)

func TestFilter(t *testing.T) {

	// create object to work with
	config := config.Config{
		Input: config.InputConfig{
			Terraform: config.Terraform{
				Backend: config.TerraformBackend{
					Storage:   "tfstatestorage",
					Group:     "terraformstate",
					Container: "tfstate",
				},
			},
			Business: config.Business{
				Company: "MyCompany",
				Domain:  "AI",
			},
		},
	}

	// create a table of tests, that include different paths to filter
	tables := []struct {
		filters []string
		test    string
		msg     string
	}{
		{
			[]string{"input.business.company"},
			"",
			"Nothing should be exported as the filter is invalid for the object",
		},
		{
			[]string{"business.company"},
			`^business:\s*company:\sMyCompany`,
			"Only the company name should be output:\n %s",
		},
		{
			[]string{"business.company", "terraform.backend.storage"},
			`^business:\s*company:\sMyCompany\sterraform:\s*backend:\s*storage:\stfstatestorage`,
			"Terraform storage and company name should be output: \n %s",
		},
	}

	// iterate around the table of tests
	for _, table := range tables {

		// create a new filter instance
		filter := New()
		filter.Filter(config.Input, table.filters)

		// get the result of the filter and compare against the expected value
		result := filter.String()

		found, err := regexp.MatchString(table.test, result)

		if err != nil {
			log.Fatal(err)
		}

		if !found {
			t.Errorf(table.msg, result)
		}
	}
}
