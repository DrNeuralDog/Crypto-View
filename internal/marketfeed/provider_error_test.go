package marketfeed

import (
	"errors"
	"testing"
)

func TestProviderErrorError(t *testing.T) {
	t.Run("nil receiver", func(t *testing.T) {
		var err *ProviderError
		if got := err.Error(); got != "provider error" {
			t.Fatalf("expected default for nil, got %q", got)
		}
	})

	t.Run("without wrapped err", func(t *testing.T) {
		err := &ProviderError{Provider: "cg", Kind: FailureKindRateLimit}
		if got := err.Error(); got != "cg: rate_limit" {
			t.Fatalf("expected kind message, got %q", got)
		}
	})

	t.Run("with wrapped err", func(t *testing.T) {
		inner := errors.New("connection refused")
		err := &ProviderError{Provider: "cg", Kind: FailureKindNetwork, Err: inner}
		if got := err.Error(); got != "cg: connection refused" {
			t.Fatalf("expected wrapped message, got %q", got)
		}
	})
}

func TestProviderErrorUnwrap(t *testing.T) {
	inner := errors.New("inner")
	err := &ProviderError{Provider: "cg", Err: inner}
	if got := err.Unwrap(); got != inner {
		t.Fatalf("expected Unwrap to return inner error, got %v", got)
	}

	var nilErr *ProviderError
	if nilErr.Unwrap() != nil {
		t.Fatal("expected nil Unwrap for nil receiver")
	}
}

func TestErrorsAsProviderError(t *testing.T) {
	err := &ProviderError{Provider: "cg", Kind: FailureKindRateLimit, StatusCode: 429}
	var target *ProviderError
	if !errors.As(err, &target) {
		t.Fatal("expected errors.As to match ProviderError")
	}
	if target.Provider != "cg" || target.Kind != FailureKindRateLimit {
		t.Fatalf("expected target to match, got %+v", target)
	}
}
