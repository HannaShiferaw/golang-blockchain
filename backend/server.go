package main

import (
	"coffee-demo/blockchain"
	"coffee-demo/contract"
	"coffee-demo/pki"
	"coffee-demo/state"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	r := gin.Default()

	// Create genesis block
	blockchain.CreateGenesisBlock()

	// Seed actors
	for _, role := range []string{"Exporter", "Buyer", "Bank", "Customs"} {
		kp, _ := pki.GenerateKeyPair()
		state.AddActor(&state.Actor{
			ID:         role,
			Name:       role,
			Role:       role,
			PublicKey:  kp.PublicKey,
			PrivateKey: kp.PrivateKey,
		})
	}

	r.GET("/api/actors", func(c *gin.Context) {
		c.JSON(200, state.GetAllActors())
	})

	r.GET("/api/batches", func(c *gin.Context) {
		c.JSON(200, state.GetAllBatches())
	})

	r.GET("/api/ledger", func(c *gin.Context) {
		c.JSON(200, blockchain.GetLedger())
	})

	r.POST("/api/:actor/action", func(c *gin.Context) {
		actorRole := c.Param("actor")
		var body map[string]interface{}
		c.BindJSON(&body)
		batchID := body["batchID"].(string)

		if !contract.ValidateTransition(batchID, actorRole) {
			c.JSON(400, gin.H{"error": "Invalid transition"})
			return
		}

		// Get actor
		actor := state.GetActorByRole(actorRole)

		txDataBytes, _ := json.Marshal(body)
		txData := string(txDataBytes)
		sig, _ := pki.Sign(actor.PrivateKey.(*pki.KeyPair).PrivateKey, txData)

		// Previous hash
		prevBlock := blockchain.GetLedger()[len(blockchain.GetLedger())-1]
		newBlock := blockchain.Block{
			Timestamp:       time.Now().String(),
			Actor:           actorRole,
			BatchID:         batchID,
			TransactionData: txData,
			PreviousHash:    prevBlock.Hash,
			Signature:       sig,
		}
		blockchain.AddBlock(newBlock)

		// Update batch status
		batch := state.GetBatch(batchID)
		if batch == nil {
			batch = &state.Batch{BatchID: batchID}
		}
		batch.Status = contract.GetNextStatus(actorRole)
		state.UpsertBatch(batch)

		c.JSON(200, gin.H{"success": true, "hash": newBlock.Hash})
	})

	r.POST("/api/tamper", func(c *gin.Context) {
		var body struct {
			Index   int    `json:"index"`
			NewData string `json:"newData"`
		}
		c.BindJSON(&body)
		blockchain.TamperBlock(body.Index, body.NewData)
		c.JSON(200, gin.H{"success": true})
	})

	fmt.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}