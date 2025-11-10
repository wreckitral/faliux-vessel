package main

import (
	"context"
	"log"
	pb "github.com/wreckitral/faliux-vessel/services/consignment/generated/consignment/v1"
	vesselProto "github.com/wreckitral/faliux-vessel/services/vessel/generated/vessel/v1"
	"go-micro.dev/v5"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type consignmentService struct {
	repo         repository
	vesselClient vesselProto.VesselService  // Add vessel client
}

func (s *consignmentService) CreateConsignment(ctx context.Context, req *pb.Consignment,
	res *pb.Response) error {

	// Create specification from consignment
	spec := &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	}

	// Call vessel service to find available vessel
	vesselResponse, err := s.vesselClient.FindAvailable(ctx, spec)
	if err != nil {
		return err
	}

	// Set the vessel ID on the consignment
	req.VesselId = vesselResponse.Vessel.Id

	// Create the consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *consignmentService) GetConsignments(ctx context.Context, req *pb.GetRequest,
	res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &Repository{}

	service := micro.NewService(
		micro.Name("faliux.service.consignment"),
	)
	service.Init()

	// Create vessel service client
	vesselClient := vesselProto.NewVesselService("faliux.service.vessel", service.Client())

	if err := pb.RegisterShippingServiceHandler(service.Server(),
		&consignmentService{
			repo:         repo,
			vesselClient: vesselClient,  // Inject vessel client
		}); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
