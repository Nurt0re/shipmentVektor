package domain

type Shipment struct {
	ID    string
	ReferenceNumber int32
	Origin string
	Destination string
	Status Status
	Events []Event
	Cost float64
	DriverRevenue float64
}