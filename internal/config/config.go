package config

import "github.com/BurntSushi/toml"

type Config struct {
	BindAddr       string `toml:"bind_addr"`
	MaxHeaderBytes int    `toml:"max_header_bytes"`
	ReadTimeout    int    `toml:"read_timeout"`
	WriteTimeout   int    `toml:"write_timeout"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr:       ":8080",
		MaxHeaderBytes: 1,
		ReadTimeout:    10,
		WriteTimeout:   10,
	}
}

func (c *Config) Load(configPath string) error {
	_, err := toml.DecodeFile(configPath, c)
	return err
}
