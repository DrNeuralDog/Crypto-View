package marketfeed

import (
	"context"
	"errors"
	"testing"
	"time"

	"cryptoview/internal/model"
	"cryptoview/internal/ui/i18n"
)

type fakeMarketProvider struct {
	name      string
	calls     int
	fetchFunc func(context.Context) (MarketSnapshot, error)
}

func (p *fakeMarketProvider) Name() string { return p.name }

func (p *fakeMarketProvider) FetchUSD(ctx context.Context) (MarketSnapshot, error) {
	p.calls++
	if p.fetchFunc == nil {
		return MarketSnapshot{}, nil
	}
	return p.fetchFunc(ctx)
}

type fakeFXProvider struct {
	calls     int
	fetchFunc func(context.Context) (FXSnapshot, error)
}

func (p *fakeFXProvider) Name() string { return "fakefx" }

func (p *fakeFXProvider) FetchRates(ctx context.Context) (FXSnapshot, error) {
	p.calls++
	if p.fetchFunc == nil {
		return FXSnapshot{}, nil
	}
	return p.fetchFunc(ctx)
}

func TestFeedFallbackOn429UsesNextProvider(t *testing.T) {
	p1 := &fakeMarketProvider{
		name: "cg",
		fetchFunc: func(context.Context) (MarketSnapshot, error) {
			return MarketSnapshot{}, &ProviderError{Provider: "cg", Kind: FailureKindRateLimit, StatusCode: 429}
		},
	}
	p2 := &fakeMarketProvider{
		name: "coincap",
		fetchFunc: func(context.Context) (MarketSnapshot, error) {
			return snapshotWithBTC("coincap", 200), nil
		},
	}
	fx := &fakeFXProvider{
		fetchFunc: func(context.Context) (FXSnapshot, error) {
			return FXSnapshot{
				Base: "USD",
				Rates: map[i18n.FiatCurrency]float64{
					i18n.FiatUSD: 1,
					i18n.FiatRUB: 90,
				},
			}, nil
		},
	}
	var gotCoins []model.Coin
	var gotStatus StatusEvent
	feed := New([]MarketProvider{p1, p2}, fx, Callbacks{
		OnMarketUpdate: func(coins []model.Coin) { gotCoins = coins },
		OnStatus:       func(event StatusEvent) { gotStatus = event },
	})

	feed.runFXCycle()
	feed.runMarketCycle()

	if p1.calls != 1 || p2.calls != 1 {
		t.Fatalf("expected fallback chain calls 1/1, got %d/%d", p1.calls, p2.calls)
	}
	if len(gotCoins) == 0 {
		t.Fatal("expected market update from fallback provider")
	}
	if gotStatus.Kind != StatusKindWarning || gotStatus.Code != StatusCodeFallback {
		t.Fatalf("expected fallback warning status, got %+v", gotStatus)
	}
}

func TestFeedRateLimitCooldownSkipsProvider(t *testing.T) {
	p1 := &fakeMarketProvider{
		name: "cg",
		fetchFunc: func(context.Context) (MarketSnapshot, error) {
			return MarketSnapshot{}, &ProviderError{Provider: "cg", Kind: FailureKindRateLimit, StatusCode: 429}
		},
	}
	p2 := &fakeMarketProvider{
		name: "coincap",
		fetchFunc: func(context.Context) (MarketSnapshot, error) {
			return snapshotWithBTC("coincap", 123), nil
		},
	}
	fx := &fakeFXProvider{fetchFunc: func(context.Context) (FXSnapshot, error) {
		return FXSnapshot{Base: "USD", Rates: map[i18n.FiatCurrency]float64{i18n.FiatUSD: 1}}, nil
	}}
	feed := New([]MarketProvider{p1, p2}, fx, Callbacks{})

	feed.runMarketCycle()
	feed.runMarketCycle()

	if p1.calls != 1 {
		t.Fatalf("expected rate-limited provider to be skipped on second cycle, got %d calls", p1.calls)
	}
	if p2.calls != 2 {
		t.Fatalf("expected fallback provider to continue serving, got %d calls", p2.calls)
	}
}

