package marketfeed

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"cryptoview/internal/model"
	"cryptoview/internal/ui/i18n"
)

const (
	defaultMarketPollInterval = 2 * time.Second
	defaultFXPollInterval     = 30 * time.Second
)

type StatusKind string

const (
	StatusKindLoading StatusKind = "loading"
	StatusKindOK      StatusKind = "ok"
	StatusKindWarning StatusKind = "warning"
	StatusKindError   StatusKind = "error"
)

type StatusCode string

const (
	StatusCodeRateLimited StatusCode = "rate_limited"
	StatusCodeOffline     StatusCode = "offline_cached"
	StatusCodeFallback    StatusCode = "fallback_active"
	StatusCodeNoData      StatusCode = "no_data"
)

type StatusEvent struct {
	Kind     StatusKind
	Code     StatusCode
	Provider string
	Err      error
}

type Callbacks struct {
	OnMarketUpdate func([]model.Coin)
	OnStatus       func(StatusEvent)
}

type MarketProvider interface {
	Name() string
	FetchUSD(ctx context.Context) (MarketSnapshot, error)
}

type FXProvider interface {
	Name() string
	FetchRates(ctx context.Context) (FXSnapshot, error)
}

type CoinQuoteUSD struct {
	ID         string
	Name       string
	Ticker     string
	PriceUSD   float64
	Change24h  *float64
	LastUpdate time.Time
}

type MarketSnapshot struct {
	Provider  string
	FetchedAt time.Time
	Coins     map[string]CoinQuoteUSD
}

type FXSnapshot struct {
	Base      string
	FetchedAt time.Time
	Rates     map[i18n.FiatCurrency]float64
}

type providerState struct {
	cooldownUntil       time.Time
	consecutiveFailures int
}

type attemptFailure struct {
	err error
}

type Feed struct {
	mu sync.RWMutex

	providers   []MarketProvider
	fxProvider  FXProvider
	callbacks   Callbacks
	currentFiat i18n.FiatCurrency

	lastMarket *MarketSnapshot
	lastFX     *FXSnapshot
	state      map[string]*providerState

	marketPollInterval time.Duration
	fxPollInterval     time.Duration

	stopCh   chan struct{}
	wg       sync.WaitGroup
	started  bool
	stopOnce sync.Once
}

func NewDefault(callbacks Callbacks) *Feed {
	return New(
		[]MarketProvider{
			NewCoinGeckoProvider(1 * time.Second),
			NewCryptoCompareProvider(3 * time.Second),
			NewCoinLoreProvider(3 * time.Second),
		},
		NewOpenExchangeRatesProvider(1*time.Second),
		callbacks,
	)
}

func New(providers []MarketProvider, fxProvider FXProvider, callbacks Callbacks) *Feed {
	if len(providers) == 0 {
		panic("marketfeed: at least one market provider is required")
	}
	if fxProvider == nil {
		panic("marketfeed: fx provider is required")
	}

	f := &Feed{
		providers:          providers,
		fxProvider:         fxProvider,
		callbacks:          callbacks,
		currentFiat:        i18n.FiatUSD,
		state:              make(map[string]*providerState, len(providers)),
		marketPollInterval: defaultMarketPollInterval,
		fxPollInterval:     defaultFXPollInterval,
		stopCh:             make(chan struct{}),
	}
	for _, p := range providers {
		f.state[p.Name()] = &providerState{}
	}
	f.lastFX = &FXSnapshot{
		Base:      "USD",
		FetchedAt: time.Time{},
		Rates: map[i18n.FiatCurrency]float64{
			i18n.FiatUSD: 1,
		},
	}
	return f
}

func (f *Feed) Start() {
	f.mu.Lock()
	if f.started {
		f.mu.Unlock()
		return
	}
	f.started = true
	f.mu.Unlock()

	f.emitStatus(StatusEvent{Kind: StatusKindLoading})

	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		f.runLoop()
	}()
}

func (f *Feed) Stop() {
	f.stopOnce.Do(func() {
		close(f.stopCh)
	})
	f.wg.Wait()
}

func (f *Feed) SetFiat(currency i18n.FiatCurrency) {
	if _, ok := i18n.ParseFiatCurrency(string(currency)); !ok {
		return
	}
	f.mu.Lock()
	f.currentFiat = currency
	coins, ok := f.buildDisplayCoinsLocked()
	f.mu.Unlock()
	if ok {
		f.emitMarketUpdate(coins)
	}
}

func (f *Feed) runLoop() {
	fxTicker := time.NewTicker(f.fxPollInterval)
	defer fxTicker.Stop()
	marketTicker := time.NewTicker(f.marketPollInterval)
	defer marketTicker.Stop()

	f.runFXCycle()
	f.runMarketCycle()

	for {
		select {
		case <-marketTicker.C:
			f.runMarketCycle()
		case <-fxTicker.C:
			f.runFXCycle()
		case <-f.stopCh:
			return
		}
	}
}

