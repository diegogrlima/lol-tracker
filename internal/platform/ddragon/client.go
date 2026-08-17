package ddragonadapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/diegogrlima/lol-tracker/internal/champion"
)

type Client struct {
	baseURL    string
	version    string
	locale     string
	httpClient *http.Client
}

type championResponse struct {
	Data map[string]championData `json:"data"`
}

type championData struct {
	ID    string   `json:"id"`
	Key   string   `json:"key"`
	Name  string   `json:"name"`
	Title string   `json:"title"`
	Blurb string   `json:"blurb"`
	Tags  []string `json:"tags"`
	Image struct {
		Full string `json:"full"`
	} `json:"image"`
}

func NewClient(baseURL, version, locale string) (*Client, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	version = strings.TrimSpace(version)
	locale = strings.TrimSpace(locale)

	if baseURL == "" {
		return nil, errors.New("Data Dragon base URL is required")
	}
	if version == "" {
		return nil, errors.New("Data Dragon version is required")
	}
	if locale == "" {
		return nil, errors.New("Data Dragon locale is required")
	}

	return &Client{
		baseURL: baseURL,
		version: version,
		locale:  locale,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *Client) List(ctx context.Context) ([]champion.Champion, error) {
	endpoint := fmt.Sprintf(
		"%s/cdn/%s/data/%s/champion.json",
		c.baseURL,
		c.version,
		c.locale,
	)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create Data Dragon request: %w", err)
	}
	request.Header.Set("Accept", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request Data Dragon: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("Data Dragon returned status %d", response.StatusCode)
	}

	var payload championResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode Data Dragon champions: %w", err)
	}

	champions := make([]champion.Champion, 0, len(payload.Data))
	for _, item := range payload.Data {
		champions = append(champions, champion.Champion{
			ID:       item.ID,
			Key:      item.Key,
			Name:     item.Name,
			Title:    item.Title,
			Blurb:    item.Blurb,
			Tags:     item.Tags,
			ImageURL: fmt.Sprintf("%s/cdn/%s/img/champion/%s", c.baseURL, c.version, item.Image.Full),
		})
	}

	return champions, nil
}
