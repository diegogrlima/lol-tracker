package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServerAddress string
	RedisAddress  string
	RedisPassword string
	RedisDB       int
	RiotAPIKey    string
	RiotRegion    string
	CacheTTL      time.Duration
}

func Load() (Config, error) {
	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil || redisDB < 0 {
		return Config{}, errors.New("REDIS_DB must be a non-negative integer")
	}

	cacheTTL, err := time.ParseDuration(getEnv("CACHE_TTL", "5m"))
	if err != nil || cacheTTL <= 0 {
		return Config{}, errors.New("CACHE_TTL must be a positive duration")
	}

	serverPort := getEnv("SERVER_PORT", "8080")
	if _, err := strconv.ParseUint(serverPort, 10, 16); err != nil {
		return Config{}, fmt.Errorf("SERVER_PORT must be a valid port: %w", err)
	}

	cfg := Config{
		ServerAddress: ":" + serverPort,
		RedisAddress:  getEnv("REDIS_ADDRESS", "localhost:6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
		RiotAPIKey:    strings.TrimSpace(os.Getenv("RIOT_API_KEY")),
		RiotRegion:    getEnv("RIOT_REGION", "americas"),
		CacheTTL:      cacheTTL,
	}

	if cfg.RiotAPIKey == "" {
		return Config{}, errors.New("RIOT_API_KEY is required")
	}

	return cfg, nil
}

func LoadMatchServerAddress() (string, error) {
	port := getEnv("MATCH_SERVER_PORT", "8081")
	if _, err := strconv.ParseUint(port, 10, 16); err != nil {
		return "", fmt.Errorf("MATCH_SERVER_PORT must be a valid port: %w", err)
	}

	return ":" + port, nil
}

func getEnv(name, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}

	return value
}
