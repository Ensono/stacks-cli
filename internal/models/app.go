package models

import (
	"fmt"
	"os"

	"github.com/Ensono/stacks-cli/internal/constants"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// App contains the objects from which the application will work,
// such as the application logger
type App struct {
	Logger *log.Logger
	Help   Help
}

func (app *App) LoadHelp(data []byte) {
	err := yaml.Unmarshal(data, &app.Help)
	if err != nil {
		log.Printf("Unable to parse help data, help messages will be missing: %s", err.Error())
	}
}

// ConfigureLogging sets up the logging for the application
// When running Debug Mode is will be configured to add the caller to the message
func (app *App) ConfigureLogging(logging Log) {

	// Initialise the Logger as a new Logger
	app.Logger = log.New()

	// Set the default to std out
	app.Logger.Out = os.Stdout

	// Attempt to set the logging level
	ll, err := log.ParseLevel(logging.Level)
	if err != nil {
		ll = log.DebugLevel
		app.HandleError(err, "fatal", "Unable to set logging level")
	}

	// if the log level is set to debug, add the caller to the messages
	if ll == log.TraceLevel {
		app.Logger.SetReportCaller(true)
	}

	// set the logging level
	app.Logger.SetLevel(ll)

	// set the format of the log messages
	switch logging.Format {
	case "json":
		app.Logger.SetFormatter(&log.JSONFormatter{
			TimestampFormat: constants.LoggingTimestamp,
		})
	default:
		app.Logger.SetFormatter(&log.TextFormatter{
			ForceColors:     logging.Colour,
			FullTimestamp:   false,
			TimestampFormat: constants.LoggingTimestamp,
		})

		app.Logger.SetOutput(colorable.NewColorableStdout())
	}

	// if a log file has been set redirect all logs to the file
	if logging.File != "" {
		file, err := os.OpenFile(logging.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			// app.Logger.Out = file
			app.Logger.SetOutput(file)
		} else {
			app.Log("LOG001", "warn", err.Error())
		}
	}
}

// HandleError handles errors from the application
// This is to ensure that all errors are handled in the same way
func (app *App) HandleError(err error, errorType string, msg ...string) {

	// if no messages have been add the default message
	if len(msg) == 0 {
		msg = append(msg, constants.DefaultErrorMessage)
	}

	var fields map[string]interface{}
	app.HandleErrorWithFields(err, errorType, msg[0], fields)
}

func (app *App) HandleErrorWithFields(err error, errorType string, msg string, fields map[string]interface{}) {

	var message string

	if err != nil {

		if fields == nil {
			fields = make(map[string]interface{})
		}

		// set the message
		if msg != "" {
			message = msg
		} else {
			message = constants.DefaultErrorMessage
		}

		// set the fields that need to be set in the message
		if errorType == "error" ||
			errorType == "fatal" {

			fields["error"] = err
		}

		switch errorType {
		case "error":
			app.Logger.WithFields(fields).Error(message)
		case "fatal":
			app.Logger.WithFields(fields).Fatal(message)
		case "warn":
			app.Logger.WithFields(fields).Warn(message)
		}
	}
}

func (app *App) Log(data string, level string, replacements ...interface{}) {

	// get the message from help, if it does not exist use the data as is
	message := app.Help.GetMessage(data)
	if message == "" {
		message = data
	}

	// perform the subsitutions on the message
	message = fmt.Sprintf(message, replacements...)

	// output the message with the correct level
	switch level {
	case "info":
		app.Logger.Info(message)
	case "error":
		app.Logger.Error(message)
	case "warn":
		app.Logger.Warn(message)
	case "fatal":
		app.Logger.Fatal(message)
	case "debug":
		app.Logger.Debug(message)
	}
}
