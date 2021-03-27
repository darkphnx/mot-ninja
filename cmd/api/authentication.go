package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/darkphnx/vehiclemanager/internal/models"
	"golang.org/x/crypto/bcrypt"
)

const jwtCookieName = "jwt"

type loginPayload struct {
	Email    string
	Password string
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var payload loginPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		renderError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := models.GetUser(s.Database, payload.Email)
	if err != nil {
		renderBadUsernamePassword(w)
		return
	}

	if !checkPassword(payload.Password, user.HashedPassword) {
		renderBadUsernamePassword(w)
		return
	}

	accessToken, err := s.AuthService.GenerateAccessToken(user)
	if err != nil {
		renderError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     jwtCookieName,
		Value:    accessToken,
		HttpOnly: true,
	})

	renderOkay(w, http.StatusOK)
}

func renderBadUsernamePassword(w http.ResponseWriter) {
	renderError(w, "Incorrect email or password", http.StatusForbidden)
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     jwtCookieName,
		Value:    "",
		HttpOnly: true,
	})

	renderOkay(w, http.StatusOK)
}

func (s *Server) AuthJwtTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err := r.Cookie(jwtCookieName)

		if err != nil {
			renderError(w, "Missing JWT token", http.StatusForbidden)
			return
		}

		jwtClaim, err := s.AuthService.VerifyAccessToken(jwtCookie.Value)
		if err != nil {
			renderError(w, "Invalid JWT token", http.StatusForbidden)
			return
		}

		user, err := models.GetUser(s.Database, jwtClaim.UserID)
		if err != nil {
			renderError(w, "Could not find user", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *models.User {
	return r.Context().Value("user").(*models.User)
}
