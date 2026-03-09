package couchdb

import (
	"context"
	"encoding/json"
	"fmt"
)

type StateStore struct {
	c *Client
}

func NewStateStore(c *Client) *StateStore {
	return &StateStore{c: c}
}

func (s *StateStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
	var m map[string]any
	found, _, err := s.c.Get(ctx, key, &m)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}
	delete(m, "_id")
	delete(m, "_rev")
	raw, err := json.Marshal(m)
	if err != nil {
		return nil, false, fmt.Errorf("marshal state: %w", err)
	}
	return raw, true, nil
}

func (s *StateStore) Put(ctx context.Context, key string, raw []byte) error {
	// Merge current rev (if exists) to make PUT idempotent.
	var existing map[string]any
	found, _, err := s.c.Get(ctx, key, &existing)
	if err != nil {
		return err
	}

	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return fmt.Errorf("invalid state json: %w", err)
	}
	m["_id"] = key
	if found {
		if rev, ok := existing["_rev"].(string); ok && rev != "" {
			m["_rev"] = rev
		}
	}

	_, err = s.c.Put(ctx, key, m)
	return err
}

