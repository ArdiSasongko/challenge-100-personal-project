package userservice

import (
	"context"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/configs"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/usermodel"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/repository/contentrepository"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/repository/userrepository"
)

type service struct {
	repo userrepository.Repository
	cr   contentrepository.Repository
	cfg  *configs.Config
}

func NewUserService(repo userrepository.Repository, cr contentrepository.Repository, cfg *configs.Config) *service {
	return &service{
		repo: repo,
		cfg:  cfg,
		cr:   cr,
	}
}

type Service interface {
	CreateUser(ctx context.Context, req usermodel.UserRequest) error
	LoginUser(ctx context.Context, req usermodel.LoginRequest) (*usermodel.LoginResponse, error)
	CreateComment(ctx context.Context, req usermodel.CommentRequest, contentID, userID int64, username string) error
}
