package contenthandler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/api"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/internal/comments"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/internal/contents"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/internal/usersactivities"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type handler struct {
	*gin.Engine
	s  contents.Service
	cs comments.Service
	ua usersactivities.Service
	*validator.Validate
}

func NewHandler(api *gin.Engine, s contents.Service, cs comments.Service, ua usersactivities.Service, v *validator.Validate) *handler {
	return &handler{
		Engine:   api,
		s:        s,
		cs:       cs,
		ua:       ua,
		Validate: v,
	}
}

func (h *handler) RegisterRouter(router *gin.RouterGroup) {
	content := router.Group("/content/")
	content.GET("", h.GetContents)
	content.Use(middleware.AuthMiddleware())
	content.POST(":content_id/comment", h.InsertComment)
	content.PUT(":content_id/liked", h.UpsertUserActivities)
	content.PUT(":content_id/saved", h.UpsertUserActivities)
	content.POST("/upload", h.UploadContent)
	content.GET(":content_id", h.GetContent)
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
	if err := h.s.CreateContent(ctx, userID, username, *request); err != nil {
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

func (h *handler) InsertComment(c *gin.Context) {
	ctx := c.Request.Context()

	request := new(comments.CommentRequest)
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
	username := c.GetString("username")
	contentIDStr := c.Param("content_id")
	contentID, err := strconv.Atoi(contentIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "INTERNAL SERVER ERROR",
			Error:      err.Error(),
		})
		return
	}

	if err := h.cs.InsertComment(ctx, userID, int64(contentID), username, *request); err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "INTERNAL SERVER ERROR",
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

func (h *handler) GetContents(c *gin.Context) {
	ctx := c.Request.Context()

	pI := c.Query("page_index")
	pS := c.Query("page_size")

	if pS == "" {
		pS = "5"
	}

	pageIndex, err := strconv.Atoi(pI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "INTERNAL SERVER ERROR",
			Error:      err.Error(),
		})
		return
	}

	pageSize, err := strconv.Atoi(pS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "INTERNAL SERVER ERROR",
			Error:      err.Error(),
		})
		return
	}

	response, err := h.s.GetContents(ctx, int64(pageSize), int64(pageIndex))
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

func (h *handler) UpsertUserActivities(c *gin.Context) {
	ctx := c.Request.Context()

	request := new(usersactivities.UserActivitiesRequest)
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
	username := c.GetString("username")
	contentIDStr := c.Param("content_id")
	contentID, err := strconv.Atoi(contentIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "INTERNAL SERVER ERROR",
			Error:      err.Error(),
		})
		return
	}

	err = h.ua.UpsertUserActivities(ctx, userID, username, int64(contentID), *request)
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
		Data:       nil,
	})
}

func (h *handler) GetContent(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetInt64("id")
	contentIDStr := c.Param("content_id")
	contentID, err := strconv.Atoi(contentIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "INTERNAL SERVER ERROR",
			Error:      err.Error(),
		})
		return
	}

	response, err := h.s.GetContent(ctx, userID, int64(contentID))
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, sql.ErrNoRows) {
			status = http.StatusNotFound
		}
		c.JSON(http.StatusBadRequest, api.ResponseError{
			StatusCode: int64(status),
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
