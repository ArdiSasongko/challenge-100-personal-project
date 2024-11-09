package users

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
)

type repository struct {
}

func NewRepository() *repository {
	return &repository{}
}

type Repository interface {
	InsertUser(ctx context.Context, tx *sql.Tx, model UserModel) error
	GetUser(ctx context.Context, tx *sql.Tx, id int64, username, email string) (*UserModel, error)
}

func (r *repository) InsertUser(ctx context.Context, tx *sql.Tx, model UserModel) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO users(name, username, email, password, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := tx.QueryRowContext(ctx, query, &model.Name, &model.Username, &model.Email, &model.Password, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy).Scan(&model.ID)

	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return nil
	}

	return nil
}

func (r *repository) GetUser(ctx context.Context, tx *sql.Tx, id int64, username, email string) (*UserModel, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, name, username, email, password, created_at, updated_at, created_by, updated_by FROM users WHERE id = $1 OR username = $2 OR email = $3`

	row := tx.QueryRowContext(ctx, query, id, username, email)
	var model UserModel

	err := row.Scan(&model.ID, &model.Name, &model.Username, &model.Email, &model.Password, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return nil, err
	}

	return &model, nil
}
