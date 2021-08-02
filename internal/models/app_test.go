package models

import (
	"errors"
	"testing"

	"github.com/amido/stacks-cli/internal/constants"
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
