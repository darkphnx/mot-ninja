package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/darkphnx/vehiclemanager/cmd/api"
	"github.com/darkphnx/vehiclemanager/cmd/background"
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

	vesapiClient := vesapi.NewClient(*vesapiKey, "")
	mothistoryClient := mothistoryapi.NewClient(*mothistoryapiKey, "")

	backgroundTasks := background.Task{
		Database:                 database,
		VehicleEnquiryServiceAPI: vesapiClient,
		MotHistoryAPI:            mothistoryClient,
	}
	go backgroundTasks.Begin()

	apiServer := api.Server{
		Database:                 database,
		VehicleEnquiryServiceAPI: vesapiClient,
		MotHistoryAPI:            mothistoryClient,
	}

	mux := mux.NewRouter()

	mux.Use(api.LoggingMiddleware)

	mux.HandleFunc("/vehicles/{registration}", apiServer.VehicleShow).Methods("GET")
	mux.HandleFunc("/vehicles/{registration}", apiServer.VehicleDelete).Methods("DELETE")
	mux.HandleFunc("/vehicles", apiServer.VehicleList).Methods("GET")
	mux.HandleFunc("/vehicles", apiServer.VehicleCreate).Methods("POST")

	mux.Handle("/", http.FileServer(http.Dir("./ui/build")))

	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
