package scanner

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Result struct {
	err error
}

type Transformer func(Yaml) Yaml

type Yaml map[interface{}]interface{}

// merge adds src's Yaml structure over to dst.
func merge(dst, src Yaml) {
	for k, v := range src {
		dst[k] = v
	}
}

// shortImageName returns a best guess at the application's name.
func shortImageName(longname string) string {
	tokens := strings.Split(longname, "/")
	return tokens[len(tokens)-1]
}

// ParseYamlFile parses a file's contents and returns a Yaml structure.
func ParseYamlFile(filename string) (Yaml, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseYaml(data)
}

// ParseYaml parses a byte stream and returns a Yaml structure.
func ParseYaml(data []byte) (Yaml, error) {
	m := make(Yaml)
	err := yaml.Unmarshal(data, &m)
	return m, err
}

// DefaultPrometheusConfiguration returns a sane, default Prometheus
// configuration section that can be inserted into a Yaml structure.
func DefaultPrometheusConfiguration(app string) (Yaml, error) {
	template := `
prometheus:
  enabled: true
  path: /%s/internal/metrics
`
	data := fmt.Sprintf(template, app)
	return ParseYaml([]byte(data))
}

func InsertPrometheus(src Yaml) Yaml {
	name := shortImageName(src["image"].(string))
	result := make(Yaml)
	merge(result, src)
	defaultConfig, _ := DefaultPrometheusConfiguration(name)
	merge(result, defaultConfig)
	return result
}

func AppConfig(yml string) Result {
	return Result{}
}
