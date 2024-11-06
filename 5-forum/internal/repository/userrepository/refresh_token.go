package userrepository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/usermodel"
	"github.com/sirupsen/logrus"
)

func (r *repository) InsertToken(ctx context.Context, model usermodel.RefreshTokenModel) error {
	query := `INSERT INTO refresh_token(user_id, refresh_token, expired_at, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, &model.UserID, &model.RefreshToken, &model.ExpiredAt, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.CreatedBy).Scan(&model.ID)

	if err != nil {
		logrus.WithField(
			"error", "error failed insert refresh token",
		).Error(err.Error())
		return err
	}

	return nil
}

func (r *repository) GetToken(ctx context.Context, userID int64, now time.Time) (*usermodel.RefreshTokenModel, error) {
	query := `SELECT id, user_id, refresh_token, expired_at, created_at, updated_at, created_by, updated_by FROM refresh_token WHERE user_id = $1 AND expired_at >= $2`

	var token usermodel.RefreshTokenModel
	row := r.db.QueryRowContext(ctx, query, userID, now)
	err := row.Scan(&token.ID, &token.UserID, &token.RefreshToken, &token.ExpiredAt, &token.CreatedAt, &token.UpdatedAt, &token.CreatedBy, &token.UpdatedBy)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithField("error", "failed scan row token").Error(err.Error())
	}

	return &token, nil
}
