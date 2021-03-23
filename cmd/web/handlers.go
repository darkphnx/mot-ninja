package main

import (
	"encoding/json"
	"net/http"

	"github.com/darkphnx/vehiclemanager/internal/models"
	"github.com/darkphnx/vehiclemanager/internal/usecases"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func vehicleDelete(server *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		vehicleID, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		vehicle := models.Vehicle{ID: vehicleID}

		err = models.DeleteVehicle(server.Database, &vehicle)
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
