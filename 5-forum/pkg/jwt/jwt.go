package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func CreateJWT(id int64, username, email, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       id,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(10 * time.Minute).Unix(),
	})

	key := []byte(secretKey)
	validToken, err := token.SignedString(key)

	if err != nil {
		logrus.WithField("error", "failed create jwt token").Error(err.Error())
		return "", err
	}

	return validToken, nil
}

func ValidateToken(tokenStr, secretKey string) (*ClaimsToken, error) {
	key := []byte(secretKey)
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		logrus.WithField("error", "failed get claims value").Error(err.Error())
		return nil, errors.New("failed get claims")
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return &ClaimsToken{
		ID:       int64(claims["id"].(float64)),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
	}, nil
}
