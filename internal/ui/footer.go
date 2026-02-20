package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewFooter() fyne.CanvasObject {
	indicator := canvas.NewCircle(color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF})
	indicatorWrap := container.NewGridWrap(fyne.NewSize(8, 8), indicator)

	statusLabel := widget.NewLabel("Status:")
	okText := canvas.NewText("OK", color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF})
	okText.TextSize = 16

	row := container.NewHBox(
		container.NewCenter(indicatorWrap),
		container.NewCenter(statusLabel),
		container.NewCenter(okText),
	)
	return container.NewPadded(row)
}