func TestFeedOfflineFiatRecalculationUsesCachedFX(t *testing.T) {
	p1 := &fakeMarketProvider{
		name: "cg",
		fetchFunc: func(context.Context) (MarketSnapshot, error) {
			return snapshotWithBTC("cg", 100), nil
		},
	}
	fx := &fakeFXProvider{
		fetchFunc: func(context.Context) (FXSnapshot, error) {
			return FXSnapshot{
				Base: "USD",
				Rates: map[i18n.FiatCurrency]float64{
					i18n.FiatUSD: 1,
					i18n.FiatEUR: 0.9,
					i18n.FiatRUB: 90,
				},
			}, nil
		},
	}
	var updates [][]model.Coin
	feed := New([]MarketProvider{p1}, fx, Callbacks{
		OnMarketUpdate: func(coins []model.Coin) {
			cp := make([]model.Coin, len(coins))
			copy(cp, coins)
			updates = append(updates, cp)
		},
	})

	feed.runFXCycle()
	feed.runMarketCycle()
	if len(updates) == 0 {
		t.Fatal("expected initial update")
	}
	usdPrice := firstBTCPrice(t, updates[len(updates)-1])

	p1.fetchFunc = func(context.Context) (MarketSnapshot, error) {
		return MarketSnapshot{}, &ProviderError{Provider: "cg", Kind: FailureKindNetwork, Err: errors.New("offline")}
	}
	beforeCalls := p1.calls
	feed.SetFiat(i18n.FiatRUB)
	if p1.calls != beforeCalls {
		t.Fatalf("expected no HTTP/provider call on fiat switch, got %d -> %d", beforeCalls, p1.calls)
	}
	if len(updates) < 2 {
		t.Fatal("expected local recalculation update on fiat switch")
	}
	rubPrice := firstBTCPrice(t, updates[len(updates)-1])
	if rubPrice != usdPrice*90 {
		t.Fatalf("expected offline RUB recalc %.2f, got %.2f", usdPrice*90, rubPrice)
	}
}

func TestFeedUsesCachedDataWarningWhenAllProvidersFail(t *testing.T) {
	p1 := &fakeMarketProvider{
		name: "cg",
		fetchFunc: func(context.Context) (MarketSnapshot, error) {
			return snapshotWithBTC("cg", 111), nil
		},
	}
	fx := &fakeFXProvider{fetchFunc: func(context.Context) (FXSnapshot, error) {
		return FXSnapshot{Base: "USD", Rates: map[i18n.FiatCurrency]float64{i18n.FiatUSD: 1}}, nil
	}}
	var lastStatus StatusEvent
	feed := New([]MarketProvider{p1}, fx, Callbacks{
		OnStatus: func(event StatusEvent) { lastStatus = event },
	})

	feed.runFXCycle()
	feed.runMarketCycle()

	p1.fetchFunc = func(context.Context) (MarketSnapshot, error) {
		return MarketSnapshot{}, &ProviderError{Provider: "cg", Kind: FailureKindRateLimit, StatusCode: 429}
	}
	feed.runMarketCycle()

	if lastStatus.Kind != StatusKindWarning || lastStatus.Code != StatusCodeRateLimited {
		t.Fatalf("expected rate-limited warning with cached data, got %+v", lastStatus)
	}
}

func snapshotWithBTC(provider string, price float64) MarketSnapshot {
	change := 1.25
	return MarketSnapshot{
		Provider:  provider,
		FetchedAt: time.Unix(1700000000, 0),
		Coins: map[string]CoinQuoteUSD{
			"bitcoin": {
				ID:         "bitcoin",
				Name:       "Bitcoin",
				Ticker:     "BTC",
				PriceUSD:   price,
				Change24h:  &change,
				LastUpdate: time.Unix(1700000000, 0),
			},
		},
	}
}

func firstBTCPrice(t *testing.T, coins []model.Coin) float64 {
	t.Helper()
	for _, coin := range coins {
		if coin.ID == "bitcoin" {
			return coin.Price
		}
	}
	t.Fatal("bitcoin not found")
	return 0
}
