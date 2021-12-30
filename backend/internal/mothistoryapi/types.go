package mothistoryapi

import "time"

func convertJSONToTimestamp(dateFormat string, data []byte) (time.Time, error) {
	// strip the quotes away from the data
	date := string(data)
	date = date[1 : len(date)-1]

	return time.Parse(dateFormat, date)
}

// DottedDate parses a date like 2020.10.31 and adds a time of midnight
type DottedDate struct {
	time.Time
}

// UnmarshalJSON decodes a date like 2020.10.31
func (d *DottedDate) UnmarshalJSON(data []byte) error {
	time, _ := convertJSONToTimestamp("2006.01.02", data)
	d.Time = time

	return nil
}

// DottedTime parses a timestamp like 2020.10.31 11:28:31
type DottedTime struct {
	time.Time
}

// UnmarshalJSON decodes a timestamp like 2020.10.31 11:28:31
func (d *DottedTime) UnmarshalJSON(data []byte) error {
	time, _ := convertJSONToTimestamp("2006.01.02 15:04:05", data)
	d.Time = time

	return nil
}

// Vehicles is a list of Vehicle returned from the API
type Vehicles []Vehicle

// Vehicle contains the MOT history for one vehicle
type Vehicle struct {
	Registration     string     `json:"registration"`
	Make             string     `json:"make"`
	Model            string     `json:"model"`
	FirstUsedDate    DottedDate `json:"firstUsedDate"`
	FuelType         string     `json:"fuelType"`
	PrimaryColour    string     `json:"primaryColour"`
	VehicleID        string     `json:"vehicleId"`
	RegistrationDate DottedDate `json:"registrationDate"`
	ManufactureDate  DottedDate `json:"manufactureDate"`
	EngineSize       int        `json:"engineSize,string"`
	MotTests         []MotTest  `json:"motTests"`
}

// MotTest contains a single MOT test for a vehicle
type MotTest struct {
	CompletedDate      DottedTime      `json:"completedDate"`
	TestResult         string          `json:"testResult"`
	ExpiryDate         DottedDate      `json:"expiryDate,omitempty"`
	OdometerValue      int             `json:"odometerValue,string"`
	OdometerUnit       string          `json:"odometerUnit,omitempty"`
	MotTestNumber      int             `json:"motTestNumber,string"`
	OdometerResultType string          `json:"odometerResultType"`
	RfrAndComments     []RfrAndComment `json:"rfrAndComments"`
}

// RfrAndComment contains the reason for failure, advisories and any comments from an MotTest
type RfrAndComment struct {
	Text      string `json:"text"`
	Type      string `json:"type"`
	Dangerous bool   `json:"dangerous"`
}
