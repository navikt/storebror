package scanner

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Result struct {
	Description string
	Err         error
}

type ResultSet []Result

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

// WriteYamlFile writes a Yaml structure to a file.
func WriteYamlFile(yml Yaml, filename string) error {
	data, err := yaml.Marshal(yml)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0)
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

// InsertPrometheusConfig inserts a valid Prometheus configuration section into
// an existing Yaml structure.
func InsertPrometheusConfig(src Yaml) Yaml {
	name := shortImageName(src["image"].(string))
	result := make(Yaml)
	merge(result, src)
	defaultConfig, _ := DefaultPrometheusConfiguration(name)
	merge(result, defaultConfig)
	return result
}

// HasValidPrometheusConfig returns true if the Yaml structure has
// .prometheus.enabled: true.
func HasValidPrometheusConfig(src Yaml) bool {
	if _, ok := src["prometheus"]; !ok {
		return false
	}
	metricsEnabled := src["prometheus"].(Yaml)["enabled"]
	switch x := metricsEnabled.(type) {
	case bool:
		return x
	default:
		return false
	}
}

// ProcessFile checks if a file has the required Prometheus configuration, and
// writes additional config back to the file if needed.
func ProcessFile(path string) *Result {
	src, err := ParseYamlFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return &Result{Err: err}
	}
	if HasValidPrometheusConfig(src) {
		return nil
	}
	dest := InsertPrometheusConfig(src)
	err = WriteYamlFile(dest, path)
	return &Result{"add Prometheus configuration section", err}
}

// Process runs ProcessFile on all the eligible files in the repository.
func Process(repository string) ResultSet {
	results := make([]Result, 0)
	for _, file := range []string{"app-config.yaml", "app-config-sbs.yaml", "app-config-fss.yaml"} {
		path := path.Join(repository, file)
		result := ProcessFile(path)
		if result == nil {
			continue
		}
		result.Description = fmt.Sprintf("%s: %s", file, result.Description)
		results = append(results, *result)
	}
	return results
}

// Description returns a newline-delimited collection of descriptions in the
// result set.
func (rs ResultSet) Description() string {
	ds := make([]string, len(rs))
	for i, d := range rs {
		ds[i] = d.Description
	}
	return strings.TrimSpace(strings.Join(ds, "\n"))
}
