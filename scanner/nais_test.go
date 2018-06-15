package scanner_test

import (
	"io/ioutil"
	"testing"

	"github.com/navikt/storebror/scanner"
	"github.com/stretchr/testify/assert"
)

func TestParseYaml(t *testing.T) {
	t.Run("Parse YAML file correctly", func(t *testing.T) {
		data, err := ioutil.ReadFile("fixtures/app-config/base.yml")
		assert.Nil(t, err)
		yml, err := scanner.ParseYaml(data)
		assert.Equal(t, "foo", yml["image"])
		assert.Equal(t, true, yml["prometheus"].(scanner.Yaml)["enabled"])
	})
}

func TestPrometheus(t *testing.T) {
	t.Run("Default Prometheus configuration is sane", func(t *testing.T) {
		yml, err := scanner.DefaultPrometheusConfiguration("foo-app")
		assert.Nil(t, err)
		assert.Equal(t, true, yml["prometheus"].(scanner.Yaml)["enabled"])
		assert.Equal(t, "/foo-app/internal/metrics", yml["prometheus"].(scanner.Yaml)["path"])
	})

	t.Run("Insert Prometheus configuration if absent", func(t *testing.T) {
		src, err := scanner.ParseYamlFile("fixtures/app-config/prometheus.input.yml")
		assert.Nil(t, err)
		expected, err := scanner.ParseYamlFile("fixtures/app-config/prometheus.output.yml")
		assert.Nil(t, err)
		dst := scanner.InsertPrometheusConfig(src)
		assert.Equal(t, expected, dst)
	})
}
