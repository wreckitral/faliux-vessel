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
    port = ":30021"
)

type repository interface {
    Create(*pb.Consignment) (*pb.Consignment, error)
}

// Repository - simulate the need of database
type Repository struct {
    mu sync.RWMutex
    consignments []*pb.Consignment
}

// Create a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
    repo.mu.Lock()
    defer repo.mu.Unlock()

    updated := append(repo.consignments, consignment)
    repo.consignments = updated

    return consignment, nil
}

// service should implement all of the methods to satisfy the service
// we defined in protobuf, check the interface in the generated code for the
// method signature
type service struct {
    repo repository
	pb.UnimplementedShippingServer
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main() {
	repo := &Repository{}

	// grpc server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterShippingServer(s, &service{repo: repo})

	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
