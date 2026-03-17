package outbound

import (
	"context"
	"shipment/internal/domain"
)

type ShipmentRepository interface {
	Save(ctx context.Context, shipment *domain.Shipment) error
	GetByID(ctx context.Context, id string) (*domain.Shipment, error)
}
