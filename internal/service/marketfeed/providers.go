package marketfeed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cryptoview/internal/api"
	"cryptoview/internal/ui/i18n"
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

type CoinGeckoProvider struct {
	client *api.Client
}

func NewCoinGeckoProvider(timeout time.Duration) *CoinGeckoProvider {
	return &CoinGeckoProvider{client: api.NewClient(timeout)}
}

func (p *CoinGeckoProvider) Name() string { return "coingecko" }

func (p *CoinGeckoProvider) FetchUSD(ctx context.Context) (MarketSnapshot, error) {
	markets, err := p.client.GetMarkets(ctx, "usd")
	if err != nil {
		var statusErr *api.StatusError
		if errors.As(err, &statusErr) {
			kind := FailureKindOther
			if statusErr.StatusCode == http.StatusTooManyRequests {
				kind = FailureKindRateLimit
			}
			return MarketSnapshot{}, &ProviderError{
				Provider:   p.Name(),
				Kind:       kind,
				StatusCode: statusErr.StatusCode,
				RetryAfter: statusErr.RetryAfter,
				Err:        err,
			}
		}
		return MarketSnapshot{}, wrapNetworkError(p.Name(), err)
	}

	now := time.Now()
	coins := make(map[string]CoinQuoteUSD, len(markets))
	for _, m := range markets {
		id := normalizeCanonicalID(m.ID, strings.ToUpper(m.Symbol))
		if id == "" {
			continue
		}
		change := m.PriceChangePercentage24h
		lastUpdate := now
		if parsed, err := time.Parse(time.RFC3339Nano, m.LastUpdated); err == nil {
			lastUpdate = parsed
		} else if parsed, err := time.Parse(time.RFC3339, m.LastUpdated); err == nil {
			lastUpdate = parsed
		}
		coins[id] = CoinQuoteUSD{
			ID:         id,
			Name:       m.Name,
			Ticker:     strings.ToUpper(m.Symbol),
			PriceUSD:   m.CurrentPrice,
			Change24h:  &change,
			LastUpdate: lastUpdate,
		}
	}
	return MarketSnapshot{Provider: p.Name(), FetchedAt: now, Coins: coins}, nil
}

type CoinCapProvider struct {
	httpClient *http.Client
	baseURL    string
}

func NewCoinCapProvider(timeout time.Duration) *CoinCapProvider {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &CoinCapProvider{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    "https://api.coincap.io/v2",
	}
}

func (p *CoinCapProvider) Name() string { return "coincap" }

