package usersactivities

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
)

type repository struct{}

func NewRepository() *repository {
	return &repository{}
}

type Repository interface {
	InsertUserActivities(ctx context.Context, tx *sql.Tx, model UserActivitiesModel) error
	GetUserActivites(ctx context.Context, tx *sql.Tx, model UserActivitiesModel) (*UserActivitiesModel, error)
	UpdatedUserActivities(ctx context.Context, tx *sql.Tx, model UserActivitiesModel) error
	CountLikes(ctx context.Context, tx *sql.Tx, contentID int64) (int, error)
}

func (r *repository) InsertUserActivities(ctx context.Context, tx *sql.Tx, model UserActivitiesModel) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO users_activities(user_id, content_id, is_liked, is_saved, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := tx.QueryRowContext(ctx, query, &model.UserID, &model.ContentID, &model.IsLiked, &model.IsSaved, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy).Scan(&model.ID)

	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return err
	}

	return nil
}

func (r *repository) GetUserActivites(ctx context.Context, tx *sql.Tx, model UserActivitiesModel) (*UserActivitiesModel, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, user_id, content_id, is_liked, is_saved, created_at, updated_at, created_at, updated_at FROM users_activities WHERE content_id = $1 AND user_id = $2`

	var response UserActivitiesModel
	row := tx.QueryRowContext(ctx, query, model.ContentID, model.UserID)
	err := row.Scan(&response.ID, &response.UserID, &response.ContentID, &response.IsLiked, &response.IsSaved, &response.CreatedAt, &response.UpdatedAt, &response.CreatedBy, &response.UpdatedBy)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return nil, err
	}

	return &response, nil
}

func (r *repository) UpdatedUserActivities(ctx context.Context, tx *sql.Tx, model UserActivitiesModel) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `UPDATE users_activities SET is_liked = $1, is_saved = $2, updated_at = $3, updated_by = $4 WHERE content_id = $5 AND user_id = $6`

	_, err := tx.ExecContext(ctx, query, model.IsLiked, model.IsSaved, model.UpdatedAt, model.UpdatedBy, model.ContentID, model.UserID)
	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return err
	}

	return nil
}

func (r *repository) CountLikes(ctx context.Context, tx *sql.Tx, contentID int64) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT COUNT(id) FROM users_activities WHERE content_id = $1 AND is_liked = TRUE`

	var likes int
	err := tx.QueryRowContext(ctx, query, contentID).Scan(&likes)

	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return 0, err
	}

	return likes, nil
}
