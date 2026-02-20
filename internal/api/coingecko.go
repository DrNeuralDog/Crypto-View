package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"cryptoview/internal/model"
)

const trackedCoinIDs = "bitcoin,ethereum,the-open-network,solana,dogecoin,ripple,litecoin"

func (c *Client) GetMarkets(ctx context.Context, fiat string) ([]model.CoinGeckoMarket, error) {
	normalized, err := normalizeFiatCurrency(fiat)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("vs_currency", normalized)
	params.Set("ids", trackedCoinIDs)
	params.Set("order", "market_cap_desc")
	params.Set("sparkline", "false")
	params.Set("price_change_percentage", "24h")

	endpoint := fmt.Sprintf("%s/coins/markets?%s", c.baseURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("coingecko status: %d", resp.StatusCode)
	}

	var markets []model.CoinGeckoMarket
	if err := json.NewDecoder(resp.Body).Decode(&markets); err != nil {
		return nil, err
	}

	return markets, nil
}

func normalizeFiatCurrency(fiat string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(fiat)) {
	case "usd", "eur", "rub":
		return strings.ToLower(strings.TrimSpace(fiat)), nil
	default:
		return "", fmt.Errorf("unsupported fiat currency: %s", fiat)
	}
}
