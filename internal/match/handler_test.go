package match

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeMatchQuery struct {
	matchIDs []string
	match    *Match
	err      error
	puuid    string
	matchID  string
	options  ListOptions
}

func (f *fakeMatchQuery) ListIDsByPUUID(
	_ context.Context,
	puuid string,
	options ListOptions,
) ([]string, error) {
	f.puuid = puuid
	f.options = options
	return f.matchIDs, f.err
}

func (f *fakeMatchQuery) GetByID(
	_ context.Context,
	matchID string,
) (*Match, error) {
	f.matchID = matchID
	return f.match, f.err
}

func TestHandlerRoutesListMatchIDs(t *testing.T) {
	query := &fakeMatchQuery{matchIDs: []string{"BR1_123", "BR1_456"}}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewHandler(query, logger)
	request := httptest.NewRequest(
		http.MethodGet,
		"/by-puuid/player-puuid?start=10&count=5",
		nil,
	)
	recorder := httptest.NewRecorder()

	handler.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if query.puuid != "player-puuid" {
		t.Fatalf("PUUID = %q, want %q", query.puuid, "player-puuid")
	}
	if query.options != (ListOptions{Start: 10, Count: 5}) {
		t.Fatalf("options = %#v, want start 10 and count 5", query.options)
	}
	if !strings.Contains(recorder.Body.String(), `"BR1_123"`) {
		t.Fatalf("body = %q, want match IDs", recorder.Body.String())
	}
}

func TestHandlerRoutesGetMatchByID(t *testing.T) {
	query := &fakeMatchQuery{match: &Match{
		Metadata: Metadata{MatchID: "BR1_123"},
	}}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewHandler(query, logger)
	request := httptest.NewRequest(http.MethodGet, "/BR1_123", nil)
	recorder := httptest.NewRecorder()

	handler.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if query.matchID != "BR1_123" {
		t.Fatalf("match ID = %q, want %q", query.matchID, "BR1_123")
	}
	if !strings.Contains(recorder.Body.String(), `"matchId":"BR1_123"`) {
		t.Fatalf("body = %q, want match details", recorder.Body.String())
	}
}
