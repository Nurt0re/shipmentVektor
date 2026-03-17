package domain

import "time"

type Shipment struct {
	ID              string
	ReferenceNumber int32
	Origin          string
	Destination     string
	Status          Status
	Events          []Event
	Cost            float64
	DriverRevenue   float64
}

func NewShipment(id string, referenceNumber int32, origin, destination string, cost float64) *Shipment {
	return &Shipment{
		ID:              id,
		ReferenceNumber: referenceNumber,
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
