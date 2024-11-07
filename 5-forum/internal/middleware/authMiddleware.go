package middleware

import (
	"net/http"
	"strings"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/configs"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	secretKey := configs.GetConfig().Service.SecretJWT
	return func(ctx *gin.Context) {
		header := ctx.Request.Header.Get("Authorization")

		header = strings.TrimSpace(header)
		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing header",
			})
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing header",
			})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		token = strings.TrimSpace(token)

		tokenClaims, err := jwt.ValidateToken(token, secretKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		}

		ctx.Set("user_id", tokenClaims.ID)
		ctx.Set("username", tokenClaims.Username)
		ctx.Set("email", tokenClaims.Email)

		ctx.Next()
	}
}
