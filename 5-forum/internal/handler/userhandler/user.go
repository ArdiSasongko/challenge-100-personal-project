package userhandler

import (
	"net/http"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/usermodel"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *handler) RegisterUser(c *gin.Context) {
	ctx := c.Request.Context()

	request := new(usermodel.UserRequest)
	if err := c.ShouldBind(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "BAD REQUEST",
			"error":      err.Error(),
		})
		return
	}

	if err := h.v.Struct(request); err != nil {
		if errorValidate, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, e := range errorValidate {
				errors[e.Field()] = e.Error()
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"message":    "BAD REQUEST",
				"error":      errors,
			})
			return
		}
	}

	err := h.service.CreateUser(ctx, *request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "BAD REQUEST",
			"error":      err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"message":    "CREATED",
	})
}

func (h *handler) LoginUser(c *gin.Context) {
	ctx := c.Request.Context()

	request := new(usermodel.LoginRequest)
	if err := c.ShouldBind(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "BAD REQUEST",
			"error":      err.Error(),
		})
		return
	}

	h.v.RegisterStructValidation(usermodel.MustHaveUsernameOrEmail, usermodel.LoginRequest{})
	if err := h.v.Struct(request); err != nil {
		if errorValidate, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, e := range errorValidate {
				errors[e.Field()] = e.Error()
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"message":    "BAD REQUEST",
				"error":      errors,
			})
			return
		}
	}

	result, err := h.service.LoginUser(ctx, *request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "BAD REQUEST",
			"error":      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Success Login",
		"data":       result,
	})
}
