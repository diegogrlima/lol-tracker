package match

import "errors"

var (
	ErrInvalidPUUID      = errors.New("invalid PUUID")
	ErrInvalidMatchID    = errors.New("invalid match ID")
	ErrMatchNotFound     = errors.New("match not found")
	ErrRateLimited       = errors.New("riot API rate limit exceeded")
	ErrCacheUnavailable  = errors.New("match cache unavailable")
	ErrInvalidPagination = errors.New("invalid pagination")
)
