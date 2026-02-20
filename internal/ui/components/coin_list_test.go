package components

import (
	"image/color"
	"strings"
	"testing"

	"cryptoview/internal/model"
	"cryptoview/internal/ui/i18n"
	"fyne.io/fyne/v2/test"
)

func TestNewCoinList(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	data := model.GetMockCoins()
	list := NewCoinList(data, i18n.NewTranslator(i18n.LangEN))
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

	controller := NewCoinList(model.GetMockCoins(), i18n.NewTranslator(i18n.LangEN))
	item := controller.Widget().CreateItem()
	row := item.(*coinListItem)

	controller.Widget().UpdateItem(0, item)
	if !strings.HasPrefix(row.price.Text, "$") {
		t.Fatalf("expected USD symbol, got %q", row.price.Text)
	}

	controller.SetCurrency(i18n.FiatEUR)
	controller.Widget().UpdateItem(0, item)
	if !strings.HasPrefix(row.price.Text, "\u20ac") {
		t.Fatalf("expected EUR symbol, got %q", row.price.Text)
	}

	controller.SetCurrency(i18n.FiatRUB)
	controller.Widget().UpdateItem(0, item)
	if !strings.HasPrefix(row.price.Text, "\u20bd") {
		t.Fatalf("expected RUB symbol, got %q", row.price.Text)
	}

	controller.SetLanguage(i18n.LangRU)
	controller.SetCurrency(i18n.FiatUSD)
	controller.Widget().UpdateItem(0, item)
	if !strings.Contains(row.price.Text, "$") || !strings.Contains(row.price.Text, ",") {
		t.Fatalf("expected RU locale price with comma decimal separator, got %q", row.price.Text)
	}
}

func TestCoinListChangeColor(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	controller := NewCoinList(model.GetMockCoins(), i18n.NewTranslator(i18n.LangEN))
	item := controller.Widget().CreateItem()
	row := item.(*coinListItem)

	controller.Widget().UpdateItem(0, item)
	if got := asNRGBA(row.change.Color); got != (color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}) {
		t.Fatalf("expected positive change color green, got %+v", got)
	}

	controller.Widget().UpdateItem(1, item)
	if got := asNRGBA(row.change.Color); got != (color.NRGBA{R: 0xF4, G: 0x43, B: 0x36, A: 0xFF}) {
		t.Fatalf("expected negative change color red, got %+v", got)
	}
}

func TestCoinListLastItemSeparatorHidden(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	controller := NewCoinList(model.GetMockCoins(), i18n.NewTranslator(i18n.LangEN))
	item := controller.Widget().CreateItem()
	row := item.(*coinListItem)

	last := controller.Widget().Length() - 1
	controller.Widget().UpdateItem(last, item)
	if row.separator.Visible() {
		t.Fatal("expected last item separator to be hidden")
	}
}

func TestReplaceData(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	controller := NewCoinList(model.GetMockCoins(), i18n.NewTranslator(i18n.LangEN))
	newData := []model.Coin{
		{ID: "bitcoin", Name: "Bitcoin", Ticker: "BTC", Price: 1, Change24h: 0.1, LastUpdateTime: "10:11:12"},
	}

	controller.ReplaceData(newData)
	controller.SetCurrency(i18n.FiatUSD)

	if got := controller.Widget().Length(); got != 1 {
		t.Fatalf("expected list length 1 after replace, got %d", got)
	}
}

func asNRGBA(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}
