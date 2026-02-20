package components

import (
	"strings"
	"testing"

	"cryptoview/internal/model"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestNewCoinList(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	data := model.GetMockCoins()
	list := NewCoinList(data)
	if list == nil {
		t.Fatal("expected list to be non-nil")
	}

	if got := list.Widget().Length(); got != len(data) {
		t.Fatalf("expected list length %d, got %d", len(data), got)
	}

	item := list.Widget().CreateItem()
	if item == nil {
		t.Fatal("expected list item to be non-nil")
	}

	list.Widget().UpdateItem(0, item)
}

func TestCoinListCurrencySwitch(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	data := model.GetMockCoins()
	controller := NewCoinList(data)
	item := controller.Widget().CreateItem()

	controller.Widget().UpdateItem(0, item)
	priceLabel := item.(*fyne.Container).Objects[3].(*widget.Label)
	if !strings.HasPrefix(priceLabel.Text, "$") {
		t.Fatalf("expected USD symbol, got %q", priceLabel.Text)
	}

	controller.SetCurrency(FiatEUR)
	controller.Widget().UpdateItem(0, item)
	if !strings.HasPrefix(priceLabel.Text, "€") {
		t.Fatalf("expected EUR symbol, got %q", priceLabel.Text)
	}

	controller.SetCurrency(FiatRUB)
	controller.Widget().UpdateItem(0, item)
	if !strings.HasPrefix(priceLabel.Text, "₽") {
		t.Fatalf("expected RUB symbol, got %q", priceLabel.Text)
	}
}

func TestReplaceData(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	controller := NewCoinList(model.GetMockCoins())
	newData := []model.Coin{
		{ID: "bitcoin", Name: "Bitcoin", Ticker: "BTC", Price: 1, Change24h: 0.1, LastUpdateTime: "10:11:12"},
	}

	controller.ReplaceData(newData)
	controller.SetCurrency(FiatUSD)

	if got := controller.Widget().Length(); got != 1 {
		t.Fatalf("expected list length 1 after replace, got %d", got)
	}
}
