package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/wreckitral/faliux-vessel/faliux-vessel-service-consignment/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":22120"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

type Repository struct {
	mu				sync.RWMutex
	consignments 	[]*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()

	updated := append(repo.consignments, consignment)
	repo.consignments = updated

	repo.mu.Unlock()

	return consignment, nil
}

type service struct {
	repo repository
	pb.UnimplementedShippingServiceServer
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (
									*pb.Response, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main() {
	repo := &Repository{}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// register service with the gRPC server
	pb.RegisterShippingServiceServer(s, &service{repo: repo})

	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
