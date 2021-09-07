package config

type Stacks struct {
	Dotnet Dotnet `mapstructure:"dotnet"`
	Java   Java   `mapstructure:"java"`
	NodeJS NodeJS `mapstructure:"nodejs"`
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
