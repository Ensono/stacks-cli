package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func CallHTTPAPI(url string, token string) (map[string]interface{}, error) {

	// create the data map to hold the information
	var data map[string]interface{}
	var err error

	// create a client to make the http request
	// this so the headers can be added if required
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, fmt.Errorf("unable to create HTTP request for '%s': %s", url, err.Error())
	}

	// if the token is not null, add the headers
	if token != "" {
		req.Header = http.Header{
			"Authorization": []string{fmt.Sprintf("token %s", token)},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return data, fmt.Errorf("unable to access requested API:\n\tURL: %s\n\t%s", url, err.Error())
	}

	// read all of the data returned in the call
	body, _ := io.ReadAll(resp.Body)

	// unmarshal the data into the map
	err = json.Unmarshal(body, &data)

	if err != nil {
		return data, fmt.Errorf("unable to read the data from the API: %s", err.Error())
	}

	if resp.StatusCode == 403 {

		// forbidden likely suggests that API access has been rate limited for the hour
		err = fmt.Errorf("error from HTTP endpoint: %s", data["message"])
	}

	return data, err
}

func GetGitHubArchiveUrl(path string, token string) (string, error) {

	var result string

	// call the function to get information from the API
	res, err := CallHTTPAPI(path, token)

	// ensure that a zipball_url exists in the map
	value, containsKey := res["zipball_url"]
	if containsKey {
		result = value.(string)
	}

	return result, err
}

func BuildGitHubAPIUrl(repoUrl string, ref string, trunk string, archive bool, token string) string {

	var url string

	// if the token is empty, then force the use of archive
	if token == "" {
		archive = true
	}

	// get the owner and repo name from the repoUrl
	ownerRepoName := strings.Replace(repoUrl, "https://github.com/", "", -1)

	// if hte ref has been set as latest or not set at all, build url to get the latest
	// release of the repository
	if archive {
		if ref == "latest" || ref == "" {
			ref = trunk
		}

		url = strings.Join([]string{strings.TrimSuffix(repoUrl, ".git"), fmt.Sprintf("archive/%s.zip", ref)}, "/")
	} else {
		if ref == "latest" || ref == "" {
			url = fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", ownerRepoName)
		} else {
			url = fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/%s", ownerRepoName, ref)
		}
	}

	return url
}
