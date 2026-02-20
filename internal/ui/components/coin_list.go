package components

import (
	"fmt"
	"image/color"
	"sync"

	"cryptoview/internal/model"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FiatCurrency string

const (
	FiatUSD FiatCurrency = "USD"
	FiatEUR FiatCurrency = "EUR"
	FiatRUB FiatCurrency = "RUB"
)

type CoinListController struct {
	list     *widget.List
	data     []model.Coin
	currency FiatCurrency
	mu       sync.RWMutex
}

func NewCoinList(data []model.Coin) *CoinListController {
	controller := &CoinListController{
		data:     data,
		currency: FiatUSD,
	}

	controller.list = widget.NewList(
		func() int {
			controller.mu.RLock()
			defer controller.mu.RUnlock()
			return len(controller.data)
		},
		func() fyne.CanvasObject {
			icon := widget.NewLabel("[icon]")
			nameTicker := widget.NewLabel("")
			price := widget.NewLabel("")
			change := canvas.NewText("", color.NRGBA{R: 128, G: 128, B: 128, A: 255})
			lastUpdate := widget.NewLabel("")

			return container.NewHBox(
				icon,
				nameTicker,
				layout.NewSpacer(),
				price,
				change,
				lastUpdate,
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			controller.mu.RLock()
			if id < 0 || int(id) >= len(controller.data) {
				controller.mu.RUnlock()
				return
			}
			coin := controller.data[id]
			currency := controller.currency
			controller.mu.RUnlock()

			row := item.(*fyne.Container)
			icon := row.Objects[0].(*widget.Label)
			nameTicker := row.Objects[1].(*widget.Label)
			price := row.Objects[3].(*widget.Label)
			change := row.Objects[4].(*canvas.Text)
			lastUpdate := row.Objects[5].(*widget.Label)

			icon.SetText(fmt.Sprintf("[%s]", coin.Ticker))
			nameTicker.SetText(fmt.Sprintf("%s | %s", coin.Name, coin.Ticker))
			price.SetText(formatPriceByCurrency(coin.Price, currency))
			change.Text = fmt.Sprintf("%.2f%%", coin.Change24h)
			change.Color = changeColor(coin.Change24h)
			change.Refresh()
			lastUpdate.SetText(coin.LastUpdateTime)
		},
	)

	return controller
}

func (c *CoinListController) Widget() *widget.List {
	return c.list
}

func (c *CoinListController) SetCurrency(currency FiatCurrency) {
	if _, ok := parseFiatCurrency(string(currency)); !ok {
		return
	}
	c.mu.Lock()
	c.currency = currency
	c.mu.Unlock()
	fyne.Do(func() {
		c.list.Refresh()
	})
}

func (c *CoinListController) ReplaceData(coins []model.Coin) {
	c.mu.Lock()
	c.data = coins
	c.mu.Unlock()
	fyne.Do(func() {
		c.list.Refresh()
	})
}

func parseFiatCurrency(raw string) (FiatCurrency, bool) {
	switch FiatCurrency(raw) {
	case FiatUSD, FiatEUR, FiatRUB:
		return FiatCurrency(raw), true
	default:
		return "", false
	}
}

func (f FiatCurrency) APIValue() (string, bool) {
	switch f {
	case FiatUSD:
		return "usd", true
	case FiatEUR:
		return "eur", true
	case FiatRUB:
		return "rub", true
	default:
		return "", false
	}
}

func formatPriceByCurrency(price float64, currency FiatCurrency) string {
	symbol := "$"
	switch currency {
	case FiatEUR:
		symbol = "€"
	case FiatRUB:
		symbol = "₽"
	}
	return fmt.Sprintf("%s%.2f", symbol, price)
}

func changeColor(change24h float64) color.Color {
	if change24h > 0 {
		return color.NRGBA{R: 34, G: 139, B: 34, A: 255}
	}
	if change24h < 0 {
		return color.NRGBA{R: 220, G: 20, B: 60, A: 255}
	}
	return color.NRGBA{R: 128, G: 128, B: 128, A: 255}
}
