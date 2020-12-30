package main

import (
	"encoding/json"
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

		err = models.CreateVehicle(server.Database, &vehicle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		json.NewEncoder(w).Encode(vehicle)
	}
}

func vehicleList(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vehicles, err := models.GetVehicles(server.Database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(vehicles)
	}
}
