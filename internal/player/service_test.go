package player

import (
	"context"
	"errors"
	"testing"
	"time"
)

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
	service := NewService(provider, cache, time.Minute)

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
	service := NewService(provider, cache, time.Minute)

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
	service := NewService(provider, cache, time.Minute)

	_, err := service.GetByRiotID(context.Background(), "Unknown", "BR1")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetByRiotID() error = %v, want ErrNotFound", err)
	}
}
