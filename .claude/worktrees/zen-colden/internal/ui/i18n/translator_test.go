package i18n

import "testing"

func TestTranslatorT(t *testing.T) {
	tr := NewTranslator(LangEN)
	if got := tr.T("status.loading"); got != "Loading..." {
		t.Fatalf("expected Loading..., got %q", got)
	}

	tr.SetLanguage(LangRU)
	if got := tr.T("status.loading"); got != "Загрузка..." {
		t.Fatalf("expected Загрузка..., got %q", got)
	}
}

func TestTranslatorFallbackForUnknownKey(t *testing.T) {
	tr := NewTranslator(LangEN)
	key := "missing.key"
	if got := tr.T(key); got != key {
		t.Fatalf("expected fallback to key %q, got %q", key, got)
	}
}

func TestTranslatorNoDataStatus(t *testing.T) {
	tr := NewTranslator(LangEN)
	if got := tr.T("status.error.no_data"); got != "No market data available" {
		t.Fatalf("expected EN no-data text, got %q", got)
	}

	tr.SetLanguage(LangRU)
	if got := tr.T("status.error.no_data"); got != "Нет данных рынка" {
		t.Fatalf("expected RU no-data text, got %q", got)
	}
}
