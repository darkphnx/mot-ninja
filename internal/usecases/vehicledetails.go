package usecases

import (
	"fmt"
	"time"

	"github.com/darkphnx/vehiclemanager/internal/models"
	"github.com/darkphnx/vehiclemanager/internal/mothistoryapi"
	"github.com/darkphnx/vehiclemanager/internal/vesapi"
)

// VehicleDetails is a wrapper for the APIs necessary to fetch vehicle details
type VehicleDetails struct {
	VehicleEnquiryServiceAPI *vesapi.Client
	MotHistoryAPI            *mothistoryapi.Client
}

// Fetch accesses the mothistory and vesapi and returns a populated vehicle
func (a *VehicleDetails) Fetch(registrationNumber string) (*models.Vehicle, error) {
	vehicleStatus, err := a.VehicleEnquiryServiceAPI.GetVehicleStatus(registrationNumber)
	if err != nil {
		return nil, err
	}

	vehicleHistory, err := a.MotHistoryAPI.GetVehicleHistory(registrationNumber)
	if err != nil {
		return nil, err
	}

	var motHistory []models.MOTTest
	for _, apiTest := range vehicleHistory.MotTests {
		var comments []models.RfrAndComments
		for _, apiComment := range apiTest.RfrAndComments {
			comment := models.RfrAndComments{
				Comment: apiComment.Text,
				Type:    apiComment.Type,
			}
			comments = append(comments, comment)
		}

		test := models.MOTTest{
			TestNumber:      apiTest.MotTestNumber,
			Passed:          apiTest.TestResult == "PASSED",
			CompletedDate:   apiTest.CompletedDate.Time,
			ExpiryDate:      apiTest.ExpiryDate.Time,
			OdometerReading: fmt.Sprintf("%d %s", apiTest.OdometerValue, apiTest.OdometerUnit),
			RfrAndComments:  comments,
		}

		motHistory = append(motHistory, test)
	}

	vehicle := models.Vehicle{
		RegistrationNumber: vehicleStatus.RegistrationNumber,
		Manufacturer:       vehicleHistory.Make,
		Model:              vehicleHistory.Model,
		MotDue:             vehicleHistory.MotTests[0].ExpiryDate.Time,
		VEDDue:             vehicleStatus.TaxDueDate.Time,
		MOTHistory:         motHistory,
		LastFetchedAt:      time.Now(),
	}

	return &vehicle, nil
}
