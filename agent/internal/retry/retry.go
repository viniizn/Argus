package retry

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"time"
)

type Config struct {
	MaxAttempts int
	InitialDelay time.Duration
	MaxDelay time.Duration
	Multiplier float64
}

func DefaultConfig() Config {
	return Config {
		MaxAttempts: 8,
		InitialDelay: 1 * time.Second,
		MaxDelay: 2 * time.Minute,
		Multiplier: 2.0,
	}
}

type NonRetryableError struct {
	Cause error
}

func (e * NonRetryableError) Error() string {
	return fmt.Sprintf("non-retryable: %v", e.Cause)
}

func (e *NonRetryableError) Unwrap() error {
	return e.Cause
}

func Do(ctx context.Context, cfg Config, logger * slog.Logger, fn func() error) error {
	var err error

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		var nonRetryable *NonRetryableError
		if errors.As(err, &nonRetryable) {
			return nonRetryable.Cause
		}

		if attempt == cfg.MaxAttempts {
			break
		}

		delay := calculateDelay(cfg, attempt)

		logger.Warn("operation failed, retrying",
			"attempt", attempt,
			"maxAttempts", cfg.MaxAttempts,
			"retryIn", delay.String(),
			"error", err,		
		)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		}
	}

	return fmt.Errorf("all %d attempts failed: %w", ctx.Err())
}

func calculateDelay(cfg Config, attempt int) time.Duration {
	base := float64(cfg.InitialDelay) * math.Pow(cfg.Multiplier, float64(attempt-1))
	if base > float64(cfg.MaxDelay) {
		base = float64(cfg.MaxDelay)
	}

	jitter := randomJitterFraction() * base * 0.3
	delay := time.Duration(base + jitter)

	return delay
}

func randomJitterFraction() float64 {
	var n uint64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &n); err != nil {
		return 0
	}
	return float64(n) / float64(math.MaxUint64)
}