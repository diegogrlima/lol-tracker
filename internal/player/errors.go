package player

import "errors"

var (
	ErrCacheMiss        = errors.New("player cache miss")
	ErrCacheUnavailable = errors.New("player cache unavailable")
	ErrNotFound         = errors.New("player not found")
	ErrRateLimited      = errors.New("riot API rate limit exceeded")
)
