package domain

import (
	"fmt"
	"time"
)

type Event struct {
	Status    Status
	Timestamp time.Time
	Err       error
}

func (s *Shipment) AddEvent(event Event) error {

	if !s.CanUpdate(event.Status) {
		return fmt.Errorf("invalid status update")
	}

	s.Status = event.Status
	s.Events = append(s.Events, event)
	return nil
}
