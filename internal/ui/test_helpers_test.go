package ui

import (
	"time"

	"cryptoview/internal/model"
)

func mockCoins() []model.Coin {
	ts := time.Date(2026, time.February, 20, 10, 11, 12, 0, time.UTC)
	return []model.Coin{
		{ID: "bitcoin", Name: "Bitcoin", Ticker: "BTC", Price: 96543.12, Change24h: 2.54, LastUpdated: ts},
		{ID: "ethereum", Name: "Ethereum", Ticker: "ETH", Price: 3421.77, Change24h: -1.23, LastUpdated: ts},
		{ID: "the-open-network", Name: "TON Coin", Ticker: "TON", Price: 5.89, Change24h: 0, LastUpdated: ts},
		{ID: "solana", Name: "Solana", Ticker: "SOL", Price: 183.45, Change24h: 5.91, LastUpdated: ts},
		{ID: "dogecoin", Name: "Dogecoin", Ticker: "DOGE", Price: 0.25, Change24h: -3.02, LastUpdated: ts},
		{ID: "ripple", Name: "Ripple", Ticker: "XRP", Price: 0.71, Change24h: 1.04, LastUpdated: ts},
		{ID: "litecoin", Name: "Litecoin", Ticker: "LTC", Price: 102.33, Change24h: -0.67, LastUpdated: ts},
	}
}
