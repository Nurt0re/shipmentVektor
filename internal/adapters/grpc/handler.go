package grpc

import (
	"context"
	"log"
	"shipment/internal/ports/inbound"
	pb "shipment/shipmentVektor/api/shipment"
)

type Handler struct {
	service inbound.ShipmentUseCase
	pb.UnimplementedShipmentServiceServer
}

func NewHandler(service inbound.ShipmentUseCase) *Handler {

	return &Handler{service: service}
}

func (h *Handler) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
	shipment, err := h.service.CreateShipment(
		ctx,
		req.Id,
		req.Reference,
		req.Origin,
		req.Destination,
		req.Cost,
		req.DriverRevenue,
	)
	if err != nil {
		log.Printf("CreateShipment error: %v", err)
		return nil, err
	}
	return &pb.CreateShipmentResponse{
		Shipment: shipmentToProto(shipment),
	}, nil
}

func (h *Handler) GetShipment(ctx context.Context, req *pb.GetShipmentRequest) (*pb.GetShipmentResponse, error) {
	shipment, err := h.service.GetShipment(ctx, req.Id)
	if err != nil {
		log.Printf("GetShipment error: %v", err)
		return nil, err
	}
	return &pb.GetShipmentResponse{
		Shipment: shipmentToProto(shipment),
	}, nil
}

func (h *Handler) AddEvent(ctx context.Context, req *pb.AddEventRequest) (*pb.AddEventResponse, error) {
	err := h.service.AddEvent(ctx, req.Id, protoToStatus(req.Status))
	if err != nil {
		log.Printf("AddEvent error: %v", err)
		return nil, err
	}
	return &pb.AddEventResponse{
		Msg: "Successfully added event, status changed",
	}, nil
}

func (h *Handler) GetHistory(ctx context.Context, req *pb.GetShipmentHistoryRequest) (*pb.GetShipmentHistoryResponse, error) {
	events, err := h.service.GetHistory(ctx, req.Id)
	if err != nil {
		log.Printf("GetHistory error: %v", err)
		return nil, err
	}

	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvents[i] = eventToProto(event)
	}

	return &pb.GetShipmentHistoryResponse{
		Events: protoEvents,
	}, nil
}
