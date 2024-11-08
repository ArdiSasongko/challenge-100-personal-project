package contentrepository

import (
	"context"
	"database/sql"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/contentmodel"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

type Repository interface {
	CreateContent(ctx context.Context, model contentmodel.ContentModel) error
	GetContent(ctx context.Context, contentID, userID int64) (*contentmodel.GetContentResponse, error)
	GetContents(ctx context.Context, limit, offset int64) (*contentmodel.GetContents, error)
}
