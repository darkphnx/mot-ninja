package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/darkphnx/vehiclemanager/internal/authservice"
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
	AuthService              *authservice.AuthService
}

type vehicleCreatePayload struct {
	RegistrationNumber string
}

func (vcp *vehicleCreatePayload) Validate(db *models.Database, user *models.User) []string {
	var errors []string

	registrationNumber := strings.ReplaceAll(strings.ToUpper(vcp.RegistrationNumber), " ", "")
	validRegistration, _ := regexp.MatchString(`^[A-z0-9]{2,7}$`, registrationNumber)
	if !validRegistration {
		errors = append(errors, "Registration Number must be valid")
	}

	vehicleExists := models.UserVehicleExists(db, user.ID, registrationNumber)
	if vehicleExists {
		errors = append(errors, "Vehicle is already added to your account")
	}

	if len(errors) == 0 {
		return nil
	} else {
		return errors
	}
}

// VehicleCreate handles an add vehicle request
func (s *Server) VehicleCreate(w http.ResponseWriter, r *http.Request) {
	var payload vehicleCreatePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		renderError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := getUserFromContext(r)

	validationErrors := payload.Validate(s.Database, user)
	if validationErrors != nil {
		renderError(w, validationErrors, http.StatusUnprocessableEntity)
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

	vehicle.UserID = user.ID

	err = models.CreateVehicle(s.Database, vehicle)
	if err != nil {
		renderError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	renderJSON(w, vehicle, http.StatusCreated)
}

// VehicleList returns a list of all vehicles
func (s *Server) VehicleList(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	vehicles, err := models.GetUserVehicles(s.Database, user.ID)
	if err != nil {
		renderError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, vehicles, http.StatusOK)
}

func (s *Server) VehicleShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := getUserFromContext(r)

	vehicle, err := models.GetUserVehicle(s.Database, user.ID, vars["registration"])
	if err != nil {
		renderError(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, vehicle, http.StatusOK)
}

// VehicleDelete deletes a vehicle from the database
func (s *Server) VehicleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := getUserFromContext(r)

	vehicle, err := models.GetUserVehicle(s.Database, user.ID, vars["registration"])
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
