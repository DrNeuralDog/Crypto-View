package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c := NewClient(5 * time.Second)
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.httpClient == nil {
		t.Fatal("expected http client to be set")
	}
	if c.httpClient.Timeout != 5*time.Second {
		t.Fatalf("expected timeout 5s, got %v", c.httpClient.Timeout)
	}
}

func TestNewClient_ZeroTimeoutUsesDefault(t *testing.T) {
	c := newClient("https://example.com", 0)
	if c.httpClient.Timeout != defaultTimeout {
		t.Fatalf("expected default timeout %v, got %v", defaultTimeout, c.httpClient.Timeout)
	}
}

func TestNewClient_NegativeTimeoutUsesDefault(t *testing.T) {
	c := newClient("https://example.com", -1*time.Second)
	if c.httpClient.Timeout != defaultTimeout {
		t.Fatalf("expected default timeout %v, got %v", defaultTimeout, c.httpClient.Timeout)
	}
}

func TestClient_TimeoutRespected(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	client := newClient(srv.URL, 10*time.Millisecond)
	ctx := context.Background()
	_, err := client.GetMarkets(ctx, "usd")
	if err == nil {
		t.Fatal("expected timeout error")
	}
}
