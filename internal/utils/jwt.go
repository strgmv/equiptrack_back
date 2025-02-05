package utils

import (
	"equiptrack/config"
	"equiptrack/internal/models"
	"errors"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/rand"
)

// чо сделать
// код для генерации generateJWT(data, TTL)
// генерировать сразу пару токенов
// код для генерации refresh token (спизди с гитхаба)
// refresh+userID в редис
// mw переделывать вроде не надо
// /auth/refresh для обновления токенов

type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

// Generate new JWT Token
func GenerateJWTToken(user *models.User, config *config.Config) (string, error) {
	// Register the JWT claims, which includes the username and expiry time
	claims := Claims{
		ID: user.UserID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * config.Server.JwtMaxAge)),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Register the JWT string
	tokenString, err := token.SignedString([]byte(config.Server.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// вроде нигде не используется
// удолить
func ExtractJWTFromRequest(r *http.Request) (map[string]interface{}, error) {
	// Get the JWT string
	tokenString := ExtractBearerToken(r)

	// Initialize a new instance of `Claims` (here using Claims map)
	claims := jwt.MapClaims{}

	// Parse the JWT string and repositories the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (jwtKey interface{}, err error) {
		return jwtKey, err
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token ")
	}

	return claims, nil
}

// удолить тоже
func ExtractBearerToken(r *http.Request) string {
	headerAuthorization := r.Header.Get("Authorization")
	bearerToken := strings.Split(headerAuthorization, " ")
	return html.EscapeString(bearerToken[1])
}

func NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	s := rand.NewSource(uint64(time.Now().UnixNano()))
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
