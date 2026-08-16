package match

import "context"

// Repository contracts will be defined from the needs of the first use case.

type ListOptions struct {
	Start int
	Count int
}

type Repository interface {
	ListIDsByPUUID(
		ctx context.Context,
		puuid string,
		options ListOptions,
	) ([]string, error)
}
