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
			TestNumber:     apiTest.MotTestNumber,
			Passed:         apiTest.TestResult == "PASSED",
			CompletedDate:  apiTest.CompletedDate.Time,
			RfrAndComments: comments,
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
	}

	err = models.CreateVehicle(&vehicle)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", prettyPrint(vehicle))
}
