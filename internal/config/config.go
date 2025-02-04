package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	// не самый лучший вариант, но ради тестов
	DatabaseURI string `env:"DATABASE_DSN" envDefault:"postgres://postgres:somepostgres@localhost:5432/postgres?sslmode=disable"`
}

func GetConfig() (Config, error) {
	log.Println("reading config")
	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		log.Println("parsing env error", err)
		return Config{}, err
	}

	return cfg, nil
}
