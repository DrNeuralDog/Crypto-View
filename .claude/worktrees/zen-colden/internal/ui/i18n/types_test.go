package i18n

import "testing"

func TestParseAppLanguage(t *testing.T) {
	tests := []struct {
		raw    string
		want   AppLanguage
		ok     bool
	}{
		{"EN", LangEN, true},
		{"en", LangEN, true},
		{"  EN  ", LangEN, true},
		{"RU", LangRU, true},
		{"ru", LangRU, true},
		{"  RU  ", LangRU, true},
		{"", "", false},
		{"DE", "", false},
		{"xx", "", false},
	}
	for _, tt := range tests {
		got, ok := ParseAppLanguage(tt.raw)
		if ok != tt.ok || got != tt.want {
			t.Errorf("ParseAppLanguage(%q) = (%q, %v), want (%q, %v)", tt.raw, got, ok, tt.want, tt.ok)
		}
	}
}

func TestParseFiatCurrency(t *testing.T) {
	tests := []struct {
		raw  string
		want FiatCurrency
		ok   bool
	}{
		{"USD", FiatUSD, true},
		{"usd", FiatUSD, true},
		{"  USD  ", FiatUSD, true},
		{"EUR", FiatEUR, true},
		{"eur", FiatEUR, true},
		{"RUB", FiatRUB, true},
		{"rub", FiatRUB, true},
		{"", "", false},
		{"GBP", "", false},
		{"xxx", "", false},
	}
	for _, tt := range tests {
		got, ok := ParseFiatCurrency(tt.raw)
		if ok != tt.ok || got != tt.want {
			t.Errorf("ParseFiatCurrency(%q) = (%q, %v), want (%q, %v)", tt.raw, got, ok, tt.want, tt.ok)
		}
	}
}

func TestFiatCurrency_APIValue(t *testing.T) {
	tests := []struct {
		fiat FiatCurrency
		want string
		ok   bool
	}{
		{FiatUSD, "usd", true},
		{FiatEUR, "eur", true},
		{FiatRUB, "rub", true},
		{FiatCurrency("GBP"), "", false},
		{FiatCurrency(""), "", false},
	}
	for _, tt := range tests {
		got, ok := tt.fiat.APIValue()
		if ok != tt.ok || got != tt.want {
			t.Errorf("FiatCurrency(%q).APIValue() = (%q, %v), want (%q, %v)", tt.fiat, got, ok, tt.want, tt.ok)
		}
	}
}
