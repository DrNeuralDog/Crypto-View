package ui

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"cryptoview/internal/api"
	"cryptoview/internal/model"
	"cryptoview/internal/ui/components"
	uitheme "cryptoview/internal/ui/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func BuildMainWindow(a fyne.App, data []model.Coin) fyne.Window {
	a.Settings().SetTheme(uitheme.NewForMode(uitheme.ModeSystem))

	w := a.NewWindow("CryptoView")
	w.Resize(fyne.NewSize(450, 400))

	apiClient := api.NewClient(10 * time.Second)
	coinList := components.NewCoinList(data)
	var requestID int64

	fetchAndApply := func(currency components.FiatCurrency) {
		apiValue, ok := currency.APIValue()
		if !ok {
			return
		}

		id := atomic.AddInt64(&requestID, 1)
		go func(localID int64) {
			ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
			defer cancel()

			markets, err := apiClient.GetMarkets(ctx, apiValue)
			if err != nil {
				log.Printf("fetch markets failed for %s: %v", apiValue, err)
				return
			}

			coins := make([]model.Coin, 0, len(markets))
			for _, market := range markets {
				coins = append(coins, model.ToCoin(market))
			}

			if atomic.LoadInt64(&requestID) != localID {
				return
			}
			coinList.ReplaceData(coins)
		}(id)
	}

	header := components.NewToolbar(a, func(currency components.FiatCurrency) {
		coinList.SetCurrency(currency)
		fetchAndApply(currency)
	}, nil, nil)
	footer := NewFooter()

	content := container.NewBorder(header.CanvasObject(), footer, nil, nil, coinList.Widget())
	w.SetContent(content)
	fetchAndApply(components.FiatUSD)

	return w
}
