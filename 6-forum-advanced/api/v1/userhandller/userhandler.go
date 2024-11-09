package userhandler

import (
	"net/http"

	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/api"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/internal/users"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type handler struct {
	*gin.Engine
	service users.Service
	*validator.Validate
}

func NewHandler(api *gin.Engine, service users.Service, v *validator.Validate) *handler {
	return &handler{
		Engine:   api,
		service:  service,
		Validate: v,
	}
}

func (h *handler) RegisterRouter(router *gin.RouterGroup) {
	users := router.Group("/user")
	users.POST("/register", h.Register)
	users.POST("/login", h.Login)

	users.Use(middleware.RefreshToken())
	users.POST("/token", h.GetRefreshToken)
}

func (h *handler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	request := new(users.UserRequest)
	if err := c.ShouldBind(request); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "INTERNAL SERVER ERROR",
			Error:      err.Error(),
		})
		return
	}

	if err := h.Validate.Struct(request); err != nil {
		if errorValidate, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, e := range errorValidate {
				errors[e.Field()] = e.Error()
			}
			c.JSON(http.StatusBadRequest, api.ResponseError{
				StatusCode: http.StatusBadRequest,
				Message:    "BAD REQUEST",
				Error:      errors,
			})
			return
		}
	}

	if err := h.service.Register(ctx, *request); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "BAD REQUEST",
			Error:      err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, api.ResponseSuccess{
		StatusCode: http.StatusCreated,
		Message:    "CREATED",
		Data:       nil,
	})
}

func (h *handler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	request := new(users.LoginRequest)
	if err := c.ShouldBind(request); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "INTERNAL SERVER ERROR",
			Error:      err.Error(),
		})
		return
	}

	if err := h.Validate.Struct(request); err != nil {
		if errorValidate, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, e := range errorValidate {
				errors[e.Field()] = e.Error()
			}
			c.JSON(http.StatusBadRequest, api.ResponseError{
				StatusCode: http.StatusBadRequest,
				Message:    "BAD REQUEST",
				Error:      errors,
			})
			return
		}
	}

	response, err := h.service.Login(ctx, *request)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "BAD REQUEST",
			Error:      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, api.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    "OK",
		Data:       response,
	})
}

func (h *handler) GetRefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	request := new(users.TokenRequest)
	if err := c.ShouldBind(request); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "INTERNAL SERVER ERROR",
			Error:      err.Error(),
		})
		return
	}

	if err := h.Validate.Struct(request); err != nil {
		if errorValidate, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, e := range errorValidate {
				errors[e.Field()] = e.Error()
			}
			c.JSON(http.StatusBadRequest, api.ResponseError{
				StatusCode: http.StatusBadRequest,
				Message:    "BAD REQUEST",
				Error:      errors,
			})
			return
		}
	}

	userID := c.GetInt64("id")

	token, err := h.service.GetRefreshToken(ctx, userID, *request)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "BAD REQUEST",
			Error:      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, api.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    "OK",
		Data:       token,
	})
}
