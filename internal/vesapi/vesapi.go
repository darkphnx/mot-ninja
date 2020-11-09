package vesapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	defaultHost = "https://driver-vehicle-licensing.api.gov.uk"
)

// Client is an API Client for Vehicle Enquiry Service API
type Client struct {
	apiKey   string
	baseHost string
	client   *http.Client
}

// Date is a special Time which only has date components
type Date struct {
	time.Time
}

// UnmarshalJSON parses an rfc8601 date into a Date
func (d *Date) UnmarshalJSON(data []byte) error {
	// strip the quotes away from the data
	date := string(data)
	date = date[1 : len(date)-1]

	rfc3339Timestamp := fmt.Sprintf("%sT00:00:00Z", date)
	time, _ := time.Parse(time.RFC3339, rfc3339Timestamp)
	d.Time = time

	return nil
}

// VehicleStatus contains the data returned from the API for a successful query
type VehicleStatus struct {
	ArtEndDate               string `json:"artEndDate"`
	Co2Emissions             int    `json:"co2Emissions"`
	Colour                   string `json:"colour"`
	EngineCapacity           int    `json:"engineCapacity"`
	FuelType                 string `json:"fuelType"`
	Make                     string `json:"make"`
	MarkedForExport          bool   `json:"markedForExport"`
	MonthOfFirstRegistration string `json:"monthOfFirstRegistration"`
	MotStatus                string `json:"motStatus"`
	RegistrationNumber       string `json:"registrationNumber"`
	RevenueWeight            int    `json:"revenueWeight"`
	TaxDueDate               Date   `json:"taxDueDate"`
	TaxStatus                string `json:"taxStatus"`
	TypeApproval             string `json:"typeApproval"`
	Wheelplan                string `json:"wheelplan"`
	YearOfManufacture        int    `json:"yearOfManufacture"`
	EuroStatus               string `json:"euroStatus"`
	RealDrivingEmissions     string `json:"realDrivingEmissions"`
	DateOfLastV5CIssued      Date   `json:"dateOfLastV5CIssued"`
}

// NewClient returns a new Vehicle Enquiry Service API Client
func NewClient(apiKey, host string) *Client {
	if host == "" {
		host = defaultHost
	}

	return &Client{
		apiKey:   apiKey,
		baseHost: host,
		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// GetVehicleStatus fetches the details from the VES API for the given vehicle
func (c *Client) GetVehicleStatus(registrationNumber string) (*VehicleStatus, error) {
	requestBody := fmt.Sprintf("{\"registrationNumber\":\"%s\"}", registrationNumber)

	req, err := http.NewRequest("POST", c.baseHost+"/vehicle-enquiry/v1/vehicles", strings.NewReader(requestBody))

	if err != nil {
		return nil, err
	}

	res := VehicleStatus{}
	err = c.sendRequest(req, &res)

	return &res, err
}

func (c *Client) sendRequest(req *http.Request, vehicleStatus *VehicleStatus) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", c.apiKey)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, body)
	}

	return json.NewDecoder(res.Body).Decode(&vehicleStatus)
}
