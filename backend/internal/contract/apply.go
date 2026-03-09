package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"coffee-consortium/backend/internal/domain"
	"coffee-consortium/backend/internal/ledger"
)

func Apply(ctx context.Context, st StateStore, tx ledger.Transaction) error {
	switch tx.Type {
	case TxCreateOrder:
		return applyCreateOrder(ctx, st, tx)
	case TxAcceptOrder:
		return applyAcceptOrder(ctx, st, tx)
	case TxIssueLC:
		return applyIssueLC(ctx, st, tx)
	case TxApproveCustoms:
		return applyApproveCustoms(ctx, st, tx)
	case TxCreateShipment:
		return applyCreateShipment(ctx, st, tx)
	case TxUpdateShipment:
		return applyUpdateShipment(ctx, st, tx)
	case TxConfirmDelivery:
		return applyConfirmDelivery(ctx, st, tx)
	case TxReleasePayment:
		return applyReleasePayment(ctx, st, tx)
	default:
		return fmt.Errorf("unknown tx type: %s", tx.Type)
	}
}

func applyCreateOrder(ctx context.Context, st StateStore, tx ledger.Transaction) error {
	if tx.Actor.Role != domain.RoleExporter {
		return fmt.Errorf("only exporter can create order")
	}
	var p CreateOrderPayload
	if err := json.Unmarshal(tx.Payload, &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	if p.BuyerID == "" || p.CoffeeGrade == "" || p.QuantityKg <= 0 || p.UnitPriceUSD <= 0 {
		return fmt.Errorf("invalid order fields")
	}

	now := time.Now().UTC()
	id := tx.ID
	order := ExportOrder{
		ID:           id,
		ExporterID:   tx.Actor.IdentityID,
		BuyerID:      p.BuyerID,
		CoffeeGrade:  p.CoffeeGrade,
		QuantityKg:   p.QuantityKg,
		UnitPriceUSD: p.UnitPriceUSD,
		TotalUSD:     float64(p.QuantityKg) * p.UnitPriceUSD,
		Status:       OrderCreated,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	raw, _ := json.Marshal(order)
	return st.Put(ctx, OrderKey(id), raw)
}

func applyAcceptOrder(ctx context.Context, st StateStore, tx ledger.Transaction) error {
	if tx.Actor.Role != domain.RoleBuyer {
		return fmt.Errorf("only buyer can accept order")
	}
	var p AcceptOrderPayload
	if err := json.Unmarshal(tx.Payload, &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	order, err := getOrder(ctx, st, p.OrderID)
	if err != nil {
		return err
	}
	if order.BuyerID != tx.Actor.IdentityID {
		return fmt.Errorf("buyer mismatch")
	}
	if order.Status != OrderCreated {
		return fmt.Errorf("order must be CREATED")
	}
	order.Status = OrderAccepted
	order.UpdatedAt = time.Now().UTC()
	raw, _ := json.Marshal(order)
	return st.Put(ctx, OrderKey(order.ID), raw)
}

func applyIssueLC(ctx context.Context, st StateStore, tx ledger.Transaction) error {
	if tx.Actor.Role != domain.RoleBank {
		return fmt.Errorf("only bank can issue LC")
	}
	var p IssueLCPayload
	if err := json.Unmarshal(tx.Payload, &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	order, err := getOrder(ctx, st, p.OrderID)
	if err != nil {
		return err
	}
	if order.Status != OrderAccepted {
		return fmt.Errorf("order must be ACCEPTED")
	}
	if p.AmountUSD <= 0 || p.AmountUSD < order.TotalUSD {
		return fmt.Errorf("LC amount must cover order total")
	}

	lc := LetterOfCredit{
		ID:        tx.ID,
		OrderID:   order.ID,
		BankID:    tx.Actor.IdentityID,
		AmountUSD: p.AmountUSD,
		IssuedAt:  time.Now().UTC(),
	}
	rawLC, _ := json.Marshal(lc)
	if err := st.Put(ctx, LCKey(lc.ID), rawLC); err != nil {
		return err
	}

	order.Status = OrderLCIssued
	order.UpdatedAt = time.Now().UTC()
	raw, _ := json.Marshal(order)
	return st.Put(ctx, OrderKey(order.ID), raw)
}

func applyApproveCustoms(ctx context.Context, st StateStore, tx ledger.Transaction) error {
	if tx.Actor.Role != domain.RoleCustoms {
		return fmt.Errorf("only customs can approve clearance")
	}
	var p ApproveCustomsPayload
	if err := json.Unmarshal(tx.Payload, &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	order, err := getOrder(ctx, st, p.OrderID)
	if err != nil {
		return err
	}
	if order.Status != OrderLCIssued {
		return fmt.Errorf("order must be LC_ISSUED")
	}

	cc := CustomsClearance{
		ID:         tx.ID,
		OrderID:    order.ID,
		CustomsID:  tx.Actor.IdentityID,
		ApprovedAt: time.Now().UTC(),
		Notes:      p.Notes,
	}
	rawCC, _ := json.Marshal(cc)
	if err := st.Put(ctx, CustomsKey(cc.ID), rawCC); err != nil {
		return err
	}

	order.Status = OrderCustomsCleared
	order.UpdatedAt = time.Now().UTC()
	raw, _ := json.Marshal(order)
	return st.Put(ctx, OrderKey(order.ID), raw)
}

func applyCreateShipment(ctx context.Context, st StateStore, tx ledger.Transaction) error {
	if tx.Actor.Role != domain.RoleShipment {
		return fmt.Errorf("only shipment can create shipment")
	}
	var p CreateShipmentPayload
	if err := json.Unmarshal(tx.Payload, &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	order, err := getOrder(ctx, st, p.OrderID)
	if err != nil {
		return err
	}
	if order.Status != OrderCustomsCleared {
		return fmt.Errorf("order must be CUSTOMS_CLEARED")
	}

	sh := Shipment{
		ID:         tx.ID,
		OrderID:    order.ID,
		ShipperID:  tx.Actor.IdentityID,
		Status:     ShipmentCreated,
		UpdatedAt:  time.Now().UTC(),
		Location:   "ADDIS_ABABA",
		TrackingNo: p.TrackingNo,
	}
	rawSh, _ := json.Marshal(sh)
	if err := st.Put(ctx, ShipmentKey(sh.ID), rawSh); err != nil {
		return err
	}

	order.Status = OrderInTransit
	order.UpdatedAt = time.Now().UTC()
	raw, _ := json.Marshal(order)
	return st.Put(ctx, OrderKey(order.ID), raw)
}

func applyUpdateShipment(ctx context.Context, st StateStore, tx ledger.Transaction) error {
	if tx.Actor.Role != domain.RoleShipment {
		return fmt.Errorf("only shipment can update shipment")
	}
	var p UpdateShipmentPayload
	if err := json.Unmarshal(tx.Payload, &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	raw, found, err := st.Get(ctx, ShipmentKey(p.ShipmentID))
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("shipment not found")
	}
	var sh Shipment
	if err := json.Unmarshal(raw, &sh); err != nil {
		return fmt.Errorf("corrupt shipment state")
	}
	if sh.ShipperID != tx.Actor.IdentityID {
		return fmt.Errorf("shipper mismatch")
	}

	switch ShipmentStatus(p.Status) {
	case ShipmentPickedUp, ShipmentExported, ShipmentArrived, ShipmentDelivered:
		sh.Status = ShipmentStatus(p.Status)
	default:
		return fmt.Errorf("invalid shipment status")
	}
	if p.Location != "" {
		sh.Location = p.Location
	}
	sh.UpdatedAt = time.Now().UTC()

	raw2, _ := json.Marshal(sh)
	return st.Put(ctx, ShipmentKey(sh.ID), raw2)
}

func applyConfirmDelivery(ctx context.Context, st StateStore, tx ledger.Transaction) error {
	if tx.Actor.Role != domain.RoleBuyer {
		return fmt.Errorf("only buyer can confirm delivery")
	}
	var p ConfirmDeliveryPayload
	if err := json.Unmarshal(tx.Payload, &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	order, err := getOrder(ctx, st, p.OrderID)
	if err != nil {
		return err
	}
	if order.BuyerID != tx.Actor.IdentityID {
		return fmt.Errorf("buyer mismatch")
	}
	if order.Status != OrderInTransit {
		return fmt.Errorf("order must be IN_TRANSIT")
	}
	order.Status = OrderDelivered
	order.UpdatedAt = time.Now().UTC()
	raw, _ := json.Marshal(order)
	return st.Put(ctx, OrderKey(order.ID), raw)
}

func applyReleasePayment(ctx context.Context, st StateStore, tx ledger.Transaction) error {
	if tx.Actor.Role != domain.RoleBank {
		return fmt.Errorf("only bank can release payment")
	}
	var p ReleasePaymentPayload
	if err := json.Unmarshal(tx.Payload, &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	order, err := getOrder(ctx, st, p.OrderID)
	if err != nil {
		return err
	}
	if order.Status != OrderDelivered {
		return fmt.Errorf("order must be DELIVERED")
	}

	pay := PaymentRelease{
		ID:         tx.ID,
		OrderID:    order.ID,
		BankID:     tx.Actor.IdentityID,
		ReleasedAt: time.Now().UTC(),
	}
	rawPay, _ := json.Marshal(pay)
	if err := st.Put(ctx, PaymentKey(pay.ID), rawPay); err != nil {
		return err
	}

	order.Status = OrderSettled
	order.UpdatedAt = time.Now().UTC()
	raw, _ := json.Marshal(order)
	return st.Put(ctx, OrderKey(order.ID), raw)
}

func getOrder(ctx context.Context, st StateStore, orderID string) (ExportOrder, error) {
	if orderID == "" {
		return ExportOrder{}, fmt.Errorf("missing orderId")
	}
	raw, found, err := st.Get(ctx, OrderKey(orderID))
	if err != nil {
		return ExportOrder{}, err
	}
	if !found {
		return ExportOrder{}, fmt.Errorf("order not found")
	}
	var o ExportOrder
	if err := json.Unmarshal(raw, &o); err != nil {
		return ExportOrder{}, fmt.Errorf("corrupt order state")
	}
	return o, nil
}

