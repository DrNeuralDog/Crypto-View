package providers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cryptoview/internal/catalog"
	"cryptoview/internal/marketfeed"
	"cryptoview/internal/ui/i18n"
)

func TestCoinGeckoProviderFetchUSDBuildsQueryAndDecodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/coins/markets" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("vs_currency"); got != "usd" {
			t.Fatalf("unexpected vs_currency: %s", got)
		}
		if got := r.URL.Query().Get("ids"); got != catalog.IDsCSV() {
			t.Fatalf("unexpected ids: %s", got)
		}
		if got := r.URL.Query().Get("price_change_percentage"); got != "24h" {
			t.Fatalf("unexpected price_change_percentage: %s", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":"bitcoin","symbol":"btc","name":"Bitcoin","current_price":100.5,"price_change_percentage_24h":1.2,"last_updated":"2026-02-20T10:11:12Z"}]`))
	}))
	defer server.Close()

	provider := &CoinGeckoProvider{
		httpClient: &http.Client{Timeout: time.Second},
		baseURL:    server.URL,
	}

	snapshot, err := provider.FetchUSD(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snapshot.Provider != provider.Name() {
		t.Fatalf("expected provider name %s, got %s", provider.Name(), snapshot.Provider)
	}
	quote, ok := snapshot.Coins["bitcoin"]
	if !ok {
		t.Fatal("expected bitcoin quote")
	}
	if quote.PriceUSD != 100.5 {
		t.Fatalf("expected price 100.5, got %f", quote.PriceUSD)
	}
	if quote.LastUpdate.IsZero() {
		t.Fatal("expected parsed last update time")
	}
}

func TestCoinGeckoProviderFetchUSDStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "7")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	provider := &CoinGeckoProvider{
		httpClient: &http.Client{Timeout: time.Second},
		baseURL:    server.URL,
	}

	_, err := provider.FetchUSD(context.Background())
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}

	var providerErr *marketfeed.ProviderError
	if !errors.As(err, &providerErr) {
		t.Fatalf("expected ProviderError, got %T", err)
	}
	if providerErr.Kind != marketfeed.FailureKindRateLimit {
		t.Fatalf("expected rate-limit error, got %s", providerErr.Kind)
	}
	if providerErr.RetryAfter != 7*time.Second {
		t.Fatalf("expected retry-after 7s, got %v", providerErr.RetryAfter)
	}
}

func TestOpenExchangeRatesProviderFiltersInvalidRates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"rates":{"EUR":0,"RUB":999999},"time_last_update_unix":1700000000}`))
	}))
	defer server.Close()

	provider := &OpenExchangeRatesProvider{
		httpClient: &http.Client{Timeout: time.Second},
		baseURL:    server.URL,
	}

	snapshot, err := provider.FetchRates(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snapshot.Rates[i18n.FiatUSD] != 1 {
		t.Fatalf("expected USD base rate, got %v", snapshot.Rates[i18n.FiatUSD])
	}
	if _, ok := snapshot.Rates[i18n.FiatEUR]; ok {
		t.Fatal("expected invalid EUR rate to be dropped")
	}
	if _, ok := snapshot.Rates[i18n.FiatRUB]; ok {
		t.Fatal("expected invalid RUB rate to be dropped")
	}
}
