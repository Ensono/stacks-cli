package config

type Stacks struct {
	Dotnet Dotnet `mapstructure:"dotnet"`
	Java   Java   `mapstructure:"java"`
	NodeJS NodeJS `mapstructure:"nodejs"`
	Infra  Infra  `mapstructure:"infra"`
}

type Dotnet struct {
	Webapi RepoInfo `mapstructure:"webapi"`
	CQRS   RepoInfo `mapstructure:"cqrs"`
	Events RepoInfo `mapstructure:"events"`
}

type Java struct {
	Webapi RepoInfo `mapstructure:"webapi"`
	CQRS   RepoInfo `mapstructure:"cqrs"`
	Events RepoInfo `mapstructure:"events"`
}

type NodeJS struct {
	CSR RepoInfo `mapstructure:"csr"`
	SSR RepoInfo `mapstructure:"ssr"`
}

type Infra struct {
	AKS RepoInfo `mapstructure:"aks"`
}

type RepoInfo struct {
	URL   string `mapstructure:"url"`
	Trunk string `mapstructure:"trunk"`
}

// GetSrcURLMap returns a map of the source control repositores
func (stacks *Stacks) GetSrcURLMap() map[string]RepoInfo {

	srcUrls := map[string]RepoInfo{
		"dotnet_webapi": stacks.Dotnet.Webapi,
		"dotnet_cqrs":   stacks.Dotnet.CQRS,
		"dotnet_events": stacks.Dotnet.Events,
		"java_webapi":   stacks.Java.Webapi,
		"java_cqrs":     stacks.Java.CQRS,
		"java_events":   stacks.Java.Events,
		"nodejs_csr":    stacks.NodeJS.CSR,
		"nodejs_ssr":    stacks.NodeJS.SSR,
		"infra_aks":     stacks.Infra.AKS,
	}

	return srcUrls
}

func (stacks *Stacks) GetSrcURL(key string) RepoInfo {
	srcUrls := stacks.GetSrcURLMap()
	return srcUrls[key]
}
