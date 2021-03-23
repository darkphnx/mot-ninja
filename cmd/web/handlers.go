package main

import (
	"encoding/json"
	"net/http"

	"github.com/darkphnx/vehiclemanager/internal/models"
	"github.com/darkphnx/vehiclemanager/internal/usecases"
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

		vehicleDetails := usecases.VehicleDetails{
			VehicleEnquiryServiceAPI: server.VehicleEnquiryServiceAPI,
			MotHistoryAPI:            server.MotHistoryAPI,
		}
		vehicle, err := vehicleDetails.Fetch(payload.RegistrationNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		err = models.CreateVehicle(server.Database, vehicle)
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
