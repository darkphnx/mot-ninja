package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/darkphnx/vehiclemanager/internal/models"
	"github.com/darkphnx/vehiclemanager/internal/mothistoryapi"
	"github.com/darkphnx/vehiclemanager/internal/usecases"
	"github.com/darkphnx/vehiclemanager/internal/vesapi"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Server contains items request handling
type Server struct {
	Database                 *models.Database
	VehicleEnquiryServiceAPI *vesapi.Client
	MotHistoryAPI            *mothistoryapi.Client
}

type vehicleRequestPayload struct {
	RegistrationNumber string
}

// VehicleCreate handles an add vehicle request
func (s *Server) VehicleCreate(w http.ResponseWriter, r *http.Request) {
	var payload vehicleRequestPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vehicleDetails := usecases.VehicleDetails{
		VehicleEnquiryServiceAPI: s.VehicleEnquiryServiceAPI,
		MotHistoryAPI:            s.MotHistoryAPI,
	}
	vehicle, err := vehicleDetails.Fetch(payload.RegistrationNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = models.CreateVehicle(s.Database, vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vehicle)
}

// VehicleList returns a list of all vehicles
func (s *Server) VehicleList(w http.ResponseWriter, r *http.Request) {
	vehicles, err := models.GetAllVehicles(s.Database)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vehicles)
}

// VehicleDelete deletes a vehicle from the database
func (s *Server) VehicleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vehicleID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vehicle := models.Vehicle{ID: vehicleID}

	err = models.DeleteVehicle(s.Database, &vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// StaticFiles serves up anything in the UI directory
func (s *Server) StaticFiles() http.Handler {
	return http.FileServer(http.Dir("./ui/build"))
}
