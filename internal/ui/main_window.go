package ui

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"cryptoview/internal/api"
	"cryptoview/internal/model"
	"cryptoview/internal/ui/assets"
	"cryptoview/internal/ui/components"
	"cryptoview/internal/ui/i18n"
	uitheme "cryptoview/internal/ui/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func BuildMainWindow(a fyne.App, data []model.Coin) fyne.Window {
	a.Settings().SetTheme(uitheme.NewForMode(uitheme.ModeSystem))

	translator := i18n.NewTranslator(i18n.LangEN)
	w := a.NewWindow(translator.T("app.title"))
	w.Resize(fyne.NewSize(450, 480))
	w.SetFixedSize(true)
	appIcon := assets.LoadResource("resources/Logo/CryptoView Icon.png")
	if appIcon == nil {
		appIcon = theme.FyneLogo()
	}
	w.SetIcon(appIcon)

	apiClient := api.NewClient(10 * time.Second)
	coinList := components.NewCoinList(data, translator)
	footer := NewFooterController(translator)

	currentCurrency := i18n.FiatUSD
	currentLanguage := i18n.LangEN
	var requestID int64
	var fetchInFlight atomic.Bool
	var header *components.Toolbar

	triggerFetch := func(reason string, force bool) {
		_ = reason
		if !force && !fetchInFlight.CompareAndSwap(false, true) {
			return
		}
		if force {
			fetchInFlight.Store(true)
		}

		apiValue, ok := currentCurrency.APIValue()
		if !ok {
			fetchInFlight.Store(false)
			return
		}

		fyne.Do(func() {
			footer.SetLoading()
		})

		id := atomic.AddInt64(&requestID, 1)
		go func(localID int64) {
			defer fetchInFlight.Store(false)

			ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
			defer cancel()

			markets, err := apiClient.GetMarkets(ctx, apiValue)
			if err != nil {
				log.Printf("fetch markets failed for %s: %v", apiValue, err)
				if atomic.LoadInt64(&requestID) != localID {
					return
				}
				fyne.Do(func() {
					footer.SetError(translator.T("status.error.network"))
				})
				return
			}

			coins := make([]model.Coin, 0, len(markets))
			for _, market := range markets {
				coins = append(coins, model.ToCoin(market))
			}

			if atomic.LoadInt64(&requestID) != localID {
				return
			}

			fyne.Do(func() {
				coinList.ReplaceData(coins)
				footer.SetOK()
			})
		}(id)
	}

	header = components.NewToolbar(
		a,
		translator,
		func(currency i18n.FiatCurrency) {
			currentCurrency = currency
			coinList.SetCurrency(currency)
			triggerFetch("currency", true)
		},
		nil,
		func(language i18n.AppLanguage) {
			currentLanguage = language
			translator.SetLanguage(language)
			coinList.SetLanguage(language)
			if header != nil {
				header.SetLanguage(language)
			}
			footer.SetLanguage(language)
			w.SetTitle(translator.T("app.title"))
		},
		func() {
			triggerFetch("manual", true)
		},
	)

	content := container.NewBorder(header.CanvasObject(), footer.CanvasObject(), nil, nil, coinList.Widget())
	w.SetContent(content)
	coinList.SetCurrency(currentCurrency)
	coinList.SetLanguage(currentLanguage)
	footer.SetLoading()
	triggerFetch("initial", true)

	refreshTicker := time.NewTicker(60 * time.Second)
	stopCh := make(chan struct{})
	var stopOnce sync.Once

	go func() {
		for {
			select {
			case <-refreshTicker.C:
				triggerFetch("auto", false)
			case <-stopCh:
				return
			}
		}
	}()

	w.SetCloseIntercept(func() {
		stopOnce.Do(func() {
			refreshTicker.Stop()
			close(stopCh)
		})
		w.Close()
	})

	return w
}
