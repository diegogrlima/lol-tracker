package champion

import (
	"context"
	"fmt"
	"sort"
)

type Service struct {
	champions Repository
}

func NewService(champions Repository) *Service {
	return &Service{champions: champions}
}

func (s *Service) List(ctx context.Context) ([]Champion, error) {
	champions, err := s.champions.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list champions: %w", err)
	}

	sort.Slice(champions, func(i, j int) bool {
		return champions[i].Name < champions[j].Name
	})

	return champions, nil
}
