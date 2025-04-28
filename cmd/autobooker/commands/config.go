package commands

import "github.com/caarlos0/env/v11"

type Config struct {
	UserID       string `env:"USER_ID"`
	UserPassword string `env:"USER_PW"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
