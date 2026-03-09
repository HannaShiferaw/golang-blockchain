package ledger

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	Index     int           `json:"index"`
	PrevHash  string        `json:"prevHash"`
	CreatedAt time.Time     `json:"createdAt"`
	Tx        []Transaction `json:"tx"`
	Hash      string        `json:"hash"`
}

func (b *Block) ComputeHash() error {
	type hashable struct {
		Index     int           `json:"index"`
		PrevHash  string        `json:"prevHash"`
		CreatedAt time.Time     `json:"createdAt"`
		TxHashes  []string      `json:"txHashes"`
	}
	txh := make([]string, 0, len(b.Tx))
	for _, t := range b.Tx {
		if t.TxHash == "" {
			return fmt.Errorf("tx missing hash: %s", t.ID)
		}
		txh = append(txh, t.TxHash)
	}
	payload, err := json.Marshal(hashable{
		Index:     b.Index,
		PrevHash:  b.PrevHash,
		CreatedAt: b.CreatedAt.UTC(),
		TxHashes:  txh,
	})
	if err != nil {
		return err
	}
	sum := sha256.Sum256(payload)
	b.Hash = hex.EncodeToString(sum[:])
	return nil
}

