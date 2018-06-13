package scanner

import (
	nais "github.com/nais/naisd/api"
	yaml "gopkg.in/yaml.v2"
)

type Result struct {
	err error
}

type Yaml map[interface{}]interface{}

func ParseYaml(data []byte) (nais.NaisManifest, error) {
	manifest := nais.NaisManifest{}
	err := yaml.Unmarshal(data, &manifest)
	return manifest, err
}

func AppConfig(yml string) Result {
	return Result{}
}
