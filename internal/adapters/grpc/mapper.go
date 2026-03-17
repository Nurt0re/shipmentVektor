package grpc

import (
	"shipment/internal/domain"
	pb "shipment/shipmentVektor/api/shipment"
)

func shipmentToProto(shipment *domain.Shipment) *pb.Shipment {
	return &pb.Shipment{
		Id:            shipment.ID,
		Reference:     shipment.Reference,
		Origin:        shipment.Origin,
		Destination:   shipment.Destination,
		Status:        statusToProto(shipment.Status),
		Cost:          shipment.Cost,
		DriverRevenue: shipment.DriverRevenue,
	}
}

func eventToProto(event domain.Event) *pb.Event {
	return &pb.Event{
		Status:    statusToProto(event.Status),
		Timestamp: event.Timestamp.Unix(),
	}
}

func statusToProto(status domain.Status) pb.Status {
	switch status {
	case domain.StatusPending:
		return pb.Status_PENDING
	case domain.StatusShipped:
		return pb.Status_SHIPPED
	case domain.StatusOnTheWay:
		return pb.Status_ON_THE_WAY
	case domain.StatusDelivered:
		return pb.Status_DELIVERED
	default:
		return pb.Status_UNKNOWN
	}
}

func protoToStatus(status pb.Status) domain.Status {
	switch status {
	case pb.Status_PENDING:
		return domain.StatusPending
	case pb.Status_SHIPPED:
		return domain.StatusShipped
	case pb.Status_ON_THE_WAY:
		return domain.StatusOnTheWay
	case pb.Status_DELIVERED:
		return domain.StatusDelivered
	default:
		return ""
	}
}
