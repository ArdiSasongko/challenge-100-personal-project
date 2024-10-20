package auth

import (
	"basic-rest-api/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(secret []byte, userID int) (string, error) {
	exp := time.Second * time.Duration(config.Envs.JWTExpired)
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
