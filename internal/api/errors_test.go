package api

import "testing"

func TestStatusError_Error(t *testing.T) {
	t.Run("nil receiver", func(t *testing.T) {
		var e *StatusError
		if got := e.Error(); got != "http status error" {
			t.Fatalf("expected default message for nil, got %q", got)
		}
	})

	t.Run("with status code", func(t *testing.T) {
		e := &StatusError{StatusCode: 429}
		if got := e.Error(); got != "coingecko status: 429" {
			t.Fatalf("expected status message, got %q", got)
		}
	})

	t.Run("with 404", func(t *testing.T) {
		e := &StatusError{StatusCode: 404}
		if got := e.Error(); got != "coingecko status: 404" {
			t.Fatalf("expected status message, got %q", got)
		}
	})
}
