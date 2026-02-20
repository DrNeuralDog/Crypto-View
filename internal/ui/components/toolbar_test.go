package components

import (
	"testing"

	"cryptoview/internal/ui/i18n"
	"fyne.io/fyne/v2/test"
)

func TestToolbarCurrencyAndLanguageCallbacks(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	translator := i18n.NewTranslator(i18n.LangEN)
	var gotCurrency i18n.FiatCurrency
	var gotLanguage i18n.AppLanguage
	refreshClicks := 0

	toolbar := NewToolbar(a, translator, func(currency i18n.FiatCurrency) {
		gotCurrency = currency
	}, nil, func(language i18n.AppLanguage) {
		gotLanguage = language
	}, func() {
		refreshClicks++
	})

	toolbar.CurrencySelect().SetSelected(string(i18n.FiatEUR))
	if gotCurrency != i18n.FiatEUR {
		t.Fatalf("expected currency callback %q, got %q", i18n.FiatEUR, gotCurrency)
	}

	toolbar.LanguageSelect().SetSelected("RU")
	if gotLanguage != i18n.LangRU {
		t.Fatalf("expected language callback RU, got %q", gotLanguage)
	}

	test.Tap(toolbar.RefreshButton())
	if refreshClicks != 1 {
		t.Fatalf("expected refresh callback to be called once, got %d", refreshClicks)
	}
}

func TestToolbarThemeButtonTogglesTheme(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	toolbar := NewToolbar(a, i18n.NewTranslator(i18n.LangEN), nil, nil, nil, nil)
	before := a.Settings().Theme()
	beforeIcon := toolbar.ThemeButton().Icon

	test.Tap(toolbar.ThemeButton())

	after := a.Settings().Theme()
	afterIcon := toolbar.ThemeButton().Icon

	if before == after {
		t.Fatal("expected theme to change after tapping theme button")
	}
	if beforeIcon == afterIcon {
		t.Fatal("expected theme action icon resource to change after toggle")
	}
	if toolbar.ThemeMode() == "system" {
		t.Fatal("expected theme mode to leave system after first toggle")
	}
}
