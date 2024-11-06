package userrepository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/usermodel"
)

type repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

type Repository interface {
	CreateUser(ctx context.Context, model *usermodel.UserModel) error
	GetUser(ctx context.Context, userID int64, username, email string) (*usermodel.UserModel, error)
	InsertToken(ctx context.Context, model usermodel.RefreshTokenModel) error
	GetToken(ctx context.Context, userID int64, now time.Time) (*usermodel.RefreshTokenModel, error)
}
