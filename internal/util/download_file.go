package util

import (
	"io"
	"net/http"
)

// DownloadFile provides a wrapper function to download a file into memory
// and return the []byte array of the downloaded item
func DownloadFile(url string) ([]byte, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
