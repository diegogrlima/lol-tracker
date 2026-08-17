package player

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakePlayerFinder struct {
	player *Player
	err    error
	calls  int
}

func (f *fakePlayerFinder) GetByRiotID(
	context.Context,
	string,
	string,
) (*Player, error) {
	f.calls++
	return f.player, f.err
}

func TestHandlerReturnsPlayer(t *testing.T) {
	finder := &fakePlayerFinder{player: &Player{
		PUUID:    "puuid",
		GameName: "Player",
		TagLine:  "BR1",
	}}
	handler := NewHandler(finder, discardLogger())
	recorder := httptest.NewRecorder()

	handler.GetByRiotID(recorder, requestWithRiotID("Player", "BR1"))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Body.String(), `"puuid":"puuid"`) {
		t.Fatalf("body = %q, want player payload", recorder.Body.String())
	}
}

func TestHandlerRejectsEmptyRiotID(t *testing.T) {
	finder := &fakePlayerFinder{}
	handler := NewHandler(finder, discardLogger())
	recorder := httptest.NewRecorder()

	handler.GetByRiotID(recorder, requestWithRiotID(" ", "BR1"))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	if finder.calls != 0 {
		t.Fatalf("finder calls = %d, want 0", finder.calls)
	}
}

func TestHandlerMapsDomainErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "not found", err: ErrNotFound, wantStatus: http.StatusNotFound},
		{name: "rate limit", err: ErrRateLimited, wantStatus: http.StatusServiceUnavailable},
		{name: "cache unavailable", err: ErrCacheUnavailable, wantStatus: http.StatusServiceUnavailable},
		{name: "upstream failure", err: errors.New("unavailable"), wantStatus: http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(&fakePlayerFinder{err: tt.err}, discardLogger())
			recorder := httptest.NewRecorder()

			handler.GetByRiotID(recorder, requestWithRiotID("Player", "BR1"))

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
		})
	}
}

func requestWithRiotID(gameName, tagLine string) *http.Request {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("gameName", gameName)
	routeContext.URLParams.Add("tagLine", tagLine)

	return request.WithContext(context.WithValue(
		request.Context(),
		chi.RouteCtxKey,
		routeContext,
	))
}
