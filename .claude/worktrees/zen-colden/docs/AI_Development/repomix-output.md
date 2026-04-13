# CryptoView - Packed Codebase (AI-friendly)

> Generated for NotebookLM context. Project: CryptoView | Tech: Go 1.22, Fyne v2, CoinGecko API

---

## go.mod

```go
module cryptoview

go 1.22

require fyne.io/fyne/v2 v2.7.2

require (
	fyne.io/systray v1.12.0 // indirect
	github.com/BurntSushi/toml v1.5.0 // indirect
	// ... (see full go.mod for indirect deps)
)
```

---

## Makefile

```makefile
.PHONY: build run clean

build:
	mkdir -p bin
	go build -o bin/cryptoview ./cmd/cryptoview

run:
	go run ./cmd/cryptoview

clean:
	rm -rf bin
```

---

## cmd/cryptoview/main.go

```go
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
```

---

## internal/model/coin.go

```go
package model

type Coin struct {
	ID             string
	Name           string
	Ticker         string
	Price          float64
	Change24h      float64
	LastUpdateTime string
	IconPath       string
}

func GetMockCoins() []Coin {
	return []Coin{
		{ID: "bitcoin", Name: "Bitcoin", Ticker: "BTC", Price: 96543.12, Change24h: 2.54, ...},
		{ID: "ethereum", Name: "Ethereum", Ticker: "ETH", Price: 3421.77, Change24h: -1.23, ...},
		{ID: "toncoin", Name: "TON Coin", Ticker: "TON", Price: 5.89, Change24h: 0.00, ...},
		{ID: "solana", Name: "Solana", Ticker: "SOL", Price: 183.45, Change24h: 5.91, ...},
		{ID: "dogecoin", Name: "Dogecoin", Ticker: "DOGE", Price: 0.25, Change24h: -3.02, ...},
		{ID: "ripple", Name: "Ripple", Ticker: "XRP", Price: 0.71, Change24h: 1.04, ...},
		{ID: "litecoin", Name: "Litecoin", Ticker: "LTC", Price: 102.33, Change24h: -0.67, ...},
	}
}
```

---

## internal/model/coin_test.go

```go
package model

import "testing"

func TestGetMockCoins(t *testing.T) {
	coins := GetMockCoins()
	if len(coins) != 7 {
		t.Fatalf("expected 7 coins, got %d", len(coins))
	}
	for i, coin := range coins {
		if coin.Ticker == "" { t.Fatalf("coin[%d] has empty ticker", i) }
		if coin.Name == "" { t.Fatalf("coin[%d] has empty name", i) }
		if coin.Price <= 0 { t.Fatalf("coin[%d] has non-positive price: %f", i, coin.Price) }
	}
}
```

---

## internal/ui/main_window.go

```go
package ui

import (
	"cryptoview/internal/model"
	"cryptoview/internal/ui/components"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func BuildMainWindow(a fyne.App, data []model.Coin) fyne.Window {
	w := a.NewWindow("CryptoView")
	w.Resize(fyne.NewSize(900, 600))
	header := components.NewToolbar()
	coinList := components.NewCoinList(data)
	footer := NewFooter()
	content := container.NewBorder(header, footer, nil, nil, coinList)
	w.SetContent(content)
	return w
}
```

---

## internal/ui/footer.go

```go
package ui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewFooter() fyne.CanvasObject {
	return container.NewHBox(widget.NewLabel("Status: OK"))
}
```

---

## internal/ui/components/toolbar.go

```go
package components

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func NewToolbar() fyne.CanvasObject {
	currencySelect := widget.NewSelect([]string{"USD", "EUR", "RUB"}, func(string) {})
	currencySelect.SetSelected("USD")
	langButton := widget.NewButton("Lang", func() {})
	themeButton := widget.NewButton("Theme", func() {})
	return container.NewHBox(
		widget.NewLabel("CryptoView"),
		layout.NewSpacer(),
		currencySelect,
		langButton,
		themeButton,
	)
}
```

---

## internal/ui/components/coin_list.go

```go
package components

import (
	"fmt"
	"image/color"
	"cryptoview/internal/model"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func NewCoinList(data []model.Coin) *widget.List {
	return widget.NewList(
		func() int { return len(data) },
		func() fyne.CanvasObject {
			icon := widget.NewLabel("[icon]")
			nameTicker := widget.NewLabel("")
			price := widget.NewLabel("")
			change := canvas.NewText("", color.NRGBA{R: 128, G: 128, B: 128, A: 255})
			lastUpdate := widget.NewLabel("")
			return container.NewHBox(icon, nameTicker, layout.NewSpacer(), price, change, lastUpdate)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			row := item.(*fyne.Container)
			coin := data[id]
			icon := row.Objects[0].(*widget.Label)
			nameTicker := row.Objects[1].(*widget.Label)
			price := row.Objects[3].(*widget.Label)
			change := row.Objects[4].(*canvas.Text)
			lastUpdate := row.Objects[5].(*widget.Label)
			icon.SetText(fmt.Sprintf("[%s]", coin.Ticker))
			nameTicker.SetText(fmt.Sprintf("%s | %s", coin.Name, coin.Ticker))
			price.SetText(fmt.Sprintf("$%.2f", coin.Price))
			change.Text = fmt.Sprintf("%.2f%%", coin.Change24h)
			change.Color = changeColor(coin.Change24h)
			change.Refresh()
			lastUpdate.SetText(coin.LastUpdateTime)
		},
	)
}

func changeColor(change24h float64) color.Color {
	if change24h > 0 { return color.NRGBA{R: 34, G: 139, B: 34, A: 255} }
	if change24h < 0 { return color.NRGBA{R: 220, G: 20, B: 60, A: 255} }
	return color.NRGBA{R: 128, G: 128, B: 128, A: 255}
}
```

---

## internal/ui/components/coin_list_test.go

```go
package components

import (
	"testing"
	"cryptoview/internal/model"
	"fyne.io/fyne/v2/test"
)

func TestNewCoinList(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()
	data := model.GetMockCoins()
	list := NewCoinList(data)
	if list == nil { t.Fatal("expected list to be non-nil") }
	if got := list.Length(); got != len(data) {
		t.Fatalf("expected list length %d, got %d", len(data), got)
	}
	item := list.CreateItem()
	if item == nil { t.Fatal("expected list item to be non-nil") }
	list.UpdateItem(0, item)
}
```

---

## docs/PRD.md (Summary)

**Project:** CryptoView — Desktop app for crypto monitoring.  
**Stack:** Go 1.22+, Fyne, CoinGecko API.  
**Features:** Scrollable list (BTC, ETH, TON, SOL, DOGE, XRP, LTC), currency selector (USD/EUR/RUB), theme/lang switch, 60s auto-refresh, network error handling.  
**Layout:** `layout.Border` — header (toolbar), center (widget.List), footer (status).  
**Architecture:** cmd/internal/model, internal/api, internal/ui, resources.

---

## docs/AI_Development/CryptoView_Implementation.md (Summary)

**Current:** Stage 2 — Core Features (MVP).  
**Done:** Go module, dir structure, Fyne, minimal window, Makefile, widget.List with mock data.  
**Todo:** JSON structs, HTTP client, CoinGecko provider, connect real API, currency selector.  
**Next:** Define JSON structs for API response (model/coin.go).

---

## docs/AI_Development/project_structure.md (Summary)

**Structure:** cmd/cryptoview, internal/model, internal/api, internal/service, internal/ui, pkg, resources, docs.  
**Build:** `go build ./cmd/cryptoview`, `make build`, `make run`.