func (f *Feed) runFXCycle() {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	snapshot, err := f.fxProvider.FetchRates(ctx)
	if err != nil {
		log.Printf("fx fetch failed via %s: %v", f.fxProvider.Name(), err)
		return
	}
	if snapshot.Rates == nil {
		return
	}
	if _, ok := snapshot.Rates[i18n.FiatUSD]; !ok {
		snapshot.Rates[i18n.FiatUSD] = 1
	}

	f.mu.Lock()
	f.lastFX = &snapshot
	f.mu.Unlock()
}

func (f *Feed) runMarketCycle() {
	now := time.Now()
	failures := make([]attemptFailure, 0, len(f.providers))
	attemptedProviders := 0

	for idx, provider := range f.providers {
		if !f.providerAvailable(provider.Name(), now) {
			remaining := f.providerCooldownRemaining(provider.Name(), now)
			log.Printf("marketfeed: skip provider=%s reason=cooldown remaining=%s", provider.Name(), remaining.Round(time.Second))
			continue
		}
		log.Printf("marketfeed: fetch attempt provider=%s", provider.Name())
		attemptedProviders++
		snapshot, err := f.fetchProvider(now, provider)
		if err != nil {
			log.Printf("marketfeed: fetch failed provider=%s err=%v", provider.Name(), err)
			failures = append(failures, attemptFailure{err: err})
			continue
		}
		log.Printf("marketfeed: fetch success provider=%s coins=%d", provider.Name(), len(snapshot.Coins))

		f.mu.Lock()
		f.mergeMissingChangesLocked(&snapshot)
		f.lastMarket = &snapshot
		coins, ok := f.buildDisplayCoinsLocked()
		f.mu.Unlock()

		if ok {
			f.emitMarketUpdate(coins)
		}

		if idx == 0 {
			f.emitStatus(StatusEvent{Kind: StatusKindOK, Provider: provider.Name()})
		} else {
			f.emitStatus(StatusEvent{
				Kind:     StatusKindWarning,
				Code:     StatusCodeFallback,
				Provider: provider.Name(),
			})
		}
		return
	}

	f.mu.RLock()
	coins, hasCache := f.buildDisplayCoinsLocked()
	f.mu.RUnlock()
	if hasCache {
		log.Printf("marketfeed: all providers failed, using cached market snapshot")
		f.emitMarketUpdate(coins)
		if hasRateLimitFailure(failures) {
			f.emitStatus(StatusEvent{Kind: StatusKindWarning, Code: StatusCodeRateLimited})
		} else {
			f.emitStatus(StatusEvent{Kind: StatusKindWarning, Code: StatusCodeOffline})
		}
		return
	}
	if attemptedProviders == 0 {
		log.Printf("marketfeed: all providers are cooling down; waiting for next window")
		f.emitStatus(StatusEvent{Kind: StatusKindLoading})
		return
	}

	combinedErr := combineFailures(failures)
	f.emitStatus(StatusEvent{
		Kind: StatusKindError,
		Code: StatusCodeNoData,
		Err:  combinedErr,
	})
	log.Printf("marketfeed: no provider data and no cache available err=%v", combinedErr)
}

func (f *Feed) providerAvailable(name string, now time.Time) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	st := f.state[name]
	if st == nil {
		return true
	}
	return !now.Before(st.cooldownUntil)
}

func (f *Feed) providerCooldownRemaining(name string, now time.Time) time.Duration {
	f.mu.RLock()
	defer f.mu.RUnlock()
	st := f.state[name]
	if st == nil || st.cooldownUntil.IsZero() {
		return 0
	}
	if now.After(st.cooldownUntil) {
		return 0
	}
	return st.cooldownUntil.Sub(now)
}

func (f *Feed) fetchProvider(now time.Time, provider MarketProvider) (MarketSnapshot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	snapshot, err := provider.FetchUSD(ctx)
	if err != nil {
		f.recordProviderFailure(now, provider.Name(), err)
		return MarketSnapshot{}, err
	}
	f.recordProviderSuccess(provider.Name())
	return snapshot, nil
}

func (f *Feed) recordProviderSuccess(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	st := f.state[name]
	if st == nil {
		st = &providerState{}
		f.state[name] = st
	}
	st.consecutiveFailures = 0
	st.cooldownUntil = time.Time{}
}

