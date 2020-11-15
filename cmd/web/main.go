package main

import (
	"encoding/json"
	"flag"
	"log"

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

	log.Printf("%+v\n", vehicleStatus)
	log.Printf("%s", prettyPrint(vehicleHistory))
}
