package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func NewToolbar(onCurrencyChanged func(FiatCurrency)) fyne.CanvasObject {
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

	langButton := widget.NewButton("Lang", func() {})
	themeButton := widget.NewButton("Theme", func() {})

	return container.NewHBox(
		widget.NewLabel("CryptoView"),
		layout.NewSpacer(),
		currencySelect,
		langButton,
		themeButton,
	)
}
