package util

import (
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
)

func setupDownloadFileTestCase(t *testing.T) (func(t *testing.T), string) {
	// create a temporary directory
	tempDir := t.TempDir()

	deferFunc := func(t *testing.T) {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Logf("[ERROR] Unable to remove dir: %v", err)
		}
	}

	return deferFunc, tempDir
}

func TestDownloadFile(t *testing.T) {

	urlLocation := "https://stacks-cli.ensonodigital.com/files/config.yml"

	// create a mock for the HTTP methods
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// register a responder for the call to http.Get
	httpmock.RegisterResponder(
		"GET",
		urlLocation,
		httpmock.NewStringResponder(200, `
stacks:
  components:
    dotnet_webapi:
    group: dotnet
    name: webapi
    package:
      name: Amido.Stacks.Templates
      type: nuget
      id: stacks-dotnet
    `),
	)

	cleanup, tempDir := setupDownloadFileTestCase(t)
	defer cleanup(t)

	// parse the rul;
	u, _ := url.ParseRequestURI(urlLocation)

	// call the download function
	file, err := DownloadFile(u, tempDir)

	if err != nil {
		t.Errorf("Unable to download file: %s", err.Error())
	}

	if !Exists(file) {
		t.Errorf("File does not exist: %s", file)
	}

}
