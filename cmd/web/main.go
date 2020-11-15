package main

import (
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/darkphnx/vehiclemanager/internal/models"
	"github.com/darkphnx/vehiclemanager/internal/mothistoryapi"
	"github.com/darkphnx/vehiclemanager/internal/vesapi"
)

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func main() {
	vesapiKey := flag.String("vesapi-key", "", "Vehicle Enquiry Service API Key")
	mothistoryapiKey := flag.String("mothistoryapi-key", "", "MOT History API Key")
	regNumber := flag.String("reg-number", "", "Registration Number to Lookup")
	flag.Parse()

	err := models.InitDB("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	vesClient := vesapi.NewClient(*vesapiKey, "")

	vehicleStatus, err := vesClient.GetVehicleStatus(*regNumber)
	if err != nil {
		log.Fatal(err)
	}

	motHistoryClient := mothistoryapi.NewClient(*mothistoryapiKey, "")

	vehicleHistory, err := motHistoryClient.GetVehicleHistory(*regNumber)
	if err != nil {
		log.Fatal(err)
	}

	vehicle := models.Vehicle{
		RegistrationNumber: vehicleStatus.RegistrationNumber,
		Manufacturer:       vehicleHistory.Make,
		Model:              vehicleHistory.Model,
		MotDue:             time.Unix(vehicleHistory.MotTests[0].ExpiryDate.Unix(), 0),
		VEDDue:             time.Unix(vehicleStatus.TaxDueDate.Unix(), 0),
	}

	log.Printf("%s", prettyPrint(vehicle))

	err = models.CreateVehicle(&vehicle)
	if err != nil {
		log.Fatal(err)
	}

}
