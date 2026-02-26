package marketfeed

import (
	"errors"
	"testing"
)

func TestProviderError_Error(t *testing.T) {
	t.Run("nil receiver", func(t *testing.T) {
		var e *ProviderError
		if got := e.Error(); got != "provider error" {
			t.Fatalf("expected default for nil, got %q", got)
		}
	})

	t.Run("without wrapped err", func(t *testing.T) {
		e := &ProviderError{Provider: "cg", Kind: FailureKindRateLimit}
		if got := e.Error(); got != "cg: rate_limit" {
			t.Fatalf("expected kind message, got %q", got)
		}
	})

	t.Run("with wrapped err", func(t *testing.T) {
		inner := errors.New("connection refused")
		e := &ProviderError{Provider: "cg", Kind: FailureKindNetwork, Err: inner}
		if got := e.Error(); got != "cg: connection refused" {
			t.Fatalf("expected wrapped message, got %q", got)
		}
	})
}

func TestProviderError_Unwrap(t *testing.T) {
	inner := errors.New("inner")
	e := &ProviderError{Provider: "cg", Err: inner}
	if got := e.Unwrap(); got != inner {
		t.Fatalf("expected Unwrap to return inner error, got %v", got)
	}

	var nilErr *ProviderError
	if nilErr.Unwrap() != nil {
		t.Fatal("expected nil Unwrap for nil receiver")
	}
}

func TestErrorsAs_ProviderError(t *testing.T) {
	pe := &ProviderError{Provider: "cg", Kind: FailureKindRateLimit, StatusCode: 429}
	var target *ProviderError
	if !errors.As(pe, &target) {
		t.Fatal("expected errors.As to match ProviderError")
	}
	if target.Provider != "cg" || target.Kind != FailureKindRateLimit {
		t.Fatalf("expected target to match, got %+v", target)
	}
}
