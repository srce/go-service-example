package config

import (
	"context"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"github.com/dzyanis/go-service-example/internal/api"
	"github.com/dzyanis/go-service-example/pkg/database"
)

type Config struct {
	Postgres database.Config
	API      api.Config
}

func LoadEnv(ctx context.Context, filename string) (*Config, error) {
	cnf := Config{}

	if err := godotenv.Load(filename); err != nil {
		return &cnf, fmt.Errorf("loading env: %w", err)
	}

	if err := env.Parse(&cnf); err != nil {
		return &cnf, fmt.Errorf("parsing env: %w", err)
	}

	return &cnf, nil
}
