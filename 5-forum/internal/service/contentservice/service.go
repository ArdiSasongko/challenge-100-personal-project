package contentservice

import (
	"context"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/contentmodel"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/repository/contentrepository"
)

type service struct {
	repo contentrepository.Repository
}

func NewService(repo contentrepository.Repository) *service {
	return &service{
		repo: repo,
	}
}

type Service interface {
	CreateContent(ctx context.Context, req contentmodel.ContentRequest, userID int64, username string) error
}
