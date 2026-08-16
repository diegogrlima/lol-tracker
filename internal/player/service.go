package player

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrCacheMiss   = errors.New("player cache miss")
	ErrNotFound    = errors.New("player not found")
	ErrRateLimited = errors.New("riot API rate limit exceeded")
)

type Player struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type AccountProvider interface {
	GetByRiotID(
		ctx context.Context,
		gameName string,
		tagLine string,
	) (*Player, error)
}

type Cache interface {
	Get(
		ctx context.Context,
		gameName string,
		tagLine string,
	) (*Player, error)

	Set(
		ctx context.Context,
		player *Player,
		ttl time.Duration,
	) error
}

type Service struct {
	accounts AccountProvider
	cache    Cache
	cacheTTL time.Duration
}

func NewService(
	accounts AccountProvider,
	cache Cache,
	cacheTTL time.Duration,
) *Service {
	return &Service{
		accounts: accounts,
		cache:    cache,
		cacheTTL: cacheTTL,
	}
}

func (s *Service) GetByRiotID(
	ctx context.Context,
	gameName string,
	tagLine string,
) (*Player, error) {
	cachedPlayer, err := s.cache.Get(ctx, gameName, tagLine)
	if err == nil {
		return cachedPlayer, nil
	}

	player, err := s.accounts.GetByRiotID(ctx, gameName, tagLine)
	if err != nil {
		return nil, fmt.Errorf("get Riot account: %w", err)
	}

	// Cache failures must not turn a successful Riot response into an HTTP error.
	_ = s.cache.Set(ctx, player, s.cacheTTL)

	return player, nil
}
