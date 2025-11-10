package main

import (
	"context"
	"errors"
	"log"

	pb "github.com/wreckitral/faliux-vessel/services/vessel/generated/vessel/v1"
	"go-micro.dev/v5"
)

type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type VesselRepository struct {
	vessels []*pb.Vessel
}

func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}

	return nil, errors.New("No vessel found by that specification")
}

type vesselService struct {
	repo Repository
}

func (s *vesselService) FindAvailable(ctx context.Context, req *pb.Specification,
										res *pb.Response) error {
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return err
	}

	res.Vessel = vessel

	return nil
}

func main() {
	vessels := []*pb.Vessel{
		&pb.Vessel{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}

	repo := &VesselRepository{vessels}

	service := micro.NewService(micro.Name("faliux.service.vessel"))

	service.Init()

	if err := pb.RegisterVesselServiceHandler(service.Server(), &vesselService{repo}); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
