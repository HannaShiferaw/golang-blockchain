package contract

import "context"

type StateStore interface {
	// Get returns raw JSON for a key.
	// found=false means key does not exist.
	Get(ctx context.Context, key string) (raw []byte, found bool, err error)
	Put(ctx context.Context, key string, raw []byte) error
}

