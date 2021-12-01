package util

import (
	"fmt"
	"testing"
)

// declare the repoUrl that will be used for all the tests
var repoUrl string = "https://github.com/amido/stacks-dotnet"

func TestBuildGitHubAPIUrl(t *testing.T) {

	// build test table
	tables := []struct {
		ref     string
		test    string
		msg     string
		archive bool
	}{
		{
			"",
			"https://api.github.com/repos/amido/stacks-dotnet/releases/latest",
			"An empty ref should return the latest release URL",
			false,
		},
		{
			"latest",
			"https://api.github.com/repos/amido/stacks-dotnet/releases/latest",
			"Specifying latest ref should return the latest release URL",
			false,
		},
		{
			"v3.0.232",
			"https://api.github.com/repos/amido/stacks-dotnet/releases/tags/v3.0.232",
			"A specified tag should return the release for that tag",
			false,
		},
		{
			"feature/dotnet-6",
			fmt.Sprintf("%s/archive/feature/dotnet-6.zip", repoUrl),
			"A branch can be specified if the archive flag is used",
			true,
		},
	}

	// iterate around the test table
	for _, table := range tables {

		// get the ghUrl from the method
		ghUrl := BuildGitHubAPIUrl(repoUrl, table.ref, table.archive)

		if ghUrl != table.test {
			t.Error(table.msg)
		}
	}
}
