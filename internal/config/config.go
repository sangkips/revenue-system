package config

import (
	"os"

	"github.com/rs/zerolog/log"
)

type Config struct {
	DBURL     string
	Port      string
	JWTSecret string
}

func Load() *Config {
	cfg := &Config{
		DBURL:     os.Getenv("DB_URL"),
		Port:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	if cfg.DBURL == "" {
		log.Fatal().Msg("DB_URL not set")
	}

	if cfg.JWTSecret == "" {
		log.Fatal().Msg("JWT_SECRET not set")
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	return cfg
}
