package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// GitClone uses standard network library to fetch a defined commit and avoids bloating the binary
func GitClone(repoUrl, commitHash, tmpPath, zipPath string) error {

	resp, err := http.Get(ArchiveUrl(repoUrl, commitHash))
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

	if err := os.WriteFile(zipPath, zip, os.FileMode(0777)); err != nil {
		return err
	}

	return nil

}

// ArchiveUrl returns the archive url for the repo at a given commit hash or branch or v release
func ArchiveUrl(repoUrl, commitHash string) string {
	return strings.Join([]string{strings.TrimSuffix(repoUrl, ".git"), fmt.Sprintf("archive/%s.zip", commitHash)}, "/")
}

func Unzip(scr, dest string) {

}
