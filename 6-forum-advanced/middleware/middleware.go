package middleware

import (
	"net/http"
	"strings"

	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/api"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/config"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware() gin.HandlerFunc {
	secretKey := config.GetConfig().Service.SecretJWT
	return func(ctx *gin.Context) {
		header := ctx.Request.Header.Get("Authorization")
		header = strings.TrimSpace(header)

		if header == "" {
			logrus.WithField("get header", "Missing Header").Error("Missing Header")
			ctx.JSON(http.StatusUnauthorized, api.ResponseError{
				StatusCode: http.StatusUnauthorized,
				Message:    "UNAUTHORIZED",
				Error:      "Missing Header",
			})
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			logrus.WithField("get header", "Missing Header").Error("Missing Header")
			ctx.JSON(http.StatusUnauthorized, api.ResponseError{
				StatusCode: http.StatusUnauthorized,
				Message:    "UNAUTHORIZED",
				Error:      "Missing Header",
			})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		tokenClaims, err := jwt.ValidateToken(token, secretKey)
		if err != nil {
			logrus.WithField("validaye token", err.Error()).Error("failed validate token")
			ctx.JSON(http.StatusUnauthorized, api.ResponseError{
				StatusCode: http.StatusUnauthorized,
				Message:    "UNAUTHORIZED",
				Error:      err.Error(),
			})
			return
		}

		ctx.Set("id", tokenClaims.ID)
		ctx.Set("username", tokenClaims.Username)
		ctx.Set("email", tokenClaims.Email)

		ctx.Next()
	}
}

func RefreshToken() gin.HandlerFunc {
	secretKey := config.GetConfig().Service.SecretJWT
	return func(ctx *gin.Context) {
		header := ctx.Request.Header.Get("Authorization")
		header = strings.TrimSpace(header)

		if header == "" {
			logrus.WithField("get header", "Missing Header").Error("Missing Header")
			ctx.JSON(http.StatusUnauthorized, api.ResponseError{
				StatusCode: http.StatusUnauthorized,
				Message:    "UNAUTHORIZED",
				Error:      "Missing Header",
			})
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			logrus.WithField("get header", "Missing Header").Error("Missing Header")
			ctx.JSON(http.StatusUnauthorized, api.ResponseError{
				StatusCode: http.StatusUnauthorized,
				Message:    "UNAUTHORIZED",
				Error:      "Missing Header",
			})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		tokenClaims, err := jwt.ValidateTokenWithOutExpired(token, secretKey)
		if err != nil {
			logrus.WithField("validaye token", err.Error()).Error("failed validate token")
			ctx.JSON(http.StatusUnauthorized, api.ResponseError{
				StatusCode: http.StatusUnauthorized,
				Message:    "UNAUTHORIZED",
				Error:      err.Error(),
			})
			return
		}

		ctx.Set("id", tokenClaims.ID)
		ctx.Set("username", tokenClaims.Username)
		ctx.Set("email", tokenClaims.Email)

		ctx.Next()
	}
}
