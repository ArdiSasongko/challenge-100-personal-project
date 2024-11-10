package comments

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
	InsertComment(ctx context.Context, tx *sql.Tx, model CommentModel) error
	GetCommentsByContent(ctx context.Context, tx *sql.Tx, content_id int64) (*[]CommentsResponse, error)
}

func (r *repository) InsertComment(ctx context.Context, tx *sql.Tx, model CommentModel) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO comments(user_id, content_id, comment_body, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := tx.QueryRowContext(ctx, query, &model.UserID, &model.ContentID, &model.CommentBody, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy).Scan(&model.ID)

	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return err
	}

	return nil
}

func (r *repository) GetCommentsByContent(ctx context.Context, tx *sql.Tx, content_id int64) (*[]CommentsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, user_id, content_id, comment_body, updated_at, created_by FROM comments WHERE content_id = $1`

	rows, err := tx.QueryContext(ctx, query, &content_id)
	if err != nil {
		logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var (
		comment  CommentsResponse
		comments = make([]CommentsResponse, 0)
	)
	for rows.Next() {
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.ContentID, &comment.UpdatedAt, &comment.CreatedBy)
		if err != nil {
			logrus.WithContext(ctx).WithField("error", err.Error()).Error(err.Error())
			return nil, err
		}
		comments = append(comments, CommentsResponse{
			ID:          comment.ID,
			UserID:      comment.UserID,
			ContentID:   comment.ContentID,
			CommentBody: comment.CommentBody,
			UpdatedAt:   comment.UpdatedAt,
			CreatedBy:   comment.CreatedBy,
		})
	}

	return &comments, nil
}
