package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName     string
	AppEnv      string
	HTTPAddress string
	DatabaseURL string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		AppName:     os.Getenv("APP_NAME"),
		AppEnv:      os.Getenv("APP_ENV"),
		HTTPAddress: os.Getenv("HTTP_ADDRESS"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	if cfg.AppName == "" {
		cfg.AppName = "samatalk-api"
	}

	if cfg.AppEnv == "" {
		cfg.AppEnv = "local"
	}

	if cfg.HTTPAddress == "" {
		cfg.HTTPAddress = ":8080"
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) Validate() error {
	if cfg.DatabaseURL == "" {
		return fmt.Errorf(
			"DATABASE_URL is required",
		)
	}

	return nil
}
