package config

import (
	"fmt"

	"amido.com/stacks-cli/internal/config/static"
	"gopkg.in/yaml.v2"
)

// ParseSpecific source files
func ParseLocalSpecific(platform, deployment, _type string) (TypeDetail, error) {
	t := TypeDetail{}
	data := static.Config(fmt.Sprintf("%s_%s_%s", platform, deployment, _type))
	err := yaml.Unmarshal([]byte(data), &t)
	return t, err
}

func ParseLocalShared() (TypeDetail, error) {
	t := TypeDetail{}
	data := static.Config("shared")
	err := yaml.Unmarshal([]byte(data), &t)
	return t, err
}

// TODO: generata config
// Will Generate a config yaml from
func Generate(config Config) {

}
