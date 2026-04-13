package i18n

import "strings"

type AppLanguage string

const (
	LangEN AppLanguage = "EN"
	LangRU AppLanguage = "RU"
)

func ParseAppLanguage(raw string) (AppLanguage, bool) {
	switch AppLanguage(strings.ToUpper(strings.TrimSpace(raw))) {
	case LangEN:
		return LangEN, true
	case LangRU:
		return LangRU, true
	default:
		return "", false
	}
}

type FiatCurrency string

const (
	FiatUSD FiatCurrency = "USD"
	FiatEUR FiatCurrency = "EUR"
	FiatRUB FiatCurrency = "RUB"
)

func ParseFiatCurrency(raw string) (FiatCurrency, bool) {
	switch FiatCurrency(strings.ToUpper(strings.TrimSpace(raw))) {
	case FiatUSD:
		return FiatUSD, true
	case FiatEUR:
		return FiatEUR, true
	case FiatRUB:
		return FiatRUB, true
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
