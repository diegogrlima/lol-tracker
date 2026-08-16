package player

import (
	"context"
	"time"
)

// AccountProvider retrieves player accounts from an external provider.
type AccountProvider interface {
	GetByRiotID(ctx context.Context, gameName string, tagLine string) (*Player, error)
}

// Cache stores player accounts for a limited period.
type Cache interface {
	Get(ctx context.Context, gameName string, tagLine string) (*Player, error)
	Set(ctx context.Context, player *Player, ttl time.Duration) error
}
