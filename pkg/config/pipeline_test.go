package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupPipelineTests(t *testing.T, name string) (func(t *testing.T), string) {

	// create a temporary directory
	tempDir := t.TempDir()

	// create a build file to work with
	pipeline := `
stages:
- stage: Build
	variables:
	- group: amido-stacks-infra-credentials-nonprod
	- group: stacks-credentials-nonprod-kv
	- group: amido-stacks-webapp
	`

	// write the file to the tempDir
	err := os.WriteFile(filepath.Join(tempDir, name), []byte(pipeline), 0666)
	if err != nil {
		t.Logf("[ERROR] Error creating file: %v", err)
	}

	deferFunc := func(t *testing.T) {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Logf("[ERROR] Unable to remove dir: %v", err)
		}
	}

	return deferFunc, tempDir
}

func TestSupportedPipelines(t *testing.T) {

	// create the pipeline object
	pipeline := Pipeline{}

	// get a list of the supported pipelines
	pipelines := pipeline.GetSupported()

	// check the length of the slice
	assert.Equal(t, 1, len(pipelines))
}

func TestAzDoPipelineIsValid(t *testing.T) {

	// create a pipeline object
	pipeline := Pipeline{}

	// test that AzDo is a valid pipeline
	actual := pipeline.IsSupported("AzDo")

	assert.Equal(t, true, actual)
}

func TestGHAPipelineIsValid(t *testing.T) {

	// create a pipeline object
	pipeline := Pipeline{}

	// test that AzDo is a valid pipeline
	actual := pipeline.IsSupported("gha")

	assert.Equal(t, true, actual)
}

func TestPipelineIsNotValid(t *testing.T) {

	// create a pipeline object
	pipeline := Pipeline{}

	// test that AzDo is a valid pipeline
	actual := pipeline.IsSupported("new-pipeline")

	assert.Equal(t, false, actual)
}

func TestReplacePatterns(t *testing.T) {

	// set the name of the build file
	name := "build.yml"

	// setup the environment
	cleanup, dir := setupPipelineTests(t, name)
	defer cleanup(t)

	buildFile := filepath.Join(dir, name)

	// create the replacements that need to be performed
	replacements := make([]PipelineReplacement, 1)
	replacements[0] = PipelineReplacement{
		Pattern: `^.*stacks-credentials-nonprod-kv$`,
		Value:   "",
	}

	// create the file list
	list := make([]PipelineFile, 1)
	list[0] = PipelineFile{
		Name: "build",
		Path: name,
	}

	// create the pipeline settings
	pipeline := Pipeline{
		File:         list,
		Replacements: replacements,
	}

	// call the function
	errs := pipeline.ReplacePatterns(dir)

	// set the expected value of the contents of the file
	expected := `
stages:
- stage: Build
	variables:
	- group: amido-stacks-infra-credentials-nonprod

	- group: amido-stacks-webapp
	`

	// read in the contents of the build file, which should have been modified
	actual, _ := ioutil.ReadFile(buildFile)

	// check that there are no errors
	for _, err := range errs {
		assert.Equal(t, nil, err)
	}
	assert.Equal(t, expected, string(actual))
}

func TestReplacePatternsWithNoReplace(t *testing.T) {

	// set the name of the build file
	name := "build.yml"

	// setup the environment
	cleanup, dir := setupPipelineTests(t, name)
	defer cleanup(t)

	buildFile := filepath.Join(dir, name)

	// create the replacements that need to be performed
	replacements := make([]PipelineReplacement, 1)
	replacements[0] = PipelineReplacement{
		Pattern: `^.*stacks-credentials-nonprod-kv$`,
		Value:   "",
	}

	// create the file list
	list := make([]PipelineFile, 1)
	list[0] = PipelineFile{
		Name:      "build",
		Path:      name,
		NoReplace: true,
	}

	// create the pipeline settings
	pipeline := Pipeline{
		File:         list,
		Replacements: replacements,
	}

	// call the function
	errs := pipeline.ReplacePatterns(dir)

	// set the expected value of the contents of the file
	expected := `
stages:
- stage: Build
	variables:
	- group: amido-stacks-infra-credentials-nonprod
	- group: stacks-credentials-nonprod-kv
	- group: amido-stacks-webapp
	`

	// read in the contents of the build file, which should have been modified
	actual, _ := ioutil.ReadFile(buildFile)

	// check that there are no errors
	for _, err := range errs {
		assert.Equal(t, nil, err)
	}
	assert.Equal(t, expected, string(actual))
}
