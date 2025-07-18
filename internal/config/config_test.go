package config

import (
	"math/rand/v2"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Ensures that marshalling an incomplete config file does not affect values
// of missing attributes.
func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	type TestConfig struct {
		AAA int `yaml:"aaa"`
		BBB int `yaml:"bbb"`
	}

	defaultConfig := TestConfig{
		AAA: rand.Int(),
		BBB: rand.Int(),
	}

	expectedConfig := defaultConfig
	expectedConfig.BBB = rand.Int()
	yamlConfig := "bbb: " + strconv.Itoa(expectedConfig.BBB) + "\n"

	actualConfig := defaultConfig
	err := yaml.Unmarshal([]byte(yamlConfig), &actualConfig)
	require.NoError(t, err)

	require.EqualValues(t, expectedConfig.AAA, actualConfig.AAA)
	require.EqualValues(t, expectedConfig.BBB, actualConfig.BBB)
}
