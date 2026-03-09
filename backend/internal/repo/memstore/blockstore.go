package memstore

import (
	"context"
	"sync"
	"time"

	"coffee-consortium/backend/internal/ledger"
)

type BlockStore struct {
	mu     sync.RWMutex
	blocks []ledger.Block
}

func NewBlockStore() *BlockStore {
	return &BlockStore{blocks: make([]ledger.Block, 0, 64)}
}

func (bs *BlockStore) Append(_ context.Context, tx ledger.Transaction) (ledger.Block, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	prev := ""
	if n := len(bs.blocks); n > 0 {
		prev = bs.blocks[n-1].Hash
	}

	b := ledger.Block{
		Index:     len(bs.blocks),
		PrevHash:  prev,
		CreatedAt: time.Now().UTC(),
		Tx:        []ledger.Transaction{tx},
	}
	if err := b.ComputeHash(); err != nil {
		return ledger.Block{}, err
	}
	bs.blocks = append(bs.blocks, b)
	return b, nil
}

func (bs *BlockStore) List(_ context.Context, limit int) ([]ledger.Block, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	if limit <= 0 || limit > len(bs.blocks) {
		limit = len(bs.blocks)
	}
	out := make([]ledger.Block, limit)
	copy(out, bs.blocks[len(bs.blocks)-limit:])
	return out, nil
}

