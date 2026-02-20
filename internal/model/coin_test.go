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
