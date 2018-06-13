package scanner

import yaml "gopkg.in/yaml.v2"

type Result struct {
	err error
}

type Yaml map[interface{}]interface{}

func ParseYaml(data []byte) (Yaml, error) {
	m := make(Yaml)
	err := yaml.Unmarshal(data, &m)
	return m, err
}

func AppConfig(yml string) Result {
	return Result{}
}
