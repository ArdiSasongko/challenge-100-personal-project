package userrepository

import (
	"context"
	"database/sql"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/usermodel"
)

func (r *repository) InsertLikeAndSave(ctx context.Context, model usermodel.UserActivities) error {
	query := `INSERT INTO users_activities(user_id, content_id, is_liked, is_saved, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, &model.UserID, &model.ContentID, &model.IsLiked, &model.IsSaved, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy).Scan(&model.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetUserActivity(ctx context.Context, model usermodel.UserActivities) (*usermodel.UserActivities, error) {
	query := `SELECT id, user_id, content_id, is_liked, is_saved, created_at, updated_at, created_by, updated_by FROM users_activities WHERE content_id = $1 AND user_id = $2`

	var response usermodel.UserActivities
	row := r.db.QueryRowContext(ctx, query, model.ContentID, model.UserID)

	err := row.Scan(&response.ID, &response.UserID, &response.ContentID, &response.IsLiked, &response.IsSaved, &response.CreatedAt, &response.UpdatedAt, &response.CreatedBy, &response.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &response, nil
}

func (r *repository) UpdateUserActivities(ctx context.Context, model usermodel.UserActivities) error {
	query := `UPDATE users_activities SET is_liked = $1, is_saved = $2, updated_at = $3, updated_by = $4 WHERE content_id = $5 AND user_id = $6`

	_, err := r.db.ExecContext(ctx, query, model.IsLiked, model.IsSaved, model.UpdatedAt, model.UpdatedBy, model.ContentID, model.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CountLikes(ctx context.Context, content_id int64) (int, error) {
	query := `SELECT COUNT(id) FROM users_activities WHERE content_id = $1 AND is_liked = TRUE`
	var result int

	row := r.db.QueryRowContext(ctx, query, content_id)
	err := row.Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}
