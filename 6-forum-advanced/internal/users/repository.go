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
	InsertToken(ctx context.Context, tx *sql.Tx, model RefreshTokenModel) error
	GetToken(ctx context.Context, tx *sql.Tx, userID int64, now time.Time) (*RefreshTokenModel, error)
	UpdateToken(ctx context.Context, tx *sql.Tx, model RefreshTokenModel) error
}

func (r *repository) InsertUser(ctx context.Context, tx *sql.Tx, model UserModel) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO users(name, username, email, password, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := tx.QueryRowContext(ctx, query, &model.Name, &model.Username, &model.Email, &model.Password, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy).Scan(&model.ID)

	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return err
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

func (r *repository) InsertToken(ctx context.Context, tx *sql.Tx, model RefreshTokenModel) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO refresh_token(user_id, token, expired_at, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := tx.QueryRowContext(ctx, query, &model.UserID, &model.Token, &model.ExpiredAt, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy).Scan(&model.ID)

	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return err
	}

	return nil
}

func (r *repository) GetToken(ctx context.Context, tx *sql.Tx, userID int64, now time.Time) (*RefreshTokenModel, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, user_id, token, expired_at, created_at, updated_at, created_by, updated_by FROM refresh_token WHERE user_id = $1 AND expired_at >= $2`

	row := tx.QueryRowContext(ctx, query, userID, now)

	var model RefreshTokenModel
	err := row.Scan(&model.ID, &model.UserID, &model.Token, &model.ExpiredAt, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return nil, err
	}

	return &model, nil
}

func (r *repository) UpdateToken(ctx context.Context, tx *sql.Tx, model RefreshTokenModel) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `UPDATE refresh_token SET token = $1, updated_at = $2, updated_by = $3 WHERE user_id = $4 AND expired_at <= $5`

	_, err := tx.ExecContext(ctx, query, model.Token, model.UpdatedAt, model.UpdatedBy, model.UserID)

	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return err
	}

	return nil
}
