package ui

import (
	"strings"
	"sync"
	"sync/atomic"

	"cryptoview/internal/marketfeed"
	"cryptoview/internal/model"
	"cryptoview/internal/ui/assets"
	"cryptoview/internal/ui/components"
	"cryptoview/internal/ui/i18n"
	uitheme "cryptoview/internal/ui/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

type marketFeed interface {
	Start()
	Stop()
	SetFiat(i18n.FiatCurrency)
	SetCallbacks(marketfeed.Callbacks)
}

type noopFeed struct{}

func (noopFeed) Start()                            {}
func (noopFeed) Stop()                             {}
func (noopFeed) SetFiat(i18n.FiatCurrency)         {}
func (noopFeed) SetCallbacks(marketfeed.Callbacks) {}

func BuildMainWindow(a fyne.App, data []model.Coin, feed marketFeed) fyne.Window {
	if feed == nil {
		feed = noopFeed{}
	}

	a.Settings().SetTheme(uitheme.NewForMode(uitheme.ModeSystem))

	translator := i18n.NewTranslator(i18n.LangEN)
	w := a.NewWindow(translator.T("app.title"))
	w.Resize(fyne.NewSize(450, 480))
	w.SetFixedSize(true)

	appIcon := assets.LoadAppIcon()
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

	feed.SetCallbacks(marketfeed.Callbacks{
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
				case marketfeed.StatusKindError:
					footer.SetError(errorStatusMessage(translator, event))
				default:
					footer.SetError(errorStatusMessage(translator, event))
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
	w.SetOnClosed(func() {
		stopOnce.Do(func() {
			feed.Stop()
		})
	})
	w.SetCloseIntercept(func() {
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
	return base + " \u2022 " + name
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

func errorStatusMessage(translator *i18n.Translator, event marketfeed.StatusEvent) string {
	if translator == nil {
		return "Network error"
	}
	if event.Code == marketfeed.StatusCodeNoData {
		return translator.T("status.error.no_data")
	}
	return translator.T("status.error.network")
}
