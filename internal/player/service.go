package player

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

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

	playerAccount, err := s.accounts.GetByRiotID(ctx, gameName, tagLine)
	if err != nil {
		return nil, fmt.Errorf("get Riot account: %w", err)
	}

	if err := s.cache.Set(ctx, playerAccount, s.cacheTTL); err != nil {
		s.logger.WarnContext(
			ctx,
			"failed to cache player",
			"error", err,
		)
	}

	return playerAccount, nil
}
