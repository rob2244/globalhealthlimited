package model

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "docker"
	dbName   = "gbhl_device"
)

// DeviceDataClient is a data client api for the device database
type DeviceDataClient struct {
	*sql.DB
}

// NewDeviceDataClient is the constructor function for a device data client
func NewDeviceDataClient() (*DeviceDataClient, error) {
	var conn *sql.DB
	var err error

	if conn, err = sql.Open("postgres", buildConnString()); err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return &DeviceDataClient{conn}, nil
}

// GetDeviceKey gets a random device key from the database
func (c *DeviceDataClient) GetDeviceKey() (string, error) {
	rows, err := c.Query("select device_key from public.device order by random() LIMIT 1;")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	rows.Next()
	var deviceID string
	rows.Scan(&deviceID)

	return strings.ReplaceAll(deviceID, "-", "_"), nil
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
		hst, port, user, pswd, dbName)
}
