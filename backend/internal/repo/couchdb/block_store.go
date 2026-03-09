package couchdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"coffee-consortium/backend/internal/ledger"
)

type BlockStore struct {
	c *Client
}

func NewBlockStore(c *Client) *BlockStore {
	return &BlockStore{c: c}
}

func (bs *BlockStore) Append(ctx context.Context, tx ledger.Transaction) (ledger.Block, error) {
	const metaID = "meta:ledger"

	for attempt := 0; attempt < 3; attempt++ {
		var meta map[string]any
		found, _, err := bs.c.Get(ctx, metaID, &meta)
		if err != nil {
			return ledger.Block{}, err
		}
		if !found {
			meta = map[string]any{"_id": metaID, "lastIndex": -1, "lastHash": ""}
		}

		lastIdx := toInt(meta["lastIndex"], -1)
		lastHash, _ := meta["lastHash"].(string)

		b := ledger.Block{
			Index:     lastIdx + 1,
			PrevHash:  lastHash,
			CreatedAt: time.Now().UTC(),
			Tx:        []ledger.Transaction{tx},
		}
		if err := b.ComputeHash(); err != nil {
			return ledger.Block{}, err
		}

		// Store block doc.
		blockID := fmt.Sprintf("block:%09d", b.Index)
		doc := map[string]any{
			"_id":       blockID,
			"index":     b.Index,
			"prevHash":  b.PrevHash,
			"createdAt": b.CreatedAt,
			"tx":        b.Tx,
			"hash":      b.Hash,
		}
		if _, err := bs.c.Put(ctx, blockID, doc); err != nil {
			// If block already exists, meta is stale (rare); retry.
			if strings.Contains(err.Error(), "409") {
				continue
			}
			return ledger.Block{}, err
		}

		// Update meta.
		meta["lastIndex"] = b.Index
		meta["lastHash"] = b.Hash
		if _, err := bs.c.Put(ctx, metaID, meta); err != nil {
			// If meta update conflicts, retry (block is already stored; we can advance meta next attempt).
			if strings.Contains(err.Error(), "409") {
				continue
			}
			return ledger.Block{}, err
		}

		return b, nil
	}

	return ledger.Block{}, fmt.Errorf("append block failed after retries")
}

func (bs *BlockStore) List(ctx context.Context, limit int) ([]ledger.Block, error) {
	if limit <= 0 {
		limit = 50
	}
	q := url.Values{}
	q.Set("include_docs", "true")
	q.Set("descending", "true")
	q.Set("limit", strconv.Itoa(limit))
	q.Set("startkey", `"block:\ufff0"`)
	q.Set("endkey", `"block:"`)

	var res struct {
		Rows []struct {
			ID  string          `json:"id"`
			Doc json.RawMessage `json:"doc"`
		} `json:"rows"`
	}
	if err := bs.c.AllDocs(ctx, q, &res); err != nil {
		return nil, err
	}

	out := make([]ledger.Block, 0, len(res.Rows))
	for _, row := range res.Rows {
		var b ledger.Block
		if err := json.Unmarshal(row.Doc, &b); err != nil {
			continue
		}
		out = append(out, b)
	}
	return out, nil
}

func toInt(v any, def int) int {
	switch x := v.(type) {
	case float64:
		return int(x)
	case int:
		return x
	case int64:
		return int(x)
	default:
		return def
	}
}

