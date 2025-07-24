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

	MetricsPort: 9090,

	Database: dbConfig{
		Host: "localhost",
		Port: 5432,
		// User: "",
		// Password: "",
		// Name: "",
	},
}

type dbConfig struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Config struct {
	// Define config structure here.
	// Do not forget the `yaml` or `json` tags.

	// Port used to expose Prometheus metrics.
	MetricsPort uint16 `yaml:"metricsPort"`

	Database dbConfig `yaml:"database"`
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
