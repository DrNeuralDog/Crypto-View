package main

import (
	"log/slog"
	"time"

	"cryptoview/internal/marketfeed"
	"cryptoview/internal/providers"
	"cryptoview/internal/ui"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()

	feed, err := marketfeed.New(
		[]marketfeed.MarketProvider{
			providers.NewCoinGeckoProvider(1 * time.Second),
			providers.NewCryptoCompareProvider(3 * time.Second),
			providers.NewCoinLoreProvider(3 * time.Second),
		},
		providers.NewOpenExchangeRatesProvider(1*time.Second),
		marketfeed.Callbacks{},
		marketfeed.WithLogger(slog.Default()),
	)
	if err != nil {
		slog.Error("failed to initialize market feed", "error", err)
		return
	}

	w := ui.BuildMainWindow(a, nil, feed)
	w.ShowAndRun()
}
