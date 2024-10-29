package userrepository

import (
	"context"
	"database/sql"
	"restapi/internal/model/domain"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, tx *sql.Tx, user domain.User) (*domain.User, error)
	Update(ctx context.Context, tx *sql.Tx, user domain.User, userID int) (*domain.User, error)
	Delete(ctx context.Context, tx *sql.Tx, userID int) error
	FindByID(ctx context.Context, tx *sql.Tx, userID int) (*domain.User, error)
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*domain.User, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]*domain.User, error)
}
