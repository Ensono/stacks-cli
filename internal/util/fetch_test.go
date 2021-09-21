package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Define a base URL that all tests will use
var repoUrl string = "https://github.com/org/repo"

func TestGitCloneWithoutRef(t *testing.T) {
	// set the values for the repoUrl and the ref
	ref := ""

	// define the expected value
	expected := fmt.Sprintf("%s/archive/master.zip", repoUrl)

	// get the actualUrl
	actual := ArchiveUrl(repoUrl, ref)

	assert.Equal(t, expected, actual)
}

func TestGitCloneWithSHARef(t *testing.T) {
	// set the values for the repoUrl and the has
	ref := "DED8A6B16DE379DDB54F242C930F1E8650308888"

	// define the expected value
	expected := fmt.Sprintf("%s/archive/%s.zip", repoUrl, ref)

	// get the actual calculated Url
	actual := ArchiveUrl(repoUrl, ref)

	assert.Equal(t, expected, actual)

}

func TestGitCloneWithBranch(t *testing.T) {
	// set the values for the repoUrl and the has
	ref := "feature/my-new-one"

	// define the expected value
	expected := fmt.Sprintf("%s/archive/%s.zip", repoUrl, ref)

	// get the actual calculated Url
	actual := ArchiveUrl(repoUrl, ref)

	assert.Equal(t, expected, actual)

}

func TestGitCloneWithTag(t *testing.T) {
	// set the values for the repoUrl and the has
	ref := "v0.0.200"

	// define the expected value
	expected := fmt.Sprintf("%s/archive/%s.zip", repoUrl, ref)

	// get the actual calculated Url
	actual := ArchiveUrl(repoUrl, ref)

	assert.Equal(t, expected, actual)

}
