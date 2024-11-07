package contentrepository

import (
	"context"

	"github.com/ArdiSasongko/challenge-100-personal-project/5-forum/internal/model/contentmodel"
	"github.com/sirupsen/logrus"
)

func (r *repository) CreateContent(ctx context.Context, model contentmodel.ContentModel) error {
	query := `INSERT INTO contents(user_id, content_title, content, content_hastags, created_at, updated_at, created_by, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, &model.UserID, &model.ContentTitle, &model.Content, &model.ContentHastags, &model.CreatedAt, &model.UpdatedAt, &model.CreatedBy, &model.UpdatedBy).Scan(&model.ID)

	if err != nil {
		logrus.WithField(
			"error", "error failed created content",
		).Error(err.Error())
		return err
	}

	return nil
}
