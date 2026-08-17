package ddragonadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	itemdomain "github.com/diegogrlima/lol-tracker/internal/item"
)

type itemResponse struct {
	Data map[string]itemData `json:"data"`
}

type itemData struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Plaintext   string   `json:"plaintext"`
	Tags        []string `json:"tags"`
	Gold        struct {
		Base        int  `json:"base"`
		Total       int  `json:"total"`
		Sell        int  `json:"sell"`
		Purchasable bool `json:"purchasable"`
	} `json:"gold"`
	Image struct {
		Full string `json:"full"`
	} `json:"image"`
}

func (c *Client) ListItems(ctx context.Context) ([]itemdomain.Item, error) {
	endpoint := fmt.Sprintf(
		"%s/cdn/%s/data/%s/item.json",
		c.baseURL,
		c.version,
		c.locale,
	)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create Data Dragon item request: %w", err)
	}
	request.Header.Set("Accept", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request Data Dragon items: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("Data Dragon items returned status %d", response.StatusCode)
	}

	var payload itemResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode Data Dragon items: %w", err)
	}

	items := make([]itemdomain.Item, 0, len(payload.Data))
	for itemID, data := range payload.Data {
		items = append(items, mapItem(itemID, data, c.baseURL, c.version))
	}

	return items, nil
}

func (c *Client) GetItemByID(
	ctx context.Context,
	itemID string,
) (*itemdomain.Item, error) {
	itemID = strings.TrimSpace(itemID)
	if itemID == "" {
		return nil, itemdomain.ErrInvalidItemID
	}

	items, err := c.ListItems(ctx)
	if err != nil {
		return nil, err
	}

	for i := range items {
		if items[i].ID == itemID {
			return &items[i], nil
		}
	}

	return nil, fmt.Errorf("%w: %s", itemdomain.ErrItemNotFound, itemID)
}

func mapItem(
	itemID string,
	data itemData,
	baseURL string,
	version string,
) itemdomain.Item {
	return itemdomain.Item{
		ID:          itemID,
		Name:        data.Name,
		Description: data.Description,
		Plaintext:   data.Plaintext,
		Tags:        data.Tags,
		BasePrice:   data.Gold.Base,
		TotalPrice:  data.Gold.Total,
		SellPrice:   data.Gold.Sell,
		Purchasable: data.Gold.Purchasable,
		ImageURL: fmt.Sprintf(
			"%s/cdn/%s/img/item/%s",
			baseURL,
			version,
			data.Image.Full,
		),
	}
}
