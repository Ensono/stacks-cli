package models

import (
	"errors"
	"testing"

	"github.com/Ensono/stacks-cli/internal/constants"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestHandleErrorWithDefaultMessage(t *testing.T) {

	logger, hook := test.NewNullLogger()

	app := App{
		Logger: logger,
	}

	// create a new error to send to the handler
	err := errors.New("Unit testing error message")

	app.HandleError(err, "error")

	assert.Equal(t, 1, len(hook.Entries), "Should be 1 error message")
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, constants.DefaultErrorMessage, hook.LastEntry().Message)

	assert.Contains(t, hook.LastEntry().Data, "error")

	hook.Reset()

}

func TestConfigureLogging(t *testing.T) {

	// create the configuration logging object
	logging := Log{
		Level: "info",
	}

	logger, hook := test.NewNullLogger()

	app := App{
		Logger: logger,
	}

	app.ConfigureLogging(logging)

	assert.Equal(t, "info", app.Logger.GetLevel().String())

	hook.Reset()
}

func TestHelpMessage(t *testing.T) {

	// Create YAML file which will be loaded into the app
	help_data := `
help:
  - name: TEST001
    value: This is a test help message
  - name: GEN001
    value: "Error in %s: %s"
`

	// build up the test tables
	tables := []struct {
		code     string
		subs     []interface{}
		expected string
	}{
		{
			"TEST001",
			[]interface{}{},
			"This is a test help message",
		},
		{
			"GEN001",
			[]interface{}{"export", "missing"},
			"Error in export: missing",
		},
	}

	// Create an App object
	app := App{}

	// load the help data into the app
	app.LoadHelp([]byte(help_data))

	// iterate around the test tables
	for _, table := range tables {
		message := app.Help.GetMessage(table.code, table.subs...)

		assert.Equal(t, table.expected, message, "The message should be '%s'", table.expected)
	}

}
