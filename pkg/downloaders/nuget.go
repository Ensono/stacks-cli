package downloaders

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/amido/stacks-cli/internal/models"
	"github.com/amido/stacks-cli/internal/util"
)

var nuget = &Nuget{}

type Nuget struct {
	Name             string
	ID               string
	Version          string
	FrameworkVersion string
	TempDir          string
	CacheDir         string

	// define private properties
	url    string
	latest bool
}

type NugetResponse struct {
	Items          []NugetResponseItem "json:`items`"
	PackageContent string              "json:`packageContent`"
}

type NugetResponseItem struct {
	Count int                      "json:`count`"
	Items []NugetResponsePageItems "json:`items`"
}

type NugetResponsePageItems struct {
	CatalogEntry NugetItemCatalogEntry "json:`catalogEntry`"
}

type NugetItemCatalogEntry struct {
	ID             string "json:`id`"
	PackageContent string "json:`packageContent`"
	Version        string "json:`version`"
}

func NewNugetDownloader(name string, id string, version string, frameworkVersion string, cacheDir string, tempDir string) *Nuget {
	nuget.Name = name
	nuget.ID = id
	nuget.Version = version
	nuget.FrameworkVersion = frameworkVersion
	nuget.CacheDir = cacheDir
	nuget.TempDir = tempDir

	if nuget.Version == "" || strings.ToLower(nuget.Version) == "latest" {
		nuget.latest = true
	}

	return nuget
}

// Get downloads the specified, or latest, version of the named package from Nuget
func (n *Nuget) Get() (string, error) {

	// define function variables
	var dir string
	var err error

	// initalise an API call
	ac := models.NewAPICall("", "")

	// if a version has not been set, e.g. the version is latest then get the latest version and
	// set as the version on the model
	_ = n.getLatestVersion(ac)

	// Set the URL to use to get package and version information
	n.setURL()
	ac.UpdateURL(n.url)

	// get the data from the Nuget API
	err, statusCode := ac.Do("GET")
	if err != nil {
		return dir, fmt.Errorf("problem calling the Nuget api: %s", err.Error())
	}

	// check the statuscode of the response
	if statusCode > 299 {
		return "", fmt.Errorf("error downloading package: %d", statusCode)
	}

	uri, _ := n.getPackageURL(ac)
	n.url = uri
	ac.UpdateURL(uri)

	// get the name of the file from the URL
	u, _ := url.Parse(uri)
	_, filename := path.Split(u.Path)

	// determine the path for the downloaded file, this will go into the cachedir
	downloadPath := filepath.Join(n.CacheDir, filename)

	// if the file does not exist download it
	if !util.Exists(downloadPath) {
		err = ac.Download(downloadPath)
		if err != nil {
			return dir, err
		}
	}

	// Unpack the Nuget package into the TempDir
	// but ensure that there is a top level dir to work with in the tempDir as
	// all projects get unpacked into here
	unpackDir := filepath.Join(n.TempDir, strings.ToLower(n.Name))
	_, err = util.Unzip(downloadPath, unpackDir)
	if err != nil {
		return "", err
	}

	// The package will unpack into a different folder structure than if had been retrieved
	// from github. This nested path needs to be set on the returned dir so that the CLI
	// can find the settings file for the project
	templateDir := filepath.Join(unpackDir, "content", "templates", n.ID)

	return templateDir, err
}

func (n *Nuget) PackageURL() string {
	return n.url
}

func (n *Nuget) setURL() {

	// based on the Version that has been specified, set the name of the file that should be
	// downloaded
	var file string = "index.json"
	if !n.latest {
		file = fmt.Sprintf("%s.json", n.Version)
	}

	n.url = fmt.Sprintf("https://api.nuget.org/v3/registration5-semver1/%s/%s", strings.ToLower(n.Name), file)
}

// getPackageURL gets the URL for the package according to the version that has been requested
func (n *Nuget) getPackageURL(ac *models.APICall) (string, error) {

	var url string

	// unmarshal the raw data into the NugetResponse
	var resp NugetResponse
	err := json.Unmarshal(ac.Raw(), &resp)

	if err != nil {
		return url, fmt.Errorf("issue reading response body from Nuget: %s", err.Error())
	}

	// if getting the latest version, get the last item in the Items[].Items array
	// otherwise get the package content
	if n.latest {
		url = resp.Items[0].Items[len(resp.Items[0].Items)-1].CatalogEntry.PackageContent
	} else {
		url = resp.PackageContent
	}

	return url, err
}

func (n *Nuget) getLatestVersion(ac *models.APICall) error {

	var err error
	var data map[string][]string

	// return of not looking for latest
	if !nuget.latest {
		return err
	}

	// set the url to get the list of versions for the package
	ac.UpdateURL(fmt.Sprintf("https://api.nuget.org/v3-flatcontainer/%s/index.json", strings.ToLower(n.Name)))

	err, _ = ac.Do("GET")
	if err != nil {
		return fmt.Errorf("problem calling the Nuget api: %s", err.Error())
	}

	// unmarshal the data from the raw data into a map
	err = json.Unmarshal(ac.Raw(), &data)
	if err != nil {
		return fmt.Errorf("problem retrieving list of versions: %s", err.Error())
	}

	// set the version to the last in the list
	n.Version = data["versions"][len(data["versions"])-1]
	n.latest = false

	return err
}
