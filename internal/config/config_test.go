package config

import "testing"

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("SERVER_PORT", "")
	t.Setenv("REDIS_ADDRESS", "")
	t.Setenv("RIOT_API_KEY", "test-key")
	t.Setenv("RIOT_REGION", "")

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
	if cfg.RiotRegion != "americas" {
		t.Errorf("RiotRegion = %q, want %q", cfg.RiotRegion, "americas")
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
	t.Setenv("RIOT_API_KEY", "test-key")
	t.Setenv("RIOT_REGION", "europe")

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
	if cfg.RiotAPIKey != "test-key" {
		t.Error("RiotAPIKey does not match the environment value")
	}
	if cfg.RiotRegion != "europe" {
		t.Errorf("RiotRegion = %q, want %q", cfg.RiotRegion, "europe")
	}
}
