package util

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// DownloadFile attempts to download the specified file to the temporary directory
// and return the path
func DownloadFile(url *url.URL, dir string) (string, error) {

	var err error

	// determine the filename of the downloaded filepath
	filepath := filepath.Join(dir, filepath.Base(url.Path))

	// get the data from the specified url
	resp, err := http.Get(url.String())
	if err != nil {
		return filepath, err
	}
	defer resp.Body.Close()

	// write the data to the file
	out, err := os.Create(filepath)
	if err != nil {
		return filepath, err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return filepath, err
}
