package main

import (
	"errors"
	"fmt"

	context "golang.org/x/net/context"

	pb "github.com/barcatfigaro/botto/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
)

// Repository is an interface for repos
type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

// VesselRepository simulates a repository
type VesselRepository struct {
	vessels []*pb.Vessel
}

// FindAvailable checks a specification against a map of vessels
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}

	return nil, errors.New("no vessel found by that spec")
}

type service struct {
	repo Repository
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return err
	}

	res.Vessel = vessel
	return nil
}

func main() {
	vessels := []*pb.Vessel{
		&pb.Vessel{Id: "vessel001", Name: "Kane's Salty Secret", MaxWeight: 200000, Capacity: 500},
	}
	repo := &VesselRepository{vessels}

	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)

	srv.Init()

	pb.RegisterVesselServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
