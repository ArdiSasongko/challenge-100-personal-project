package contenthandler

import (
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/middleware"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/service/contentservice"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type handler struct {
	*gin.Engine
	service contentservice.Service
	*validator.Validate
}

func NewHandler(api *gin.Engine, service contentservice.Service, v *validator.Validate) *handler {
	return &handler{
		Engine:   api,
		service:  service,
		Validate: v,
	}
}

func (h *handler) RegisterRouter() {
	router := h.Group("/content")
	router.Use(middleware.AuthMiddleware())
	router.POST("/", h.CreateContent)
}
