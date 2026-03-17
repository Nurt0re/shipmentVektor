package domain

type Status string

const (
	StatusPending   Status = "pending"
	StatusShipped   Status = "shipped"
	StatusDelivered Status = "delivered"
)