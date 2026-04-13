package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"cryptoview/internal/catalog"
	"cryptoview/internal/marketfeed"
	"cryptoview/internal/ui/i18n"
)

const maxResponseBodySize = 5 * 1024 * 1024 // 5 MB

type CoinGeckoProvider struct {
	httpClient *http.Client
	baseURL    string
}

func NewCoinGeckoProvider(timeout time.Duration) *CoinGeckoProvider {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &CoinGeckoProvider{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    "https://api.coingecko.com/api/v3",
	}
}

func (p *CoinGeckoProvider) Name() string { return "coingecko" }

func (p *CoinGeckoProvider) FetchUSD(ctx context.Context) (marketfeed.MarketSnapshot, error) {
	values := url.Values{}
	values.Set("vs_currency", "usd")
	values.Set("ids", catalog.IDsCSV())
	values.Set("order", "market_cap_desc")
	values.Set("sparkline", "false")
	values.Set("price_change_percentage", "24h")

	endpoint := p.baseURL + "/coins/markets?" + values.Encode()

	body, _, err := doJSONRequest(ctx, p.httpClient, p.Name(), endpoint)
	if err != nil {
		return marketfeed.MarketSnapshot{}, err
	}

	var payload []struct {
		ID                       string  `json:"id"`
		Symbol                   string  `json:"symbol"`
		Name                     string  `json:"name"`
		CurrentPrice             float64 `json:"current_price"`
		PriceChangePercentage24h float64 `json:"price_change_percentage_24h"`
		LastUpdated              string  `json:"last_updated"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return marketfeed.MarketSnapshot{}, &marketfeed.ProviderError{Provider: p.Name(), Kind: marketfeed.FailureKindOther, Err: err}
	}

	now := time.Now()
	coins := make(map[string]marketfeed.CoinQuoteUSD, len(payload))
	for _, item := range payload {
		id := catalog.CanonicalID(item.ID, item.Symbol)
		if id == "" {
			continue
		}
		if !isFinitePositive(item.CurrentPrice) {
			continue
		}

		change := item.PriceChangePercentage24h
		if !isFinite(change) {
			change = 0
		}

		lastUpdate := now
		if parsed, err := time.Parse(time.RFC3339Nano, item.LastUpdated); err == nil {
			lastUpdate = parsed
		} else if parsed, err := time.Parse(time.RFC3339, item.LastUpdated); err == nil {
			lastUpdate = parsed
		}

		coins[id] = marketfeed.CoinQuoteUSD{
			ID:         id,
			Name:       sanitizeDisplayString(item.Name, 100),
			Ticker:     sanitizeDisplayString(strings.ToUpper(item.Symbol), 20),
			PriceUSD:   item.CurrentPrice,
			Change24h:  float64Ptr(change),
			LastUpdate: lastUpdate,
		}
	}

	return marketfeed.MarketSnapshot{Provider: p.Name(), FetchedAt: now, Coins: coins}, nil
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

func (p *CoinCapProvider) FetchUSD(ctx context.Context) (marketfeed.MarketSnapshot, error) {
	values := url.Values{}
	values.Set("ids", strings.Join([]string{
		"bitcoin", "ethereum", "toncoin", "solana", "dogecoin", "xrp", "litecoin",
	}, ","))
	endpoint := p.baseURL + "/assets?" + values.Encode()

	body, _, err := doJSONRequest(ctx, p.httpClient, p.Name(), endpoint)
	if err != nil {
		return marketfeed.MarketSnapshot{}, err
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
		return marketfeed.MarketSnapshot{}, &marketfeed.ProviderError{Provider: p.Name(), Kind: marketfeed.FailureKindOther, Err: err}
	}

	now := time.Now()
	coins := make(map[string]marketfeed.CoinQuoteUSD, len(payload.Data))
	for _, item := range payload.Data {
		id := catalog.CanonicalID(item.ID, item.Symbol)
		if id == "" {
			continue
		}

		price, err := strconv.ParseFloat(item.PriceUSD, 64)
		if err != nil || !isFinitePositive(price) {
			continue
		}

		coin := marketfeed.CoinQuoteUSD{
			ID:         id,
			Name:       sanitizeDisplayString(item.Name, 100),
			Ticker:     sanitizeDisplayString(strings.ToUpper(item.Symbol), 20),
			PriceUSD:   price,
			LastUpdate: now,
		}
		if item.ChangePercent24h != "" {
			if change, err := strconv.ParseFloat(item.ChangePercent24h, 64); err == nil && isFinite(change) {
				coin.Change24h = float64Ptr(change)
			}
		}
		coins[id] = coin
	}

	return marketfeed.MarketSnapshot{Provider: p.Name(), FetchedAt: now, Coins: coins}, nil
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

func (p *CoinPaprikaProvider) FetchUSD(ctx context.Context) (marketfeed.MarketSnapshot, error) {
	endpoint := p.baseURL + "/tickers?quotes=USD"
	body, _, err := doJSONRequest(ctx, p.httpClient, p.Name(), endpoint)
	if err != nil {
		return marketfeed.MarketSnapshot{}, err
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
		return marketfeed.MarketSnapshot{}, &marketfeed.ProviderError{Provider: p.Name(), Kind: marketfeed.FailureKindOther, Err: err}
	}

	now := time.Now()
	coins := make(map[string]marketfeed.CoinQuoteUSD, len(catalog.Tracked()))
	for _, item := range payload {
		id := catalog.CanonicalID(item.ID, item.Symbol)
		if id == "" || !isFinitePositive(item.Quotes.USD.Price) {
			continue
		}

		lastUpdate := now
		if ts, err := time.Parse(time.RFC3339Nano, item.LastUpdated); err == nil {
			lastUpdate = ts
		} else if ts, err := time.Parse(time.RFC3339, item.LastUpdated); err == nil {
			lastUpdate = ts
		}

		change := item.Quotes.USD.PercentChange24
		if !isFinite(change) {
			change = 0
		}

		coins[id] = marketfeed.CoinQuoteUSD{
			ID:         id,
			Name:       sanitizeDisplayString(item.Name, 100),
			Ticker:     sanitizeDisplayString(strings.ToUpper(item.Symbol), 20),
			PriceUSD:   item.Quotes.USD.Price,
			Change24h:  float64Ptr(change),
			LastUpdate: lastUpdate,
		}
		if len(coins) == len(catalog.Tracked()) {
			break
		}
	}

	return marketfeed.MarketSnapshot{Provider: p.Name(), FetchedAt: now, Coins: coins}, nil
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

func (p *CryptoCompareProvider) FetchUSD(ctx context.Context) (marketfeed.MarketSnapshot, error) {
	values := url.Values{}
	values.Set("fsyms", "BTC,ETH,TON,SOL,DOGE,XRP,LTC")
	values.Set("tsyms", "USD")
	endpoint := p.baseURL + "?" + values.Encode()

	body, _, err := doJSONRequest(ctx, p.httpClient, p.Name(), endpoint)
	if err != nil {
		return marketfeed.MarketSnapshot{}, err
	}

	var payload struct {
		RAW map[string]map[string]struct {
			Price          float64 `json:"PRICE"`
			ChangePct24h   float64 `json:"CHANGEPCT24HOUR"`
			LastUpdateUnix int64   `json:"LASTUPDATE"`
		} `json:"RAW"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return marketfeed.MarketSnapshot{}, &marketfeed.ProviderError{Provider: p.Name(), Kind: marketfeed.FailureKindOther, Err: err}
	}

	now := time.Now()
	coins := make(map[string]marketfeed.CoinQuoteUSD, len(payload.RAW))
	for symbol, byFiat := range payload.RAW {
		usd, ok := byFiat["USD"]
		if !ok || !isFinitePositive(usd.Price) {
			continue
		}

		id := catalog.CanonicalID(symbol)
		if id == "" {
			continue
		}

		change := usd.ChangePct24h
		if !isFinite(change) {
			change = 0
		}

		lastUpdate := now
		if usd.LastUpdateUnix > 0 {
			lastUpdate = time.Unix(usd.LastUpdateUnix, 0)
		}

		meta, _ := catalog.Lookup(id)
		coins[id] = marketfeed.CoinQuoteUSD{
			ID:         id,
			Name:       sanitizeDisplayString(meta.Name, 100),
			Ticker:     sanitizeDisplayString(strings.ToUpper(symbol), 20),
			PriceUSD:   usd.Price,
			Change24h:  float64Ptr(change),
			LastUpdate: lastUpdate,
		}
	}

	return marketfeed.MarketSnapshot{Provider: p.Name(), FetchedAt: now, Coins: coins}, nil
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

func (p *BinanceProvider) FetchUSD(ctx context.Context) (marketfeed.MarketSnapshot, error) {
	values := url.Values{}
	values.Set("symbols", `["BTCUSDT","ETHUSDT","TONUSDT","SOLUSDT","DOGEUSDT","XRPUSDT","LTCUSDT"]`)
	endpoint := p.baseURL + "?" + values.Encode()

	body, _, err := doJSONRequest(ctx, p.httpClient, p.Name(), endpoint)
	if err != nil {
		return marketfeed.MarketSnapshot{}, err
	}

	var payload []struct {
		Symbol             string `json:"symbol"`
		LastPrice          string `json:"lastPrice"`
		PriceChangePercent string `json:"priceChangePercent"`
		CloseTime          int64  `json:"closeTime"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return marketfeed.MarketSnapshot{}, &marketfeed.ProviderError{Provider: p.Name(), Kind: marketfeed.FailureKindOther, Err: err}
	}

	now := time.Now()
	coins := make(map[string]marketfeed.CoinQuoteUSD, len(payload))
	for _, item := range payload {
		id := catalog.CanonicalID(strings.TrimSuffix(strings.ToUpper(item.Symbol), "USDT"))
		if id == "" {
			continue
		}

		price, err := strconv.ParseFloat(item.LastPrice, 64)
		if err != nil || !isFinitePositive(price) {
			continue
		}

		lastUpdate := now
		if item.CloseTime > 0 {
			lastUpdate = time.UnixMilli(item.CloseTime)
		}

		meta, _ := catalog.Lookup(id)
		coin := marketfeed.CoinQuoteUSD{
			ID:         id,
			Name:       sanitizeDisplayString(meta.Name, 100),
			Ticker:     sanitizeDisplayString(meta.Ticker, 20),
			PriceUSD:   price,
			LastUpdate: lastUpdate,
		}
		if item.PriceChangePercent != "" {
			if change, err := strconv.ParseFloat(item.PriceChangePercent, 64); err == nil && isFinite(change) {
				coin.Change24h = float64Ptr(change)
			}
		}
		coins[id] = coin
	}

	return marketfeed.MarketSnapshot{Provider: p.Name(), FetchedAt: now, Coins: coins}, nil
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

func (p *CoinLoreProvider) FetchUSD(ctx context.Context) (marketfeed.MarketSnapshot, error) {
	values := url.Values{}
	values.Set("start", "0")
	values.Set("limit", "100")
	endpoint := p.baseURL + "?" + values.Encode()

	body, _, err := doJSONRequest(ctx, p.httpClient, p.Name(), endpoint)
	if err != nil {
		return marketfeed.MarketSnapshot{}, err
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
		return marketfeed.MarketSnapshot{}, &marketfeed.ProviderError{Provider: p.Name(), Kind: marketfeed.FailureKindOther, Err: err}
	}

	now := time.Now()
	coins := make(map[string]marketfeed.CoinQuoteUSD, len(payload.Data))
	for _, item := range payload.Data {
		id := catalog.CanonicalID(item.Symbol)
		if id == "" {
			continue
		}

		price, err := strconv.ParseFloat(item.PriceUSD, 64)
		if err != nil || !isFinitePositive(price) {
			continue
		}

		meta, _ := catalog.Lookup(id)
		coin := marketfeed.CoinQuoteUSD{
			ID:         id,
			Name:       sanitizeDisplayString(chooseString(item.Name, meta.Name), 100),
			Ticker:     sanitizeDisplayString(chooseString(meta.Ticker, strings.ToUpper(item.Symbol)), 20),
			PriceUSD:   price,
			LastUpdate: now,
		}
		if item.PercentChange24 != "" {
			if change, err := strconv.ParseFloat(item.PercentChange24, 64); err == nil && isFinite(change) {
				coin.Change24h = float64Ptr(change)
			}
		}
		coins[id] = coin
		if len(coins) == len(catalog.Tracked()) {
			break
		}
	}

	return marketfeed.MarketSnapshot{Provider: p.Name(), FetchedAt: now, Coins: coins}, nil
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

func (p *OpenExchangeRatesProvider) FetchRates(ctx context.Context) (marketfeed.FXSnapshot, error) {
	body, _, err := doJSONRequest(ctx, p.httpClient, p.Name(), p.baseURL)
	if err != nil {
		return marketfeed.FXSnapshot{}, err
	}

	var payload struct {
		Result string             `json:"result"`
		Rates  map[string]float64 `json:"rates"`
		Time   int64              `json:"time_last_update_unix"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return marketfeed.FXSnapshot{}, &marketfeed.ProviderError{Provider: p.Name(), Kind: marketfeed.FailureKindOther, Err: err}
	}
	if len(payload.Rates) == 0 {
		return marketfeed.FXSnapshot{}, &marketfeed.ProviderError{Provider: p.Name(), Kind: marketfeed.FailureKindOther, Err: fmt.Errorf("empty fx rates")}
	}

	snapshot := marketfeed.FXSnapshot{
		Base:      "USD",
		FetchedAt: time.Now(),
		Rates: map[i18n.FiatCurrency]float64{
			i18n.FiatUSD: 1,
		},
	}
	if payload.Time > 0 {
		snapshot.FetchedAt = time.Unix(payload.Time, 0)
	}
	if value := payload.Rates["EUR"]; isFinitePositive(value) && value > 0.1 && value < 10 {
		snapshot.Rates[i18n.FiatEUR] = value
	}
	if value := payload.Rates["RUB"]; isFinitePositive(value) && value > 1 && value < 10000 {
		snapshot.Rates[i18n.FiatRUB] = value
	}
	return snapshot, nil
}

func doJSONRequest(ctx context.Context, client *http.Client, providerName, endpoint string) ([]byte, http.Header, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, &marketfeed.ProviderError{Provider: providerName, Kind: marketfeed.FailureKindOther, Err: err}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "CryptoView/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, wrapNetworkError(providerName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))

		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		kind := marketfeed.FailureKindOther
		if resp.StatusCode == http.StatusTooManyRequests {
			kind = marketfeed.FailureKindRateLimit
		}
		return nil, resp.Header, &marketfeed.ProviderError{
			Provider:   providerName,
			Kind:       kind,
			StatusCode: resp.StatusCode,
			RetryAfter: retryAfter,
			Err:        fmt.Errorf("%s status: %d", providerName, resp.StatusCode),
		}
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodySize))
	if err != nil {
		return nil, resp.Header, wrapNetworkError(providerName, err)
	}
	return body, resp.Header, nil
}

func wrapNetworkError(provider string, err error) error {
	var netErr net.Error
	if errors.As(err, &netErr) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return &marketfeed.ProviderError{Provider: provider, Kind: marketfeed.FailureKindNetwork, Err: err}
	}
	return &marketfeed.ProviderError{Provider: provider, Kind: marketfeed.FailureKindOther, Err: err}
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
		if delay := time.Until(when); delay > 0 {
			return delay
		}
	}
	return 0
}

func isFinitePositive(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value > 0
}

func isFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func sanitizeDisplayString(value string, maxLen int) string {
	if !utf8.ValidString(value) {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(value))

	for _, r := range value {
		if unicode.IsControl(r) {
			continue
		}
		if unicode.Is(unicode.Cf, r) {
			continue
		}
		if unicode.Is(unicode.Co, r) {
			continue
		}

		builder.WriteRune(r)
		if maxLen > 0 && builder.Len() >= maxLen {
			break
		}
	}

	return strings.TrimSpace(builder.String())
}

func chooseString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func float64Ptr(value float64) *float64 {
	copyValue := value
	return &copyValue
}
