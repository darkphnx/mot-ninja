package vesapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	// BaseURL contains the request endpoint for the VES API
	BaseURL = "https://driver-vehicle-licensing.api.gov.uk/vehicle-enquiry/v1/vehicles"
)

// Client is an API Client for Vehicle Enquiry Service API
type Client struct {
	apiKey     string
	BaseURL    string
	HTTPClient *http.Client
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
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		BaseURL: BaseURL,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// GetVehicleStatus fetches the details from the VES API for the given vehicle
func (c *Client) GetVehicleStatus(registrationNumber string) (*VehicleStatus, error) {
	requestBody := fmt.Sprintf("{\"registrationNumber\":\"%s\"}", registrationNumber)

	req, err := http.NewRequest("POST", c.BaseURL, strings.NewReader(requestBody))

	if err != nil {
		return nil, err
	}

	res := VehicleStatus{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) sendRequest(req *http.Request, vehicleStatus *VehicleStatus) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", c.apiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, res.Body)
	}

	if err = json.NewDecoder(res.Body).Decode(&vehicleStatus); err != nil {
		return err
	}

	return nil
}
