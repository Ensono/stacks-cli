package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func CallGitHubAPI(url string) (map[string]interface{}, error) {

	// create the data map to hold the information
	var data map[string]interface{}
	var err error

	// attempt to call the API to get the requested information
	resp, err := http.Get(url)
	if err != nil {
		return data, fmt.Errorf("unable to access requested GitHub API:\n\tURL: %s\n\t%s", url, err.Error())
	}

	// read all of the data returned in the call
	body, _ := ioutil.ReadAll(resp.Body)

	// unmarshal the data into the map
	err = json.Unmarshal(body, &data)

	if err != nil {
		return data, fmt.Errorf("unable to read the data from the API: %s", err.Error())
	}

	if resp.StatusCode == 403 {

		// forbidden likely suggests that API access has been rate limited for the hour
		err = fmt.Errorf("error from GitHub: %s", data["message"])
	}

	return data, err
}

func GetGitHubArchiveUrl(path string) (string, error) {

	var result string

	// call the function to get information from the API
	res, err := CallGitHubAPI(path)

	// ensure that a zipball_url exists in the map
	value, containsKey := res["zipball_url"]
	if containsKey {
		result = value.(string)
	}

	return result, err
}

func BuildGitHubAPIUrl(repoUrl string, ref string, archive bool) string {

	var url string

	// get the owner and repo name from the repoUrl
	ownerRepoName := strings.Replace(repoUrl, "https://github.com/", "", -1)

	// if hte ref has been set as latest or not set at all, build url to get the latest
	// release of the repository
	if archive {
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
