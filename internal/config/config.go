package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type PlayerConfig struct {
	ServerAddress string
	RedisAddress  string
	RedisPassword string
	RedisDB       int
	RiotAPIKey    string
	RiotRegion    string
	CacheTTL      time.Duration
}

type MatchConfig struct {
	ServerAddress  string
	RedisAddress   string
	RedisPassword  string
	RedisDB        int
	RiotAPIKey     string
	RiotRegion     string
	IDsCacheTTL    time.Duration
	DetailCacheTTL time.Duration
}

type GameDataConfig struct {
	ServerAddress     string
	DataDragonBaseURL string
	DataDragonVersion string
	DataDragonLocale  string
}

func LoadPlayer() (PlayerConfig, error) {
	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil || redisDB < 0 {
		return PlayerConfig{}, errors.New("REDIS_DB must be a non-negative integer")
	}

	cacheTTL, err := time.ParseDuration(getEnv("CACHE_TTL", "5m"))
	if err != nil || cacheTTL <= 0 {
		return PlayerConfig{}, errors.New("CACHE_TTL must be a positive duration")
	}

	serverPort := getEnv("SERVER_PORT", "8080")
	if _, err := strconv.ParseUint(serverPort, 10, 16); err != nil {
		return PlayerConfig{}, fmt.Errorf("SERVER_PORT must be a valid port: %w", err)
	}

	cfg := PlayerConfig{
		ServerAddress: ":" + serverPort,
		RedisAddress:  getEnv("REDIS_ADDRESS", "localhost:6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
		RiotAPIKey:    strings.TrimSpace(os.Getenv("RIOT_API_KEY")),
		RiotRegion:    getEnv("RIOT_REGION", "americas"),
		CacheTTL:      cacheTTL,
	}

	if cfg.RiotAPIKey == "" {
		return PlayerConfig{}, errors.New("RIOT_API_KEY is required")
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

func LoadMatch() (MatchConfig, error) {
	serverAddress, err := LoadMatchServerAddress()
	if err != nil {
		return MatchConfig{}, err
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil || redisDB < 0 {
		return MatchConfig{}, errors.New("REDIS_DB must be a non-negative integer")
	}

	idsCacheTTL, err := time.ParseDuration(getEnv("MATCH_IDS_CACHE_TTL", "5m"))
	if err != nil || idsCacheTTL <= 0 {
		return MatchConfig{}, errors.New("MATCH_IDS_CACHE_TTL must be a positive duration")
	}

	detailCacheTTL, err := time.ParseDuration(getEnv("MATCH_DETAIL_CACHE_TTL", "24h"))
	if err != nil || detailCacheTTL <= 0 {
		return MatchConfig{}, errors.New("MATCH_DETAIL_CACHE_TTL must be a positive duration")
	}

	cfg := MatchConfig{
		ServerAddress:  serverAddress,
		RedisAddress:   getEnv("REDIS_ADDRESS", "localhost:6379"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		RedisDB:        redisDB,
		RiotAPIKey:     strings.TrimSpace(os.Getenv("RIOT_API_KEY")),
		RiotRegion:     getEnv("RIOT_REGION", "americas"),
		IDsCacheTTL:    idsCacheTTL,
		DetailCacheTTL: detailCacheTTL,
	}

	if cfg.RiotAPIKey == "" {
		return MatchConfig{}, errors.New("RIOT_API_KEY is required")
	}

	return cfg, nil
}

func LoadGameData() (GameDataConfig, error) {
	port := getEnv(
		"GAME_DATA_SERVER_PORT",
		getEnv("CHAMPION_SERVER_PORT", getEnv("GAME_SERVER_PORT", "8082")),
	)
	if _, err := strconv.ParseUint(port, 10, 16); err != nil {
		return GameDataConfig{}, fmt.Errorf("GAME_DATA_SERVER_PORT must be a valid port: %w", err)
	}

	cfg := GameDataConfig{
		ServerAddress:     ":" + port,
		DataDragonBaseURL: strings.TrimRight(getEnv("DATA_DRAGON_BASE_URL", "https://ddragon.leagueoflegends.com"), "/"),
		DataDragonVersion: getEnv("DATA_DRAGON_VERSION", "16.1.1"),
		DataDragonLocale:  getEnv("DATA_DRAGON_LOCALE", "pt_BR"),
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
