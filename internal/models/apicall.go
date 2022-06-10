package models

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type APICall struct {
	url          string
	token        string
	raw          []byte
	downloadPath string
}

var apiCall = &APICall{}

func NewAPICall(url string, token string) *APICall {
	apiCall.url = url
	apiCall.token = token

	return apiCall
}

func (ac *APICall) Do(method string) (error, int) {
	var err error

	// create a client to make the http request
	// this so the headers can be added if required
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	req, err := http.NewRequest(method, ac.url, nil)
	if err != nil {
		return fmt.Errorf("unable to create HTTP request for '%s': %s", ac.url, err.Error()), 0
	}

	// if the token is not null, add the headers
	if ac.token != "" {
		req.Header = http.Header{
			"Authorization": []string{fmt.Sprintf("token %s", ac.token)},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to access requested API:\n\tURL: %s\n\t%s", ac.url, err.Error()), 0
	}
	defer resp.Body.Close()

	// read all of the data returned in the call
	// and save into a file if it has been specified
	if ac.downloadPath != "" {
		file, err := os.Create(ac.downloadPath)
		if err != nil {
			return err, 0
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err, 0
		}

	} else {
		ac.raw, _ = ioutil.ReadAll(resp.Body)
	}

	if err != nil {
		return fmt.Errorf("unable to read the data from the API: %s", err.Error()), 0
	}

	/*
		if resp.StatusCode == 403 {

			// forbidden likely suggests that API access has been rate limited for the hour
			err = fmt.Errorf("error from HTTP endpoint: %s", data["message"])
		}
	*/

	return err, resp.StatusCode
}

func (ac *APICall) Raw() []byte {
	return ac.raw
}

func (ac *APICall) Download(filepath string) error {
	var err error

	ac.downloadPath = filepath
	err, _ = ac.Do("GET")
	if err != nil {
		return err
	}

	// restet the downloadPath to null
	ac.downloadPath = ""

	return err
}

func (ac *APICall) UpdateURL(url string) {
	ac.url = url
}
