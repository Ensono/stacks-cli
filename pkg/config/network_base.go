package config

type NetworkBase struct {
	Domain DomainType `mapstructure:"domain"`
}

type DomainType struct {
	Internal string `mapstructure:"internal"`
	External string `mapstructure:"external"`
}
