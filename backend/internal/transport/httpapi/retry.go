package httpapi

import (
	"context"
	"time"
)

func retry(ctx context.Context, attempts int, delay time.Duration, fn func() error) error {
	var last error
	for i := 0; i < attempts; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := fn(); err == nil {
			return nil
		} else {
			last = err
		}
		time.Sleep(delay)
	}
	return last
}

