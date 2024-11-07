package contentservice

import (
	"context"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/contentmodel"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/repository/contentrepository"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/repository/userrepository"
)

type service struct {
	repo contentrepository.Repository
	ur   userrepository.Repository
}

func NewService(repo contentrepository.Repository, ur userrepository.Repository) *service {
	return &service{
		repo: repo,
		ur:   ur,
	}
}

type Service interface {
	CreateContent(ctx context.Context, req contentmodel.ContentRequest, userID int64, username string) error
	GetContent(ctx context.Context, contentId int64) (*contentmodel.GetResponse, error)
}
