package domain

import "time"

type Event struct {
	Status Status
	timestamp time.Time
	err error
}