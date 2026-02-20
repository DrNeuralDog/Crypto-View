package components

import (
	"cryptoview/internal/ui/assets"
	"cryptoview/internal/ui/i18n"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Toolbar struct {
	root           fyne.CanvasObject
	title          *widget.Label
	refreshButton  *widget.Button
	themeButton    *widget.Button
	themeControl   *ThemeController
	currencySelect *widget.Select
	langSelect     *widget.Select
	translator     *i18n.Translator
}

func NewToolbar(
	app fyne.App,
	translator *i18n.Translator,
	onCurrencyChanged func(i18n.FiatCurrency),
	onThemeChanged func(),
	onLanguageChanged func(i18n.AppLanguage),
	onRefreshRequested func(),
) *Toolbar {
	if translator == nil {
		translator = i18n.NewTranslator(i18n.LangEN)
	}

	title := widget.NewLabel(translator.T("app.title"))
	title.TextStyle = fyne.TextStyle{Bold: true}
	logoResource := assets.LoadResource("resources/Logo/CryptoView Icon.png")
	if logoResource == nil {
		logoResource = theme.FyneLogo()
	}
	logo := widget.NewIcon(logoResource)
	logoWrap := container.NewGridWrap(fyne.NewSize(28, 28), logo)

	currencySelect := widget.NewSelect(
		[]string{string(i18n.FiatUSD), string(i18n.FiatEUR), string(i18n.FiatRUB)},
		func(selected string) {
			if onCurrencyChanged == nil {
				return
			}
			currency, ok := i18n.ParseFiatCurrency(selected)
			if !ok {
				return
			}
			onCurrencyChanged(currency)
		})
	currencySelect.SetSelected(string(i18n.FiatUSD))

	langSelect := widget.NewSelect([]string{translator.T("toolbar.lang.en"), translator.T("toolbar.lang.ru")}, func(selected string) {
		if onLanguageChanged != nil {
			language, ok := i18n.ParseAppLanguage(selected)
			if !ok {
				return
			}
			onLanguageChanged(language)
		}
	})
	langSelect.SetSelected(string(i18n.LangEN))

	refreshButton := widget.NewButtonWithIcon(translator.T("toolbar.refresh.tooltip"), theme.ViewRefreshIcon(), func() {
		if onRefreshRequested != nil {
			onRefreshRequested()
		}
	})
	refreshButton.Importance = widget.LowImportance
	refreshButtonWrap := container.NewGridWrap(fyne.NewSize(96, 40), refreshButton)

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
	right := container.NewHBox(currencySelect, langSelect, refreshButtonWrap, themeButtonWrap)
	header := container.NewBorder(nil, canvas.NewLine(theme.Color(theme.ColorNameSeparator)), left, right)

	return &Toolbar{
		root:           header,
		title:          title,
		refreshButton:  refreshButton,
		themeButton:    themeButton,
		themeControl:   themeControl,
		currencySelect: currencySelect,
		langSelect:     langSelect,
		translator:     translator,
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

func (t *Toolbar) RefreshButton() *widget.Button {
	return t.refreshButton
}

func (t *Toolbar) SetLanguage(language i18n.AppLanguage) {
	t.translator.SetLanguage(language)
	t.title.SetText(t.translator.T("app.title"))
	t.refreshButton.SetText(t.translator.T("toolbar.refresh.tooltip"))
	t.langSelect.Options = []string{t.translator.T("toolbar.lang.en"), t.translator.T("toolbar.lang.ru")}
	t.langSelect.Refresh()
}
