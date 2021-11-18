package models

type Command struct {
	Framework       string
	Binary          string
	VersionFound    string
	VersionRequired string
	Message         string
}
