package config

import "time"

type Config struct {
	RunAddress      string
	InactiveTimeout time.Duration
}

func GetConfig() *Config {
	return &Config{
		":8080",
		15 * time.Minute,
	}
}
