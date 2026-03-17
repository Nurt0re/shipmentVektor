package application

import (
	"shipment/internal/domain"
	"shipment/internal/ports/outbound"
	"time"
)

type ShipmentService struct {
	repo outbound.ShipmentRepository
}

func NewShipmentService(repo outbound.ShipmentRepository) *ShipmentService {
	return &ShipmentService{repo: repo}
}

func (s *ShipmentService) CreateShipment(
	id string,
	reference int32,
	origin string,
	destination string,
	amount float64,
	revenue float64,
) (*domain.Shipment, error) {

	shipment := domain.NewShipment(id, reference, origin, destination, amount, revenue)
	err := s.repo.Save(shipment)
	if err != nil {
		return nil, err
	}
	return shipment, nil
}

func (s *ShipmentService) GetShipment(id string) (*domain.Shipment, error) {
	return s.repo.GetByID(id)
}

func (s *ShipmentService) AddEvent(id string, status domain.Status) error {
	shipment, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	event := domain.Event{
		Status:    status,
		Timestamp: time.Now(),
	}

	err = shipment.ApplyEvent(event)
	if err != nil {
		return err
	}

	return s.repo.Save(shipment)
}

func (s *ShipmentService) GetHistory(id string) ([]domain.Event, error) {
	shipment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return shipment.Events, nil
}
