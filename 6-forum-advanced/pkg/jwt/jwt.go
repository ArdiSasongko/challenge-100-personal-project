package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func GeneratedToken(claims ClaimsToken, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       claims.ID,
		"username": claims.Username,
		"email":    claims.Email,
		"exp":      time.Now().Add(10 * time.Minute).Unix(),
	})

	validToken, err := token.SignedString([]byte(secretKey))

	if err != nil {
		logrus.WithField("generate token", err.Error()).Error(err.Error())
		return "", err
	}

	return validToken, nil
}

func ValidateToken(tokenStr, secretKey string) (*ClaimsToken, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		logrus.WithField("validate token", err.Error()).Error(err.Error())
		return nil, err
	}

	if !token.Valid {
		logrus.WithField("validate token", "token invalid").Error("token invalid")
		return nil, err
	}

	return &ClaimsToken{
		ID:       int64(claims["id"].(float64)),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
	}, nil
}

func GeneratedRefreshToken() string {
	random := make([]byte, 18)
	_, err := rand.Read(random)
	if err != nil {
		logrus.WithField("refresh token", err.Error()).Error(err.Error())
		return ""
	}

	return hex.EncodeToString(random)
}

func ValidateTokenWithOutExpired(tokenStr, secretKey string) (*ClaimsToken, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		logrus.WithField("validate token", err.Error()).Error(err.Error())
		return nil, err
	}

	if !token.Valid {
		logrus.WithField("validate token", "token invalid").Error("token invalid")
		return nil, err
	}

	return &ClaimsToken{
		ID:       int64(claims["id"].(float64)),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
	}, nil
}
