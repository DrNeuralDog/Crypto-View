package api

import (
	"fmt"
	"time"
)

type StatusError struct {
	StatusCode int
	RetryAfter time.Duration
}

func (e *StatusError) Error() string {
	if e == nil {
		return "http status error"
	}
	return fmt.Sprintf("coingecko status: %d", e.StatusCode)
}
