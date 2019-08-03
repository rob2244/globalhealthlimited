package data

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	_ "github.com/lib/pq"
)

const (
	host         = "localhost"
	port         = 5432
	user         = "postgres"
	password     = "docker"
	dbName       = "gbhl_device"
	deviceInsert = "insert into public.device(device_key, publish_date, brand, company, model_number) values ($1, $2, $3, $4, $5)  ON CONFLICT DO NOTHING"
	identInsert  = "insert into public.identifier(device_id, device_type, issuing_agency, package_quantity, discontinue_date, device_key) values ($1, $2, $3, $4, $5, $6)  ON CONFLICT DO NOTHING"
)

// Devices holds a list of all the device data parsed from the xml file
type Devices struct {
	Devices []DeviceData `xml:"device"`
}

// DeviceData holds parsed device data parsed from an xml file
type DeviceData struct {
	DeviceKey   string        `xml:"publicDeviceRecordKey"`
	PublishDate string        `xml:"devicePublishDate"`
	Brand       string        `xml:"brandName"`
	Company     string        `xml:"companyName"`
	ModelNumber string        `xml:"versionModelNumber"`
	Identifiers []Identifiers `xml:"identifiers>identifier"`
}

// Identifiers holds the identifier information of individual devices parsed from the xml file
type Identifiers struct {
	DeviceID        string `xml:"deviceId"`
	DeviceType      string `xml:"deviceType"`
	IssuingAgency   string `xml:"deviceIdIssuingAgency"`
	PackageQuantity string `xml:"pkgQuantity"`
	DiscontinueDate string `xml:"pkgDiscontinueDate"`
}

var conn *sql.DB

func init() {
	var err error
	conn, err = sql.Open("postgres", buildConnString())

	if err != nil {
		panic(err)
	}

	if err = conn.Ping(); err != nil {
		panic(err)
	}
}

// Parse parses the given xml file into a struct of devices
func Parse(file io.Reader) (*Devices, error) {
	devices := &Devices{}
	if err := xml.NewDecoder(file).Decode(devices); err != nil {
		return nil, err
	}

	return devices, nil
}

// SaveDeviceData saves device data to database
func SaveDeviceData(devices Devices) error {
	tran, err := conn.Begin()
	if err != nil {
		return err
	}

	for _, dev := range devices.Devices {
		_, err := tran.Exec(deviceInsert, dev.DeviceKey, dev.PublishDate,
			dev.Brand, dev.Company, dev.ModelNumber)

		if err != nil {
			tranErr := tran.Rollback()
			if tranErr != nil {
				panic(tranErr)
			}

			return err
		}

		for _, ident := range dev.Identifiers {
			_, err := tran.Exec(identInsert, ident.DeviceID, ident.DeviceType,
				ident.IssuingAgency, ident.PackageQuantity, ident.DiscontinueDate, dev.DeviceKey)

			if err != nil {
				tranErr := tran.Rollback()
				if tranErr != nil {
					panic(tranErr)
				}

				return err
			}
		}

	}

	err = tran.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Close closes an existing open DB connection
func Close() {
	conn.Close()
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
