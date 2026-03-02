package ui

import (
	"testing"

	"cryptoview/internal/model"
	"cryptoview/internal/service/marketfeed"
	"cryptoview/internal/ui/i18n"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

type fakeFeed struct {
	callbacks marketfeed.Callbacks
	started   bool
	stopCalls int
	lastFiat  i18n.FiatCurrency
}

func newFakeFeed(callbacks marketfeed.Callbacks) *fakeFeed {
	return &fakeFeed{callbacks: callbacks}
}

func (f *fakeFeed) Start() {
	f.started = true
}

func (f *fakeFeed) Stop() {
	f.stopCalls++
}

func (f *fakeFeed) SetFiat(currency i18n.FiatCurrency) {
	f.lastFiat = currency
}

func (f *fakeFeed) EmitStatus(event marketfeed.StatusEvent) {
	if f.callbacks.OnStatus != nil {
		f.callbacks.OnStatus(event)
	}
}

func (f *fakeFeed) EmitMarketUpdate(coins []model.Coin) {
	if f.callbacks.OnMarketUpdate != nil {
		f.callbacks.OnMarketUpdate(coins)
	}
}

func TestBuildMainWindow_IntegrationLifecycle_NoNetwork(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	var feed *fakeFeed
	w := buildMainWindowWithFeedFactory(a, nil, func(callbacks marketfeed.Callbacks) marketFeed {
		feed = newFakeFeed(callbacks)
		return feed
	})

	if w == nil {
		t.Fatal("expected non-nil window")
	}
	if w.Content() == nil {
		t.Fatal("expected non-nil window content")
	}
	if feed == nil {
		t.Fatal("expected injected fake feed to be created")
	}
	if !feed.started {
		t.Fatal("expected feed.Start to be called during window setup")
	}
	list := findFirstList(w.Content())
	if list == nil {
		t.Fatal("expected coin list widget to be present")
	}
	if list.Length() != 0 {
		t.Fatalf("expected empty coin list on startup without mocks, got %d items", list.Length())
	}

	feed.EmitStatus(marketfeed.StatusEvent{
		Kind:     marketfeed.StatusKindOK,
		Provider: "coingecko",
	})
	feed.EmitStatus(marketfeed.StatusEvent{
		Kind: marketfeed.StatusKindWarning,
		Code: marketfeed.StatusCodeRateLimited,
	})
	feed.EmitMarketUpdate([]model.Coin{
		{
			ID:             "bitcoin",
			Name:           "Bitcoin",
			Ticker:         "BTC",
			Price:          100000,
			Change24h:      2.5,
			LastUpdateTime: "10:11:12",
			IconPath:       model.IconPathForID("bitcoin"),
		},
	})
	fyne.DoAndWait(func() {})
	if list.Length() != 1 {
		t.Fatalf("expected 1 coin item after market update, got %d", list.Length())
	}

	feed.EmitStatus(marketfeed.StatusEvent{
		Kind: marketfeed.StatusKindError,
		Code: marketfeed.StatusCodeNoData,
	})
	fyne.DoAndWait(func() {})
	if !containsCanvasText(w.Content(), "No market data available") {
		t.Fatal("expected EN no-data status message in footer")
	}

	langSelect := findLanguageSelect(w.Content())
	if langSelect == nil {
		t.Fatal("expected language selector to be present")
	}
	langSelect.SetSelected("RU")
	fyne.DoAndWait(func() {})

	feed.EmitStatus(marketfeed.StatusEvent{
		Kind: marketfeed.StatusKindError,
		Code: marketfeed.StatusCodeNoData,
	})
	fyne.DoAndWait(func() {})
	if !containsCanvasText(w.Content(), "Нет данных рынка") {
		t.Fatal("expected RU no-data status message in footer")
	}

	// Window should remain healthy after callback-driven UI updates.
	if w.Content() == nil {
		t.Fatal("expected content to remain valid after feed callbacks")
	}

	w.Close()
	if feed.stopCalls != 1 {
		t.Fatalf("expected feed.Stop to be called once on close, got %d", feed.stopCalls)
	}
}

func findFirstList(obj fyne.CanvasObject) *widget.List {
	switch current := obj.(type) {
	case *widget.List:
		return current
	case *fyne.Container:
		for _, child := range current.Objects {
			if found := findFirstList(child); found != nil {
				return found
			}
		}
	}
	return nil
}

func findLanguageSelect(obj fyne.CanvasObject) *widget.Select {
	switch current := obj.(type) {
	case *widget.Select:
		hasEN := false
		hasRU := false
		for _, option := range current.Options {
			if option == "EN" {
				hasEN = true
			}
			if option == "RU" {
				hasRU = true
			}
		}
		if hasEN && hasRU {
			return current
		}
	case *fyne.Container:
		for _, child := range current.Objects {
			if found := findLanguageSelect(child); found != nil {
				return found
			}
		}
	}
	return nil
}

func containsCanvasText(obj fyne.CanvasObject, expected string) bool {
	switch current := obj.(type) {
	case *canvas.Text:
		return current.Text == expected
	case *fyne.Container:
		for _, child := range current.Objects {
			if containsCanvasText(child, expected) {
				return true
			}
		}
	}
	return false
}

