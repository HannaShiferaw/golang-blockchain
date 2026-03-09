package domain

import "fmt"

type Role string

const (
	RoleExporter Role = "EXPORTER"
	RoleBuyer    Role = "BUYER"
	RoleBank     Role = "BANK"
	RoleCustoms  Role = "CUSTOMS"
	RoleShipment Role = "SHIPMENT"
)

func ParseRole(s string) (Role, error) {
	switch Role(s) {
	case RoleExporter, RoleBuyer, RoleBank, RoleCustoms, RoleShipment:
		return Role(s), nil
	default:
		return "", fmt.Errorf("invalid role: %q", s)
	}
}

