package i18n

import "testing"

func TestFormatPriceEN(t *testing.T) {
	got := FormatPrice(12345.67, FiatUSD, LangEN)
	if got != "$12,345.67" {
		t.Fatalf("expected $12,345.67, got %q", got)
	}
}

func TestFormatPriceRU(t *testing.T) {
	got := FormatPrice(12345.67, FiatRUB, LangRU)
	if got != "12 345,67 ₽" {
		t.Fatalf("expected 12 345,67 ₽, got %q", got)
	}
}

func TestFormatTime(t *testing.T) {
	if got := FormatTime("12:34:56", LangEN); got != "12:34:56" {
		t.Fatalf("expected valid time to remain unchanged, got %q", got)
	}
	if got := FormatTime("invalid", LangRU); got != "--:--:--" {
		t.Fatalf("expected invalid time fallback, got %q", got)
	}
}
