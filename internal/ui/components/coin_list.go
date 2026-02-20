package components

import (
	"fmt"
	"image/color"
	"sync"

	"cryptoview/internal/model"
	"cryptoview/internal/ui/assets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
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
	icons    map[string]fyne.Resource
}

func NewCoinList(data []model.Coin) *CoinListController {
	controller := &CoinListController{
		data:     data,
		currency: FiatUSD,
		icons:    make(map[string]fyne.Resource),
	}

	controller.list = widget.NewList(
		func() int {
			controller.mu.RLock()
			defer controller.mu.RUnlock()
			return len(controller.data)
		},
		func() fyne.CanvasObject {
			return newCoinListItem()
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			controller.mu.RLock()
			if id < 0 || int(id) >= len(controller.data) {
				controller.mu.RUnlock()
				return
			}
			coin := controller.data[id]
			currency := controller.currency
			isLast := int(id) == len(controller.data)-1
			controller.mu.RUnlock()

			row := item.(*coinListItem)
			row.applyCoin(
				coin,
				formatPriceByCurrency(coin.Price, currency),
				changeColor(coin.Change24h),
				controller.iconForCoin(coin),
				isLast,
			)
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

func (c *CoinListController) iconForCoin(coin model.Coin) fyne.Resource {
	if coin.IconPath == "" {
		return nil
	}

	c.mu.RLock()
	cached, ok := c.icons[coin.IconPath]
	c.mu.RUnlock()
	if ok {
		return cached
	}

	resource := assets.LoadResource(coin.IconPath)
	if resource == nil {
		return nil
	}

	c.mu.Lock()
	c.icons[coin.IconPath] = resource
	c.mu.Unlock()
	return resource
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
		symbol = "\u20ac"
	case FiatRUB:
		symbol = "\u20bd"
	}
	return fmt.Sprintf("%s%.2f", symbol, price)
}

func changeColor(change24h float64) color.Color {
	if change24h > 0 {
		return color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
	}
	if change24h < 0 {
		return color.NRGBA{R: 0xF4, G: 0x43, B: 0x36, A: 0xFF}
	}
	return color.NRGBA{R: 128, G: 128, B: 128, A: 255}
}

type coinListItem struct {
	widget.BaseWidget

	root      *fyne.Container
	icon      *canvas.Image
	ticker    *widget.Label
	updatedAt *canvas.Text
	name      *widget.Label
	price     *widget.Label
	change    *canvas.Text
	chartIcon *widget.Icon
	separator *widget.Separator
}

func newCoinListItem() *coinListItem {
	icon := canvas.NewImageFromResource(theme.BrokenImageIcon())
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSize(26, 26))

	ticker := widget.NewLabel("BTC")
	ticker.TextStyle = fyne.TextStyle{Bold: true}

	updatedAt := canvas.NewText("--:--:--", theme.Color(theme.ColorNamePlaceHolder))
	updatedAt.TextSize = theme.CaptionTextSize()
	updatedAt.Alignment = fyne.TextAlignLeading

	name := widget.NewLabel("Bitcoin | BTC")

	price := widget.NewLabel("$0.00")
	price.TextStyle = fyne.TextStyle{Bold: true}

	change := canvas.NewText("+0.00%", color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF})
	change.TextStyle = fyne.TextStyle{Bold: true}
	change.TextSize = theme.TextSize()

	chartIcon := widget.NewIcon(theme.NewThemedResource(theme.ListIcon()))

	mainInfo := container.NewHBox(
		ticker,
		spacerX(8),
		name,
	)
	//timeRow := container.NewHBox(spacerX(1), updatedAt)
	meta := container.New(&vGapLayout{gap: 0.9}, mainInfo, updatedAt)
	row := container.NewHBox(
		container.NewCenter(icon),
		spacerX(8),
		meta,
		layout.NewSpacer(),
		container.NewCenter(price),
		spacerX(8),
		container.NewCenter(change),
		spacerX(8),
		container.NewCenter(chartIcon),
	)
	separator := widget.NewSeparator()
	content := container.NewVBox(container.NewPadded(container.NewPadded(row)), separator)

	item := &coinListItem{
		root:      content,
		icon:      icon,
		ticker:    ticker,
		updatedAt: updatedAt,
		name:      name,
		price:     price,
		change:    change,
		chartIcon: chartIcon,
		separator: separator,
	}
	item.ExtendBaseWidget(item)
	return item
}

func (i *coinListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(i.root)
}

func (i *coinListItem) applyCoin(
	coin model.Coin,
	price string,
	changeColor color.Color,
	iconResource fyne.Resource,
	isLast bool,
) {
	i.ticker.SetText(coin.Ticker)
	i.updatedAt.Text = coin.LastUpdateTime
	i.updatedAt.Color = theme.Color(theme.ColorNamePlaceHolder)
	i.updatedAt.Refresh()

	i.name.SetText(fmt.Sprintf("%s | %s", coin.Name, coin.Ticker))
	i.price.SetText(price)
	i.change.Text = fmt.Sprintf("%+.2f%%", coin.Change24h)
	i.change.Color = changeColor
	i.change.Refresh()

	if iconResource != nil {
		i.icon.Resource = iconResource
	} else {
		i.icon.Resource = theme.BrokenImageIcon()
	}
	i.icon.Refresh()

	if isLast {
		i.separator.Hide()
	} else {
		i.separator.Show()
	}
}

func spacerX(width float32) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(width, 1))
	return rect
}

type vGapLayout struct {
	gap float32
}

func (l *vGapLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	y := float32(0)
	for idx, obj := range objects {
		min := obj.MinSize()
		obj.Resize(fyne.NewSize(size.Width, min.Height))
		obj.Move(fyne.NewPos(0, y))
		y += min.Height
		if idx < len(objects)-1 {
			y += l.gap
		}
	}
}

func (l *vGapLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var width float32
	var height float32
	for idx, obj := range objects {
		min := obj.MinSize()
		if min.Width > width {
			width = min.Width
		}
		height += min.Height
		if idx < len(objects)-1 {
			height += l.gap
		}
	}
	return fyne.NewSize(width, height)
}
