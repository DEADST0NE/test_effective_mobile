package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	API struct {
		Port int `envconfig:"APP_API_PORT" default:"8080"`
	}
	DB struct {
		Host     string `envconfig:"APP_DB_HOST" default:"pg"`
		Name     string `envconfig:"APP_DB_NAME" default:"postgres"`
		User     string `envconfig:"APP_DB_USER" default:"postgres"`
		Schema   string `envconfig:"APP_DB_SCHEMA" default:"public"`
		Port     string `envconfig:"APP_DB_PORT" default:"5432"`
		Password string `envconfig:"APP_DB_PASSWORD" default:"postgres"`
	}
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}
