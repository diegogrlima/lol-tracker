package item

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeItemQuery struct {
	items  []Item
	item   *Item
	err    error
	itemID string
}

func (f *fakeItemQuery) ListStoreItems(context.Context) ([]Item, error) {
	return f.items, f.err
}

func (f *fakeItemQuery) GetByID(_ context.Context, itemID string) (*Item, error) {
	f.itemID = itemID
	return f.item, f.err
}

func TestHandlerRoutesListItems(t *testing.T) {
	query := &fakeItemQuery{items: []Item{{ID: "1001", Name: "Botas"}}}
	handler := NewHandler(query, itemTestLogger())
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	handler.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Body.String(), `"items":[{"id":"1001"`) {
		t.Fatalf("body = %q, want items payload", recorder.Body.String())
	}
}

func TestHandlerRoutesGetItemByID(t *testing.T) {
	query := &fakeItemQuery{item: &Item{ID: "1001", Name: "Botas"}}
	handler := NewHandler(query, itemTestLogger())
	request := httptest.NewRequest(http.MethodGet, "/1001", nil)
	recorder := httptest.NewRecorder()

	handler.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if query.itemID != "1001" {
		t.Fatalf("item ID = %q, want %q", query.itemID, "1001")
	}
	if !strings.Contains(recorder.Body.String(), `"item":{"id":"1001"`) {
		t.Fatalf("body = %q, want item payload", recorder.Body.String())
	}
}

func TestHandlerMapsItemNotFound(t *testing.T) {
	query := &fakeItemQuery{err: errors.Join(ErrItemNotFound, errors.New("1001"))}
	handler := NewHandler(query, itemTestLogger())
	request := httptest.NewRequest(http.MethodGet, "/1001", nil)
	recorder := httptest.NewRecorder()

	handler.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func itemTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
