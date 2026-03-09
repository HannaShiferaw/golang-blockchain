package contract

const (
	orderKeyPrefix    = "state:order:"
	lcKeyPrefix       = "state:lc:"
	customsKeyPrefix  = "state:customs:"
	shipmentKeyPrefix = "state:shipment:"
	paymentKeyPrefix  = "state:payment:"
)

func OrderKey(orderID string) string    { return orderKeyPrefix + orderID }
func LCKey(id string) string           { return lcKeyPrefix + id }
func CustomsKey(id string) string      { return customsKeyPrefix + id }
func ShipmentKey(id string) string     { return shipmentKeyPrefix + id }
func PaymentKey(id string) string      { return paymentKeyPrefix + id }

