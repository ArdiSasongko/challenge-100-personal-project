package main

import (
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/api/v1/contenthandler"
	userhandler "github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/api/v1/userhandller"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/config"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/internal/comments"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/internal/contents"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/internal/users"
	"github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/internal/usersactivities"
	cld "github.com/ArdiSasongko/challenge-100-personal-project/6-forum-advanced/pkg/cloudinary"
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

	cld, err := cld.Init(cfg.Service.CloudInaryURL)
	if err != nil {
		logrus.Error(err.Error())
	}
	logrus.Println(cld)

	r := gin.Default()
	v := validator.New()
	v1 := r.Group("api/v1")

	// repository
	userRepo := users.NewRepository()
	contentRepo := contents.NewRepository()
	commentRepo := comments.NewRepository()
	usersActivitiesRepo := usersactivities.NewRepository()

	// service
	userService := users.NewService(db, userRepo, cfg)
	contentService := contents.NewService(contentRepo, usersActivitiesRepo, commentRepo, cld, db)
	commentService := comments.NewService(commentRepo, db)
	usersActivitesService := usersactivities.NewService(usersActivitiesRepo, db)

	// handler
	userHandler := userhandler.NewHandler(r, userService, v)
	contentHandler := contenthandler.NewHandler(r, contentService, commentService, usersActivitesService, v)

	// register route
	userHandler.RegisterRouter(v1)
	contentHandler.RegisterRouter(v1)

	r.Run(cfg.Service.Port)
}
