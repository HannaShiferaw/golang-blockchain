package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

type Block struct {
	BlockIndex      int
	Timestamp       string
	Actor           string
	BatchID         string
	TransactionData string
	PreviousHash    string
	Hash            string
	Signature       string
}

var Ledger []Block

func CalculateHash(block Block) string {
	record := string(block.BlockIndex) + block.Timestamp + block.Actor + block.TransactionData + block.PreviousHash
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

func AddBlock(newBlock Block) (Block, error) {
	if len(Ledger) > 0 {
		lastBlock := Ledger[len(Ledger)-1]
		if newBlock.PreviousHash != lastBlock.Hash {
			return Block{}, errors.New("previous hash does not match")
		}
		newBlock.BlockIndex = lastBlock.BlockIndex + 1
	} else {
		newBlock.BlockIndex = 1
		newBlock.PreviousHash = "0"
	}

	newBlock.Hash = CalculateHash(newBlock)
	Ledger = append(Ledger, newBlock)
	return newBlock, nil
}

func GetLedger() []Block {
	return Ledger
}

// Tamper demo
func TamperBlock(index int, newData string) error {
	if index <= 0 || index > len(Ledger) {
		return errors.New("invalid block index")
	}
	Ledger[index-1].TransactionData = newData
	Ledger[index-1].Hash = CalculateHash(Ledger[index-1])
	return nil
}

// Create genesis block
func CreateGenesisBlock() {
	genesis := Block{
		BlockIndex:      0,
		Timestamp:       time.Now().String(),
		Actor:           "Genesis",
		BatchID:         "",
		TransactionData: "Genesis Block",
		PreviousHash:    "0",
		Signature:       "",
		Hash:            "",
	}
	genesis.Hash = CalculateHash(genesis)
	Ledger = append(Ledger, genesis)
}