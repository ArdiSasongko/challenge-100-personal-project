package contenthandler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/usermodel"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *handler) UpserUserActivities(c *gin.Context) {
	ctx := c.Request.Context()

	request := new(usermodel.UserActivitiesRequest)
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

	contentIDStr := c.Param("content_id")
	contentID, err := strconv.Atoi(contentIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.New("post_id not valid").Error(),
		})
		return
	}

	err = h.us.UpsertUserActivities(ctx, userID, int64(contentID), username, *request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "BAD REQUEST",
			"error":      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}
