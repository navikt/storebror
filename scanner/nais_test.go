package scanner_test

import (
	"io/ioutil"
	"testing"

	"github.com/navikt/storebror/scanner"
	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	t.Run("Parse YAML file correctly", func(t *testing.T) {
		data, err := ioutil.ReadFile("fixtures/app-config/base.yml")
		if err != nil {
			t.Fatalf("fixture is missing")
		}
		yml, err := scanner.ParseYaml(data)
		assert.Equal(t, "foo", yml["image"])
		assert.Equal(t, true, yml["prometheus"].(scanner.Yaml)["enabled"])
	})
}
