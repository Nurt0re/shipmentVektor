package domain

import "time"

type Shipment struct {
	ID              string
	Reference int32
	Origin          string
	Destination     string
	Status          Status
	Events          []Event
	Cost            float64
	DriverRevenue   float64
}

func NewShipment(id string, reference int32, origin, destination string, cost, revenue float64) *Shipment {
	return &Shipment{
		ID:              id,
		Reference: reference,
		Origin:          origin,
		Destination:     destination,
		Status:          StatusPending,
		Cost:            cost,
		Events: []Event{
			{
				Status:    StatusPending,
				Timestamp: time.Now(),
			},
		},
		DriverRevenue: revenue,
	}
}

func (s *Shipment) CanUpdate(newStatus Status) bool {

	updates := map[Status][]Status{
		StatusPending:  {StatusShipped},
		StatusShipped:  {StatusOnTheWay},
		StatusOnTheWay: {StatusDelivered},
	}

	allowedStatuses, ok := updates[s.Status]
	if !ok {
		return false
	}

	for _, allowedStatus := range allowedStatuses {
		if newStatus == allowedStatus {
			return true
		}
	}

	return false
}
