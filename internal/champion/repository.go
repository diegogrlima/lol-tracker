package champion

import "context"

type Repository interface {
	List(ctx context.Context) ([]Champion, error)
}
