package redisadapter

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	matchdomain "github.com/diegogrlima/lol-tracker/internal/match"
	"github.com/redis/go-redis/v9"
)

type fakeCacheClient struct {
	values   map[string][]byte
	getErr   error
	setErr   error
	setCalls int
	lastTTL  time.Duration
}

func (f *fakeCacheClient) Get(_ context.Context, key string) *redis.StringCmd {
	if f.getErr != nil {
		return redis.NewStringResult("", f.getErr)
	}

	value, ok := f.values[key]
	if !ok {
		return redis.NewStringResult("", redis.Nil)
	}

	return redis.NewStringResult(string(value), nil)
}

func (f *fakeCacheClient) Set(
	_ context.Context,
	key string,
	value any,
	expiration time.Duration,
) *redis.StatusCmd {
	f.setCalls++
	f.lastTTL = expiration
	if f.values == nil {
		f.values = make(map[string][]byte)
	}
	if data, ok := value.([]byte); ok {
		f.values[key] = data
	}
	return redis.NewStatusResult("OK", f.setErr)
}

type fakeMatchSource struct {
	matchIDs    []string
	match       *matchdomain.Match
	err         error
	listCalls   int
	detailCalls int
}

func (f *fakeMatchSource) ListIDsByPUUID(
	context.Context,
	string,
	matchdomain.ListOptions,
) ([]string, error) {
	f.listCalls++
	return f.matchIDs, f.err
}

func (f *fakeMatchSource) GetByID(context.Context, string) (*matchdomain.Match, error) {
	f.detailCalls++
	return f.match, f.err
}

func TestCachedMatchRepositoryCachesMatchIDs(t *testing.T) {
	cache := &fakeCacheClient{values: make(map[string][]byte)}
	source := &fakeMatchSource{matchIDs: []string{"BR1_123"}}
	repository := newTestCachedMatchRepository(cache, source)
	options := matchdomain.ListOptions{Start: 0, Count: 20}

	for range 2 {
		result, err := repository.ListIDsByPUUID(context.Background(), "player-puuid", options)
		if err != nil {
			t.Fatalf("ListIDsByPUUID() returned an error: %v", err)
		}
		if len(result) != 1 || result[0] != "BR1_123" {
			t.Fatalf("match IDs = %#v, want cached IDs", result)
		}
	}

	if source.listCalls != 1 {
		t.Fatalf("source list calls = %d, want 1", source.listCalls)
	}
	if cache.setCalls != 1 || cache.lastTTL != 5*time.Minute {
		t.Fatalf("cache writes = %d with TTL %v, want one write with 5m", cache.setCalls, cache.lastTTL)
	}
}

func TestCachedMatchRepositoryCachesMatchDetails(t *testing.T) {
	cache := &fakeCacheClient{values: make(map[string][]byte)}
	source := &fakeMatchSource{match: &matchdomain.Match{
		Metadata: matchdomain.Metadata{MatchID: "BR1_123"},
	}}
	repository := newTestCachedMatchRepository(cache, source)

	for range 2 {
		result, err := repository.GetByID(context.Background(), "BR1_123")
		if err != nil {
			t.Fatalf("GetByID() returned an error: %v", err)
		}
		if result.Metadata.MatchID != "BR1_123" {
			t.Fatalf("match ID = %q, want cached match", result.Metadata.MatchID)
		}
	}

	if source.detailCalls != 1 {
		t.Fatalf("source detail calls = %d, want 1", source.detailCalls)
	}
	if cache.lastTTL != 24*time.Hour {
		t.Fatalf("cache TTL = %v, want 24h", cache.lastTTL)
	}
}

func TestCachedMatchRepositoryFailsClosedOnRedisReadError(t *testing.T) {
	cache := &fakeCacheClient{getErr: errors.New("connection refused")}
	source := &fakeMatchSource{matchIDs: []string{"unexpected"}}
	repository := newTestCachedMatchRepository(cache, source)

	_, err := repository.ListIDsByPUUID(
		context.Background(),
		"player-puuid",
		matchdomain.ListOptions{Start: 0, Count: 20},
	)
	if !errors.Is(err, matchdomain.ErrCacheUnavailable) {
		t.Fatalf("error = %v, want ErrCacheUnavailable", err)
	}
	if source.listCalls != 0 {
		t.Fatalf("source calls = %d, want 0", source.listCalls)
	}
}

func TestCachedMatchRepositoryReturnsSourceResultOnRedisWriteError(t *testing.T) {
	cache := &fakeCacheClient{
		values: make(map[string][]byte),
		setErr: errors.New("connection refused"),
	}
	source := &fakeMatchSource{matchIDs: []string{"BR1_123"}}
	repository := newTestCachedMatchRepository(cache, source)

	result, err := repository.ListIDsByPUUID(
		context.Background(),
		"player-puuid",
		matchdomain.ListOptions{Start: 0, Count: 20},
	)
	if err != nil {
		t.Fatalf("ListIDsByPUUID() returned an error: %v", err)
	}
	if len(result) != 1 || result[0] != "BR1_123" {
		t.Fatalf("match IDs = %#v, want source result", result)
	}
}

func newTestCachedMatchRepository(
	cache cacheClient,
	source matchdomain.Repository,
) *CachedMatchRepository {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewCachedMatchRepository(cache, source, 5*time.Minute, 24*time.Hour, logger)
}
