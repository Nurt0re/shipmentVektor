package main

import (
	"log"
	"net"
	"shipment/internal/adapters/grpc"
	"shipment/internal/adapters/postgres"
	"shipment/internal/application"
	"shipment/internal/pkg"
	pb "shipment/shipmentVektor/api/shipment"

	grpcserver "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting gRPC shipment service...")

	cfg := pkg.LoadConfig()
	db, err := pkg.NewPostgresDB(cfg.DBConn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repo := postgres.NewRepository(db)
	service := application.NewShipmentService(repo)
	handler := grpc.NewHandler(service)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	grpcServer := grpcserver.NewServer()
	pb.RegisterShipmentServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on port %s", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
