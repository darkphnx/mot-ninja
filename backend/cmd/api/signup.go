package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/darkphnx/vehiclemanager/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type signupPayload struct {
	Email              string
	Password           string
	PasswordConfirm    string
	TermsAndConditions bool
}

func (sp *signupPayload) Validate(db *models.Database) []string {
	var errors []string

	validEmail, _ := regexp.MatchString(`^.+?@.+?\..+?$`, sp.Email)
	if !validEmail {
		errors = append(errors, "E-mail address is not valid")
	}

	emailExists := models.UserExists(db, sp.Email)
	if emailExists {
		errors = append(errors, "E-mail address is already registered")
	}

	validPassword, _ := regexp.MatchString(`^.{6,64}$`, sp.Password)
	if !validPassword {
		errors = append(errors, "Password must be between 6 and 64 characters in length")
	}

	if sp.Password != sp.PasswordConfirm {
		errors = append(errors, "Password and confirmation must be the same")
	}

	if !sp.TermsAndConditions {
		errors = append(errors, "Terms and Conditions must be agreed to")
	}

	if len(errors) == 0 {
		return nil
	} else {
		return errors
	}
}

func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
	var payload signupPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		renderError(w, err.Error(), http.StatusBadRequest)
		return
	}

	errors := payload.Validate(s.Database)
	if errors != nil {
		renderError(w, errors, http.StatusUnprocessableEntity)
		return
	}

	hashedPassword, err := hashPassword(payload.Password)
	if err != nil {
		renderError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := models.User{
		Email:          payload.Email,
		HashedPassword: hashedPassword,
		VehicleLimit:   5,
	}

	err = models.CreateUser(s.Database, &user)
	if err != nil {
		renderError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	renderJSON(w, &user, http.StatusCreated)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
