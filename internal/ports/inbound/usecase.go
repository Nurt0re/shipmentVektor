package inbound

import (
	"context"
	"shipment/internal/domain"
)

type ShipmentUseCase interface {
	CreateShipment(
		ctx context.Context,
		id string,
		reference int32,
		origin string,
		destination string,
		cost float64,
		revenue float64,
	) (*domain.Shipment, error)

	GetShipment(ctx context.Context, id string) (*domain.Shipment, error)

	AddEvent(ctx context.Context, id string, status domain.Status) error

	GetHistory(ctx context.Context, id string) ([]domain.Event, error)
}