func (p *CoinCapProvider) FetchUSD(ctx context.Context) (MarketSnapshot, error) {
	values := url.Values{}
	values.Set("ids", strings.Join([]string{
		"bitcoin", "ethereum", "toncoin", "solana", "dogecoin", "xrp", "litecoin",
	}, ","))
	endpoint := p.baseURL + "/assets?" + values.Encode()

	body, _, err := doJSONRequest(ctx, p.httpClient, p.Name(), endpoint)
	if err != nil {
		return MarketSnapshot{}, err
	}
	var payload struct {
		Data []struct {
			ID               string `json:"id"`
			Symbol           string `json:"symbol"`
			Name             string `json:"name"`
			PriceUSD         string `json:"priceUsd"`
			ChangePercent24h string `json:"changePercent24Hr"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return MarketSnapshot{}, &ProviderError{Provider: p.Name(), Kind: FailureKindOther, Err: err}
	}
	now := time.Now()
	coins := make(map[string]CoinQuoteUSD, len(payload.Data))
	for _, item := range payload.Data {
		id := normalizeCanonicalID(item.ID, strings.ToUpper(item.Symbol))
		if id == "" {
			continue
		}
		price, err := strconv.ParseFloat(item.PriceUSD, 64)
		if err != nil || price <= 0 {
			continue
		}
		var changePtr *float64
		if item.ChangePercent24h != "" {
			if c, err := strconv.ParseFloat(item.ChangePercent24h, 64); err == nil {
				changeCopy := c
				changePtr = &changeCopy
			}
		}
		coins[id] = CoinQuoteUSD{
			ID:         id,
			Name:       item.Name,
			Ticker:     strings.ToUpper(item.Symbol),
			PriceUSD:   price,
			Change24h:  changePtr,
			LastUpdate: now,
		}
	}
	return MarketSnapshot{Provider: p.Name(), FetchedAt: now, Coins: coins}, nil
}

type CoinPaprikaProvider struct {
	httpClient *http.Client
	baseURL    string
}

func NewCoinPaprikaProvider(timeout time.Duration) *CoinPaprikaProvider {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return &CoinPaprikaProvider{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    "https://api.coinpaprika.com/v1",
	}
}

func (p *CoinPaprikaProvider) Name() string { return "coinpaprika" }

func (p *CoinPaprikaProvider) FetchUSD(ctx context.Context) (MarketSnapshot, error) {
	endpoint := p.baseURL + "/tickers?quotes=USD"
	body, _, err := doJSONRequest(ctx, p.httpClient, p.Name(), endpoint)
	if err != nil {
		return MarketSnapshot{}, err
	}

	var payload []struct {
		ID          string `json:"id"`
		Symbol      string `json:"symbol"`
		Name        string `json:"name"`
		LastUpdated string `json:"last_updated"`
		Quotes      struct {
			USD struct {
				Price           float64 `json:"price"`
				PercentChange24 float64 `json:"percent_change_24h"`
			} `json:"USD"`
		} `json:"quotes"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return MarketSnapshot{}, &ProviderError{Provider: p.Name(), Kind: FailureKindOther, Err: err}
	}

	now := time.Now()
	coins := make(map[string]CoinQuoteUSD, len(trackedOrder))
	for _, item := range payload {
		id := normalizeCanonicalID(item.ID, strings.ToUpper(item.Symbol))
		if id == "" {
			continue
		}
		if item.Quotes.USD.Price <= 0 {
			continue
		}
		lastUpdate := now
		if ts, err := time.Parse(time.RFC3339Nano, item.LastUpdated); err == nil {
			lastUpdate = ts
		} else if ts, err := time.Parse(time.RFC3339, item.LastUpdated); err == nil {
			lastUpdate = ts
		}
		change := item.Quotes.USD.PercentChange24
		coins[id] = CoinQuoteUSD{
			ID:         id,
			Name:       item.Name,
			Ticker:     strings.ToUpper(item.Symbol),
			PriceUSD:   item.Quotes.USD.Price,
			Change24h:  &change,
			LastUpdate: lastUpdate,
		}
		if len(coins) == len(trackedOrder) {
			break
		}
	}
	return MarketSnapshot{Provider: p.Name(), FetchedAt: now, Coins: coins}, nil
}

type CryptoCompareProvider struct {
	httpClient *http.Client
	baseURL    string
}

func NewCryptoCompareProvider(timeout time.Duration) *CryptoCompareProvider {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &CryptoCompareProvider{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    "https://min-api.cryptocompare.com/data/pricemultifull",
	}
}

func (p *CryptoCompareProvider) Name() string { return "cryptocompare" }

func (p *CryptoCompareProvider) FetchUSD(ctx context.Context) (MarketSnapshot, error) {
	values := url.Values{}
	values.Set("fsyms", "BTC,ETH,TON,SOL,DOGE,XRP,LTC")
	values.Set("tsyms", "USD")
	endpoint := p.baseURL + "?" + values.Encode()

	body, _, err := doJSONRequest(ctx, p.httpClient, p.Name(), endpoint)
	if err != nil {
		return MarketSnapshot{}, err
	}

	var payload struct {
		RAW map[string]map[string]struct {
			Price          float64 `json:"PRICE"`
			ChangePct24h   float64 `json:"CHANGEPCT24HOUR"`
			LastUpdateUnix int64   `json:"LASTUPDATE"`
		} `json:"RAW"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return MarketSnapshot{}, &ProviderError{Provider: p.Name(), Kind: FailureKindOther, Err: err}
	}

	now := time.Now()
	coins := make(map[string]CoinQuoteUSD, len(payload.RAW))
	for symbol, byFiat := range payload.RAW {
		usd, ok := byFiat["USD"]
		if !ok || usd.Price <= 0 {
			continue
		}
		id := normalizeCanonicalID("", symbol)
		if id == "" {
			continue
		}
		change := usd.ChangePct24h
		lastUpdate := now
		if usd.LastUpdateUnix > 0 {
			lastUpdate = time.Unix(usd.LastUpdateUnix, 0)
		}
		coins[id] = CoinQuoteUSD{
			ID:         id,
			Name:       defaultCoinNames[id],
			Ticker:     strings.ToUpper(symbol),
			PriceUSD:   usd.Price,
			Change24h:  &change,
			LastUpdate: lastUpdate,
		}
	}

	return MarketSnapshot{Provider: p.Name(), FetchedAt: now, Coins: coins}, nil
}

type BinanceProvider struct {
	httpClient *http.Client
	baseURL    string
}

func NewBinanceProvider(timeout time.Duration) *BinanceProvider {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &BinanceProvider{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    "https://api.binance.com/api/v3/ticker/24hr",
	}
}

func (p *BinanceProvider) Name() string { return "binance" }

func (p *BinanceProvider) FetchUSD(ctx context.Context) (MarketSnapshot, error) {
	symbols := `["BTCUSDT","ETHUSDT","TONUSDT","SOLUSDT","DOGEUSDT","XRPUSDT","LTCUSDT"]`
	values := url.Values{}
	values.Set("symbols", symbols)
	endpoint := p.baseURL + "?" + values.Encode()

	body, _, err := doJSONRequest(ctx, p.httpClient, p.Name(), endpoint)
	if err != nil {
		return MarketSnapshot{}, err
	}

	var payload []struct {
		Symbol             string `json:"symbol"`
		LastPrice          string `json:"lastPrice"`
		PriceChangePercent string `json:"priceChangePercent"`
		CloseTime          int64  `json:"closeTime"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return MarketSnapshot{}, &ProviderError{Provider: p.Name(), Kind: FailureKindOther, Err: err}
	}

	now := time.Now()
	coins := make(map[string]CoinQuoteUSD, len(payload))
	for _, item := range payload {
		id := normalizeCanonicalID("", strings.TrimSuffix(strings.ToUpper(item.Symbol), "USDT"))
		if id == "" {
			continue
		}
		price, err := strconv.ParseFloat(item.LastPrice, 64)
		if err != nil || price <= 0 {
			continue
		}
		var changePtr *float64
		if item.PriceChangePercent != "" {
			if c, err := strconv.ParseFloat(item.PriceChangePercent, 64); err == nil {
				changeCopy := c
				changePtr = &changeCopy
			}
		}
		lastUpdate := now
		if item.CloseTime > 0 {
			lastUpdate = time.UnixMilli(item.CloseTime)
		}
		coins[id] = CoinQuoteUSD{
			ID:         id,
			Name:       defaultCoinNames[id],
			Ticker:     defaultTickers[id],
			PriceUSD:   price,
			Change24h:  changePtr,
			LastUpdate: lastUpdate,
		}
	}

	return MarketSnapshot{Provider: p.Name(), FetchedAt: now, Coins: coins}, nil
}

type CoinLoreProvider struct {
	httpClient *http.Client
	baseURL    string
}

func NewCoinLoreProvider(timeout time.Duration) *CoinLoreProvider {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &CoinLoreProvider{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    "https://api.coinlore.net/api/tickers/",
	}
}

func (p *CoinLoreProvider) Name() string { return "coinlore" }

func (p *CoinLoreProvider) FetchUSD(ctx context.Context) (MarketSnapshot, error) {
	values := url.Values{}
	values.Set("start", "0")
	values.Set("limit", "100")
	endpoint := p.baseURL + "?" + values.Encode()

	body, _, err := doJSONRequest(ctx, p.httpClient, p.Name(), endpoint)
	if err != nil {
		return MarketSnapshot{}, err
	}

	var payload struct {
		Data []struct {
			Symbol          string `json:"symbol"`
			Name            string `json:"name"`
			PriceUSD        string `json:"price_usd"`
			PercentChange24 string `json:"percent_change_24h"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return MarketSnapshot{}, &ProviderError{Provider: p.Name(), Kind: FailureKindOther, Err: err}
	}

	now := time.Now()
	coins := make(map[string]CoinQuoteUSD, len(payload.Data))
	for _, item := range payload.Data {
		id := normalizeCanonicalID("", item.Symbol)
		if id == "" {
			continue
		}
		price, err := strconv.ParseFloat(item.PriceUSD, 64)
		if err != nil || price <= 0 {
			continue
		}
		var changePtr *float64
		if item.PercentChange24 != "" {
			if c, err := strconv.ParseFloat(item.PercentChange24, 64); err == nil {
				changeCopy := c
				changePtr = &changeCopy
			}
		}
		coins[id] = CoinQuoteUSD{
			ID:         id,
			Name:       chooseString(item.Name, defaultCoinNames[id]),
			Ticker:     chooseString(defaultTickers[id], strings.ToUpper(item.Symbol)),
			PriceUSD:   price,
			Change24h:  changePtr,
			LastUpdate: now,
		}
		if len(coins) == len(trackedOrder) {
			break
		}
	}

	return MarketSnapshot{Provider: p.Name(), FetchedAt: now, Coins: coins}, nil
}

type OpenExchangeRatesProvider struct {
	httpClient *http.Client
	baseURL    string
}

func NewOpenExchangeRatesProvider(timeout time.Duration) *OpenExchangeRatesProvider {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &OpenExchangeRatesProvider{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    "https://open.er-api.com/v6/latest/USD",
	}
}

func (p *OpenExchangeRatesProvider) Name() string { return "open-er-api" }

func (p *OpenExchangeRatesProvider) FetchRates(ctx context.Context) (FXSnapshot, error) {
	body, _, err := doJSONRequest(ctx, p.httpClient, p.Name(), p.baseURL)
	if err != nil {
		return FXSnapshot{}, err
	}
	var payload struct {
		Result string             `json:"result"`
		Rates  map[string]float64 `json:"rates"`
		Time   int64              `json:"time_last_update_unix"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return FXSnapshot{}, &ProviderError{Provider: p.Name(), Kind: FailureKindOther, Err: err}
	}
	if len(payload.Rates) == 0 {
		return FXSnapshot{}, &ProviderError{Provider: p.Name(), Kind: FailureKindOther, Err: fmt.Errorf("empty fx rates")}
	}
	snapshot := FXSnapshot{
		Base:      "USD",
		FetchedAt: time.Now(),
		Rates: map[i18n.FiatCurrency]float64{
			i18n.FiatUSD: 1,
		},
	}
	if payload.Time > 0 {
		snapshot.FetchedAt = time.Unix(payload.Time, 0)
	}
	if v := payload.Rates["EUR"]; v > 0 {
		snapshot.Rates[i18n.FiatEUR] = v
	}
	if v := payload.Rates["RUB"]; v > 0 {
		snapshot.Rates[i18n.FiatRUB] = v
	}
	return snapshot, nil
}

func doJSONRequest(ctx context.Context, client *http.Client, providerName, endpoint string) ([]byte, http.Header, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, &ProviderError{Provider: providerName, Kind: FailureKindOther, Err: err}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "CryptoView/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, wrapNetworkError(providerName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		kind := FailureKindOther
		if resp.StatusCode == http.StatusTooManyRequests {
			kind = FailureKindRateLimit
		}
		return nil, resp.Header, &ProviderError{
			Provider:   providerName,
			Kind:       kind,
			StatusCode: resp.StatusCode,
			RetryAfter: retryAfter,
			Err:        fmt.Errorf("%s status: %d", providerName, resp.StatusCode),
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.Header, wrapNetworkError(providerName, err)
	}
	return body, resp.Header, nil
}

func wrapNetworkError(provider string, err error) error {
	var netErr net.Error
	if errors.As(err, &netErr) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return &ProviderError{Provider: provider, Kind: FailureKindNetwork, Err: err}
	}
	return &ProviderError{Provider: provider, Kind: FailureKindOther, Err: err}
}

func parseRetryAfter(raw string) time.Duration {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}
	if secs, err := strconv.Atoi(raw); err == nil && secs > 0 {
		return time.Duration(secs) * time.Second
	}
	if when, err := http.ParseTime(raw); err == nil {
		if d := time.Until(when); d > 0 {
			return d
		}
	}
	return 0
}

func normalizeCanonicalID(providerID, symbol string) string {
	switch strings.ToLower(strings.TrimSpace(providerID)) {
	case "bitcoin", "btc-bitcoin":
		return "bitcoin"
	case "ethereum", "eth-ethereum":
		return "ethereum"
	case "the-open-network", "toncoin", "ton-toncoin", "toncoin-toncoin":
		return "the-open-network"
	case "solana", "sol-solana":
		return "solana"
	case "dogecoin", "doge-dogecoin":
		return "dogecoin"
	case "ripple", "xrp", "xrp-xrp":
		return "ripple"
	case "litecoin", "ltc-litecoin":
		return "litecoin"
	}

	switch strings.ToUpper(strings.TrimSpace(symbol)) {
	case "BTC":
		return "bitcoin"
	case "ETH":
		return "ethereum"
	case "TON":
		return "the-open-network"
	case "SOL":
		return "solana"
	case "DOGE":
		return "dogecoin"
	case "XRP":
		return "ripple"
	case "LTC":
		return "litecoin"
	}

	return ""
}
