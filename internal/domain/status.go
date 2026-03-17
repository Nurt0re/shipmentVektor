package domain

type Status string

const (
	StatusPending   Status = "pending"
	StatusShipped   Status = "shipped"
	StatusOnTheWay Status = "on_the_way"
	StatusDelivered Status = "delivered"
)