package domain

import (
	"testing"
	"time"
)

func TestValidStatusUpdate(t *testing.T) {
	shipment := NewShipment("1", 123, "New York", "Los Angeles", 100.0, 50.0)

	err := shipment.ApplyEvent(Event{Status: StatusShipped, Timestamp: time.Now()})
	if err != nil {
		t.Errorf("expected valid status update, got error: %v", err)
	}
}

func TestInvalidStatusUpdate(t *testing.T) {
	shipment := NewShipment("1", 123, "New York", "Los Angeles", 100.0, 50.0)

	err := shipment.ApplyEvent(Event{Status: StatusDelivered, Timestamp: time.Now()})
	if err == nil {
		t.Errorf("expected error for invalid status update, got nil")
	}
}

func TestSequentialStatusUpdates(t *testing.T) {
	shipment := NewShipment("1", 123, "New York", "Los Angeles", 100.0, 50.0)

	statuses := []Status{StatusShipped, StatusOnTheWay, StatusDelivered}
	for _, status := range statuses {
		err := shipment.ApplyEvent(Event{Status: status, Timestamp: time.Now()})
		if err != nil {
			t.Errorf("failed to apply status %s: %v", status, err)
		}
	}

	if shipment.Status != StatusDelivered {
		t.Errorf("expected final status DELIVERED, got %s", shipment.Status)
	}
}

func TestNewShipmentInitialization(t *testing.T) {
	shipment := NewShipment("ship-001", 456, "NYC", "LA", 5000.0, 1000.0)

	if shipment.ID != "ship-001" {
		t.Errorf("expected ID ship-001, got %s", shipment.ID)
	}

	if shipment.Reference != 456 {
		t.Errorf("expected reference 456, got %d", shipment.Reference)
	}

	if shipment.Status != StatusPending {
		t.Errorf("expected status PENDING, got %s", shipment.Status)
	}

	if len(shipment.Events) != 1 {
		t.Errorf("expected 1 initial event, got %d", len(shipment.Events))
	}
}

func TestCanUpdate(t *testing.T) {
	shipment := NewShipment("1", 123, "New York", "Los Angeles", 100.0, 50.0)

	if !shipment.CanUpdate(StatusShipped) {
		t.Errorf("expected PENDING -> SHIPPED to be valid")
	}

	if shipment.CanUpdate(StatusOnTheWay) {
		t.Errorf("expected PENDING -> ON_THE_WAY to be invalid")
	}

	shipment.ApplyEvent(Event{Status: StatusShipped, Timestamp: time.Now()})

	if !shipment.CanUpdate(StatusOnTheWay) {
		t.Errorf("expected SHIPPED -> ON_THE_WAY to be valid")
	}

	if shipment.CanUpdate(StatusPending) {
		t.Errorf("expected SHIPPED -> PENDING to be invalid")
	}
}
