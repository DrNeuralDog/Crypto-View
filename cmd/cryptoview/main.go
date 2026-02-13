package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("CryptoView")
	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}
