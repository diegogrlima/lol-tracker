package redisadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/diegogrlima/lol-tracker/internal/player"
	"github.com/redis/go-redis/v9"
)

type PlayerCache struct {
	client *redis.Client
}

func NewPlayerCache(client *redis.Client) *PlayerCache {
	return &PlayerCache{client: client}
}

func (c *PlayerCache) Get(
	ctx context.Context,
	gameName string,
	tagLine string,
) (*player.Player, error) {
	value, err := c.client.Get(ctx, playerKey(gameName, tagLine)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, player.ErrCacheMiss
		}
		return nil, fmt.Errorf("get cached player: %w", err)
	}

	var cachedPlayer player.Player
	if err := json.Unmarshal(value, &cachedPlayer); err != nil {
		return nil, fmt.Errorf("decode cached player: %w", err)
	}

	return &cachedPlayer, nil
}

func (c *PlayerCache) Set(
	ctx context.Context,
	value *player.Player,
	ttl time.Duration,
) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode player for cache: %w", err)
	}

	if err := c.client.Set(
		ctx,
		playerKey(value.GameName, value.TagLine),
		data,
		ttl,
	).Err(); err != nil {
		return fmt.Errorf("cache player: %w", err)
	}

	return nil
}

func playerKey(gameName, tagLine string) string {
	riotID := strings.ToLower(gameName + ":" + tagLine)
	return "player:account:" + url.QueryEscape(riotID)
}
