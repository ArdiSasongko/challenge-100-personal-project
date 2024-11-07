package userrepository

import (
	"context"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/contentmodel"
	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/usermodel"
	"github.com/sirupsen/logrus"
)

func (r *repository) CreateComment(ctx context.Context, model *usermodel.CommentModel) error {
	query := `INSERT INTO comments(user_id, content_id, comment, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, &model.UserID, &model.ContentID, &model.Comment, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy).Scan(&model.ID)

	if err != nil {
		logrus.WithField(
			"error", "error failed created user",
		).Error(err.Error())
		return err
	}

	return nil
}

func (r *repository) GetCommentByContent(ctx context.Context, content_id int64) (*[]contentmodel.Comment, error) {
	query := `SELECT id, user_id, created_by, comment FROM comments WHERE content_id = $1`

	var (
		comments contentmodel.Comment
		response = make([]contentmodel.Comment, 0)
	)

	rows, err := r.db.QueryContext(ctx, query, content_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&comments.ID, &comments.UserID, &comments.Username, &comments.Comment)
		if err != nil {
			return nil, err
		}

		response = append(response, contentmodel.Comment{
			ID:       comments.ID,
			UserID:   comments.UserID,
			Username: comments.Username,
			Comment:  comments.Comment,
		})
	}

	return &response, nil
}
