package appresources

import "testing"

func TestEmbeddedResourcesAvailable(t *testing.T) {
	if AppIcon() == nil {
		t.Fatal("expected embedded app icon")
	}
	if CoinIcon("bitcoin") == nil {
		t.Fatal("expected embedded bitcoin icon")
	}
	if CoinIcon("unknown") != nil {
		t.Fatal("expected nil icon for unknown coin")
	}
}
