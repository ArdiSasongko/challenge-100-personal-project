package utils

import (
	"restapi/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "userID"

func GeneratedToken(userID int) (string, error) {
	exp := time.Second * time.Duration(config.Envs.JWTExpired)
	secret := []byte(config.Envs.JWTSecret)

	Claims := jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	validToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return validToken, nil
}
