package config

import "github.com/joho/godotenv"

type Config struct {

}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config {}

	return cfg, nil
}