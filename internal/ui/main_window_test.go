package ui

import (
	"testing"

	"cryptoview/internal/model"
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

func TestBuildMainWindow_Smoke(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	data := model.GetMockCoins()
	w := BuildMainWindow(a, data)

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
