package domain

import "time"

type Event struct {
	Status Status
	Timestamp time.Time
	Err error
}