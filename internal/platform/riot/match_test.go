package riotadapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	matchdomain "github.com/diegogrlima/lol-tracker/internal/match"
)

func TestClientListsMatchIDsByPUUID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/lol/match/v5/matches/by-puuid/player-puuid/ids" {
			t.Errorf("path = %q, want Match-V5 IDs path", r.URL.Path)
		}
		if r.URL.Query().Get("start") != "10" {
			t.Errorf("start = %q, want %q", r.URL.Query().Get("start"), "10")
		}
		if r.URL.Query().Get("count") != "5" {
			t.Errorf("count = %q, want %q", r.URL.Query().Get("count"), "5")
		}
		if r.Header.Get("X-Riot-Token") != "test-key" {
			t.Errorf("X-Riot-Token header was not set")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`["BR1_123","BR1_456"]`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", "americas")
	if err != nil {
		t.Fatalf("NewClient() returned an error: %v", err)
	}
	client.baseURL = server.URL
	client.httpClient = server.Client()

	matchIDs, err := client.ListIDsByPUUID(
		context.Background(),
		"player-puuid",
		matchdomain.ListOptions{Start: 10, Count: 5},
	)
	if err != nil {
		t.Fatalf("ListIDsByPUUID() returned an error: %v", err)
	}
	if len(matchIDs) != 2 || matchIDs[0] != "BR1_123" {
		t.Fatalf("match IDs = %#v, want Riot response", matchIDs)
	}
}

func TestClientGetsMatchByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/lol/match/v5/matches/BR1_123" {
			t.Errorf("path = %q, want Match-V5 match path", r.URL.Path)
		}
		if r.Header.Get("X-Riot-Token") != "test-key" {
			t.Error("X-Riot-Token header was not set")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"metadata":{"dataVersion":"2","matchId":"BR1_123","participants":["puuid"]},
			"info":{"gameMode":"CLASSIC","queueId":420,"participants":[{
				"puuid":"puuid","championName":"Ahri","kills":10,"deaths":2,"assists":8,"win":true
			}]}
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", "americas")
	if err != nil {
		t.Fatalf("NewClient() returned an error: %v", err)
	}
	client.baseURL = server.URL
	client.httpClient = server.Client()

	result, err := client.GetByID(context.Background(), "BR1_123")
	if err != nil {
		t.Fatalf("GetByID() returned an error: %v", err)
	}
	if result.Metadata.MatchID != "BR1_123" {
		t.Fatalf("match ID = %q, want %q", result.Metadata.MatchID, "BR1_123")
	}
	if len(result.Info.Participants) != 1 || result.Info.Participants[0].ChampionName != "Ahri" {
		t.Fatalf("participants = %#v, want decoded participant", result.Info.Participants)
	}
}
