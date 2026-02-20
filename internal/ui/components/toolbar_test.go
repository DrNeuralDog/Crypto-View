package components

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestToolbarCurrencyAndLanguageCallbacks(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	var gotCurrency FiatCurrency
	var gotLanguage string

	toolbar := NewToolbar(a, func(currency FiatCurrency) {
		gotCurrency = currency
	}, nil, func(language string) {
		gotLanguage = language
	})

	toolbar.CurrencySelect().SetSelected(string(FiatEUR))
	if gotCurrency != FiatEUR {
		t.Fatalf("expected currency callback %q, got %q", FiatEUR, gotCurrency)
	}

	toolbar.LanguageSelect().SetSelected("ENG")
	if gotLanguage != "ENG" {
		t.Fatalf("expected language callback ENG, got %q", gotLanguage)
	}
}

func TestToolbarThemeButtonTogglesTheme(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	toolbar := NewToolbar(a, nil, nil, nil)
	before := a.Settings().Theme()
	beforeLabel := toolbar.ThemeButton().Text

	test.Tap(toolbar.ThemeButton())

	after := a.Settings().Theme()
	afterLabel := toolbar.ThemeButton().Text

	if before == after {
		t.Fatal("expected theme to change after tapping theme button")
	}
	if beforeLabel == afterLabel {
		t.Fatal("expected theme action icon text to change after toggle")
	}
	if toolbar.ThemeMode() == "system" {
		t.Fatal("expected theme mode to leave system after first toggle")
	}
}
