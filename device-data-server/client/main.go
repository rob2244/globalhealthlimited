package main

import (
	"context"
	pb "globalhealthlimited/device-data-server/device-data-service"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDeviceDataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetDeviceData(ctx, &pb.DeviceDataRequest{DeviceKey: "bf88f3c5-4e6e-4829-8802-da7295cba4ed"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetBrand())

	s, err := c.GetIdentifiers(ctx, &pb.DeviceDataRequest{DeviceKey: "7776b889-8bc4-4f69-a41f-09fc260303c3"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	for {
		data, err := s.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		log.Println(data.DeviceID, data.DeviceType, data.IssuingAgency,
			data.IssuingAgency, data.PackageQuantity, data.DiscontinueDate)
	}

}
