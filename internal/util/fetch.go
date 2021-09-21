package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GitClone uses standard network library to fetch a defined commit and avoids bloating the binary
func GitClone(repoUrl, commitHash, tmpPath string) error {

	// get the URL to be used to clone the repo from
	archiveUrl := ArchiveUrl(repoUrl, commitHash)

	resp, err := http.Get(archiveUrl)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("StatusCode: %d", resp.StatusCode)
	}

	zip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// write the contents of the HTTP Get to a zip file
	zipPath := filepath.Join(os.TempDir(), RandomString(7))
	if err := os.WriteFile(zipPath, zip, os.FileMode(0777)); err != nil {
		return err
	}

	// unzip the downloaded files to the tempdir for the project
	err = Unzip(zipPath, tmpPath)
	if err != nil {
		return err
	}

	return nil

}

// ArchiveUrl returns the archive url for the repo at a given commit hash or branch or v release
func ArchiveUrl(repoUrl, commitHash string) string {

	// of the commitHash is empty, set as master
	if commitHash == "" {
		commitHash = "master"
	}

	return strings.Join([]string{strings.TrimSuffix(repoUrl, ".git"), fmt.Sprintf("archive/%s.zip", commitHash)}, "/")
}
