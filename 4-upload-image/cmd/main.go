package main

import (
	"upload_image/internal/configs"
	"upload_image/internal/handler"
	"upload_image/internal/repository"
	"upload_image/internal/service"
	"upload_image/pkg/cloudinary"
	"upload_image/pkg/sql"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	var cfg *configs.Config

	err := configs.Init(
		configs.ConfigWithFolder(
			[]string{"./internal/configs"},
		),
		configs.ConfigWithFile("config"),
		configs.ConfigWithType("yaml"),
	)

	if err != nil {
		logrus.WithField(
			"error", err.Error(),
		).Fatal(err.Error())
	}

	cfg = configs.GetConfig()
	logrus.Info(cfg.Service.Port)
	db, err := sql.Connect(cfg.Database.DataSourceName)
	if err != nil {
		logrus.WithField(
			"error", err.Error(),
		).Fatal(err.Error())
	}

	sql.InitStorage(db)

	cdi, err := cloudinary.Init(cfg.Service.CloudInaryURL)
	if err != nil {
		logrus.WithField(
			"error", err.Error(),
		).Fatal(err.Error())
	}

	r := gin.Default()
	uploadRepo := repository.NewRepository(db)
	uploadService := service.NewService(uploadRepo, cdi)
	uploadHandler := handler.NewHandler(r, uploadService)
	uploadHandler.RegisterRouter()

	r.Run(cfg.Service.Port)
}
