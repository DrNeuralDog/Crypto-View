package marketfeed

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"cryptoview/internal/catalog"
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

type FailureKind string

const (
	FailureKindRateLimit FailureKind = "rate_limit"
	FailureKindNetwork   FailureKind = "network"
	FailureKindOther     FailureKind = "other"
)

type ProviderError struct {
	Provider   string
	Kind       FailureKind
	StatusCode int
	RetryAfter time.Duration
	Err        error
}

func (e *ProviderError) Error() string {
	if e == nil {
		return "provider error"
	}
	if e.Err == nil {
		return fmt.Sprintf("%s: %s", e.Provider, e.Kind)
	}
	return fmt.Sprintf("%s: %v", e.Provider, e.Err)
}

func (e *ProviderError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

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

type Option func(*Feed)

type Feed struct {
	mu sync.RWMutex

	providers   []MarketProvider
	fxProvider  FXProvider
	callbacks   Callbacks
	currentFiat i18n.FiatCurrency

	lastMarket *MarketSnapshot
	lastFX     *FXSnapshot
	state      map[string]*providerState
	logger     *slog.Logger

	marketPollInterval time.Duration
	fxPollInterval     time.Duration
	runCtx             context.Context
	runCancel          context.CancelFunc

	stopCh   chan struct{}
	wg       sync.WaitGroup
	started  bool
	stopOnce sync.Once
}

func WithLogger(logger *slog.Logger) Option {
	return func(feed *Feed) {
		if logger != nil {
			feed.logger = logger
		}
	}
}

func New(providers []MarketProvider, fxProvider FXProvider, callbacks Callbacks, opts ...Option) (*Feed, error) {
	if len(providers) == 0 {
		return nil, errors.New("marketfeed: at least one market provider is required")
	}
	if fxProvider == nil {
		return nil, errors.New("marketfeed: fx provider is required")
	}

	runCtx, runCancel := context.WithCancel(context.Background())

	feed := &Feed{
		providers:          providers,
		fxProvider:         fxProvider,
		callbacks:          callbacks,
		currentFiat:        i18n.FiatUSD,
		state:              make(map[string]*providerState, len(providers)),
		logger:             slog.New(slog.NewTextHandler(io.Discard, nil)),
		marketPollInterval: defaultMarketPollInterval,
		fxPollInterval:     defaultFXPollInterval,
		runCtx:             runCtx,
		runCancel:          runCancel,
		stopCh:             make(chan struct{}),
	}

	for _, opt := range opts {
		opt(feed)
	}

	for _, provider := range providers {
		feed.state[provider.Name()] = &providerState{}
	}

	feed.lastFX = &FXSnapshot{
		Base:      "USD",
		FetchedAt: time.Time{},
		Rates: map[i18n.FiatCurrency]float64{
			i18n.FiatUSD: 1,
		},
	}

	return feed, nil
}

func (f *Feed) SetCallbacks(callbacks Callbacks) {
	f.mu.Lock()
	f.callbacks = callbacks
	f.mu.Unlock()
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
		if f.runCancel != nil {
			f.runCancel()
		}
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
	if f.isStopping() {
		return
	}

	ctx, cancel := context.WithTimeout(f.runCtx, 12*time.Second)
	defer cancel()

	snapshot, err := f.fxProvider.FetchRates(ctx)
	if err != nil {
		f.logger.Debug("fx fetch failed", "provider", f.fxProvider.Name(), "error", err)
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
	if f.isStopping() {
		return
	}

	now := time.Now()
	failures := make([]attemptFailure, 0, len(f.providers))
	attemptedProviders := 0

	for idx, provider := range f.providers {
		if f.isStopping() {
			return
		}
		if !f.providerAvailable(provider.Name(), now) {
			remaining := f.providerCooldownRemaining(provider.Name(), now)
			f.logger.Debug("skip provider on cooldown", "provider", provider.Name(), "remaining", remaining.Round(time.Second))
			continue
		}

		f.logger.Debug("fetch attempt", "provider", provider.Name())
		attemptedProviders++

		snapshot, err := f.fetchProvider(now, provider)
		if err != nil {
			f.logger.Debug("fetch failed", "provider", provider.Name(), "error", err)
			failures = append(failures, attemptFailure{err: err})
			continue
		}

		f.logger.Debug("fetch success", "provider", provider.Name(), "coins", len(snapshot.Coins))

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
		f.logger.Debug("all providers failed, using cached market snapshot")
		f.emitMarketUpdate(coins)
		if hasRateLimitFailure(failures) {
			f.emitStatus(StatusEvent{Kind: StatusKindWarning, Code: StatusCodeRateLimited})
		} else {
			f.emitStatus(StatusEvent{Kind: StatusKindWarning, Code: StatusCodeOffline})
		}
		return
	}
	if attemptedProviders == 0 {
		f.logger.Debug("all providers are cooling down")
		f.emitStatus(StatusEvent{Kind: StatusKindLoading})
		return
	}

	combinedErr := combineFailures(failures)
	f.emitStatus(StatusEvent{
		Kind: StatusKindError,
		Code: StatusCodeNoData,
		Err:  combinedErr,
	})
	f.logger.Debug("no provider data and no cache available", "error", combinedErr)
}

func (f *Feed) providerAvailable(name string, now time.Time) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	state := f.state[name]
	if state == nil {
		return true
	}
	return !now.Before(state.cooldownUntil)
}

func (f *Feed) providerCooldownRemaining(name string, now time.Time) time.Duration {
	f.mu.RLock()
	defer f.mu.RUnlock()
	state := f.state[name]
	if state == nil || state.cooldownUntil.IsZero() {
		return 0
	}
	if now.After(state.cooldownUntil) {
		return 0
	}
	return state.cooldownUntil.Sub(now)
}

func (f *Feed) fetchProvider(now time.Time, provider MarketProvider) (MarketSnapshot, error) {
	ctx, cancel := context.WithTimeout(f.runCtx, 12*time.Second)
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

	state := f.state[name]
	if state == nil {
		state = &providerState{}
		f.state[name] = state
	}
	state.consecutiveFailures = 0
	state.cooldownUntil = time.Time{}
}

func (f *Feed) recordProviderFailure(now time.Time, name string, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	state := f.state[name]
	if state == nil {
		state = &providerState{}
		f.state[name] = state
	}
	state.consecutiveFailures++
	cooldown := failureCooldown(state.consecutiveFailures, err)
	if cooldown > 0 {
		state.cooldownUntil = now.Add(cooldown)
	}
}

func failureCooldown(failures int, err error) time.Duration {
	var providerErr *ProviderError
	if errors.As(err, &providerErr) {
		switch providerErr.Kind {
		case FailureKindRateLimit:
			if providerErr.RetryAfter > 0 {
				if providerErr.RetryAfter > 20*time.Second {
					return 20 * time.Second
				}
				return providerErr.RetryAfter
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
		if value, ok := f.lastFX.Rates[fiat]; ok && value > 0 {
			rate = value
		} else if fiat != i18n.FiatUSD {
			return nil, false
		}
	} else if fiat != i18n.FiatUSD {
		return nil, false
	}

	coins := make([]model.Coin, 0, len(catalog.Tracked()))
	for _, meta := range catalog.Tracked() {
		quote, ok := f.lastMarket.Coins[meta.ID]
		if !ok {
			continue
		}

		change := 0.0
		if quote.Change24h != nil {
			change = *quote.Change24h
		}

		lastUpdated := quote.LastUpdate
		if lastUpdated.IsZero() {
			lastUpdated = f.lastMarket.FetchedAt
		}

		coins = append(coins, model.Coin{
			ID:          meta.ID,
			Name:        chooseString(quote.Name, meta.Name, meta.ID),
			Ticker:      chooseString(quote.Ticker, meta.Ticker),
			Price:       quote.PriceUSD * rate,
			Change24h:   change,
			LastUpdated: lastUpdated,
		})
	}
	return coins, len(coins) > 0
}

func (f *Feed) emitMarketUpdate(coins []model.Coin) {
	f.mu.RLock()
	callback := f.callbacks.OnMarketUpdate
	f.mu.RUnlock()

	if callback != nil {
		callback(coins)
	}
}

func (f *Feed) emitStatus(event StatusEvent) {
	f.mu.RLock()
	callback := f.callbacks.OnStatus
	f.mu.RUnlock()

	if callback != nil {
		callback(event)
	}
}

func hasRateLimitFailure(failures []attemptFailure) bool {
	for _, failure := range failures {
		var providerErr *ProviderError
		if errors.As(failure.err, &providerErr) && providerErr.Kind == FailureKindRateLimit {
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

func (f *Feed) isStopping() bool {
	select {
	case <-f.stopCh:
		return true
	default:
		return false
	}
}

func (f *Feed) setIntervalsForTest(market, fx time.Duration) {
	if market > 0 {
		f.marketPollInterval = market
	}
	if fx > 0 {
		f.fxPollInterval = fx
	}
}
