package scanner_test

import (
	"io/ioutil"
	"testing"

	"github.com/navikt/storebror/scanner"
	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	t.Run("Parse nais manifest from a YAML file", func(t *testing.T) {
		data, err := ioutil.ReadFile("fixtures/app-config/base.yml")
		if err != nil {
			t.Fatalf("fixture is missing")
		}
		manifest, err := scanner.ParseYaml(data)
		assert.Equal(t, "foo", manifest.Image)
		assert.Equal(t, true, manifest.Prometheus.Enabled)
	})
}
