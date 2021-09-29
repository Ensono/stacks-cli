package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupportedPipelines(t *testing.T) {

	// create the pipeline object
	pipeline := Pipeline{}

	// get a list of the supported pipelines
	pipelines := pipeline.GetSupported()

	// check the length of the slice
	assert.Equal(t, 1, len(pipelines))
}

func TestPipelineIsValid(t *testing.T) {

	// create a pipeline object
	pipeline := Pipeline{}

	// test that AzDo is a valid pipeline
	actual := pipeline.IsSupported("AzDo")

	assert.Equal(t, true, actual)
}

func TestPipelineIsNotValid(t *testing.T) {

	// create a pipeline object
	pipeline := Pipeline{}

	// test that AzDo is a valid pipeline
	actual := pipeline.IsSupported("new-pipeline")

	assert.Equal(t, false, actual)
}
