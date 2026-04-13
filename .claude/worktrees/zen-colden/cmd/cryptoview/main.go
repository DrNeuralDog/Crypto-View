package main

import (
	"cryptoview/internal/ui"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := ui.BuildMainWindow(a, nil)
	w.ShowAndRun()
}
