package api

import (
	"net/http"
	"time"
)

const defaultTimeout = 10 * time.Second

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(timeout time.Duration) *Client {
	return newClient("https://api.coingecko.com/api/v3", timeout)
}

func newClient(baseURL string, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    baseURL,
	}
}
