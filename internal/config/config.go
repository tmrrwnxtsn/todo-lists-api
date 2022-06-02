package config

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// Config represents an application configuration.
type Config struct {
	// BindAddr is the binding address. required.
	BindAddr string `yaml:"bind_addr"`
	// DSN is the data source name for connecting to the database. required.
	DSN string `yaml:"dsn"`
}

// Validate validates the application configuration.
func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.DSN, validation.Required),
		validation.Field(&c.BindAddr, validation.Required),
	)
}

// Load returns an application configuration which is populated from the given configuration file and environment variables.
func Load(filePath string) (*Config, error) {
	c := Config{}

	// load from YAML config file
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	c.DSN = getEnv("DSN", c.DSN)
	c.BindAddr = getEnv("BIND_ADDR", c.BindAddr)

	if err = c.Validate(); err != nil {
		return nil, err
	}

	return &c, nil
}

// getEnv is a simple helper function to read an environment or return a default value.
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
