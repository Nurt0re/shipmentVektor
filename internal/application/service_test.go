package application

import (
	"context"
	"errors"
	"testing"

	"shipment/internal/domain"
)

type MockRepository struct {
	ships map[string]*domain.Shipment
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		ships: make(map[string]*domain.Shipment),
	}
}

func (m *MockRepository) Save(ctx context.Context, s *domain.Shipment) error {
	m.ships[s.ID] = s
	return nil
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*domain.Shipment, error) {
	if ship, ok := m.ships[id]; ok {
		return ship, nil
	}
	return nil, errors.New("shipment not found")
}

func TestCreateShipment(t *testing.T) {
	repo := NewMockRepository()
	service := NewShipmentService(repo)

	shipment, err := service.CreateShipment(context.Background(), "ship-001", 123, "NYC", "LA", 1500.0, 300.0)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if shipment.ID != "ship-001" {
		t.Errorf("expected ID ship-001, got %s", shipment.ID)
	}

	if shipment.Status != domain.StatusPending {
		t.Errorf("expected status PENDING, got %s", shipment.Status)
	}
}

func TestGetShipment(t *testing.T) {
	repo := NewMockRepository()
	service := NewShipmentService(repo)

	service.CreateShipment(context.Background(), "ship-001", 123, "NYC", "LA", 1500.0, 300.0)

	shipment, err := service.GetShipment(context.Background(), "ship-001")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if shipment.ID != "ship-001" {
		t.Errorf("expected ID ship-001, got %s", shipment.ID)
	}
}

func TestGetShipmentNotFound(t *testing.T) {
	repo := NewMockRepository()
	service := NewShipmentService(repo)

	_, err := service.GetShipment(context.Background(), "nonexistent")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestAddEvent(t *testing.T) {
	repo := NewMockRepository()
	service := NewShipmentService(repo)

	service.CreateShipment(context.Background(), "ship-001", 123, "NYC", "LA", 1500.0, 300.0)

	err := service.AddEvent(context.Background(), "ship-001", domain.StatusShipped)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	shipment, _ := service.GetShipment(context.Background(), "ship-001")
	if shipment.Status != domain.StatusShipped {
		t.Errorf("expected status SHIPPED, got %s", shipment.Status)
	}
}

func TestGetHistory(t *testing.T) {
	repo := NewMockRepository()
	service := NewShipmentService(repo)

	service.CreateShipment(context.Background(), "ship-001", 123, "NYC", "LA", 1500.0, 300.0)
	service.AddEvent(context.Background(), "ship-001", domain.StatusShipped)

	events, err := service.GetHistory(context.Background(), "ship-001")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}
