package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type EnvConfig struct {
	DB DBConfig `envPrefix:"DB_"`
}

type DBConfig struct {
	Name     string `env:"NAME" validate:"required"`
	User     string `env:"USER" validate:"required"`
	Password string `env:"PASSWORD" validate:"omitempty"`
	Host     string `env:"HOST" default:"localhost"`
	Port     int    `env:"PORT" default:"5432"`
}

func LoadConfig() (*EnvConfig, error) {
	cfg := new(EnvConfig)

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("env parse failed: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("env validate failed: %w", err)
	}

	return cfg, nil
}
