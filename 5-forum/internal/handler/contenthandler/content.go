package contenthandler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/contentmodel"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *handler) CreateContent(c *gin.Context) {
	ctx := c.Request.Context()

	request := new(contentmodel.ContentRequest)
	if err := c.ShouldBind(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "BAD REQUEST",
			"error":      err.Error(),
		})
		return
	}

	if err := h.Validate.Struct(request); err != nil {
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

	userID := c.GetInt64("user_id")
	username := c.GetString("username")

	err := h.service.CreateContent(ctx, *request, userID, username)
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

func (h *handler) GetContent(c *gin.Context) {
	ctx := c.Request.Context()

	contentIDStr := c.Param("content_id")
	contentID, err := strconv.ParseInt(contentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.New("post_id not valid").Error(),
		})
		return
	}

	result, err := h.service.GetContent(ctx, contentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, result)
}
