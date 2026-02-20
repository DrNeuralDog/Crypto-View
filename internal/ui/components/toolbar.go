package components

import (
	"cryptoview/internal/ui/assets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Toolbar struct {
	root           fyne.CanvasObject
	themeButton    *widget.Button
	themeControl   *ThemeController
	currencySelect *widget.Select
	langSelect     *widget.Select
}

func NewToolbar(
	app fyne.App,
	onCurrencyChanged func(FiatCurrency),
	onThemeChanged func(),
	onLanguageChanged func(string),
) *Toolbar {
	title := widget.NewLabel("CryptoView")
	title.TextStyle = fyne.TextStyle{Bold: true}
	logoResource := assets.LoadResource("resources/Logo/CryptoView Icon.png")
	if logoResource == nil {
		logoResource = theme.FyneLogo()
	}
	logo := widget.NewIcon(logoResource)
	logoWrap := container.NewGridWrap(fyne.NewSize(28, 28), logo)

	currencySelect := widget.NewSelect([]string{string(FiatUSD), string(FiatEUR), string(FiatRUB)}, func(selected string) {
		if onCurrencyChanged == nil {
			return
		}
		currency, ok := parseFiatCurrency(selected)
		if !ok {
			return
		}
		onCurrencyChanged(currency)
	})
	currencySelect.SetSelected(string(FiatUSD))

	langSelect := widget.NewSelect([]string{"EN", "ENG"}, func(selected string) {
		if onLanguageChanged != nil {
			onLanguageChanged(selected)
		}
	})
	langSelect.SetSelected("EN")

	themeControl := NewThemeController(app)
	var themeButton *widget.Button
	themeButton = widget.NewButtonWithIcon("", themeControl.ActionIconResource(), func() {
		themeControl.Toggle()
		themeButton.SetIcon(themeControl.ActionIconResource())
		if onThemeChanged != nil {
			onThemeChanged()
		}
	})
	themeButton.Importance = widget.LowImportance
	themeButton.SetIcon(themeControl.ActionIconResource())
	themeButtonWrap := container.NewGridWrap(fyne.NewSize(56, 40), themeButton)

	left := container.NewHBox(logoWrap, title)
	right := container.NewHBox(currencySelect, langSelect, themeButtonWrap)
	header := container.NewBorder(nil, canvas.NewLine(theme.Color(theme.ColorNameSeparator)), left, right)

	return &Toolbar{
		root:           header,
		themeButton:    themeButton,
		themeControl:   themeControl,
		currencySelect: currencySelect,
		langSelect:     langSelect,
	}
}

func (t *Toolbar) CanvasObject() fyne.CanvasObject {
	return t.root
}

func (t *Toolbar) ThemeButton() *widget.Button {
	return t.themeButton
}

func (t *Toolbar) ThemeMode() string {
	return string(t.themeControl.Mode())
}

func (t *Toolbar) CurrencySelect() *widget.Select {
	return t.currencySelect
}

func (t *Toolbar) LanguageSelect() *widget.Select {
	return t.langSelect
}
