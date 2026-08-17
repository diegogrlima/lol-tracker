package item

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

type Service struct {
	items Repository
}

func (s *Service) GetByID(ctx context.Context, itemID string) (*Item, error) {
	itemID = strings.TrimSpace(itemID)
	if itemID == "" {
		return nil, ErrInvalidItemID
	}

	storeItem, err := s.items.GetItemByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("get item by ID: %w", err)
	}

	return storeItem, nil
}

func NewService(items Repository) *Service {
	return &Service{items: items}
}

func (s *Service) ListStoreItems(ctx context.Context) ([]Item, error) {
	items, err := s.items.ListItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("list store items: %w", err)
	}

	storeItems := make([]Item, 0, len(items))
	for _, currentItem := range items {
		if currentItem.Purchasable {
			storeItems = append(storeItems, currentItem)
		}
	}
	sort.Slice(storeItems, func(i, j int) bool {
		return storeItems[i].Name < storeItems[j].Name
	})
	return storeItems, nil
}
