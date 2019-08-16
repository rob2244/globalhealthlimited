package main

import (
	context "context"
	"database/sql"
	"fmt"
	pb "globalhealthlimited/device-data-server/device-data-service"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

const (
	port            = ":50051"
	host            = "localhost"
	dbPort          = 5432
	user            = "postgres"
	password        = "docker"
	dbName          = "gbhl_device"
	query           = "SELECT device_key, publish_date, brand, company, model_number FROM public.device where device_key = $1"
	identifierQuery = "SELECT device_id, device_type, issuing_agency, package_quantity, discontinue_date FROM public.identifier where device_key = $1"
)

var conn *sql.DB

type server struct {
}

func (S *server) GetDeviceData(ctx context.Context, in *pb.DeviceDataRequest) (*pb.DeviceDataResponse, error) {
	data, err := getDeviceData(in.GetDeviceKey())
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (S *server) GetIdentifiers(in *pb.DeviceDataRequest, stream pb.DeviceDataService_GetIdentifiersServer) error {
	results := make(chan *pb.IdentifierDataResponse, 100)
	go getIdentifierData(in.DeviceKey, results)

	for result := range results {
		if err := stream.Send(result); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var err error
	conn, err = sql.Open("postgres", buildConnString())

	if err != nil {
		log.Fatal(err)
	}

	if err = conn.Ping(); err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDeviceDataServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getDeviceData(deviceKey string) (*pb.DeviceDataResponse, error) {
	rows, err := conn.Query(query, deviceKey)
	if err != nil {
		return nil, err
	}

	result := pb.DeviceDataResponse{}
	rows.Next()
	err = rows.Scan(&result.DeviceKey, &result.PublishDate,
		&result.Brand, &result.Company, &result.ModelNumber)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getIdentifierData(deviceKey string, responses chan<- *pb.IdentifierDataResponse) {
	rows, err := conn.Query(identifierQuery, deviceKey)
	if err != nil {
		close(responses)
		log.Fatal(err)
	}

	for rows.Next() {
		result := pb.IdentifierDataResponse{}

		err = rows.Scan(&result.DeviceID, &result.DeviceType,
			&result.IssuingAgency, &result.PackageQuantity, &result.DiscontinueDate)

		if err != nil {
			close(responses)
			log.Fatal(err)
		}

		responses <- &result
	}

	close(responses)
}

func buildConnString() string {
	pswd := os.Getenv("PGSQL_PASSWORD")
	if pswd == "" {
		pswd = password
	}

	env := os.Getenv("ENVIRONMENT")
	var hst string
	if env == "" || env == "DEV" {
		hst = host
	}

	if env == "PROD" {
		hst = "postgres-postgresql"
	}

	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hst, dbPort, user, pswd, dbName)
}
