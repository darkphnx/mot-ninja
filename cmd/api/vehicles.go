package api

import (
	"encoding/json"
	"net/http"

	"github.com/darkphnx/vehiclemanager/internal/models"
	"github.com/darkphnx/vehiclemanager/internal/mothistoryapi"
	"github.com/darkphnx/vehiclemanager/internal/usecases"
	"github.com/darkphnx/vehiclemanager/internal/vesapi"
	"github.com/gorilla/mux"
)

// Server contains items request handling
type Server struct {
	Database                 *models.Database
	VehicleEnquiryServiceAPI *vesapi.Client
	MotHistoryAPI            *mothistoryapi.Client
}

type vehicleCreatePayload struct {
	RegistrationNumber string
}

// VehicleCreate handles an add vehicle request
func (s *Server) VehicleCreate(w http.ResponseWriter, r *http.Request) {
	var payload vehicleCreatePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		renderError(w, err.Error(), http.StatusBadRequest)
		return
	}

	vehicleDetails := usecases.VehicleDetails{
		VehicleEnquiryServiceAPI: s.VehicleEnquiryServiceAPI,
		MotHistoryAPI:            s.MotHistoryAPI,
	}
	vehicle, err := vehicleDetails.Fetch(payload.RegistrationNumber)
	if err != nil {
		renderError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = models.CreateVehicle(s.Database, vehicle)
	if err != nil {
		renderError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	renderJSON(w, vehicle, http.StatusCreated)
}

// VehicleList returns a list of all vehicles
func (s *Server) VehicleList(w http.ResponseWriter, r *http.Request) {
	vehicles, err := models.GetAllVehicles(s.Database)
	if err != nil {
		renderError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, vehicles, http.StatusOK)
}

func (s *Server) VehicleShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	vehicle, err := models.GetVehicle(s.Database, vars["registration"])
	if err != nil {
		renderError(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, vehicle, http.StatusOK)
}

// VehicleDelete deletes a vehicle from the database
func (s *Server) VehicleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	vehicle, err := models.GetVehicle(s.Database, vars["registration"])
	if err != nil {
		renderError(w, err.Error(), http.StatusNotFound)
		return
	}

	err = models.DeleteVehicle(s.Database, vehicle)
	if err != nil {
		renderError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderOkay(w, http.StatusOK)
}

type simpleResponse struct {
	Status string
}

type errorResponse struct {
	Error interface{}
}

func renderOkay(w http.ResponseWriter, status int) {
	renderJSON(w, simpleResponse{Status: "ok"}, status)
}

func renderError(w http.ResponseWriter, errMsg interface{}, status int) {
	err := errorResponse{Error: errMsg}
	renderJSON(w, err, status)
}

func renderJSON(w http.ResponseWriter, payload interface{}, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
