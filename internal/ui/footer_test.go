package ui

import (
	"image/color"
	"testing"

	"cryptoview/internal/ui/i18n"
	"fyne.io/fyne/v2/test"
)

func TestFooterControllerStates(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	footer := NewFooterController(i18n.NewTranslator(i18n.LangEN))

	footer.SetLoading()
	if footer.statusValue.Text != "Loading..." {
		t.Fatalf("expected loading text, got %q", footer.statusValue.Text)
	}
	if !footer.progress.Visible() {
		t.Fatal("expected progress visible in loading state")
	}

	footer.SetError("Network error")
	if footer.statusValue.Text != "Network error" {
		t.Fatalf("expected error text, got %q", footer.statusValue.Text)
	}
	if footer.progress.Visible() {
		t.Fatal("expected progress hidden in error state")
	}
	if got := asNRGBA(footer.statusValue.Color); got != (color.NRGBA{R: 0xF4, G: 0x43, B: 0x36, A: 0xFF}) {
		t.Fatalf("expected red error color, got %+v", got)
	}

	footer.SetOK()
	if footer.statusValue.Text != "OK" {
		t.Fatalf("expected OK text, got %q", footer.statusValue.Text)
	}
	if footer.progress.Visible() {
		t.Fatal("expected progress hidden in OK state")
	}
}

func TestFooterLanguageSwitch(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	footer := NewFooterController(i18n.NewTranslator(i18n.LangEN))
	footer.SetLanguage(i18n.LangRU)
	footer.SetLoading()

	if footer.statusLabel.Text != "Статус:" {
		t.Fatalf("expected localized status label, got %q", footer.statusLabel.Text)
	}
	if footer.statusValue.Text != "Загрузка..." {
		t.Fatalf("expected localized loading text, got %q", footer.statusValue.Text)
	}
}

func asNRGBA(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}
