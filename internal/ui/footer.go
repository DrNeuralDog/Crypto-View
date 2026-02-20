package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewFooter() fyne.CanvasObject {
	return container.NewHBox(widget.NewLabel("Status: OK"))
}
