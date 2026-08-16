package player

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type fakeAccountProvider struct {
	player *Player
	err    error
	calls  int
}

func (f *fakeAccountProvider) GetByRiotID(
	context.Context,
	string,
	string,
) (*Player, error) {
	f.calls++
	return f.player, f.err
}

type fakeCache struct {
	player   *Player
	getErr   error
	setErr   error
	setCalls int
}

func (f *fakeCache) Get(context.Context, string, string) (*Player, error) {
	return f.player, f.getErr
}

func (f *fakeCache) Set(context.Context, *Player, time.Duration) error {
	f.setCalls++
	return f.setErr
}

func TestServiceReturnsCachedPlayer(t *testing.T) {
	cachedPlayer := &Player{PUUID: "cached"}
	provider := &fakeAccountProvider{}
	cache := &fakeCache{player: cachedPlayer}
	service := NewService(provider, cache, time.Minute, discardLogger())

	result, err := service.GetByRiotID(context.Background(), "Player", "BR1")
	if err != nil {
		t.Fatalf("GetByRiotID() returned an error: %v", err)
	}
	if result != cachedPlayer {
		t.Fatalf("GetByRiotID() = %#v, want cached player %#v", result, cachedPlayer)
	}
	if provider.calls != 0 {
		t.Fatalf("provider calls = %d, want 0", provider.calls)
	}
}

func TestServiceFetchesAndCachesPlayerOnMiss(t *testing.T) {
	fetchedPlayer := &Player{PUUID: "fetched", GameName: "Player", TagLine: "BR1"}
	provider := &fakeAccountProvider{player: fetchedPlayer}
	cache := &fakeCache{getErr: ErrCacheMiss}
	service := NewService(provider, cache, time.Minute, discardLogger())

	result, err := service.GetByRiotID(context.Background(), "Player", "BR1")
	if err != nil {
		t.Fatalf("GetByRiotID() returned an error: %v", err)
	}
	if result != fetchedPlayer {
		t.Fatalf("GetByRiotID() = %#v, want fetched player %#v", result, fetchedPlayer)
	}
	if provider.calls != 1 {
		t.Fatalf("provider calls = %d, want 1", provider.calls)
	}
	if cache.setCalls != 1 {
		t.Fatalf("cache Set calls = %d, want 1", cache.setCalls)
	}
}

func TestServicePreservesProviderError(t *testing.T) {
	provider := &fakeAccountProvider{err: ErrNotFound}
	cache := &fakeCache{getErr: ErrCacheMiss}
	service := NewService(provider, cache, time.Minute, discardLogger())

	_, err := service.GetByRiotID(context.Background(), "Unknown", "BR1")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetByRiotID() error = %v, want ErrNotFound", err)
	}
}

func TestServiceDoesNotCallProviderWhenCacheFails(t *testing.T) {
	provider := &fakeAccountProvider{player: &Player{PUUID: "unexpected"}}
	cache := &fakeCache{getErr: errors.New("Redis connection refused")}
	service := NewService(provider, cache, time.Minute, discardLogger())

	_, err := service.GetByRiotID(context.Background(), "Player", "BR1")
	if !errors.Is(err, ErrCacheUnavailable) {
		t.Fatalf("GetByRiotID() error = %v, want ErrCacheUnavailable", err)
	}
	if provider.calls != 0 {
		t.Fatalf("provider calls = %d, want 0", provider.calls)
	}
}

func TestServiceReturnsPlayerWhenCacheWriteFails(t *testing.T) {
	fetchedPlayer := &Player{PUUID: "fetched", GameName: "Player", TagLine: "BR1"}
	provider := &fakeAccountProvider{player: fetchedPlayer}
	cache := &fakeCache{
		getErr: ErrCacheMiss,
		setErr: errors.New("Redis connection refused"),
	}
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	service := NewService(provider, cache, time.Minute, logger)

	result, err := service.GetByRiotID(context.Background(), "Player", "BR1")
	if err != nil {
		t.Fatalf("GetByRiotID() returned an error: %v", err)
	}
	if result != fetchedPlayer {
		t.Fatalf("GetByRiotID() = %#v, want fetched player %#v", result, fetchedPlayer)
	}
	if !strings.Contains(logs.String(), "failed to cache player") {
		t.Fatalf("logs = %q, want cache write warning", logs.String())
	}
}
