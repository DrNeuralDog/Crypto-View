package ui

import (
	"strings"
	"sync"
	"sync/atomic"

	"cryptoview/internal/model"
	"cryptoview/internal/service/marketfeed"
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

	coinList := components.NewCoinList(data, translator)
	footer := NewFooterController(translator)

	currentCurrency := i18n.FiatUSD
	currentLanguage := i18n.LangEN
	var header *components.Toolbar
	var statusEventID int64
	feed := marketfeed.NewDefault(marketfeed.Callbacks{
		OnMarketUpdate: func(coins []model.Coin) {
			fyne.Do(func() {
				coinList.ReplaceData(coins)
			})
		},
		OnStatus: func(event marketfeed.StatusEvent) {
			localID := atomic.AddInt64(&statusEventID, 1)
			fyne.Do(func() {
				if atomic.LoadInt64(&statusEventID) != localID {
					return
				}
				switch event.Kind {
				case marketfeed.StatusKindLoading:
					footer.SetLoading()
				case marketfeed.StatusKindOK:
					footer.SetOKWithMessage(okStatusMessage(translator, event.Provider))
				case marketfeed.StatusKindWarning:
					switch event.Code {
					case marketfeed.StatusCodeRateLimited:
						footer.SetWarning(translator.T("status.warning.rate"))
					case marketfeed.StatusCodeFallback:
						footer.SetOKWithMessage(okStatusMessage(translator, event.Provider))
					default:
						footer.SetWarning(translator.T("status.warning.cached"))
					}
				default:
					footer.SetError(translator.T("status.error.network"))
				}
			})
		},
	})

	header = components.NewToolbar(
		a,
		translator,
		func(currency i18n.FiatCurrency) {
			currentCurrency = currency
			coinList.SetCurrency(currency)
			feed.SetFiat(currency)
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
	)

	content := container.NewBorder(header.CanvasObject(), footer.CanvasObject(), nil, nil, coinList.Widget())
	w.SetContent(content)
	coinList.SetCurrency(currentCurrency)
	coinList.SetLanguage(currentLanguage)
	footer.SetLoading()
	feed.Start()

	var stopOnce sync.Once

	w.SetCloseIntercept(func() {
		stopOnce.Do(func() {
			feed.Stop()
		})
		w.Close()
	})

	return w
}

func okStatusMessage(translator *i18n.Translator, provider string) string {
	base := "OK"
	if translator != nil {
		base = translator.T("status.ok")
	}
	name := providerDisplayName(provider)
	if name == "" {
		return base
	}
	return base + " • " + name
}

func providerDisplayName(provider string) string {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "coingecko":
		return "CoinGecko"
	case "coincap":
		return "CoinCap"
	case "coinpaprika":
		return "CoinPaprika"
	case "cryptocompare":
		return "CryptoCompare"
	case "binance":
		return "Binance"
	case "coinlore":
		return "CoinLore"
	case "open-er-api":
		return "Open ER API"
	default:
		if provider == "" {
			return ""
		}
		return provider
	}
}
