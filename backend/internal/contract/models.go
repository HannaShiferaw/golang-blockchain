package contract

import (
	"time"

	"coffee-consortium/backend/internal/domain"
)

type ExportOrderStatus string

const (
	OrderCreated       ExportOrderStatus = "CREATED"
	OrderAccepted      ExportOrderStatus = "ACCEPTED"
	OrderLCIssued      ExportOrderStatus = "LC_ISSUED"
	OrderCustomsCleared ExportOrderStatus = "CUSTOMS_CLEARED"
	OrderInTransit     ExportOrderStatus = "IN_TRANSIT"
	OrderDelivered     ExportOrderStatus = "DELIVERED"
	OrderSettled       ExportOrderStatus = "SETTLED"
)

type ExportOrder struct {
	ID           string            `json:"id"`
	ExporterID   string            `json:"exporterId"`
	BuyerID      string            `json:"buyerId"`
	CoffeeGrade  string            `json:"coffeeGrade"`
	QuantityKg   int               `json:"quantityKg"`
	UnitPriceUSD float64           `json:"unitPriceUsd"`
	TotalUSD     float64           `json:"totalUsd"`
	Status       ExportOrderStatus `json:"status"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
}

type LetterOfCredit struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"orderId"`
	BankID    string    `json:"bankId"`
	AmountUSD float64   `json:"amountUsd"`
	IssuedAt  time.Time `json:"issuedAt"`
}

type CustomsClearance struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"orderId"`
	CustomsID string    `json:"customsId"`
	ApprovedAt time.Time `json:"approvedAt"`
	Notes     string    `json:"notes"`
}

type ShipmentStatus string

const (
	ShipmentCreated   ShipmentStatus = "CREATED"
	ShipmentPickedUp  ShipmentStatus = "PICKED_UP"
	ShipmentExported  ShipmentStatus = "EXPORTED"
	ShipmentArrived   ShipmentStatus = "ARRIVED"
	ShipmentDelivered ShipmentStatus = "DELIVERED"
)

type Shipment struct {
	ID         string         `json:"id"`
	OrderID    string         `json:"orderId"`
	ShipperID  string         `json:"shipperId"`
	Status     ShipmentStatus `json:"status"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	Location   string         `json:"location"`
	TrackingNo string         `json:"trackingNo"`
}

type PaymentRelease struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"orderId"`
	BankID    string    `json:"bankId"`
	ReleasedAt time.Time `json:"releasedAt"`
}

type Actor = struct {
	IdentityID string
	Role       domain.Role
}

