package config

import (
	"log"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	RunAddress string `env:"RUN_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func GetConfig() *Config {
	cfg := &Config{
		":8080",
		"127.0.0.1:8080",
	}

	err := env.Parse(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	return cfg

}
