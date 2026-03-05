package contract

import "coffee-demo/state"

func ValidateTransition(batchID string, actorRole string) bool {
	b := state.GetBatch(batchID)
	if actorRole == "Exporter" {
		return b == nil
	}
	if b == nil {
		return false
	}
	switch actorRole {
	case "Buyer":
		return b.Status == "Pending Buyer Confirmation"
	case "Bank":
		return b.Status == "Pending Bank Payment"
	case "Customs":
		return b.Status == "Pending Customs Clearance"
	default:
		return false
	}
}

func GetNextStatus(actorRole string) string {
	switch actorRole {
	case "Exporter":
		return "Pending Buyer Confirmation"
	case "Buyer":
		return "Pending Bank Payment"
	case "Bank":
		return "Pending Customs Clearance"
	case "Customs":
		return "Export Completed"
	default:
		return "Unknown"
	}
}