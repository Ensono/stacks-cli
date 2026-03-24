package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// declare the repoUrl that will be used for all the tests
var repoUrl string = "https://github.com/Ensono/stacks-dotnet"
var token string = "ghjkgjhgjhgj"

func TestBuildGitHubAPIUrl(t *testing.T) {

	// build test table
	tables := []struct {
		ref     string
		trunk   string
		test    string
		msg     string
		archive bool
		token   string
	}{
		{
			"",
			"master",
			"https://api.github.com/repos/Ensono/stacks-dotnet/releases/latest",
			"An empty ref should return the latest release URL",
			false,
			token,
		},
		{
			"",
			"master",
			fmt.Sprintf("%s/archive/master.zip", repoUrl),
			"An empty ref with no token should return the latest release URL using the archive URL",
			false,
			"",
		},
		{
			"latest",
			"",
			"https://api.github.com/repos/Ensono/stacks-dotnet/releases/latest",
			"Specifying latest ref should return the latest release URL",
			false,
			token,
		},
		{
			"v3.0.232",
			"",
			fmt.Sprintf("%s/archive/v3.0.232.zip", repoUrl),
			"A specified tag with no token should return the release for that tag using the archive URL",
			false,
			"",
		},
		{
			"v3.0.232",
			"",
			"https://api.github.com/repos/Ensono/stacks-dotnet/releases/tags/v3.0.232",
			"A specified tag should return the release for that tag",
			false,
			token,
		},
		{
			"feature/dotnet-6",
			"",
			fmt.Sprintf("%s/archive/feature/dotnet-6.zip", repoUrl),
			"A branch can be specified if the archive flag is used",
			true,
			"",
		},
	}

	// iterate around the test table
	for _, table := range tables {

		// get the ghUrl from the method
		ghUrl := BuildGitHubAPIUrl(repoUrl, table.ref, table.trunk, table.archive, table.token)

		if ghUrl != table.test {
			t.Error(table.msg)
		}
	}
}

func TestGetLatestReleaseTagFromURL(t *testing.T) {

	tables := []struct {
		tagName     string
		expected    string
		msg         string
		expectError bool
	}{
		{
			"v0.3.0",
			"0.3.0",
			"Tag with v prefix should have it stripped",
			false,
		},
		{
			"0.3.0",
			"0.3.0",
			"Tag without v prefix should be returned as-is",
			false,
		},
		{
			"v1.2.3-beta.1",
			"1.2.3-beta.1",
			"Pre-release tag with v prefix should have it stripped",
			false,
		},
	}

	for _, table := range tables {

		// create a mock server that returns the tag_name
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]interface{}{
				"tag_name": table.tagName,
				"url":      "https://api.github.com/repos/test/test/releases/1",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		version, err := getLatestReleaseTagFromURL(server.URL, "")
		if table.expectError && err == nil {
			t.Errorf("%s: expected error but got none", table.msg)
		}
		if !table.expectError && err != nil {
			t.Errorf("%s: unexpected error: %s", table.msg, err.Error())
		}
		if version != table.expected {
			t.Errorf("%s: expected '%s' but got '%s'", table.msg, table.expected, version)
		}
	}
}

func TestGetLatestReleaseTagFromURL_NoRelease(t *testing.T) {

	// create a mock server that returns no tag_name
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"message": "Not Found",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	_, err := getLatestReleaseTagFromURL(server.URL, "")
	if err == nil {
		t.Error("Expected error when no tag_name in response")
	}
}
