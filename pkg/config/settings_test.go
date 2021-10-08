package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupPipelines(t *testing.T) Settings {
	t.Log("Setting up pipeline settings")

	// setup the pipelines
	pipelines := make([]Pipeline, 1)

	pipelines[0] = Pipeline{
		Type: "azdo",
	}

	settings := Settings{}
	settings.Pipeline = pipelines

	return settings
}

func TestPipelineExists(t *testing.T) {

	// get the settings from the setup
	settings := setupPipelines(t)

	assert.NotEqual(t, Pipeline{}, settings.GetPipeline("azdo"))
}

func TestPipelineDoesNotExist(t *testing.T) {

	// get the settings from the setup
	settings := setupPipelines(t)

	assert.Equal(t, Pipeline{}, settings.GetPipeline("jenkins"))
}
