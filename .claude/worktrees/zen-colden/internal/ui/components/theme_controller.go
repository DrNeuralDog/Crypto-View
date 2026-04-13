package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

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
		if c.currentVariant() == theme.VariantDark {
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
	if c.currentVariant() == theme.VariantDark {
		return "light"
	}
	return "dark"
}

func (c *ThemeController) ActionIconResource() fyne.Resource {
	if c.currentVariant() == theme.VariantDark {
		return moonIconResource
	}
	return sunIconResource
}

func (c *ThemeController) setMode(mode uitheme.Mode) {
	c.mode = mode
	c.app.Settings().SetTheme(uitheme.NewForMode(mode))
}

func (c *ThemeController) currentVariant() fyne.ThemeVariant {
	switch c.mode {
	case uitheme.ModeDark:
		return theme.VariantDark
	case uitheme.ModeLight:
		return theme.VariantLight
	default:
		return c.app.Settings().ThemeVariant()
	}
}

var sunIconResource = fyne.NewStaticResource("sun.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
  <circle cx="12" cy="12" r="4.5" fill="#fbc02d"/>
  <g stroke="#fbc02d" stroke-width="1.8" stroke-linecap="round">
    <line x1="12" y1="1.8" x2="12" y2="4.2"/>
    <line x1="12" y1="19.8" x2="12" y2="22.2"/>
    <line x1="1.8" y1="12" x2="4.2" y2="12"/>
    <line x1="19.8" y1="12" x2="22.2" y2="12"/>
    <line x1="4.2" y1="4.2" x2="6" y2="6"/>
    <line x1="18" y1="18" x2="19.8" y2="19.8"/>
    <line x1="18" y1="6" x2="19.8" y2="4.2"/>
    <line x1="4.2" y1="19.8" x2="6" y2="18"/>
  </g>
</svg>
`))

var moonIconResource = fyne.NewStaticResource("moon.svg", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
  <path fill="#f5e6a8" d="M14.8 2.2c-4.7.9-8.2 5-8.2 9.8 0 4.5 3 8.4 7.3 9.6C8.5 21.5 4 17 4 11.4 4 6 8.4 1.6 13.8 1.6c.3 0 .7 0 1 .1z"/>
</svg>
`))
