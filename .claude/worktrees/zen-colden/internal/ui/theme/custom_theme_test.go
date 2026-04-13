package uitheme

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2/theme"
)

func TestCustomThemeDarkBackgroundColor(t *testing.T) {
	tm := NewForMode(ModeDark)
	got := toNRGBA(tm.Color(theme.ColorNameBackground, theme.VariantLight))
	want := color.NRGBA{R: 0x1E, G: 0x1E, B: 0x1E, A: 0xFF}
	if got != want {
		t.Fatalf("expected dark background %+v, got %+v", want, got)
	}
}

func TestCustomThemeLightBackgroundColor(t *testing.T) {
	tm := NewForMode(ModeLight)
	got := toNRGBA(tm.Color(theme.ColorNameBackground, theme.VariantDark))
	want := color.NRGBA{R: 0xF5, G: 0xF5, B: 0xF7, A: 0xFF}
	if got != want {
		t.Fatalf("expected light background %+v, got %+v", want, got)
	}
}

func TestCustomThemeSystemRespectsVariant(t *testing.T) {
	tm := NewForMode(ModeSystem)
	dark := toNRGBA(tm.Color(theme.ColorNameBackground, theme.VariantDark))
	light := toNRGBA(tm.Color(theme.ColorNameBackground, theme.VariantLight))
	if dark == light {
		t.Fatal("expected system mode to respect passed variant")
	}
}

func toNRGBA(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}
