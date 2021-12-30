package authservice

import (
	"errors"
	"time"

	"github.com/darkphnx/vehiclemanager/internal/models"
	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	Secret          []byte
	Issuer          string
	ExpirationHours int64
}

// TokenClaim
type TokenClaim struct {
	UserID string
	jwt.StandardClaims
}

func NewAuthService(secret string, expirationHours int64, issuer string) *AuthService {
	return &AuthService{
		Secret:          []byte(secret),
		ExpirationHours: expirationHours,
		Issuer:          issuer,
	}
}

func (as *AuthService) GenerateAccessToken(user *models.User) (string, error) {
	userID := user.Email
	expiresAt := time.Now().Add(time.Duration(as.ExpirationHours) * time.Hour)

	claim := TokenClaim{
		userID,
		jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			Issuer:    as.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claim)

	return token.SignedString(as.Secret)
}

func (as *AuthService) VerifyAccessToken(signedToken string) (*TokenClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&TokenClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return as.Secret, nil
		})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaim)
	if !ok {
		err = errors.New("Couldn't parse token")
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Unix() {
		err = errors.New("JWT is expired")
		return nil, err
	}

	return claims, nil
}
