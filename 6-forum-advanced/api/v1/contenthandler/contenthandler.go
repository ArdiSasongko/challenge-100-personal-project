package contenthandler

import (
	"net/http"

	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/api"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/internal/contents"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type handler struct {
	*gin.Engine
	service contents.Service
	*validator.Validate
}

func NewHandler(api *gin.Engine, service contents.Service, v *validator.Validate) *handler {
	return &handler{
		Engine:   api,
		service:  service,
		Validate: v,
	}
}

func (h *handler) RegisterRouter(router *gin.RouterGroup) {
	users := router.Group("/content")
	users.Use(middleware.AuthMiddleware())

	users.POST("/upload", h.UploadContent)
}

func (h *handler) UploadContent(c *gin.Context) {
	ctx := c.Request.Context()

	request := new(contents.ContentRequest)
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

	files, err := c.MultipartForm()
	if err == nil && len(files.File["file"]) > 0 {
		request.File = files.File["file"]
	}

	userID := c.GetInt64("id")
	username := c.GetString("username")
	if err := h.service.CreateContent(ctx, userID, username, *request); err != nil {
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
