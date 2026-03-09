package ledgerSvc

import (
	"context"

	"coffee-consortium/backend/internal/ledger"
)

type TxIndexer interface {
	Index(ctx context.Context, b ledger.Block, tx ledger.Transaction) error
}

