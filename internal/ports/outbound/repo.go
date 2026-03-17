package outbound

import "shipment/internal/domain"

type ShipmentRepository interface {
	Save(shipment *domain.Shipment) error
	GetByID(id string) (*domain.Shipment, error)
}
