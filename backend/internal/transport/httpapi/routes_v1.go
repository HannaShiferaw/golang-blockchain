package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"coffee-consortium/backend/internal/contract"
	"coffee-consortium/backend/internal/domain"
	ledgerSvc "coffee-consortium/backend/internal/service/ledger"
	"coffee-consortium/backend/internal/service/identity"
)

func registerV1Routes(rg *gin.RouterGroup, ids *identity.Service, led *ledgerSvc.Service) {
	rg.GET("/pki/ca", func(c *gin.Context) {
		pem, err := ids.RootCertPEM()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"rootCertPem": pem})
	})

	rg.GET("/identities", func(c *gin.Context) {
		c.JSON(200, gin.H{"items": ids.ListIdentities()})
	})

	rg.POST("/identities", func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
			Role string `json:"role"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid json"})
			return
		}
		role, err := domain.ParseRole(req.Role)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		it, err := ids.CreateIdentity(req.Name, role)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		it.PrivateKeyPEM = ""
		c.JSON(201, it)
	})

	rg.GET("/blocks", func(c *gin.Context) {
		items, err := led.Blocks(c.Request.Context(), 100)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"items": items})
	})

	// --- Order workflow endpoints (demo-friendly, role-enforced) ---

	rg.POST("/orders", func(c *gin.Context) {
		actorID := c.GetHeader("X-Actor-Id")
		if err := led.EnsureRole(actorID, domain.RoleExporter); err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		var p contract.CreateOrderPayload
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(400, gin.H{"error": "invalid json"})
			return
		}
		tx, err := led.Submit(c.Request.Context(), actorID, contract.TxCreateOrder, p)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"orderId": tx.ID, "tx": tx})
	})

	rg.POST("/orders/:orderId/accept", func(c *gin.Context) {
		actorID := c.GetHeader("X-Actor-Id")
		if err := led.EnsureRole(actorID, domain.RoleBuyer); err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		orderID := c.Param("orderId")
		tx, err := led.Submit(c.Request.Context(), actorID, contract.TxAcceptOrder, contract.AcceptOrderPayload{OrderID: orderID})
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"tx": tx})
	})

	rg.POST("/orders/:orderId/lc", func(c *gin.Context) {
		actorID := c.GetHeader("X-Actor-Id")
		if err := led.EnsureRole(actorID, domain.RoleBank); err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		var req struct {
			AmountUSD float64 `json:"amountUsd"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid json"})
			return
		}
		orderID := c.Param("orderId")
		tx, err := led.Submit(c.Request.Context(), actorID, contract.TxIssueLC, contract.IssueLCPayload{OrderID: orderID, AmountUSD: req.AmountUSD})
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"tx": tx})
	})

	rg.POST("/orders/:orderId/customs-approve", func(c *gin.Context) {
		actorID := c.GetHeader("X-Actor-Id")
		if err := led.EnsureRole(actorID, domain.RoleCustoms); err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		var req struct {
			Notes string `json:"notes"`
		}
		_ = c.ShouldBindJSON(&req)
		orderID := c.Param("orderId")
		tx, err := led.Submit(c.Request.Context(), actorID, contract.TxApproveCustoms, contract.ApproveCustomsPayload{OrderID: orderID, Notes: req.Notes})
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"tx": tx})
	})

	rg.POST("/orders/:orderId/shipments", func(c *gin.Context) {
		actorID := c.GetHeader("X-Actor-Id")
		if err := led.EnsureRole(actorID, domain.RoleShipment); err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		var req struct {
			TrackingNo string `json:"trackingNo"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid json"})
			return
		}
		orderID := c.Param("orderId")
		tx, err := led.Submit(c.Request.Context(), actorID, contract.TxCreateShipment, contract.CreateShipmentPayload{OrderID: orderID, TrackingNo: req.TrackingNo})
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"shipmentId": tx.ID, "tx": tx})
	})

	rg.POST("/shipments/:shipmentId/status", func(c *gin.Context) {
		actorID := c.GetHeader("X-Actor-Id")
		if err := led.EnsureRole(actorID, domain.RoleShipment); err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		var req struct {
			Status   string `json:"status"`
			Location string `json:"location"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid json"})
			return
		}
		shID := c.Param("shipmentId")
		tx, err := led.Submit(c.Request.Context(), actorID, contract.TxUpdateShipment, contract.UpdateShipmentPayload{
			ShipmentID: shID,
			Status:     req.Status,
			Location:   req.Location,
		})
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"tx": tx})
	})

	rg.POST("/orders/:orderId/confirm-delivery", func(c *gin.Context) {
		actorID := c.GetHeader("X-Actor-Id")
		if err := led.EnsureRole(actorID, domain.RoleBuyer); err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		orderID := c.Param("orderId")
		tx, err := led.Submit(c.Request.Context(), actorID, contract.TxConfirmDelivery, contract.ConfirmDeliveryPayload{OrderID: orderID})
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"tx": tx})
	})

	rg.POST("/orders/:orderId/release-payment", func(c *gin.Context) {
		actorID := c.GetHeader("X-Actor-Id")
		if err := led.EnsureRole(actorID, domain.RoleBank); err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		orderID := c.Param("orderId")
		tx, err := led.Submit(c.Request.Context(), actorID, contract.TxReleasePayment, contract.ReleasePaymentPayload{OrderID: orderID})
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"tx": tx})
	})

	// State queries (demo) - returns raw JSON of state docs by key.
	rg.GET("/state", func(c *gin.Context) {
		key := c.Query("key")
		if key == "" {
			c.JSON(400, gin.H{"error": "missing key"})
			return
		}
		raw, found, err := led.StateGet(c.Request.Context(), key)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if !found {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		var anyJSON any
		if err := json.Unmarshal(raw, &anyJSON); err != nil {
			c.JSON(500, gin.H{"error": "corrupt json"})
			return
		}
		c.JSON(200, anyJSON)
	})
}

