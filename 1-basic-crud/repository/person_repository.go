package repository

import (
	"basic-crud/model/domain"
	"context"
	"database/sql"
)

type PersonRepositoryInterface interface {
	Create(ctx context.Context, tx *sql.Tx, person domain.Person) (*domain.Person, error)
	Update(ctx context.Context, tx *sql.Tx, person domain.Person) (*domain.Person, error)
	Delete(ctx context.Context, tx *sql.Tx, personID int) error
	FindByID(ctx context.Context, tx *sql.Tx, personID int) (*domain.Person, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Person, error)
}
