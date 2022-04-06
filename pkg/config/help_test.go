package config

import "testing"

func TestGetUrl(t *testing.T) {

	// set the urls on the Help object and check they come back properly
	help := Help{
		Root:        "https://stacks-cli.help/root",
		Scaffold:    "https://stacks-cli.help/scaffold",
		Interactive: "https://stacks-cli.help/interactive",
		Version:     "https://stacks-cli.help/version",
	}

	tables := []struct {
		command string
		test    string
	}{
		{
			"root",
			help.Root,
		},
		{
			"scaffold",
			help.Scaffold,
		},
		{
			"interactive",
			help.Interactive,
		},
		{
			"version",
			help.Version,
		},
		{
			"notset",
			"",
		},
	}

	for _, table := range tables {

		// get the url from the command check it matches the test
		res := help.GetUrl(table.command)

		if res != table.test {
			t.Errorf("URL is not set properly: %s", res)
		}
	}
}
