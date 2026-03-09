package ledgerSvc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"coffee-consortium/backend/internal/contract"
	"coffee-consortium/backend/internal/domain"
	"coffee-consortium/backend/internal/ledger"
	"coffee-consortium/backend/internal/pki"
	"coffee-consortium/backend/internal/service/identity"
)

type Service struct {
	ids   *identity.Service
	state contract.StateStore
	bs    BlockStore
	ix    TxIndexer
}

func New(ids *identity.Service, state contract.StateStore, bs BlockStore, ix TxIndexer) *Service {
	return &Service{
		ids:   ids,
		state: state,
		bs:    bs,
		ix:    ix,
	}
}

func (s *Service) Submit(ctx context.Context, actorID string, txType string, payload any) (ledger.Transaction, error) {
	it, err := s.ids.GetIdentity(actorID)
	if err != nil {
		return ledger.Transaction{}, err
	}
	if it.PrivateKeyPEM == "" {
		return ledger.Transaction{}, fmt.Errorf("identity missing private key")
	}

	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return ledger.Transaction{}, fmt.Errorf("marshal payload: %w", err)
	}

	tx := ledger.Transaction{
		ID:        uuid.NewString(),
		Type:      txType,
		CreatedAt: time.Now().UTC(),
		Actor: ledger.Actor{
			IdentityID: it.ID,
			Name:       it.Name,
			Role:       it.Role,
			CertPEM:    it.CertPEM,
		},
		Payload: rawPayload,
	}
	if err := tx.ComputeHash(); err != nil {
		return ledger.Transaction{}, err
	}

	signBytes, err := tx.SigningBytes()
	if err != nil {
		return ledger.Transaction{}, err
	}
	sig, err := pki.SignPayload(it.PrivateKeyPEM, signBytes)
	if err != nil {
		return ledger.Transaction{}, err
	}
	tx.Signature = sig

	// Verify before apply.
	if err := pki.VerifyPayload(it.CertPEM, signBytes, sig); err != nil {
		return ledger.Transaction{}, fmt.Errorf("signature verification failed: %w", err)
	}

	// Apply contract changes.
	if err := contract.Apply(ctx, s.state, tx); err != nil {
		return ledger.Transaction{}, err
	}

	if s.bs != nil {
		b, err := s.bs.Append(ctx, tx)
		if err != nil {
			return ledger.Transaction{}, err
		}
		if s.ix != nil {
			if err := s.ix.Index(ctx, b, tx); err != nil {
				return ledger.Transaction{}, err
			}
		}
	}
	return tx, nil
}

func (s *Service) Blocks(ctx context.Context, limit int) ([]ledger.Block, error) {
	if s.bs == nil {
		return nil, nil
	}
	return s.bs.List(ctx, limit)
}

func (s *Service) StateGet(ctx context.Context, key string) ([]byte, bool, error) {
	return s.state.Get(ctx, key)
}

// Convenience for UI (not for production).
func (s *Service) EnsureRole(actorID string, r domain.Role) error {
	it, err := s.ids.GetIdentity(actorID)
	if err != nil {
		return err
	}
	if it.Role != r {
		return fmt.Errorf("actor role must be %s", r)
	}
	return nil
}

