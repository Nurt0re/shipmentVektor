package inbound

import "shipment/internal/domain"

type ShipmentUseCase interface {
	CreateShipment(
		id string,
		reference int32,
		origin string,
		destination string,
		amount float64,
		revenue float64,
	) (*domain.Shipment, error)

	GetShipment(id string) (*domain.Shipment, error)

	AddEvent(id string, status domain.Status) error

	GetHistory(id string) ([]domain.Event, error)
}
