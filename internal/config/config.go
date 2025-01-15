package config

import (
	"fmt"
	"os"
)

type Config struct {
	DB string
}

func New() (*Config, error) {
	db, exists := os.LookupEnv("DBSTRING")
	if !exists || db == "" {
		return nil, fmt.Errorf("No database string")
	}

	return &Config{
		DB: db,
	}, nil
}
