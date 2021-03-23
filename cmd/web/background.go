package main

import (
	"log"
	"time"

	"github.com/darkphnx/vehiclemanager/internal/models"
	"github.com/darkphnx/vehiclemanager/internal/usecases"
)

// BackgroundUpdateVehicles fetches new MOT data every 5 minutes
func BackgroundUpdateVehicles(server *Server) {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ticker.C:
			updateVehicles(server)
		}
	}
}

func updateVehicles(server *Server) {
	log.Println("Update Vehicles")

	timestamp := time.Now().Add(-1 * time.Hour)
	vehicles, err := models.GetVehiclesUpdatedBefore(server.Database, timestamp)
	if err != nil {
		log.Println(err)
	}

	vehicleDetails := usecases.VehicleDetails{
		VehicleEnquiryServiceAPI: server.VehicleEnquiryServiceAPI,
		MotHistoryAPI:            server.MotHistoryAPI,
	}

	for _, vehicle := range vehicles {
		log.Printf("Updating vehicle %s...\n", vehicle.RegistrationNumber)

		updatedVehicleDetails, err := vehicleDetails.Fetch(vehicle.RegistrationNumber)
		if err != nil {
			log.Println(err)
		}

		err = models.UpdateVehicle(server.Database, vehicle, updatedVehicleDetails)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("Updating Vehicles Complete")
}
