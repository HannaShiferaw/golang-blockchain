package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"coffee-consortium/backend/internal/contract"
	"coffee-consortium/backend/internal/ledger"
)

type StateReader interface {
	Get(ctx context.Context, key string) ([]byte, bool, error)
}

type Indexer struct {
	db    *DB
	state StateReader
}

func NewIndexer(db *DB, state StateReader) *Indexer {
	return &Indexer{db: db, state: state}
}

func (ix *Indexer) Index(ctx context.Context, b ledger.Block, tx ledger.Transaction) error {
	if ix == nil || ix.db == nil || ix.db.Pool == nil {
		return nil
	}
	payload := json.RawMessage(tx.Payload)
	_, err := ix.db.Pool.Exec(ctx, `
		INSERT INTO audit_tx (tx_id, tx_type, actor_id, actor_role, created_at, tx_hash, block_index, block_hash, payload)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (tx_id) DO NOTHING
	`, tx.ID, tx.Type, tx.Actor.IdentityID, string(tx.Actor.Role), tx.CreatedAt, tx.TxHash, b.Index, b.Hash, payload)
	if err != nil {
		return err
	}

	// Maintain order_index from world-state.
	if ix.state != nil {
		switch tx.Type {
		case contract.TxCreateOrder, contract.TxAcceptOrder, contract.TxIssueLC, contract.TxApproveCustoms, contract.TxCreateShipment, contract.TxConfirmDelivery, contract.TxReleasePayment:
			// Order id is the create tx id for CREATE_ORDER; for other tx types it's in payload.
			orderID, err := orderIDFromTx(tx)
			if err != nil {
				return nil
			}
			raw, found, err := ix.state.Get(ctx, contract.OrderKey(orderID))
			if err != nil || !found {
				return nil
			}
			var o contract.ExportOrder
			if err := json.Unmarshal(raw, &o); err != nil {
				return nil
			}
			return ix.upsertOrderIndex(ctx, o)
		}
	}

	return nil
}

func (ix *Indexer) upsertOrderIndex(ctx context.Context, o contract.ExportOrder) error {
	updated := o.UpdatedAt
	if updated.IsZero() {
		updated = time.Now().UTC()
	}
	_, err := ix.db.Pool.Exec(ctx, `
		INSERT INTO order_index(order_id, exporter_id, buyer_id, status, total_usd, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (order_id) DO UPDATE SET
			status=EXCLUDED.status,
			total_usd=EXCLUDED.total_usd,
			updated_at=EXCLUDED.updated_at
	`, o.ID, o.ExporterID, o.BuyerID, string(o.Status), o.TotalUSD, updated)
	return err
}

func orderIDFromTx(tx ledger.Transaction) (string, error) {
	if tx.Type == contract.TxCreateOrder {
		return tx.ID, nil
	}
	var m map[string]any
	if err := json.Unmarshal(tx.Payload, &m); err != nil {
		return "", err
	}
	if v, ok := m["orderId"].(string); ok && v != "" {
		return v, nil
	}
	return "", fmt.Errorf("no orderId in payload")
}

