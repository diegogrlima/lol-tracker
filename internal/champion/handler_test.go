package champion

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

type fakeChampionLister struct {
	champions []Champion
	err       error
}

func (f *fakeChampionLister) List(context.Context) ([]Champion, error) {
	return f.champions, f.err
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestHandlerListsChampions(t *testing.T) {
	lister := &fakeChampionLister{champions: []Champion{{ID: "Ahri", Name: "Ahri"}}}
	handler := NewHandler(lister, testLogger())
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/champions", nil)

	handler.List(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Body.String(), `"champions":[{"id":"Ahri"`) {
		t.Fatalf("body = %q, want champions payload", recorder.Body.String())
	}
}

func TestHandlerMapsProviderError(t *testing.T) {
	handler := NewHandler(
		&fakeChampionLister{err: errors.New("unavailable")},
		testLogger(),
	)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/champions", nil)

	handler.List(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadGateway)
	}
}
