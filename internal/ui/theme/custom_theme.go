package uitheme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Mode string

const (
	ModeSystem Mode = "system"
	ModeDark   Mode = "dark"
	ModeLight  Mode = "light"
)

type CustomTheme struct {
	base   fyne.Theme
	mode   Mode
	forced fyne.ThemeVariant
}

func NewForMode(mode Mode) fyne.Theme {
	t := &CustomTheme{
		base: theme.DefaultTheme(),
		mode: mode,
	}
	switch mode {
	case ModeDark:
		t.forced = theme.VariantDark
	case ModeLight:
		t.forced = theme.VariantLight
	}
	return t
}

func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	activeVariant := variant
	if t.mode == ModeDark || t.mode == ModeLight {
		activeVariant = t.forced
	}

	switch activeVariant {
	case theme.VariantDark:
		if c, ok := darkPalette(name); ok {
			return c
		}
	default:
		if c, ok := lightPalette(name); ok {
			return c
		}
	}

	return t.base.Color(name, activeVariant)
}

func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.base.Font(style)
}

func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.base.Icon(name)
}

func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.base.Size(name)
}

func darkPalette(name fyne.ThemeColorName) (color.Color, bool) {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0x1E, G: 0x1E, B: 0x1E, A: 0xFF}, true
	case theme.ColorNameInputBackground, theme.ColorNameButton, theme.ColorNameHeaderBackground:
		return color.NRGBA{R: 0x26, G: 0x27, B: 0x31, A: 0xFF}, true
	case theme.ColorNameMenuBackground:
		return color.NRGBA{R: 0x20, G: 0x21, B: 0x2B, A: 0xFF}, true
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0xE8, G: 0xEA, B: 0xF1, A: 0xFF}, true
	case theme.ColorNamePlaceHolder, theme.ColorNameDisabled:
		return color.NRGBA{R: 0x8E, G: 0x94, B: 0xA5, A: 0xFF}, true
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 0x33, G: 0x36, B: 0x45, A: 0xFF}, true
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 0x50, G: 0x53, B: 0x61, A: 0xFF}, true
	case theme.ColorNameScrollBarBackground:
		return color.NRGBA{R: 0x2B, G: 0x2D, B: 0x39, A: 0xFF}, true
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}, true
	case theme.ColorNameError:
		return color.NRGBA{R: 0xF4, G: 0x43, B: 0x36, A: 0xFF}, true
	}
	return nil, false
}

func lightPalette(name fyne.ThemeColorName) (color.Color, bool) {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0xF5, G: 0xF5, B: 0xF7, A: 0xFF}, true
	case theme.ColorNameInputBackground, theme.ColorNameButton, theme.ColorNameHeaderBackground:
		return color.NRGBA{R: 0xEA, G: 0xEB, B: 0xEF, A: 0xFF}, true
	case theme.ColorNameMenuBackground:
		return color.NRGBA{R: 0xF8, G: 0xF8, B: 0xFA, A: 0xFF}, true
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0x2D, G: 0x31, B: 0x40, A: 0xFF}, true
	case theme.ColorNamePlaceHolder, theme.ColorNameDisabled:
		return color.NRGBA{R: 0x8E, G: 0x94, B: 0xA5, A: 0xFF}, true
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 0xDB, G: 0xDE, B: 0xE7, A: 0xFF}, true
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 0xC1, G: 0xC6, B: 0xD4, A: 0xFF}, true
	case theme.ColorNameScrollBarBackground:
		return color.NRGBA{R: 0xE7, G: 0xE9, B: 0xF0, A: 0xFF}, true
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}, true
	case theme.ColorNameError:
		return color.NRGBA{R: 0xF4, G: 0x43, B: 0x36, A: 0xFF}, true
	}
	return nil, false
}
