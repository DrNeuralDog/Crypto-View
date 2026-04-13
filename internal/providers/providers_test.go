package providers

import (
	"errors"
	"testing"

	"cryptoview/internal/marketfeed"
)

func TestParseRetryAfter(t *testing.T) {
	if got := parseRetryAfter("7"); got.Seconds() != 7 {
		t.Fatalf("expected 7s retry-after, got %v", got)
	}
	if got := parseRetryAfter(""); got != 0 {
		t.Fatalf("expected zero retry-after for empty header, got %v", got)
	}
}

func TestWrapNetworkError(t *testing.T) {
	err := wrapNetworkError("cg", errors.New("boom"))
	var providerErr *marketfeed.ProviderError
	if !errors.As(err, &providerErr) {
		t.Fatalf("expected ProviderError, got %T", err)
	}
	if providerErr.Kind != marketfeed.FailureKindOther {
		t.Fatalf("expected other failure kind, got %s", providerErr.Kind)
	}
}

func TestSanitizeDisplayString(t *testing.T) {
	if got := sanitizeDisplayString("Bit\u200bcoin", 100); got != "Bitcoin" {
		t.Fatalf("expected zero-width chars to be removed, got %q", got)
	}
	if got := sanitizeDisplayString("\u202E", 100); got != "" {
		t.Fatalf("expected bidi control char to be removed, got %q", got)
	}
}
