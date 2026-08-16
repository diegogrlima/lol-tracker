package match

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const (
	DefaultCount = 20
	MaxCount     = 100
)

var (
	ErrInvalidPUUID      = errors.New("invalid PUUID")
	ErrInvalidMatchID    = errors.New("invalid match ID")
	ErrMatchNotFound     = errors.New("match not found")
	ErrRateLimited       = errors.New("riot API rate limit exceeded")
	ErrCacheUnavailable  = errors.New("match cache unavailable")
	ErrInvalidPagination = errors.New("invalid pagination")
)

type Service struct {
	matches Repository
}

func NewService(matches Repository) *Service {
	return &Service{
		matches: matches,
	}
}

func (s *Service) ListIDsByPUUID(
	ctx context.Context,
	puuid string,
	options ListOptions,
) ([]string, error) {
	puuid = strings.TrimSpace(puuid)

	if puuid == "" {
		return nil, ErrInvalidPUUID
	}

	if options.Start < 0 {
		return nil, fmt.Errorf(
			"%w: start must not be negative",
			ErrInvalidPagination,
		)
	}

	if options.Count < 1 || options.Count > MaxCount {
		return nil, fmt.Errorf(
			"%w: count must be between 1 and %d",
			ErrInvalidPagination,
			MaxCount,
		)
	}

	matchIDs, err := s.matches.ListIDsByPUUID(ctx, puuid, options)
	if err != nil {
		return nil, fmt.Errorf("list match IDs: %w", err)
	}

	return matchIDs, nil
}

func (s *Service) GetByID(ctx context.Context, matchID string) (*Match, error) {
	matchID = strings.TrimSpace(matchID)
	if matchID == "" {
		return nil, ErrInvalidMatchID
	}

	result, err := s.matches.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("get match by ID: %w", err)
	}

	return result, nil
}
