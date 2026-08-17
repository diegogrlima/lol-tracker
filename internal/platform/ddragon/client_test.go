package ddragonadapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientListsChampions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cdn/16.1.1/data/pt_BR/champion.json" {
			t.Fatalf("path = %q, want champion endpoint", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": {
				"Ahri": {
					"id": "Ahri",
					"key": "103",
					"name": "Ahri",
					"title": "a Raposa de Nove Caudas",
					"blurb": "Descrição",
					"tags": ["Mage", "Assassin"],
					"image": {"full": "Ahri.png"}
				}
			}
		}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "16.1.1", "pt_BR")
	if err != nil {
		t.Fatalf("NewClient() returned an error: %v", err)
	}

	champions, err := client.List(context.Background())
	if err != nil {
		t.Fatalf("List() returned an error: %v", err)
	}
	if len(champions) != 1 {
		t.Fatalf("len(champions) = %d, want 1", len(champions))
	}
	if champions[0].Name != "Ahri" {
		t.Fatalf("Name = %q, want Ahri", champions[0].Name)
	}
	wantImageURL := server.URL + "/cdn/16.1.1/img/champion/Ahri.png"
	if champions[0].ImageURL != wantImageURL {
		t.Fatalf("ImageURL = %q, want %q", champions[0].ImageURL, wantImageURL)
	}
}

func TestClientReturnsErrorForUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "invalid", "pt_BR")
	if err != nil {
		t.Fatalf("NewClient() returned an error: %v", err)
	}

	if _, err := client.List(context.Background()); err == nil {
		t.Fatal("List() returned nil error for status 404")
	}
}
