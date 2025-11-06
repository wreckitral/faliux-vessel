package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	pb "github.com/wreckitral/faliux-vessel/faliux-vessel-service-consignment/proto/consignment"
	"go-micro.dev/v5"
)

const (
	defaultFilename = "consignment.json"
)

// parse the data json file
func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// populate the consignment protobuf struct with data from json
	if err := json.Unmarshal(data, &consignment); err != nil {
		return nil, err
	}

	return consignment, err
}

func main() {
	// setup connection to the gRPC server
	service := micro.NewService(micro.Name("faliux.consignment.cli"))
	service.Init()

	client := pb.NewShippingService("faliux.service.consignment", service.Client())

	// if filename is not specified then use defaultFilename
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse the file: %v", err)
	}

	r, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("Could not create Consignment: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not get Consignments: %v", err)
	}

		for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
