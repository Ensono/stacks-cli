package interfaces

type Downloader interface {
	Get() (string, error)
	PackageURL() string
}
