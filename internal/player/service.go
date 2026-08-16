package player

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

var (
	ErrCacheMiss        = errors.New("player cache miss")
	ErrCacheUnavailable = errors.New("player cache unavailable")
	ErrNotFound         = errors.New("player not found")
	ErrRateLimited      = errors.New("riot API rate limit exceeded")
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
	logger   *slog.Logger
}

func NewService(
	accounts AccountProvider,
	cache Cache,
	cacheTTL time.Duration,
	logger *slog.Logger,
) *Service {
	if logger == nil {
		logger = slog.Default()
	}

	return &Service{
		accounts: accounts,
		cache:    cache,
		cacheTTL: cacheTTL,
		logger:   logger,
	}
}

func (s *Service) GetByRiotID(
	ctx context.Context,
	gameName string,
	tagLine string,
) (*Player, error) {
	cachedPlayer, err := s.cache.Get(ctx, gameName, tagLine)

	switch {
	case err == nil:
		return cachedPlayer, nil

	case errors.Is(err, ErrCacheMiss):
		// Cache miss é uma situação normal.

	default:
		return nil, fmt.Errorf(
			"%w: %v",
			ErrCacheUnavailable,
			err,
		)
	}

	result, err := s.accounts.GetByRiotID(ctx, gameName, tagLine)
	if err != nil {
		return nil, fmt.Errorf("get Riot account: %w", err)
	}

	if err := s.cache.Set(ctx, result, s.cacheTTL); err != nil {
		s.logger.WarnContext(
			ctx,
			"failed to cache player",
			"error", err,
		)
	}

	return result, nil
}
