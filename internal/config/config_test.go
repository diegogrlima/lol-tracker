package config

import (
	"testing"
	"time"
)

func setRequiredEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("RIOT_API_KEY", "test-key")
}

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("SERVER_PORT", "")
	t.Setenv("REDIS_ADDRESS", "")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "")
	t.Setenv("RIOT_REGION", "")
	t.Setenv("CACHE_TTL", "")
	setRequiredEnvironment(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}

	if cfg.ServerAddress != ":8080" {
		t.Errorf("ServerAddress = %q, want %q", cfg.ServerAddress, ":8080")
	}
	if cfg.RedisAddress != "localhost:6379" {
		t.Errorf("RedisAddress = %q, want %q", cfg.RedisAddress, "localhost:6379")
	}
	if cfg.RedisDB != 0 {
		t.Errorf("RedisDB = %d, want %d", cfg.RedisDB, 0)
	}
	if cfg.RiotRegion != "americas" {
		t.Errorf("RiotRegion = %q, want %q", cfg.RiotRegion, "americas")
	}
	if cfg.CacheTTL != 5*time.Minute {
		t.Errorf("CacheTTL = %v, want %v", cfg.CacheTTL, 5*time.Minute)
	}
}

func TestLoadRequiresRiotAPIKey(t *testing.T) {
	t.Setenv("RIOT_API_KEY", "")

	if _, err := Load(); err == nil {
		t.Fatal("Load() returned nil error without RIOT_API_KEY")
	}
}

func TestLoadUsesEnvironmentValues(t *testing.T) {
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("REDIS_ADDRESS", "redis:6379")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("REDIS_DB", "2")
	t.Setenv("RIOT_API_KEY", "test-key")
	t.Setenv("RIOT_REGION", "europe")
	t.Setenv("CACHE_TTL", "10m")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}

	if cfg.ServerAddress != ":9090" {
		t.Errorf("ServerAddress = %q, want %q", cfg.ServerAddress, ":9090")
	}
	if cfg.RedisAddress != "redis:6379" {
		t.Errorf("RedisAddress = %q, want %q", cfg.RedisAddress, "redis:6379")
	}
	if cfg.RedisPassword != "secret" {
		t.Error("RedisPassword does not match the environment value")
	}
	if cfg.RedisDB != 2 {
		t.Errorf("RedisDB = %d, want %d", cfg.RedisDB, 2)
	}
	if cfg.RiotAPIKey != "test-key" {
		t.Error("RiotAPIKey does not match the environment value")
	}
	if cfg.RiotRegion != "europe" {
		t.Errorf("RiotRegion = %q, want %q", cfg.RiotRegion, "europe")
	}
	if cfg.CacheTTL != 10*time.Minute {
		t.Errorf("CacheTTL = %v, want %v", cfg.CacheTTL, 10*time.Minute)
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{name: "server port", key: "SERVER_PORT", value: "invalid"},
		{name: "Redis database", key: "REDIS_DB", value: "-1"},
		{name: "cache TTL", key: "CACHE_TTL", value: "0s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setRequiredEnvironment(t)
			t.Setenv(tt.key, tt.value)

			if _, err := Load(); err == nil {
				t.Fatalf("Load() returned nil error for %s=%q", tt.key, tt.value)
			}
		})
	}
}
