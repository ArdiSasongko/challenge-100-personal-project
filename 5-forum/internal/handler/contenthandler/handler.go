package contenthandler

import (
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/middleware"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/service/contentservice"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/service/userservice"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type handler struct {
	*gin.Engine
	service contentservice.Service
	us      userservice.Service
	*validator.Validate
}

func NewHandler(api *gin.Engine, service contentservice.Service, us userservice.Service, v *validator.Validate) *handler {
	return &handler{
		Engine:   api,
		service:  service,
		us:       us,
		Validate: v,
	}
}

func (h *handler) RegisterRouter(router *gin.RouterGroup) {
	content := router.Group("/content")
	content.Use(middleware.AuthMiddleware())
	content.POST("/", h.CreateContent)
	content.GET("/", h.GetContents)
	content.GET("/:content_id", h.GetContent)
	content.POST("/:content_id/comment", h.CreateComment)
	content.PUT("/:content_id/activities", h.UpserUserActivities)
}
