package components

import (
	"fmt"
	"image/color"
	"sync"

	"cryptoview/internal/model"
	"cryptoview/internal/ui/assets"
	"cryptoview/internal/ui/i18n"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CoinListController struct {
	list       *widget.List
	data       []model.Coin
	currency   i18n.FiatCurrency
	language   i18n.AppLanguage
	translator *i18n.Translator
	mu         sync.RWMutex
	icons      map[string]fyne.Resource
	tickerW    float32
}

func NewCoinList(data []model.Coin, translator *i18n.Translator) *CoinListController {
	if translator == nil {
		translator = i18n.NewTranslator(i18n.LangEN)
	}
	controller := &CoinListController{
		data:       data,
		currency:   i18n.FiatUSD,
		language:   translator.Language(),
		translator: translator,
		icons:      make(map[string]fyne.Resource),
		tickerW:    maxTickerWidth(data),
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
			language := controller.language
			tickerW := controller.tickerW
			isLast := int(id) == len(controller.data)-1
			controller.mu.RUnlock()

			row := item.(*coinListItem)
			row.applyCoin(
				coin,
				i18n.FormatPrice(coin.Price, currency, language),
				i18n.FormatTime(coin.LastUpdateTime, language),
				changeColor(coin.Change24h),
				controller.iconForCoin(coin),
				tickerW,
				isLast,
			)
		},
	)
	return controller
}

func (c *CoinListController) Widget() *widget.List {
	return c.list
}

func (c *CoinListController) SetCurrency(currency i18n.FiatCurrency) {
	if _, ok := i18n.ParseFiatCurrency(string(currency)); !ok {
		return
	}
	c.mu.Lock()
	c.currency = currency
	c.mu.Unlock()
	fyne.Do(func() {
		c.list.Refresh()
	})
}

func (c *CoinListController) SetLanguage(language i18n.AppLanguage) {
	c.mu.Lock()
	c.language = language
	c.mu.Unlock()
	fyne.Do(func() {
		c.list.Refresh()
	})
}

func (c *CoinListController) ReplaceData(coins []model.Coin) {
	c.mu.Lock()
	c.data = coins
	c.tickerW = maxTickerWidth(coins)
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

func changeColor(change24h float64) color.Color {
	if change24h > 0 {
		return color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
	}
	if change24h < 0 {
		return color.NRGBA{R: 0xF4, G: 0x43, B: 0x36, A: 0xFF}
	}
	return color.NRGBA{R: 128, G: 128, B: 128, A: 255}
}

func maxTickerWidth(coins []model.Coin) float32 {
	base := widget.NewLabel("BTC")
	base.TextStyle = fyne.TextStyle{Bold: true}
	maxWidth := base.MinSize().Width
	for _, coin := range coins {
		if coin.Ticker == "" {
			continue
		}
		lbl := widget.NewLabel(coin.Ticker)
		lbl.TextStyle = fyne.TextStyle{Bold: true}
		width := lbl.MinSize().Width
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

type coinListItem struct {
	widget.BaseWidget

	root      *fyne.Container
	icon      *canvas.Image
	ticker    *widget.Label
	namePad   *canvas.Rectangle
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
	namePad := canvas.NewRectangle(color.Transparent)
	namePad.SetMinSize(fyne.NewSize(0, 1))

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
		namePad,
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
		namePad:   namePad,
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
	formattedTime string,
	changeColor color.Color,
	iconResource fyne.Resource,
	tickerW float32,
	isLast bool,
) {
	i.ticker.SetText(coin.Ticker)
	if tickerW > 0 {
		currentTickerW := i.ticker.MinSize().Width
		padW := tickerW - currentTickerW
		if padW < 0 {
			padW = 0
		}
		i.namePad.SetMinSize(fyne.NewSize(padW, 1))
		i.namePad.Refresh()
	}
	i.updatedAt.Text = formattedTime
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

	// Force parent row relayout so ticker fixed-width column changes affect name alignment.
	i.Refresh()
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
