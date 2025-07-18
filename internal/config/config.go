package config

import (
	"fmt"
	"os"

	"github.com/risingwavelabs/eris"
	"gopkg.in/yaml.v3"
)

// The program's configuration with default values. The latter ensures the
// program is working even if no file is provided or if it is incomplete.
var C = Config{
	// Initialise the default config here.
}

type Config struct {
	// Define config structure here.
	// Do not forget the `yaml` or `json` tags.
}

func (c *Config) Load(configPath string) error {
	if len(configPath) == 0 {
		return nil
	}

	raw, err := os.ReadFile(configPath)
	if err != nil {
		return eris.Wrapf(err, "failed to read config file '%s'", configPath)
	}

	err = yaml.Unmarshal(raw, c)
	if err != nil {
		return eris.Wrapf(err, "failed to unmarshal config file '%s'", configPath)
	}

	return nil
}

func (c *Config) Print() {
	yamlConfig, _ := yaml.Marshal(&C)
	fmt.Println()
	fmt.Println(string(yamlConfig))
}
