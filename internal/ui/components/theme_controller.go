package components

import (
	"fyne.io/fyne/v2"
	fynetheme "fyne.io/fyne/v2/theme"

	uitheme "cryptoview/internal/ui/theme"
)

type ThemeController struct {
	app  fyne.App
	mode uitheme.Mode
}

func NewThemeController(app fyne.App) *ThemeController {
	return &ThemeController{
		app:  app,
		mode: uitheme.ModeSystem,
	}
}

func (c *ThemeController) Mode() uitheme.Mode {
	return c.mode
}

func (c *ThemeController) Toggle() {
	switch c.mode {
	case uitheme.ModeSystem:
		if c.currentVariant() == fynetheme.VariantDark {
			c.setMode(uitheme.ModeLight)
		} else {
			c.setMode(uitheme.ModeDark)
		}
	case uitheme.ModeDark:
		c.setMode(uitheme.ModeLight)
	default:
		c.setMode(uitheme.ModeDark)
	}
}

func (c *ThemeController) ActionIconText() string {
	if c.currentVariant() == fynetheme.VariantDark {
		return "☀"
	}
	return "🌙"
}

func (c *ThemeController) setMode(mode uitheme.Mode) {
	c.mode = mode
	c.app.Settings().SetTheme(uitheme.NewForMode(mode))
}

func (c *ThemeController) currentVariant() fyne.ThemeVariant {
	switch c.mode {
	case uitheme.ModeDark:
		return fynetheme.VariantDark
	case uitheme.ModeLight:
		return fynetheme.VariantLight
	default:
		return c.app.Settings().ThemeVariant()
	}
}
