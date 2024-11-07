package userrepository

import (
	"context"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/usermodel"
	"github.com/sirupsen/logrus"
)

func (r *repository) CreateUser(ctx context.Context, model *usermodel.UserModel) error {
	query := `INSERT INTO users(name, username, email, password, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, &model.Name, &model.Username, &model.Email, &model.Password, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy).Scan(&model.ID)

	if err != nil {
		logrus.WithField(
			"error", "error failed created user",
		).Error(err.Error())
		return err
	}

	return nil
}

func (r *repository) GetUser(ctx context.Context, userID int64, username, email string) (*usermodel.UserModel, error) {
	query := `SELECT id, name, username, email, password, created_at, updated_at, created_by, updated_by FROM users WHERE id = $1 OR username = $2 OR email = $3`

	row := r.db.QueryRowContext(ctx, query, userID, username, email)

	var user usermodel.UserModel
	err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.CreatedBy, &user.UpdatedBy)
	if err != nil {
		logrus.WithField(
			"error", "error failed get user",
		).Error(err.Error())
		return nil, err
	}

	return &user, nil
}
