package riotadapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/diegogrlima/lol-tracker/internal/player"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type accountResponse struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type errorResponse struct {
	Status struct {
		Message string `json:"message"`
	} `json:"status"`
}

func NewClient(apiKey, region string) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("RIOT_API_KEY is required")
	}

	switch region {
	case "americas", "europe", "asia", "sea":
	default:
		return nil, fmt.Errorf("unsupported Riot region: %q", region)
	}

	return &Client{
		apiKey:  apiKey,
		baseURL: "https://" + region + ".api.riotgames.com",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *Client) GetByRiotID(
	ctx context.Context,
	gameName string,
	tagLine string,
) (*player.Player, error) {
	endpoint := fmt.Sprintf(
		"%s/riot/account/v1/accounts/by-riot-id/%s/%s",
		c.baseURL,
		url.PathEscape(gameName),
		url.PathEscape(tagLine),
	)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create Riot request: %w", err)
	}

	request.Header.Set("X-Riot-Token", c.apiKey)
	request.Header.Set("Accept", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request Riot API: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusNotFound:
		return nil, player.ErrNotFound
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf(
			"%w: retry after %s seconds",
			player.ErrRateLimited,
			response.Header.Get("Retry-After"),
		)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, decodeError(response)
	}

	var account accountResponse
	if err := json.NewDecoder(response.Body).Decode(&account); err != nil {
		return nil, fmt.Errorf("decode Riot account: %w", err)
	}

	return &player.Player{
		PUUID:    account.PUUID,
		GameName: account.GameName,
		TagLine:  account.TagLine,
	}, nil
}

func decodeError(response *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("Riot API returned status %d", response.StatusCode)
	}

	var apiError errorResponse
	if err := json.Unmarshal(body, &apiError); err == nil && apiError.Status.Message != "" {
		return fmt.Errorf(
			"Riot API returned status %d: %s",
			response.StatusCode,
			apiError.Status.Message,
		)
	}

	return fmt.Errorf("Riot API returned status %d", response.StatusCode)
}
