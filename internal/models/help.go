package models

import "fmt"

type Help struct {
	Help []HelpMessage `mapstructure:"help"`
}

type HelpMessage struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

func (h *Help) GetMessage(id string, replacements ...interface{}) string {

	// declare the message variable
	var message string

	// find the message in the help slice
	for _, item := range h.Help {
		if item.Name == id {
			message = item.Value
			break
		}
	}

	message = fmt.Sprintf(message, replacements...)

	return message
}
