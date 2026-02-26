package model

import "testing"

func TestGetMockCoins(t *testing.T) {
	coins := GetMockCoins()
	if len(coins) != 7 {
		t.Fatalf("expected 7 coins, got %d", len(coins))
	}

	for i, coin := range coins {
		if coin.Ticker == "" {
			t.Fatalf("coin[%d] has empty ticker", i)
		}
		if coin.Name == "" {
			t.Fatalf("coin[%d] has empty name", i)
		}
		if coin.Price <= 0 {
			t.Fatalf("coin[%d] has non-positive price: %f", i, coin.Price)
		}
	}
}

func TestToCoin(t *testing.T) {
	src := CoinGeckoMarket{
		ID:                       "bitcoin",
		Symbol:                   "btc",
		Name:                     "Bitcoin",
		CurrentPrice:             123.45,
		PriceChangePercentage24h: 2.34,
		LastUpdated:              "2026-02-20T10:11:12Z",
	}

	coin := ToCoin(src)
	if coin.Ticker != "BTC" {
		t.Fatalf("expected upper ticker BTC, got %s", coin.Ticker)
	}
	if coin.LastUpdateTime == "--:--:--" || coin.LastUpdateTime == "" {
		t.Fatalf("expected parsed last update time, got %s", coin.LastUpdateTime)
	}
	if coin.IconPath == "" {
		t.Fatal("expected icon path to be mapped by coin id")
	}
}

func TestToCoinInvalidTimestamp(t *testing.T) {
	src := CoinGeckoMarket{
		ID:                       "bitcoin",
		Symbol:                   "btc",
		Name:                     "Bitcoin",
		CurrentPrice:             123.45,
		PriceChangePercentage24h: 2.34,
		LastUpdated:              "not-a-time",
	}

	coin := ToCoin(src)
	if coin.LastUpdateTime != "--:--:--" {
		t.Fatalf("expected fallback time, got %s", coin.LastUpdateTime)
	}
}

func TestToCoinEmptyLastUpdated(t *testing.T) {
	src := CoinGeckoMarket{
		ID:                       "bitcoin",
		Symbol:                   "btc",
		Name:                     "Bitcoin",
		CurrentPrice:             123.45,
		PriceChangePercentage24h: 2.34,
		LastUpdated:              "",
	}

	coin := ToCoin(src)
	if coin.LastUpdateTime != "--:--:--" {
		t.Fatalf("expected fallback for empty last updated, got %s", coin.LastUpdateTime)
	}
}

func TestIconPathForID(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"bitcoin", "resources/coins/bitcoin.png"},
		{"ethereum", "resources/coins/ethereum.png"},
		{"the-open-network", "resources/coins/the-open-network.png"},
		{"toncoin", "resources/coins/the-open-network.png"},
		{"solana", "resources/coins/solana.png"},
		{"dogecoin", "resources/coins/dogecoin.png"},
		{"ripple", "resources/coins/ripple.png"},
		{"litecoin", "resources/coins/litecoin.png"},
		{"unknown-coin", ""},
		{"", ""},
	}
	for _, tt := range tests {
		got := IconPathForID(tt.id)
		if got != tt.want {
			t.Errorf("IconPathForID(%q) = %q, want %q", tt.id, got, tt.want)
		}
	}
}
