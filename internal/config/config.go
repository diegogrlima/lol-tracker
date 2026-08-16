package config

import (
	"errors"
	"os"
)

type Config struct {
	ServerAddress string
	RedisAddress  string
	RiotAPIKey    string
	RiotRegion    string
}

func Load() (Config, error) {
	cfg := Config{
		ServerAddress: ":" + getEnv("SERVER_PORT", "8080"),
		RedisAddress:  getEnv("REDIS_ADDRESS", "localhost:6379"),
		RiotAPIKey:    os.Getenv("RIOT_API_KEY"),
		RiotRegion:    getEnv("RIOT_REGION", "americas"),
	}

	if cfg.RiotAPIKey == "" {
		return Config{}, errors.New("RIOT_API_KEY is required")
	}

	return cfg, nil
}

func getEnv(name, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}

	return value
}
