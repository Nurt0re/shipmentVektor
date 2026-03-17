package domain

import (
	"testing"
	"time"
)

func TestValidStatusUpdate(t *testing.T) {
	shipment := NewShipment("1", 123, "New York", "Los Angeles", 100.0)

	err := shipment.AddEvent(Event{Status: StatusShipped, Timestamp: time.Now()})
	if err != nil {
		t.Errorf("expected valid status update, got error: %v", err)
	}
}

func TestInvalidStatusUpdate(t *testing.T) {
	shipment := NewShipment("1", 123, "New York", "Los Angeles", 100.0)

	err := shipment.AddEvent(Event{Status: StatusDelivered, Timestamp: time.Now()})
	if err == nil {
		t.Errorf("expected error for invalid status update, got nil")
	}
}
