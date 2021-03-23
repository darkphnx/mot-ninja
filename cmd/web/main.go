package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/darkphnx/vehiclemanager/internal/models"
	"github.com/darkphnx/vehiclemanager/internal/mothistoryapi"
	"github.com/darkphnx/vehiclemanager/internal/vesapi"
)

// Server contains items for DI into handlers
type Server struct {
	Database                 *models.Database
	VehicleEnquiryServiceAPI *vesapi.Client
	MotHistoryAPI            *mothistoryapi.Client
}

func main() {
	vesapiKey := flag.String("vesapi-key", "", "Vehicle Enquiry Service API Key")
	mothistoryapiKey := flag.String("mothistoryapi-key", "", "MOT History API Key")
	flag.Parse()

	database, err := models.InitDB("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	server := Server{
		Database:                 database,
		VehicleEnquiryServiceAPI: vesapi.NewClient(*vesapiKey, ""),
		MotHistoryAPI:            mothistoryapi.NewClient(*mothistoryapiKey, ""),
	}

	go BackgroundUpdateVehicles(&server)

	mux := http.NewServeMux()
	mux.Handle("/vehicle/create", vehicleCreate(&server))
	mux.Handle("/vehicle/delete", vehicleDelete(&server))
	mux.Handle("/vehicles", vehicleList(&server))
	mux.Handle("/", staticFiles())

	loggingMiddleware := RequestLoggingMiddleware()
	loggedMux := loggingMiddleware(mux)

	err = http.ListenAndServe(":4000", loggedMux)
	log.Fatal(err)
}
