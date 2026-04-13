package catalog

import "testing"

func TestTrackedIDsAndLookup(t *testing.T) {
	if got := IDsCSV(); got != "bitcoin,ethereum,the-open-network,solana,dogecoin,ripple,litecoin" {
		t.Fatalf("unexpected tracked ids csv: %s", got)
	}

	coin, ok := Lookup("bitcoin")
	if !ok {
		t.Fatal("expected bitcoin metadata to exist")
	}
	if coin.Ticker != "BTC" {
		t.Fatalf("expected BTC ticker, got %s", coin.Ticker)
	}
	if coin.IconAsset != "coins/bitcoin.png" {
		t.Fatalf("expected bitcoin icon asset, got %s", coin.IconAsset)
	}
}

func TestCanonicalID(t *testing.T) {
	tests := []struct {
		values []string
		want   string
	}{
		{values: []string{"bitcoin"}, want: "bitcoin"},
		{values: []string{"BTC"}, want: "bitcoin"},
		{values: []string{"btc-bitcoin"}, want: "bitcoin"},
		{values: []string{"toncoin"}, want: "the-open-network"},
		{values: []string{"TON"}, want: "the-open-network"},
		{values: []string{"xrp"}, want: "ripple"},
		{values: []string{"LTC"}, want: "litecoin"},
		{values: []string{"unknown", "ethereum"}, want: "ethereum"},
		{values: []string{"unknown"}, want: ""},
	}

	for _, tt := range tests {
		if got := CanonicalID(tt.values...); got != tt.want {
			t.Errorf("CanonicalID(%q) = %q, want %q", tt.values, got, tt.want)
		}
	}
}

func TestTrackedReturnsCopy(t *testing.T) {
	coins := Tracked()
	coins[0].Name = "Changed"

	again, ok := Lookup("bitcoin")
	if !ok {
		t.Fatal("expected bitcoin metadata to exist")
	}
	if again.Name != "Bitcoin" {
		t.Fatalf("expected tracked copy isolation, got %s", again.Name)
	}
}
