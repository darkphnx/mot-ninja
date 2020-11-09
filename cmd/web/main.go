package main

import (
	"flag"
	"log"

	"github.com/darkphnx/vehiclemanager/internal/vesapi"
)

func main() {
	vesapiKey := flag.String("vesapi-key", "", "Vehicle Enquiry Service API Key")
	regNumber := flag.String("reg-number", "", "Registration Number to Lookup")
	flag.Parse()

	client := vesapi.NewClient(*vesapiKey, "")

	vehicleStatus, err := client.GetVehicleStatus(*regNumber)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", vehicleStatus)
	log.Printf("%s\n", vehicleStatus.FuelType)
	log.Printf("%s\n", vehicleStatus.DateOfLastV5CIssued)
}
