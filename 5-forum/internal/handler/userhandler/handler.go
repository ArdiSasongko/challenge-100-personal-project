package userhandler

import (
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/service/userservice"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type handler struct {
	*gin.Engine
	v       *validator.Validate
	service userservice.Service
}

func NewHandler(api *gin.Engine, service userservice.Service, v *validator.Validate) *handler {
	return &handler{
		Engine:  api,
		service: service,
		v:       v,
	}
}

func (h *handler) RegisterRoutes() {
	router := h.Group("/user")
	router.POST("/register", h.RegisterUser)
	router.POST("/login", h.LoginUser)
}
