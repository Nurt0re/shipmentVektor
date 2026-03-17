package domain

import (
	"fmt"
	"time"
)

type Event struct {
	Status    Status
	Timestamp time.Time
}

func (s *Shipment) ApplyEvent(event Event) error {

	if !s.CanUpdate(event.Status) {
		return fmt.Errorf("invalid status update")
	}

	s.Status = event.Status
	s.Events = append(s.Events, event)
	return nil
}
