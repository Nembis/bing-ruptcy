package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURI string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DatabaseURI: os.Getenv("DATABASE_URI"),
	}

	err := cfg.validateEnv()

	return cfg, err
}

func (cfg *Config) validateEnv() error {
	errs := make([]error, 0)
	if cfg.DatabaseURI == "" {
		errs = append(errs, errors.New("missing DATABASE_URI"))
	}

	if len(errs) == 0 {
		return nil
	}

	joinedErrors := errors.Join(errs...)

	return joinedErrors
}
