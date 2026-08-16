package riotadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	matchdomain "github.com/diegogrlima/lol-tracker/internal/match"
)

func (c *Client) ListIDsByPUUID(
	ctx context.Context,
	puuid string,
	options matchdomain.ListOptions,
) ([]string, error) {
	endpoint := fmt.Sprintf(
		"%s/lol/match/v5/matches/by-puuid/%s/ids",
		c.baseURL,
		url.PathEscape(puuid),
	)

	query := url.Values{}
	query.Set("start", strconv.Itoa(options.Start))
	query.Set("count", strconv.Itoa(options.Count))
	endpoint += "?" + query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create Riot match request: %w", err)
	}

	request.Header.Set("X-Riot-Token", c.apiKey)
	request.Header.Set("Accept", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request Riot match API: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, decodeError(response)
	}

	var matchIDs []string
	if err := json.NewDecoder(response.Body).Decode(&matchIDs); err != nil {
		return nil, fmt.Errorf("decode Riot match IDs: %w", err)
	}

	return matchIDs, nil
}
