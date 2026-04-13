package i18n

import (
	"testing"
	"time"
)

func TestFormatPriceEN(t *testing.T) {
	got := FormatPrice(12345.67, FiatUSD, LangEN)
	if got != "$12,345.67" {
		t.Fatalf("expected $12,345.67, got %q", got)
	}
}

func TestFormatPriceRU(t *testing.T) {
	got := FormatPrice(12345.67, FiatRUB, LangRU)
	if got != "12 345,67 \u20bd" {
		t.Fatalf("expected 12 345,67 ₽, got %q", got)
	}
}

func TestFormatTime(t *testing.T) {
	value := time.Date(2026, time.February, 20, 12, 34, 56, 0, time.UTC)
	if got := FormatTime(value, LangEN); got != value.Local().Format("15:04:05") {
		t.Fatalf("expected valid time to be formatted in local timezone, got %q", got)
	}
	if got := FormatTime(time.Time{}, LangRU); got != "--:--:--" {
		t.Fatalf("expected zero time fallback, got %q", got)
	}
}
