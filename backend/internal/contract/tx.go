package contract

const (
	TxCreateOrder       = "CREATE_EXPORT_ORDER"
	TxAcceptOrder       = "ACCEPT_EXPORT_ORDER"
	TxIssueLC           = "ISSUE_LETTER_OF_CREDIT"
	TxApproveCustoms    = "APPROVE_CUSTOMS_CLEARANCE"
	TxCreateShipment    = "CREATE_SHIPMENT"
	TxUpdateShipment    = "UPDATE_SHIPMENT_STATUS"
	TxConfirmDelivery   = "CONFIRM_DELIVERY"
	TxReleasePayment    = "RELEASE_PAYMENT"
)

type CreateOrderPayload struct {
	BuyerID      string  `json:"buyerId"`
	CoffeeGrade  string  `json:"coffeeGrade"`
	QuantityKg   int     `json:"quantityKg"`
	UnitPriceUSD float64 `json:"unitPriceUsd"`
}

type AcceptOrderPayload struct {
	OrderID string `json:"orderId"`
}

type IssueLCPayload struct {
	OrderID   string  `json:"orderId"`
	AmountUSD float64 `json:"amountUsd"`
}

type ApproveCustomsPayload struct {
	OrderID string `json:"orderId"`
	Notes   string `json:"notes"`
}

type CreateShipmentPayload struct {
	OrderID    string `json:"orderId"`
	TrackingNo string `json:"trackingNo"`
}

type UpdateShipmentPayload struct {
	ShipmentID string `json:"shipmentId"`
	Status     string `json:"status"`
	Location   string `json:"location"`
}

type ConfirmDeliveryPayload struct {
	OrderID string `json:"orderId"`
}

type ReleasePaymentPayload struct {
	OrderID string `json:"orderId"`
}

