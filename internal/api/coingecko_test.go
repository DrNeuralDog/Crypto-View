package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetMarketsBuildsQueryAndDecodes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/coins/markets" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("vs_currency"); got != "usd" {
			t.Fatalf("unexpected vs_currency: %s", got)
		}
		if got := r.URL.Query().Get("ids"); got != trackedCoinIDs {
			t.Fatalf("unexpected ids: %s", got)
		}
		if got := r.URL.Query().Get("price_change_percentage"); got != "24h" {
			t.Fatalf("unexpected price_change_percentage: %s", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":"bitcoin","symbol":"btc","name":"Bitcoin","current_price":100.5,"price_change_percentage_24h":1.2,"last_updated":"2026-02-20T10:11:12Z"}]`))
	}))
	defer srv.Close()

	client := newClient(srv.URL, time.Second)
	markets, err := client.GetMarkets(context.Background(), "USD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(markets) != 1 {
		t.Fatalf("expected 1 market, got %d", len(markets))
	}
	if markets[0].ID != "bitcoin" {
		t.Fatalf("unexpected market id: %s", markets[0].ID)
	}
}

func TestGetMarketsStatusError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	client := newClient(srv.URL, time.Second)
	_, err := client.GetMarkets(context.Background(), "usd")
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestGetMarketsContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	client := newClient(srv.URL, time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err := client.GetMarkets(ctx, "usd")
	if err == nil {
		t.Fatal("expected timeout/cancel error")
	}
}
