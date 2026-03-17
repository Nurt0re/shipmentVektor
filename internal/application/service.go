package application

import (
	"context"
	"shipment/internal/domain"
	"shipment/internal/ports/inbound"
	"shipment/internal/ports/outbound"
	"time"
)

type ShipmentService struct {
	repo outbound.ShipmentRepository
}

func NewShipmentService(repo outbound.ShipmentRepository) *ShipmentService {
	return &ShipmentService{repo: repo}
}

var _ inbound.ShipmentUseCase = (*ShipmentService)(nil)

func (s *ShipmentService) CreateShipment(
	ctx context.Context,
	id string,
	reference int32,
	origin string,
	destination string,
	cost float64,
	revenue float64,
) (*domain.Shipment, error) {

	shipment := domain.NewShipment(id, reference, origin, destination, cost, revenue)
	err := s.repo.Save(ctx, shipment)
	if err != nil {
		return nil, err
	}
	return shipment, nil
}

func (s *ShipmentService) GetShipment(ctx context.Context, id string) (*domain.Shipment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ShipmentService) AddEvent(ctx context.Context, id string, status domain.Status) error {
	shipment, err := s.repo.GetByID(ctx, id)
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

	return s.repo.Save(ctx, shipment)
}

func (s *ShipmentService) GetHistory(ctx context.Context, id string) ([]domain.Event, error) {
	shipment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return shipment.Events, nil
}
