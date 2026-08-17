package champion

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	champions []Champion
	err       error
}

func (f *fakeRepository) List(context.Context) ([]Champion, error) {
	return f.champions, f.err
}

func TestServiceListsChampionsOrderedByName(t *testing.T) {
	repository := &fakeRepository{champions: []Champion{
		{ID: "Zed", Name: "Zed"},
		{ID: "Ahri", Name: "Ahri"},
	}}
	service := NewService(repository)

	champions, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List() returned an error: %v", err)
	}
	if len(champions) != 2 {
		t.Fatalf("len(champions) = %d, want 2", len(champions))
	}
	if champions[0].Name != "Ahri" || champions[1].Name != "Zed" {
		t.Fatalf("champions = %#v, want champions ordered by name", champions)
	}
}

func TestServicePreservesRepositoryError(t *testing.T) {
	wantErr := errors.New("Data Dragon unavailable")
	service := NewService(&fakeRepository{err: wantErr})

	_, err := service.List(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("List() error = %v, want %v", err, wantErr)
	}
}
