package redisadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"

	matchdomain "github.com/diegogrlima/lol-tracker/internal/match"
	"github.com/redis/go-redis/v9"
)

type cacheClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
}

type CachedMatchRepository struct {
	cache     cacheClient
	source    matchdomain.Repository
	idsTTL    time.Duration
	detailTTL time.Duration
	logger    *slog.Logger
}

func NewCachedMatchRepository(
	cache cacheClient,
	source matchdomain.Repository,
	idsTTL time.Duration,
	detailTTL time.Duration,
	logger *slog.Logger,
) *CachedMatchRepository {
	if logger == nil {
		logger = slog.Default()
	}

	return &CachedMatchRepository{
		cache:     cache,
		source:    source,
		idsTTL:    idsTTL,
		detailTTL: detailTTL,
		logger:    logger,
	}
}

func (r *CachedMatchRepository) ListIDsByPUUID(
	ctx context.Context,
	puuid string,
	options matchdomain.ListOptions,
) ([]string, error) {
	key := matchIDsKey(puuid, options)

	var matchIDs []string
	hit, err := r.readCachedJSON(ctx, key, &matchIDs)
	if err != nil {
		return nil, err
	}
	if hit {
		return matchIDs, nil
	}

	matchIDs, err = r.source.ListIDsByPUUID(ctx, puuid, options)
	if err != nil {
		return nil, err
	}

	r.writeCachedJSON(ctx, key, matchIDs, r.idsTTL)
	return matchIDs, nil
}

func (r *CachedMatchRepository) GetByID(
	ctx context.Context,
	matchID string,
) (*matchdomain.Match, error) {
	key := matchDetailKey(matchID)

	var cachedMatch matchdomain.Match
	hit, err := r.readCachedJSON(ctx, key, &cachedMatch)
	if err != nil {
		return nil, err
	}
	if hit {
		return &cachedMatch, nil
	}

	matchResult, err := r.source.GetByID(ctx, matchID)
	if err != nil {
		return nil, err
	}

	r.writeCachedJSON(ctx, key, matchResult, r.detailTTL)
	return matchResult, nil
}

func (r *CachedMatchRepository) readCachedJSON(
	ctx context.Context,
	key string,
	target any,
) (bool, error) {
	data, err := r.cache.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, fmt.Errorf("%w: read Redis key: %v", matchdomain.ErrCacheUnavailable, err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return false, fmt.Errorf("%w: decode Redis value: %v", matchdomain.ErrCacheUnavailable, err)
	}

	return true, nil
}

func (r *CachedMatchRepository) writeCachedJSON(
	ctx context.Context,
	key string,
	value any,
	ttl time.Duration,
) {
	data, err := json.Marshal(value)
	if err != nil {
		r.logger.WarnContext(ctx, "failed to encode match cache value", "error", err)
		return
	}

	if err := r.cache.Set(ctx, key, data, ttl).Err(); err != nil {
		r.logger.WarnContext(ctx, "failed to write match cache", "error", err)
	}
}

func matchIDsKey(puuid string, options matchdomain.ListOptions) string {
	return "match:ids:" + url.QueryEscape(strings.ToLower(puuid)) +
		":" + strconv.Itoa(options.Start) +
		":" + strconv.Itoa(options.Count)
}

func matchDetailKey(matchID string) string {
	return "match:detail:" + url.QueryEscape(strings.ToLower(matchID))
}
