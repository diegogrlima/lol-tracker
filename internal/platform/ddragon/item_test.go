package ddragonadapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientListsItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cdn/16.1.1/data/pt_BR/item.json" {
			t.Fatalf("path = %q, want item endpoint", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": {
				"1001": {
					"name": "Botas",
					"description": "Aumenta a velocidade de movimento",
					"plaintext": "Aumenta levemente a velocidade de movimento",
					"tags": ["Boots"],
					"gold": {"base": 300, "total": 300, "sell": 210, "purchasable": true},
					"image": {"full": "1001.png"}
				}
			}
		}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "16.1.1", "pt_BR")
	if err != nil {
		t.Fatalf("NewClient() returned an error: %v", err)
	}

	items, err := client.ListItems(context.Background())
	if err != nil {
		t.Fatalf("ListItems() returned an error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}

	got := items[0]
	if got.ID != "1001" || got.Name != "Botas" || !got.Purchasable {
		t.Fatalf("item = %#v, want mapped store item", got)
	}
	if got.BasePrice != 300 || got.TotalPrice != 300 || got.SellPrice != 210 {
		t.Fatalf("item prices = %d/%d/%d, want 300/300/210", got.BasePrice, got.TotalPrice, got.SellPrice)
	}

	wantImageURL := server.URL + "/cdn/16.1.1/img/item/1001.png"
	if got.ImageURL != wantImageURL {
		t.Fatalf("ImageURL = %q, want %q", got.ImageURL, wantImageURL)
	}
}

func TestClientGetsItemByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"1001":{"name":"Botas","gold":{},"image":{"full":"1001.png"}}}}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "16.1.1", "pt_BR")
	if err != nil {
		t.Fatalf("NewClient() returned an error: %v", err)
	}

	got, err := client.GetItemByID(context.Background(), "1001")
	if err != nil {
		t.Fatalf("GetItemByID() returned an error: %v", err)
	}
	if got.ID != "1001" || got.Name != "Botas" {
		t.Fatalf("item = %#v, want item 1001", got)
	}
}
