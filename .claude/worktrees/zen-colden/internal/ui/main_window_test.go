package ui

import (
	"testing"

	"cryptoview/internal/model"
	"cryptoview/internal/service/marketfeed"
	"cryptoview/internal/ui/i18n"
	"fyne.io/fyne/v2/test"
)

func TestProviderDisplayName(t *testing.T) {
	tests := []struct {
		provider string
		want     string
	}{
		{"coingecko", "CoinGecko"},
		{"coincap", "CoinCap"},
		{"coinpaprika", "CoinPaprika"},
		{"cryptocompare", "CryptoCompare"},
		{"binance", "Binance"},
		{"coinlore", "CoinLore"},
		{"open-er-api", "Open ER API"},
		{"  COINGECKO  ", "CoinGecko"},
		{"", ""},
		{"unknown", "unknown"},
	}
	for _, tt := range tests {
		got := providerDisplayName(tt.provider)
		if got != tt.want {
			t.Errorf("providerDisplayName(%q) = %q, want %q", tt.provider, got, tt.want)
		}
	}
}

func TestOkStatusMessage(t *testing.T) {
	tr := i18n.NewTranslator(i18n.LangEN)
	if got := okStatusMessage(tr, "coingecko"); got != "OK • CoinGecko" {
		t.Fatalf("expected OK • CoinGecko, got %q", got)
	}
	if got := okStatusMessage(tr, ""); got != "OK" {
		t.Fatalf("expected OK for empty provider, got %q", got)
	}
	if got := okStatusMessage(nil, "coincap"); got != "OK • CoinCap" {
		t.Fatalf("expected OK • CoinCap with nil translator, got %q", got)
	}
}

func TestErrorStatusMessage(t *testing.T) {
	tr := i18n.NewTranslator(i18n.LangEN)

	noData := errorStatusMessage(tr, marketfeed.StatusEvent{Code: marketfeed.StatusCodeNoData})
	if noData != "No market data available" {
		t.Fatalf("expected EN no-data message, got %q", noData)
	}

	network := errorStatusMessage(tr, marketfeed.StatusEvent{})
	if network != "Network error" {
		t.Fatalf("expected EN network error message, got %q", network)
	}

	tr.SetLanguage(i18n.LangRU)
	noDataRU := errorStatusMessage(tr, marketfeed.StatusEvent{Code: marketfeed.StatusCodeNoData})
	if noDataRU != "Нет данных рынка" {
		t.Fatalf("expected RU no-data message, got %q", noDataRU)
	}
}

func TestBuildMainWindow_Smoke(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	data := model.GetMockCoins()
	w := buildMainWindowWithFeedFactory(a, data, func(callbacks marketfeed.Callbacks) marketFeed {
		return newFakeFeed(callbacks)
	})

	if w == nil {
		t.Fatal("expected non-nil window")
	}
	if w.Content() == nil {
		t.Fatal("expected window to have content")
	}
	if w.Title() != "CryptoView" {
		t.Fatalf("expected title CryptoView, got %q", w.Title())
	}
	// Trigger close to stop marketfeed goroutines (SetCloseIntercept calls feed.Stop)
	w.Close()
}
