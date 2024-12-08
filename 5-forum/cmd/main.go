package main

import (
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/configs"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/handler/contenthandler"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/handler/userhandler"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/repository/contentrepository"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/repository/userrepository"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/service/contentservice"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/service/userservice"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/pkg/internalsql"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		logrus.WithField("error", err).Error(err.Error())
	}

	cfg = configs.GetConfig()

	db, err := internalsql.Database(cfg.Database.DataSource)
	if err != nil {
		logrus.Error(err.Error())
	}

	internalsql.InitStorage(db)

	r := gin.Default()
	v1 := r.Group("/v1/api")
	validate := validator.New()

	contentRepo := contentrepository.NewRepository(db)
	userrepo := userrepository.NewUserRepository(db)

	contentService := contentservice.NewService(contentRepo, userrepo)
	userservice := userservice.NewUserService(userrepo, contentRepo, cfg)

	contentHandler := contenthandler.NewHandler(r, contentService, userservice, validate)
	userhandler := userhandler.NewHandler(r, userservice, validate)

	contentHandler.RegisterRouter(v1)
	userhandler.RegisterRoutes(v1)

	r.Run(cfg.Service.Port)
}
