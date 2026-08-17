package item

import "errors"

var (
	ErrInvalidItemID = errors.New("invalid item ID")
	ErrItemNotFound  = errors.New("item not found")
)
