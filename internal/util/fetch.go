package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// GitClone uses standard network library to fetch a defined commit and avoids bloating the binary
func GitClone(repoUrl, ref, tmpPath string, token string) (string, error) {

	// get the URL to be used to clone the repo from
	archiveUrl, err := ArchiveUrl(repoUrl, ref, token)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(archiveUrl)
	if err != nil {
		return "", err
	}

	if resp.StatusCode > 299 {
		return archiveUrl, fmt.Errorf("StatusCode: %d", resp.StatusCode)
	}

	zip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// write the contents of the HTTP Get to a zip file
	zipPath := filepath.Join(os.TempDir(), RandomString(7))
	if err := os.WriteFile(zipPath, zip, os.FileMode(0777)); err != nil {
		return "", err
	}

	// unzip the downloaded files to the tempdir for the project
	tempRepoDir, err := Unzip(zipPath, tmpPath)
	if err != nil {
		return "", err
	}

	// remove the zip file
	_ = os.Remove(zipPath)

	return tempRepoDir, nil
}

// ArchiveUrl returns the archive url for the repo at a given commit hash or branch or v release
func ArchiveUrl(repoUrl, ref string, token string) (string, error) {

	var zipUrl string
	var err error

	// get the apiUrl from the method
	apiUrl := BuildGitHubAPIUrl(repoUrl, ref, false, token)

	if token == "" {
		zipUrl = apiUrl
	} else {
		// call the github api to get the url to the zip file to download
		zipUrl, err = GetGitHubArchiveUrl(apiUrl, token)
	}

	// if the zipUrl has not been found then drop back to the archive URL
	if zipUrl == "" && err == nil {
		apiUrl = BuildGitHubAPIUrl(repoUrl, ref, true, token)
		zipUrl, err = GetGitHubArchiveUrl(apiUrl, token)
	}

	return zipUrl, err
}
