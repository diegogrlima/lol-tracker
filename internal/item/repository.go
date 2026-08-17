package item

import "context"

type Repository interface {
	ListItems(ctx context.Context) ([]Item, error)
	GetItemByID(ctx context.Context, itemID string) (*Item, error)
}
