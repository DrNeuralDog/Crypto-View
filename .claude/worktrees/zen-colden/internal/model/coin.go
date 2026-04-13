package model

import (
	"strings"
	"time"
)

type Coin struct {
	ID             string
	Name           string
	Ticker         string
	Price          float64
	Change24h      float64
	LastUpdateTime string
	IconPath       string
}

var coinIconPathByID = map[string]string{
	"bitcoin":          "resources/coins/bitcoin.png",
	"ethereum":         "resources/coins/ethereum.png",
	"the-open-network": "resources/coins/the-open-network.png",
	"toncoin":          "resources/coins/the-open-network.png",
	"solana":           "resources/coins/solana.png",
	"dogecoin":         "resources/coins/dogecoin.png",
	"ripple":           "resources/coins/ripple.png",
	"litecoin":         "resources/coins/litecoin.png",
}

type CoinGeckoMarket struct {
	ID                       string  `json:"id"`
	Symbol                   string  `json:"symbol"`
	Name                     string  `json:"name"`
	CurrentPrice             float64 `json:"current_price"`
	PriceChangePercentage24h float64 `json:"price_change_percentage_24h"`
	LastUpdated              string  `json:"last_updated"`
}

func ToCoin(m CoinGeckoMarket) Coin {
	return Coin{
		ID:             m.ID,
		Name:           m.Name,
		Ticker:         strings.ToUpper(m.Symbol),
		Price:          m.CurrentPrice,
		Change24h:      m.PriceChangePercentage24h,
		LastUpdateTime: formatLastUpdated(m.LastUpdated),
		IconPath:       iconPathForID(m.ID),
	}
}

func iconPathForID(id string) string {
	if path, ok := coinIconPathByID[id]; ok {
		return path
	}
	return ""
}

func IconPathForID(id string) string {
	return iconPathForID(id)
}

func formatLastUpdated(raw string) string {
	if raw == "" {
		return "--:--:--"
	}

	if ts, err := time.Parse(time.RFC3339Nano, raw); err == nil {
		return ts.Local().Format("15:04:05")
	}

	if ts, err := time.Parse(time.RFC3339, raw); err == nil {
		return ts.Local().Format("15:04:05")
	}

	return "--:--:--"
}

func GetMockCoins() []Coin {
	return []Coin{
		{
			ID:             "bitcoin",
			Name:           "Bitcoin",
			Ticker:         "BTC",
			Price:          96543.12,
			Change24h:      2.54,
			LastUpdateTime: "12:34:56",
			IconPath:       iconPathForID("bitcoin"),
		},
		{
			ID:             "ethereum",
			Name:           "Ethereum",
			Ticker:         "ETH",
			Price:          3421.77,
			Change24h:      -1.23,
			LastUpdateTime: "12:34:56",
			IconPath:       iconPathForID("ethereum"),
		},
		{
			ID:             "the-open-network",
			Name:           "TON Coin",
			Ticker:         "TON",
			Price:          5.89,
			Change24h:      0.00,
			LastUpdateTime: "12:34:56",
			IconPath:       iconPathForID("the-open-network"),
		},
		{
			ID:             "solana",
			Name:           "Solana",
			Ticker:         "SOL",
			Price:          183.45,
			Change24h:      5.91,
			LastUpdateTime: "12:34:56",
			IconPath:       iconPathForID("solana"),
		},
		{
			ID:             "dogecoin",
			Name:           "Dogecoin",
			Ticker:         "DOGE",
			Price:          0.25,
			Change24h:      -3.02,
			LastUpdateTime: "12:34:56",
			IconPath:       iconPathForID("dogecoin"),
		},
		{
			ID:             "ripple",
			Name:           "Ripple",
			Ticker:         "XRP",
			Price:          0.71,
			Change24h:      1.04,
			LastUpdateTime: "12:34:56",
			IconPath:       iconPathForID("ripple"),
		},
		{
			ID:             "litecoin",
			Name:           "Litecoin",
			Ticker:         "LTC",
			Price:          102.33,
			Change24h:      -0.67,
			LastUpdateTime: "12:34:56",
			IconPath:       iconPathForID("litecoin"),
		},
	}
}
