package ledgerSvc

import (
	"context"

	"coffee-consortium/backend/internal/ledger"
)

type BlockStore interface {
	Append(ctx context.Context, tx ledger.Transaction) (ledger.Block, error)
	List(ctx context.Context, limit int) ([]ledger.Block, error)
}