func (f *Feed) recordProviderFailure(now time.Time, name string, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	st := f.state[name]
	if st == nil {
		st = &providerState{}
		f.state[name] = st
	}
	st.consecutiveFailures++
	cooldown := failureCooldown(st.consecutiveFailures, err)
	if cooldown > 0 {
		st.cooldownUntil = now.Add(cooldown)
	}
}

func failureCooldown(failures int, err error) time.Duration {
	var pe *ProviderError
	if errors.As(err, &pe) {
		switch pe.Kind {
		case FailureKindRateLimit:
			if pe.RetryAfter > 0 {
				if pe.RetryAfter > 20*time.Second {
					return 20 * time.Second
				}
				return pe.RetryAfter
			}
			steps := []time.Duration{5 * time.Second, 10 * time.Second, 20 * time.Second}
			return steps[minInt(failures-1, len(steps)-1)]
		case FailureKindNetwork:
			steps := []time.Duration{4 * time.Second, 8 * time.Second, 20 * time.Second}
			return steps[minInt(failures-1, len(steps)-1)]
		}
	}
	return 20 * time.Second
}

func (f *Feed) mergeMissingChangesLocked(next *MarketSnapshot) {
	if f.lastMarket == nil || next == nil {
		return
	}
	for id, quote := range next.Coins {
		if quote.Change24h != nil {
			continue
		}
		prev, ok := f.lastMarket.Coins[id]
		if !ok || prev.Change24h == nil {
			continue
		}
		prevChange := *prev.Change24h
		quote.Change24h = &prevChange
		next.Coins[id] = quote
	}
}

func (f *Feed) buildDisplayCoinsLocked() ([]model.Coin, bool) {
	if f.lastMarket == nil {
		return nil, false
	}

	fiat := f.currentFiat
	rate := 1.0
	if f.lastFX != nil {
		if r, ok := f.lastFX.Rates[fiat]; ok && r > 0 {
			rate = r
		} else if fiat != i18n.FiatUSD {
			return nil, false
		}
	} else if fiat != i18n.FiatUSD {
		return nil, false
	}

	coins := make([]model.Coin, 0, len(trackedOrder))
	for _, id := range trackedOrder {
		quote, ok := f.lastMarket.Coins[id]
		if !ok {
			continue
		}
		change := 0.0
		if quote.Change24h != nil {
			change = *quote.Change24h
		}
		lastTime := "--:--:--"
		if !quote.LastUpdate.IsZero() {
			lastTime = quote.LastUpdate.Local().Format("15:04:05")
		}
		coins = append(coins, model.Coin{
			ID:             id,
			Name:           chooseString(quote.Name, defaultCoinNames[id], id),
			Ticker:         chooseString(quote.Ticker, defaultTickers[id]),
			Price:          quote.PriceUSD * rate,
			Change24h:      change,
			LastUpdateTime: lastTime,
			IconPath:       model.IconPathForID(id),
		})
	}
	return coins, len(coins) > 0
}

func (f *Feed) emitMarketUpdate(coins []model.Coin) {
	if f.callbacks.OnMarketUpdate != nil {
		f.callbacks.OnMarketUpdate(coins)
	}
}

func (f *Feed) emitStatus(event StatusEvent) {
	if f.callbacks.OnStatus != nil {
		f.callbacks.OnStatus(event)
	}
}

func hasRateLimitFailure(failures []attemptFailure) bool {
	for _, failure := range failures {
		var pe *ProviderError
		if errors.As(failure.err, &pe) && pe.Kind == FailureKindRateLimit {
			return true
		}
	}
	return false
}

func combineFailures(failures []attemptFailure) error {
	if len(failures) == 0 {
		return nil
	}
	if len(failures) == 1 {
		return failures[0].err
	}
	return fmt.Errorf("%d provider failures; last: %w", len(failures), failures[len(failures)-1].err)
}

func chooseString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (f *Feed) setIntervalsForTest(market, fx time.Duration) {
	if market > 0 {
		f.marketPollInterval = market
	}
	if fx > 0 {
		f.fxPollInterval = fx
	}
}

var trackedOrder = []string{
	"bitcoin",
	"ethereum",
	"the-open-network",
	"solana",
	"dogecoin",
	"ripple",
	"litecoin",
}

var defaultCoinNames = map[string]string{
	"bitcoin":          "Bitcoin",
	"ethereum":         "Ethereum",
	"the-open-network": "TON Coin",
	"solana":           "Solana",
	"dogecoin":         "Dogecoin",
	"ripple":           "Ripple",
	"litecoin":         "Litecoin",
}

var defaultTickers = map[string]string{
	"bitcoin":          "BTC",
	"ethereum":         "ETH",
	"the-open-network": "TON",
	"solana":           "SOL",
	"dogecoin":         "DOGE",
	"ripple":           "XRP",
	"litecoin":         "LTC",
}
