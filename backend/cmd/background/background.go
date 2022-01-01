package background

import (
	"log"
	"time"

	"github.com/darkphnx/vehiclemanager/internal/models"
	"github.com/darkphnx/vehiclemanager/internal/mothistoryapi"
	"github.com/darkphnx/vehiclemanager/internal/usecases"
	"github.com/darkphnx/vehiclemanager/internal/vesapi"
)

// Task contains all of the external connections we need
type Task struct {
	Database                 *models.Database
	VehicleEnquiryServiceAPI *vesapi.Client
	MotHistoryAPI            *mothistoryapi.Client
}

// Begin fetches new MOT data every 5 minutes
func (bt *Task) Begin() {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ticker.C:
			bt.updateVehicles()
		}
	}
}

func (bt *Task) updateVehicles() {
	log.Println("Update Vehicles")

	timestamp := time.Now().Add(-1 * time.Hour)
	vehicles, err := models.GetVehiclesUpdatedBefore(bt.Database, timestamp)
	if err != nil {
		log.Println(err)
	}

	vehicleDetails := usecases.VehicleDetails{
		VehicleEnquiryServiceAPI: bt.VehicleEnquiryServiceAPI,
		MotHistoryAPI:            bt.MotHistoryAPI,
	}

	for _, vehicle := range vehicles {
		log.Printf("Updating vehicle %s...\n", vehicle.RegistrationNumber)

		updatedVehicleDetails, err := vehicleDetails.Fetch(vehicle.RegistrationNumber)
		if err != nil {
			log.Println(err)
		}

		vehicle.MOTHistory = updatedVehicleDetails.MOTHistory
		vehicle.MotDue = updatedVehicleDetails.MotDue
		vehicle.VEDDue = updatedVehicleDetails.VEDDue
		vehicle.LastFetchedAt = updatedVehicleDetails.LastFetchedAt

		err = models.UpdateVehicle(bt.Database, vehicle)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("Updating Vehicles Complete")
}
