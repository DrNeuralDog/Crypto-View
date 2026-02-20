package main

import (
	"cryptoview/internal/model"
	"cryptoview/internal/ui"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	data := model.GetMockCoins()
	w := ui.BuildMainWindow(a, data)
	w.ShowAndRun()
}
