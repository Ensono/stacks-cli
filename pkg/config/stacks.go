package config

type Stacks struct {
	Dotnet Dotnet `mapstructure:"dotnet"`
	Java   Java   `mapstructure:"java"`
	NodeJS NodeJS `mapstructure:"nodejs"`
	Infra  Infra  `mapstructure:"infra"`
}

type Dotnet struct {
	Webapi string `mapstructure:"webapi"`
	CQRS   string `mapstructure:"cqrs"`
	Events string `mapstructure:"events"`
}

type Java struct {
	Webapi string `mapstructure:"webapi"`
	CQRS   string `mapstructure:"cqrs"`
	Events string `mapstructure:"events"`
}

type NodeJS struct {
	CSR string `mapstructure:"csr"`
	SSR string `mapstructure:"ssr"`
}

type Infra struct {
	AKS string `mapstructure:"aks"`
}

// GetSrcURLMap returns a map of the source control repositores
func (stacks *Stacks) GetSrcURLMap() map[string]string {

	srcUrls := map[string]string{
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

func (stacks *Stacks) GetSrcURL(key string) string {
	srcUrls := stacks.GetSrcURLMap()
	return srcUrls[key]
}
