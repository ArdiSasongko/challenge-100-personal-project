package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restapi/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte(config.Envs.JWTSecret)

type contextKey string

var UserKey contextKey = "userID"

func GeneratedToken(userID int) (string, error) {
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

func WithJWT(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := getToken(r)

		token, err := validateToken(tokenStr)
		if err != nil {
			log.Printf("failed verified token, err : %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Print("token is invalid")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, _ := strconv.Atoi(str)

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, userID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getToken(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		return tokenAuth
	}

	return ""
}

func validateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}

		return secret, nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	WriteErr(w, http.StatusForbidden, "Forbiden", fmt.Errorf("permission denied"))
}

func GetUserIdfromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}
