package match

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	matchIDs []string
	err      error
	calls    int

	receivedPUUID   string
	receivedOptions ListOptions
}

func (f *fakeRepository) ListIDsByPUUID(
	_ context.Context,
	puuid string,
	options ListOptions,
) ([]string, error) {
	f.calls++
	f.receivedPUUID = puuid
	f.receivedOptions = options

	return f.matchIDs, f.err
}

func TestServiceListsMatchIDsByPUUID(t *testing.T) {
	repository := &fakeRepository{
		matchIDs: []string{
			"BR1_123456789",
			"BR1_987654321",
		},
	}

	service := NewService(repository)

	options := ListOptions{
		Start: 0,
		Count: 20,
	}

	result, err := service.ListIDsByPUUID(
		context.Background(),
		"player-puuid",
		options,
	)
	if err != nil {
		t.Fatalf("ListIDsByPUUID() returned an error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("result length = %d, want 2", len(result))
	}

	if result[0] != "BR1_123456789" {
		t.Errorf(
			"result[0] = %q, want %q",
			result[0],
			"BR1_123456789",
		)
	}

	if repository.calls != 1 {
		t.Errorf(
			"repository calls = %d, want 1",
			repository.calls,
		)
	}

	if repository.receivedPUUID != "player-puuid" {
		t.Errorf(
			"received PUUID = %q, want %q",
			repository.receivedPUUID,
			"player-puuid",
		)
	}

	if repository.receivedOptions != options {
		t.Errorf(
			"received options = %#v, want %#v",
			repository.receivedOptions,
			options,
		)
	}
}

func TestServiceRejectsInvalidPUUID(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository)

	_, err := service.ListIDsByPUUID(
		context.Background(),
		"   ",
		ListOptions{
			Start: 0,
			Count: DefaultCount,
		},
	)

	if !errors.Is(err, ErrInvalidPUUID) {
		t.Fatalf(
			"error = %v, want ErrInvalidPUUID",
			err,
		)
	}

	if repository.calls != 0 {
		t.Fatalf(
			"repository calls = %d, want 0",
			repository.calls,
		)
	}
}

func TestServiceRejectsInvalidPagination(t *testing.T) {
	tests := []struct {
		name    string
		options ListOptions
	}{
		{
			name: "negative start",
			options: ListOptions{
				Start: -1,
				Count: DefaultCount,
			},
		},
		{
			name: "zero count",
			options: ListOptions{
				Start: 0,
				Count: 0,
			},
		},
		{
			name: "count above maximum",
			options: ListOptions{
				Start: 0,
				Count: MaxCount + 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakeRepository{}
			service := NewService(repository)

			_, err := service.ListIDsByPUUID(
				context.Background(),
				"player-puuid",
				tt.options,
			)

			if !errors.Is(err, ErrInvalidPagination) {
				t.Fatalf(
					"error = %v, want ErrInvalidPagination",
					err,
				)
			}

			if repository.calls != 0 {
				t.Fatalf(
					"repository calls = %d, want 0",
					repository.calls,
				)
			}
		})
	}
}

func TestServicePreservesRepositoryError(t *testing.T) {
	repositoryError := errors.New("Riot unavailable")

	repository := &fakeRepository{
		err: repositoryError,
	}

	service := NewService(repository)

	_, err := service.ListIDsByPUUID(
		context.Background(),
		"player-puuid",
		ListOptions{
			Start: 0,
			Count: DefaultCount,
		},
	)

	if !errors.Is(err, repositoryError) {
		t.Fatalf(
			"error = %v, want repository error",
			err,
		)
	}
}
