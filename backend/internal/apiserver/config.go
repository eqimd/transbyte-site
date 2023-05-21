package apiserver

import "github.com/BurntSushi/toml"

type Config struct {
	BindAddr string `toml:"addr"`
}

func NewConfig(path string) (*Config, error) {
	cfg := new(Config)

	_, err := toml.DecodeFile(path, cfg)

	return cfg, err
}
