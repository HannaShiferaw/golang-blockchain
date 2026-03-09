package ledger

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"coffee-consortium/backend/internal/domain"
)

type Actor struct {
	IdentityID string      `json:"identityId"`
	Name       string      `json:"name"`
	Role       domain.Role `json:"role"`
	CertPEM    string      `json:"certPem"`
}

type Transaction struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	CreatedAt time.Time       `json:"createdAt"`
	Actor     Actor           `json:"actor"`
	Payload   json.RawMessage `json:"payload"`

	Signature string `json:"signature"` // base64 ASN.1 ECDSA signature of SigningBytes()
	TxHash    string `json:"txHash"`    // sha256 hex of SigningBytes()
}

func (t Transaction) SigningBytes() ([]byte, error) {
	// Keep this stable: do not include Signature or TxHash.
	type signable struct {
		ID        string          `json:"id"`
		Type      string          `json:"type"`
		CreatedAt time.Time       `json:"createdAt"`
		Actor     Actor           `json:"actor"`
		Payload   json.RawMessage `json:"payload"`
	}
	return json.Marshal(signable{
		ID:        t.ID,
		Type:      t.Type,
		CreatedAt: t.CreatedAt.UTC(),
		Actor:     t.Actor,
		Payload:   t.Payload,
	})
}

func (t *Transaction) ComputeHash() error {
	b, err := t.SigningBytes()
	if err != nil {
		return err
	}
	sum := sha256.Sum256(b)
	t.TxHash = hex.EncodeToString(sum[:])
	return nil
}

func (t Transaction) ValidateBasic() error {
	if t.ID == "" {
		return fmt.Errorf("missing tx id")
	}
	if t.Type == "" {
		return fmt.Errorf("missing tx type")
	}
	if t.CreatedAt.IsZero() {
		return fmt.Errorf("missing createdAt")
	}
	if t.Actor.IdentityID == "" || t.Actor.Name == "" || t.Actor.Role == "" {
		return fmt.Errorf("missing actor fields")
	}
	if len(t.Payload) == 0 {
		return fmt.Errorf("missing payload")
	}
	if t.Signature == "" {
		return fmt.Errorf("missing signature")
	}
	if t.TxHash == "" {
		return fmt.Errorf("missing txHash")
	}
	return nil
}

