package handler

import (
	"net/http"
	"os"
	"upload_image/internal/service"

	"github.com/gin-gonic/gin"
)

type handler struct {
	*gin.Engine
	service service.Service
}

func NewHandler(api *gin.Engine, service service.Service) *handler {
	return &handler{
		Engine:  api,
		service: service,
	}
}

func (h *handler) RegisterRouter() {
	router := h.Group("/")
	router.POST("/upload", h.UploadImage)
}

func (h *handler) UploadImage(c *gin.Context) {
	ctx := c.Request.Context()

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "no file available",
			"error":   err.Error(),
		})
		return
	}

	tempFilePath := "file." + file.Filename

	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to save file",
			"error":   err.Error(),
		})
		return
	}

	defer os.Remove(tempFilePath)

	fileData, err := h.service.UploadImage(ctx, tempFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to upload",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success upload",
		"data":    fileData,
	})
}
