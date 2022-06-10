package interfaces

import "github.com/sirupsen/logrus"

type Downloader interface {
	Get() (string, error)
	PackageURL() string
	SetLogger(logger *logrus.Logger)
}
