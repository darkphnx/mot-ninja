package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/darkphnx/vehiclemanager/internal/models"
)

type vehicleRequestPayload struct {
	RegistrationNumber string
}

func vehicleCreate(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload vehicleRequestPayload

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		vehicleStatus, err := server.VehicleEnquiryServiceAPI.GetVehicleStatus(payload.RegistrationNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		vehicleHistory, err := server.MotHistoryAPI.GetVehicleHistory(payload.RegistrationNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
		}

		err = models.CreateVehicle(server.Database, &vehicle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(vehicle)
	}
}

func vehicleList(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vehicles, err := models.GetAllVehicles(server.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(vehicles)
	}
}

type singleVehicleRequestPayload struct {
	ID string
}

func vehicleDelete(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload models.Vehicle

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = models.DeleteVehicle(server.Database, &payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func staticFiles() http.Handler {
	return http.FileServer(http.Dir("./ui/build"))
}
