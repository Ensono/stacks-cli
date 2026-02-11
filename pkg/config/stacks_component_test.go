package config

import (
	"testing"

	yaml "github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

func TestStacksComponentTemplateModeDefaults(t *testing.T) {

	// Test case 1: TemplateMode not specified in config (should default to true)
	configWithoutTemplateMode := []byte(`
stacks:
  components:
    dotnet_webapi:
      group: dotnet
      name: webapi
      package:
        url: https://github.com/example/stacks-dotnet
        version: main
        type: git
`)

	var config1 Config
	config1.Init()

	err := yaml.Unmarshal(configWithoutTemplateMode, &config1)
	assert.NoError(t, err)

	// Before calling SetDefaultValues, TemplateMode should be nil
	component := config1.Stacks.Components["dotnet_webapi"]
	assert.Nil(t, component.TemplateMode)

	// But the helper method should return true (default)
	assert.True(t, component.IsTemplateModeEnabled())

	// After calling SetDefaultValues, it should be set to true
	config1.SetDefaultValues()
	component = config1.Stacks.Components["dotnet_webapi"]
	assert.NotNil(t, component.TemplateMode)
	assert.True(t, *component.TemplateMode)
	assert.True(t, component.IsTemplateModeEnabled())

	// Test case 2: TemplateMode explicitly set to false (should remain false)
	configWithTemplateModeExplicitFalse := []byte(`
stacks:
  components:
    dotnet_webapi:
      group: dotnet
      name: webapi
      package:
        url: https://github.com/example/stacks-dotnet
        version: main
        type: git
      template_mode: false
`)

	var config2 Config
	config2.Init()

	err = yaml.Unmarshal(configWithTemplateModeExplicitFalse, &config2)
	assert.NoError(t, err)

	// Before calling SetDefaultValues, TemplateMode should be false
	component2 := config2.Stacks.Components["dotnet_webapi"]
	assert.NotNil(t, component2.TemplateMode)
	assert.False(t, *component2.TemplateMode)
	assert.False(t, component2.IsTemplateModeEnabled())

	// After calling SetDefaultValues, it should still be false
	config2.SetDefaultValues()
	component2 = config2.Stacks.Components["dotnet_webapi"]
	assert.NotNil(t, component2.TemplateMode)
	assert.False(t, *component2.TemplateMode)
	assert.False(t, component2.IsTemplateModeEnabled())

	// Test case 3: TemplateMode explicitly set to true (should remain true)
	configWithTemplateModeExplicitTrue := []byte(`
stacks:
  components:
    dotnet_webapi:
      group: dotnet
      name: webapi
      package:
        url: https://github.com/example/stacks-dotnet
        version: main
        type: git
      template_mode: true
`)

	var config3 Config
	config3.Init()

	err = yaml.Unmarshal(configWithTemplateModeExplicitTrue, &config3)
	assert.NoError(t, err)

	// Before calling SetDefaultValues, TemplateMode should be true
	component3 := config3.Stacks.Components["dotnet_webapi"]
	assert.NotNil(t, component3.TemplateMode)
	assert.True(t, *component3.TemplateMode)
	assert.True(t, component3.IsTemplateModeEnabled())

	// After calling SetDefaultValues, it should still be true
	config3.SetDefaultValues()
	component3 = config3.Stacks.Components["dotnet_webapi"]
	assert.NotNil(t, component3.TemplateMode)
	assert.True(t, *component3.TemplateMode)
	assert.True(t, component3.IsTemplateModeEnabled())
}

func TestStacksComponentIsTemplateModeEnabledHelper(t *testing.T) {

	// Test with nil TemplateMode (should return default true)
	component1 := StacksComponent{
		Group:        "dotnet",
		Name:         "webapi",
		TemplateMode: nil,
	}
	assert.True(t, component1.IsTemplateModeEnabled())

	// Test with explicit true
	trueValue := true
	component2 := StacksComponent{
		Group:        "dotnet",
		Name:         "webapi",
		TemplateMode: &trueValue,
	}
	assert.True(t, component2.IsTemplateModeEnabled())

	// Test with explicit false
	falseValue := false
	component3 := StacksComponent{
		Group:        "dotnet",
		Name:         "webapi",
		TemplateMode: &falseValue,
	}
	assert.False(t, component3.IsTemplateModeEnabled())
}
