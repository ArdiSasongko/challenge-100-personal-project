package main

import (
	handler "github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/api/v1"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/config"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/internal/users"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func main() {
	var cfg *config.Config
	err := config.Init(
		config.ConfigFolder(
			[]string{"./config"},
		),
		config.ConfigFile("config"),
		config.ConfigType("yaml"),
	)

	if err != nil {
		logrus.Error(err.Error())
	}

	cfg = config.GetConfig()

	db, err := database.Database(cfg.Database.SourceName)
	if err != nil {
		logrus.Error(err.Error())
	}

	database.IntiStorage(db)

	r := gin.Default()
	v := validator.New()
	v1 := r.Group("api/v1")

	// repository
	userRepo := users.NewRepository()

	// service
	userService := users.NewService(db, userRepo)

	// handler
	userHandler := handler.NewUserHandler(r, userService, v)
	userHandler.RegisterRouter(v1)

	r.Run(cfg.Service.Port)
}
