package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/darkphnx/vehiclemanager/internal/models"
	"github.com/darkphnx/vehiclemanager/internal/mothistoryapi"
	"github.com/darkphnx/vehiclemanager/internal/vesapi"
)

// Env contains items for DI into handlers
type Env struct {
	VehicleEnquiryServiceAPI *vesapi.Client
	MotHistoryAPI            *mothistoryapi.Client
}

func main() {
	vesapiKey := flag.String("vesapi-key", "", "Vehicle Enquiry Service API Key")
	mothistoryapiKey := flag.String("mothistoryapi-key", "", "MOT History API Key")
	flag.Parse()

	err := models.InitDB("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	env := Env{
		VehicleEnquiryServiceAPI: vesapi.NewClient(*vesapiKey, ""),
		MotHistoryAPI:            mothistoryapi.NewClient(*mothistoryapiKey, ""),
	}

	mux := http.NewServeMux()
	mux.Handle("/vehicle/create", vehicleCreate(&env))
	mux.Handle("/vehicles", vehicleList(&env))

	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
